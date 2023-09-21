//go:build wireinject
// +build wireinject

package middlewares

import (
	"example/internal/app/service"
	"example/internal/infrastructure/repos/mysql_repos_impl"
	"example/internal/infrastructure/repos/redis_repos_impl"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/google/wire"
	"net/http"
)

// NewHandleAuthVerify 创建
func NewHandleAuthVerify(dc *db.Connector, rdb *db.Redis) func(handler http.Handler) http.Handler {
	panic(wire.Build(redis_repos_impl.NewSessionRepository, service.NewSessionService, HandleAuthVerify))
}

// NewHandlePermissionVerify
func NewHandlePermissionVerify(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) func(handler http.Handler) http.Handler {
	panic(wire.Build(
		mysql_repos_impl.NewAccountRepository,
		service.NewPermissionService,
		HandlePermissionVerify,
	))
}
