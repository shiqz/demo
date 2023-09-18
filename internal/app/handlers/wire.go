//go:build wireinject
// +build wireinject

package handlers

import (
	"demo/internal/app/service"
	"demo/internal/infrastructure/repos/mysql_repos_impl"
	"demo/internal/infrastructure/repos/redis_repos_impl"
	"demo/internal/pkg/db"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// NewUserAPI 实例化用户控制器
func NewUserAPI(dc *db.Connector, rdb *redis.Client) *UserHandler {
	panic(wire.Build(
		mysql_repos_impl.NewUserRepository,
		redis_repos_impl.NewSessionRepository,
		service.NewUserService,
		service.NewSessionService,
		NewUserHandler,
	))
}

// NewAccountAPI 实例化管理员账户控制器
func NewAccountAPI(dc *db.Connector, rdb *redis.Client) *AccountHandler {
	panic(wire.Build(mysql_repos_impl.NewAccountRepository,
		redis_repos_impl.NewSessionRepository,
		service.NewAccountService,
		service.NewSessionService,
		NewAccountHandler,
	))
}
