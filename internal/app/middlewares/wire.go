//go:build wireinject
// +build wireinject

package middlewares

import (
	"demo/internal/app/service"
	"demo/internal/infrastructure/repos/redis_repos_impl"
	"demo/internal/pkg/db"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"net/http"
)

// NewHandleAuthVerify 创建
func NewHandleAuthVerify(dc *db.Connector, rdb *redis.Client) func(handler http.Handler) http.Handler {
	panic(wire.Build(redis_repos_impl.NewSessionRepository, service.NewSessionService, HandleAuthVerify))
}
