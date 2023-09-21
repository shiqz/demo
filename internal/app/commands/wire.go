//go:build wireinject
// +build wireinject

package commands

import (
	"example/internal/app/service"
	"example/internal/infrastructure/repos/mysql_repos_impl"
	"example/internal/infrastructure/repos/redis_repos_impl"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/google/wire"
)

// NewAccountCommand 实例化账号管理命令
func NewAccountCommand(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *AccountHandler {
	panic(wire.Build(
		mysql_repos_impl.NewAccountRepository,
		service.NewAccountService,
		NewAccountHandler,
	))
}

// NewUserCommand 实例化账号管理命令
func NewUserCommand(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *UserHandler {
	panic(wire.Build(
		redis_repos_impl.NewSessionRepository,
		service.NewSessionService,
		NewUserHandler,
	))
}
