// Package services 服务实现
package service

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
func (s *UserService) Login(ctx context.Context, username, pass string) (*entity.Session, error) {
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
	return entity.NewSession(types.UserSession, ug.User.UserID), nil
}

// GetUserinfo 获取用户信息
func (s *UserService) GetUserinfo(ctx context.Context) (*domain.UserAggregate, error) {
	session := ctx.Value(types.SessionFlag).(*entity.Session)
	return s.repo.GetOne(ctx, session.GetSessionID())
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(ctx context.Context, pass string) error {
	session := ctx.Value(types.SessionFlag).(*entity.Session)
	return s.repo.UpdatePass(ctx, session.GetSessionID(), pass)
}

// Users 用户列表
func (s *UserService) Users(ctx context.Context) ([]*domain.UserAggregate, error) {

	return nil, nil
}
