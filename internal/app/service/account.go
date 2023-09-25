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

// AccountService 账号服务
type AccountService struct {
	repo domain.AccountRepository
}

// NewAccountService 实例化服务
func NewAccountService(repo domain.AccountRepository) domain.AccountService {
	return &AccountService{repo: repo}
}

// Login 账号登录
func (s *AccountService) Login(ctx context.Context, email, pass string) (*entity.Session, error) {
	account, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.EcInvalidUser
		}
		return nil, err
	}
	if !account.Password.Valid(pass) {
		return nil, errs.EcInvalidUser
	}
	return entity.NewSession(types.AdminSession, account.AdminID), nil
}

// Create 创建账号
func (s *AccountService) Create(ctx context.Context, account *entity.Account) error {
	ac, err := s.repo.GetAccountByEmail(ctx, account.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if ac != nil {
		return errs.EcUserHasBeenExist
	}
	return s.repo.Save(ctx, account)
}

// UpdateRoleByEmail 根据邮箱修改角色
func (s *AccountService) UpdateRoleByEmail(ctx context.Context, email string, roles types.Roles) error {
	info, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrAccountNotFound
		}
		return err
	}
	if info.Roles.Eq(roles) {
		return nil
	}
	info.SetRole(roles)
	return s.repo.UpdateRole(ctx, info)
}

// UpdatePassByEmail 根据邮箱修改密码
func (s *AccountService) UpdatePassByEmail(ctx context.Context, email, pass string) error {
	info, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrAccountNotFound
		}
		return err
	}
	if err = info.SetPassword(pass); err != nil {
		return err
	}
	return s.repo.UpdatePass(ctx, info)
}

// GetAccounts 查询列表
func (s *AccountService) GetAccounts(ctx context.Context, filter *domain.AccountFilter) ([]*entity.Account, error) {
	return s.repo.Accounts(ctx, filter)
}

// GetAccountByEmail 通过邮箱获取账号信息
func (s *AccountService) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	account, err := s.repo.GetAccountByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}
	return account, nil
}
