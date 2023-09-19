//go:build wireinject
// +build wireinject

package middlewares

import (
	"demo/internal/app/service"
	"demo/internal/infrastructure/repos/redis_repos_impl"
	"demo/internal/pkg/db"
	"github.com/google/wire"
	"net/http"
)

// NewHandleAuthVerify 创建
func NewHandleAuthVerify(dc *db.Connector, rdb *db.Redis) func(handler http.Handler) http.Handler {
	panic(wire.Build(redis_repos_impl.NewSessionRepository, service.NewSessionService, HandleAuthVerify))
}
