// Package services 服务实现
package services

import (
	"context"
	"database/sql"
	"demo/internal/app/errs"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"github.com/pkg/errors"
)

// UserService 用户服务
type UserService struct {
	repo domain.UserRepository
}

// NewUserService 实例化用户服务
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &UserService{repo: repo}
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, ug *domain.UserAggregate) error {
	existUg, err := s.repo.GetUserByUsername(ctx, ug.User.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existUg != nil {
		return errs.EcUserHasBeenExist
	}
	return s.repo.Save(ctx, ug)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, pass string) (*domain.UserAggregate, error) {
	ug, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.EcInvalidUser
		}
		return nil, err
	}
	if !ug.User.IsValidPassword(pass) {
		return nil, errs.EcInvalidUser
	}
	ug.Session = entity.NewSession(types.UserSession, ug.User.UserID)
	if err = s.repo.SetSession(ctx, ug); err != nil {
		return nil, err
	}
	return ug, nil
}
