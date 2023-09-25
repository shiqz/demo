//go:build wireinject
// +build wireinject

package middlewares

import (
	"example/internal/app/service"
	"example/internal/infrastructure/depend"
	"example/internal/infrastructure/repos/mysqlreposimpl"
	"example/internal/infrastructure/repos/redisrepoimpl"
	"github.com/google/wire"
	"net/http"
)

// NewHandleAuthVerify 创建
func NewHandleAuthVerify(inject *depend.Injecter) func(handler http.Handler) http.Handler {
	panic(wire.Build(redisrepoimpl.NewSessionRepository, service.NewSessionService, HandleAuthVerify))
}

// NewHandlePermissionVerify
func NewHandlePermissionVerify(inject *depend.Injecter) func(handler http.Handler) http.Handler {
	panic(wire.Build(
		mysqlreposimpl.NewAccountRepository,
		service.NewPermissionService,
		HandlePermissionVerify,
	))
}
