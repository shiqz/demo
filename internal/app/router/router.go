// Package router API相关路由
package router

import (
	"example/internal/app/errs"
	"example/internal/app/handlers"
	"example/internal/app/middlewares"
	"example/internal/app/response"
	"example/internal/app/router/admin"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// HandleNotFound 响应不存在路由
func HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	response.Error(w, errs.New(errs.EcNotFound, errs.ErrAPINotFound))
}

// Init 初始化路由
func Init(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middlewares.HandleLogger(lg), middlewares.HandleRecover, middlewares.HandleCors)
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
