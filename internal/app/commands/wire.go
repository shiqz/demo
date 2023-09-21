//go:build wireinject
// +build wireinject

package commands

import (
	"example/internal/app/service"
	"example/internal/infrastructure/repos/mysqlreposimpl"
	"example/internal/infrastructure/repos/redisrepoimpl"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/google/wire"
)

// NewAccountCommand 实例化账号管理命令
func NewAccountCommand(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *AccountHandler {
	panic(wire.Build(
		mysqlreposimpl.NewAccountRepository,
		service.NewAccountService,
		NewAccountHandler,
	))
}

// NewUserCommand 实例化账号管理命令
func NewUserCommand(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *UserHandler {
	panic(wire.Build(
		redisrepoimpl.NewSessionRepository,
		service.NewSessionService,
		NewUserHandler,
	))
}
