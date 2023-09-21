//go:build wireinject
// +build wireinject

package middlewares

import (
	"example/internal/app/service"
	"example/internal/infrastructure/repos/mysqlreposimpl"
	"example/internal/infrastructure/repos/redisrepoimpl"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/google/wire"
	"net/http"
)

// NewHandleAuthVerify 创建
func NewHandleAuthVerify(dc *db.Connector, rdb *db.Redis) func(handler http.Handler) http.Handler {
	panic(wire.Build(redisrepoimpl.NewSessionRepository, service.NewSessionService, HandleAuthVerify))
}

// NewHandlePermissionVerify
func NewHandlePermissionVerify(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) func(handler http.Handler) http.Handler {
	panic(wire.Build(
		mysqlreposimpl.NewAccountRepository,
		service.NewPermissionService,
		HandlePermissionVerify,
	))
}
