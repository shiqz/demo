// Package services 服务实现
package service

import (
	"context"
	"database/sql"
	"example/internal/app/errs"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
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
func (s *UserService) Create(ctx context.Context, user *entity.User) error {
	existUg, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existUg != nil {
		return errs.EcUserHasBeenExist
	}
	return s.repo.Save(ctx, user)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, pass string) (*entity.Session, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.EcInvalidUser
		}
		return nil, err
	}
	if err = user.ValidPassword(pass); err != nil {
		return nil, err
	}
	if err = user.ValidState(); err != nil {
		return nil, err
	}
	return entity.NewSession(types.UserSession, user.UserID), nil
}

// GetUserinfo 获取用户信息
func (s *UserService) GetUserinfo(ctx context.Context) (*entity.User, error) {
	session := ctx.Value(types.SessionFlag).(*entity.Session)
	return s.repo.GetOne(ctx, session.GetSessionID())
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(ctx context.Context, id uint, pass string) error {
	user, err := s.repo.GetOne(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.New(errs.EcNotFound, domain.ErrUserNotFound)
		}
		return err
	}
	user.SetPassword(pass)
	return s.repo.UpdatePass(ctx, user)
}

// Users 用户列表
func (s *UserService) Users(ctx context.Context, filter *domain.UserFilter) ([]*entity.User, error) {
	return s.repo.Users(ctx, filter)
}

// UpdateStatus 更新用户密码
func (s *UserService) UpdateStatus(ctx context.Context, id uint, status types.UserState) error {
	user, err := s.repo.GetOne(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.New(errs.EcNotFound, domain.ErrUserNotFound)
		}
		return err
	}
	if user.Status == status {
		return nil
	}
	user.Status = status
	return s.repo.UpdateStatus(ctx, user)
}
