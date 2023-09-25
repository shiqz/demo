//go:build wireinject
// +build wireinject

package commands

import (
	"example/internal/app/service"
	"example/internal/infrastructure/depend"
	"example/internal/infrastructure/repos/mysqlreposimpl"
	"example/internal/infrastructure/repos/redisrepoimpl"
	"github.com/google/wire"
)

// NewAccountCommand 实例化账号管理命令
func NewAccountCommand(inject *depend.Injecter) *AccountHandler {
	panic(wire.Build(
		mysqlreposimpl.NewAccountRepository,
		service.NewAccountService,
		NewAccountHandler,
	))
}

// NewUserCommand 实例化账号管理命令
func NewUserCommand(inject *depend.Injecter) *UserHandler {
	panic(wire.Build(
		redisrepoimpl.NewSessionRepository,
		service.NewSessionService,
		NewUserHandler,
	))
}
