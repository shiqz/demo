// Package domain 领域层
package domain

import (
	"context"
	"demo/internal/domain/entity"
	"net/http"
)

type (
	// UserAggregate 用户聚合根
	UserAggregate struct {
		User    *entity.User
		Session *entity.Session
	}

	// UserRepository 用户仓库接口
	UserRepository interface {
		GetUserByUsername(ctx context.Context, uname string) (*UserAggregate, error)
		Save(ctx context.Context, user *UserAggregate) error
		SetSession(ctx context.Context, ug *UserAggregate) error
		GetOne(ctx context.Context, id uint) (*UserAggregate, error)
	}

	// UserService 用户服务接口
	UserService interface {
		Create(ctx context.Context, user *UserAggregate) error
		Login(ctx context.Context, username, pass string) (*UserAggregate, error)
	}

	// UserHandler 控制接口
	UserHandler interface {
		HandleRegister(w http.ResponseWriter, r *http.Request)
		HandleLogin(w http.ResponseWriter, r *http.Request)
	}
)
