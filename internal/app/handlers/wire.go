//go:build wireinject
// +build wireinject

package handlers

import (
	"example/internal/app/service"
	"example/internal/infrastructure/depend"
	"example/internal/infrastructure/repos/mysqlreposimpl"
	"example/internal/infrastructure/repos/redisrepoimpl"
	"github.com/google/wire"
)

// NewUserAPI 实例化用户控制器
func NewUserAPI(inject *depend.Injecter) *UserHandler {
	panic(wire.Build(
		mysqlreposimpl.NewUserRepository,
		redisrepoimpl.NewSessionRepository,
		service.NewUserService,
		service.NewSessionService,
		NewUserHandler,
	))
}

// NewAccountAPI 实例化管理员账户控制器
func NewAccountAPI(inject *depend.Injecter) *AccountHandler {
	panic(wire.Build(
		mysqlreposimpl.NewAccountRepository,
		redisrepoimpl.NewSessionRepository,
		service.NewAccountService,
		service.NewSessionService,
		NewAccountHandler,
	))
}
