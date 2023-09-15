// Package router API相关路由
package router

import (
	"demo/internal/app/controller"
	"demo/internal/app/middlewares"
	"demo/internal/app/router/admin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Init 初始化路由
func Init(db *sqlx.DB, rdb *redis.Client) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Recoverer, middlewares.HandleLogger, middlewares.HandleCors)

	mux.Route("/api/admin", admin.InitRoute)

	// 用户接口
	UserAPI := controller.NewUserAPI(db, rdb)
	mux.Post("/register", UserAPI.HandleRegister)
	mux.Post("/login", UserAPI.HandleLogin)
	mux.Route("/", func(auth chi.Router) {
		auth.Use()
	})
	return mux
}
