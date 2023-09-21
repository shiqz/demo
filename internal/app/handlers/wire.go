//go:build wireinject
// +build wireinject

package handlers

import (
	"example/internal/app/service"
	"example/internal/infrastructure/repos/mysqlreposimpl"
	"example/internal/infrastructure/repos/redisrepoimpl"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
	"github.com/google/wire"
)

// NewUserAPI 实例化用户控制器
func NewUserAPI(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *UserHandler {
	panic(wire.Build(
		mysqlreposimpl.NewUserRepository,
		redisrepoimpl.NewSessionRepository,
		service.NewUserService,
		service.NewSessionService,
		NewUserHandler,
	))
}

// NewAccountAPI 实例化管理员账户控制器
func NewAccountAPI(dc *db.Connector, rdb *db.Redis, lg *logger.Logger) *AccountHandler {
	panic(wire.Build(
		mysqlreposimpl.NewAccountRepository,
		redisrepoimpl.NewSessionRepository,
		service.NewAccountService,
		service.NewSessionService,
		NewAccountHandler,
	))
}
