// Package admin 管理后台路由
package admin

import (
	"demo/internal/app/handlers"
	"demo/internal/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

// InitRoute 初始化管理后台路由
func InitRoute(db *db.Connector, rdb *redis.Client, route chi.Router) {
	AccountAPI := handlers.NewAccountAPI(db, rdb)
	route.Post("/login", AccountAPI.HandleLogin)
}
