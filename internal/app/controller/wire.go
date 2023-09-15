//go:build wireinject
// +build wireinject

package controller

import (
	"demo/internal/app/controller/handlers"
	"demo/internal/app/services"
	"demo/internal/domain"
	"demo/internal/infrastructure/repos"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// NewUserAPI 实例化用户控制器
func NewUserAPI(db *sqlx.DB, rdb *redis.Client) domain.UserHandler {
	panic(wire.Build(
		repos.NewUserRepository,
		services.NewUserService,
		handlers.NewUserHandler,
	))
}
