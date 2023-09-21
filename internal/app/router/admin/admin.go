// Package admin 管理后台路由
package admin

import (
	"example/internal/app/handlers"
	"example/internal/app/middlewares"
	"example/internal/app/response"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// InitRoute 初始化管理后台路由
func InitRoute(dc *db.Connector, rdb *db.Redis, lg *logger.Logger, route chi.Router) {
	AccountAPI := handlers.NewAccountAPI(dc, rdb, lg)
	route.Post("/login", AccountAPI.HandleLogin)
	route.Route("/", func(auth chi.Router) {
		auth.Use(middlewares.NewHandleAuthVerify(dc, rdb))
		auth.Delete("/login", AccountAPI.HandleLogout)
		auth.Route("/", func(perm chi.Router) {
			perm.Use(middlewares.NewHandlePermissionVerify(dc, rdb, lg))
			UserAPI := handlers.NewUserAPI(dc, rdb, lg)
			{
				perm.Get("/users", UserAPI.HandleUsers)
				perm.Patch("/users/status", UserAPI.ChangeUserStatus)
				perm.Patch("/users/passwd", UserAPI.ResetUserPass)
			}
		})
	})
	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "hello", nil)
	})
}
