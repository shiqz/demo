//go:build wireinject
// +build wireinject

package commands

import (
	"demo/internal/app/service"
	"demo/internal/infrastructure/repos/mysql_repos_impl"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
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
