// Package domain 领域层
package domain

import (
	"context"
	"example/internal/domain/entity"
	"example/internal/domain/types"
)

const (
	DefaultUserPass = "123@Joyparty"
)

type (
	// UserFilter 筛选结构
	UserFilter struct {
		Filter
		UserID   *uint
		Nickname *string
		Gender   *types.UserGender
		Status   *types.UserState
	}

	// UserService 用户服务接口
	UserService interface {
		Login(ctx context.Context, username, pass string) (*entity.Session, error)
		GetUserinfo(ctx context.Context) (*entity.User, error)
		Users(ctx context.Context, filter *UserFilter) ([]*entity.User, error)

		Create(ctx context.Context, user *entity.User) error
		UpdatePassword(ctx context.Context, id uint, pass string) error
		UpdateStatus(ctx context.Context, id uint, status types.UserState) error
	}

	// UserRepository 用户仓库接口
	UserRepository interface {
		GetUserByUsername(ctx context.Context, uname string) (*entity.User, error)
		GetOne(ctx context.Context, id uint) (*entity.User, error)
		Users(ctx context.Context, filter *UserFilter) ([]*entity.User, error)

		Save(ctx context.Context, user *entity.User) error
		UpdatePass(ctx context.Context, user *entity.User) error
		UpdateStatus(ctx context.Context, user *entity.User) error
	}
)
