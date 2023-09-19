// Package admin 管理后台路由
package admin

import (
	"demo/internal/app/handlers"
	"demo/internal/app/middlewares"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// InitRoute 初始化管理后台路由
func InitRoute(dc *db.Connector, rdb *db.Redis, lg *logger.Logger, route chi.Router) {
	AccountAPI := handlers.NewAccountAPI(dc, rdb, lg)
	UserAPI := handlers.NewUserAPI(dc, rdb, lg)
	route.Post("/login", AccountAPI.HandleLogin)
	route.Route("/", func(auth chi.Router) {
		auth.Use(middlewares.NewHandleAuthVerify(dc, rdb))
		auth.Delete("/login", AccountAPI.HandleLogout)
		auth.Get("/users", UserAPI.HandleUsers)
	})
}
