// Package domain 领域层
package domain

import (
	"context"
	"demo/internal/domain/entity"
)

type (
	// UserAggregate 用户聚合根
	UserAggregate struct {
		User *entity.User
	}

	// UserBaseService 用户服务接口
	UserBaseService interface {
		Create(ctx context.Context, user *UserAggregate) error
	}

	// UserService 用户服务接口
	UserService interface {
		UserBaseService
		Login(ctx context.Context, username, pass string) (*entity.Session, error)
		GetUserinfo(ctx context.Context) (*UserAggregate, error)
		UpdatePassword(ctx context.Context, pass string) error
		Users(ctx context.Context) ([]*UserAggregate, error)
	}

	// UserRepository 用户仓库接口
	UserRepository interface {
		GetUserByUsername(ctx context.Context, uname string) (*UserAggregate, error)
		Save(ctx context.Context, user *UserAggregate) error
		GetOne(ctx context.Context, id uint) (*UserAggregate, error)
		UpdatePass(ctx context.Context, id uint, pass string) error
	}
)
