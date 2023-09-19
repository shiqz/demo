//go:build wireinject
// +build wireinject

package handlers

import (
	"demo/internal/app/service"
	"demo/internal/infrastructure/repos/mysql_repos_impl"
	"demo/internal/infrastructure/repos/redis_repos_impl"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	"github.com/google/wire"
)

// NewUserAPI 实例化用户控制器
func NewUserAPI(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *UserHandler {
	panic(wire.Build(
		mysql_repos_impl.NewUserRepository,
		redis_repos_impl.NewSessionRepository,
		service.NewUserService,
		service.NewSessionService,
		NewUserHandler,
	))
}

// NewAccountAPI 实例化管理员账户控制器
func NewAccountAPI(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *AccountHandler {
	panic(wire.Build(
		mysql_repos_impl.NewAccountRepository,
		redis_repos_impl.NewSessionRepository,
		service.NewAccountService,
		service.NewSessionService,
		NewAccountHandler,
	))
}
