// Package domain 领域层
package domain

import (
	"context"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
)

type (
	// UserFilter 筛选结构
	UserFilter struct {
		UserID   *uint
		Nickname *string
		Gender   *types.UserGender
		Status   *types.UserState
	}

	// UserBaseService 用户服务接口
	UserBaseService interface {
		Create(ctx context.Context, user *entity.User) error
	}

	// UserService 用户服务接口
	UserService interface {
		UserBaseService
		Login(ctx context.Context, username, pass string) (*entity.Session, error)
		GetUserinfo(ctx context.Context) (*entity.User, error)
		UpdatePassword(ctx context.Context, pass string) error
		Users(ctx context.Context, filter *UserFilter) ([]*entity.User, error)
	}

	// UserRepository 用户仓库接口
	UserRepository interface {
		GetUserByUsername(ctx context.Context, uname string) (*entity.User, error)
		Save(ctx context.Context, user *entity.User) error
		GetOne(ctx context.Context, id uint) (*entity.User, error)
		UpdatePass(ctx context.Context, id uint, pass string) error
	}
)
