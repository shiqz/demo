package middlewares

import (
	"context"
	"database/sql"
	"demo/internal/app/controller/response"
	"demo/internal/app/errs"
	"demo/internal/domain"
	"demo/internal/domain/types"
	"demo/internal/pkg/logger"
	"demo/internal/pkg/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// HandleCors 处理跨域
func HandleCors(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}

// HandleLogger 日志中间件
func HandleLogger(next http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  logger.New("Trace"),
		NoColor: false,
	})(next)
}

// HandleAuthVerify 登录校验中间件
func HandleAuthVerify(repo domain.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.GetRequestToken(r.Header.Get("Authorization"))
			if err != nil {
				response.Error(w, errs.EcUnauthorized)
				return
			}
			scene, uid, err := types.ParseToken(token)
			if err != nil {
				response.Error(w, errs.EcUnauthorized)
				return
			}
			if scene == types.UserSession {
				var ug *domain.UserAggregate
				ug, err = repo.GetOne(r.Context(), uid)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) || errors.Is(err, redis.Nil) {
						err = errs.EcUnauthorized
					}
					response.Error(w, err)
					return
				}
				if ug.Session == nil || ug.Session.IsExpired() {
					response.Error(w, errs.EcUnauthorized)
					return
				}
				r.WithContext(context.WithValue(r.Context(), types.SessionFlag, ug))
			}
			next.ServeHTTP(w, r)
		})
	}
}
