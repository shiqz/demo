// Package middlewares API中间件
package middlewares

import (
	"context"
	"example/internal/app/errs"
	"example/internal/app/response"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"example/internal/pkg/logger"
	"example/internal/pkg/utils"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net/http"
	"runtime/debug"
	"time"
)

// HandleRecover 奔溃处理
func HandleRecover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				if r.Header.Get("Connection") != "Upgrade" {
					response.Error(w, fmt.Errorf("%v", rvr))
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// HandleCors 处理跨域
func HandleCors(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}

// HandleLogger 日志中间件
func HandleLogger(lg *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		lf := &middleware.DefaultLogFormatter{
			Logger:  lg,
			NoColor: false,
		}

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := lf.NewLogEntry(r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			res := ww.Unwrap().(*response.Responser)
			t1 := time.Now()
			defer func() {
				entry.Write(res.Status(), res.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()
			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

// HandleAuthVerify 登录校验中间件
func HandleAuthVerify(srv domain.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.GetRequestToken(r.Header.Get("Authorization"))
			if err != nil {
				response.Error(w, errs.EcUnauthorized)
				return
			}
			ts, err := entity.LoadSessionByToken(token)
			if err != nil {
				response.Error(w, errs.EcUnauthorized)
				return
			}
			session, err := srv.Get(r.Context(), ts.GetScene(), ts.GetSessionID())
			if err != nil {
				if errors.Is(err, redis.Nil) {
					err = errs.EcUnauthorized
				}
				response.Error(w, err)
				return
			}
			if session == nil || session.Token != ts.Token || session.IsExpired() {
				response.Error(w, errs.EcUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), types.SessionFlag, session))
			next.ServeHTTP(w, r)
			if session.IsExpireSoon() {
				if err = srv.Refresh(r.Context(), session); err != nil {
					response.Error(w, err)
					return
				}
				w.Header().Add("Refresh-Token", session.FormatToken())
			}
		})
	}
}

// HandlePermissionVerify 路由权限校验中间件
func HandlePermissionVerify(srv domain.PermissionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := srv.CheckPermission(r.Context(), types.Route{Method: r.Method, Path: r.URL.Path}); err != nil {
				response.Error(w, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
