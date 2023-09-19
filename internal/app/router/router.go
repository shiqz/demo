// Package router API相关路由
package router

import (
	"demo/internal/app/errs"
	"demo/internal/app/handlers"
	"demo/internal/app/middlewares"
	"demo/internal/app/response"
	"demo/internal/app/router/admin"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// HandleNotFound 响应不存在路由
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	response.Error(w, errs.New(errs.EcNotFound, errs.ErrAPINotFound))
}

// Init 初始化路由
func Init(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middlewares.HandleLogger(lg), middlewares.HandleFinal, middlewares.HandleRecover, middlewares.HandleCors)
	mux.NotFound(HandleNotFound)
	mux.MethodNotAllowed(HandleNotFound)
	mux.Route("/api/admin", func(r chi.Router) {
		admin.InitRoute(dc, rdb, lg, r)
	})

	// 用户接口
	UserAPI := handlers.NewUserAPI(dc, rdb, lg)
	mux.Post("/register", UserAPI.HandleRegister)
	mux.Post("/login", UserAPI.HandleLogin)
	mux.Route("/", func(auth chi.Router) {
		auth.Use(middlewares.NewHandleAuthVerify(dc, rdb))
		auth.Delete("/login", UserAPI.HandleLogout)
		auth.Get("/my/identity", UserAPI.HandleIdentity)
		auth.Put("/my/password", UserAPI.HandleChangePass)
	})
	return mux
}
