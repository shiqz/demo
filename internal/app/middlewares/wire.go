//go:build wireinject
// +build wireinject

package middlewares

import (
	"demo/internal/infrastructure/repos"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// NewHandleAuthVerify 创建
func NewHandleAuthVerify(db *sqlx.DB, rdb *redis.Client) {
	panic(wire.Build(repos.NewUserRepository, HandleAuthVerify))
}
