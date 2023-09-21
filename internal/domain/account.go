package domain

import (
	"context"
	"example/internal/domain/entity"
	"example/internal/domain/types"
)

// 管理员账号领域核心业务
type (
	// AccountFilter 用户列表筛选
	AccountFilter struct {
		Filter
	}

	// AccountBaseService 基础服务接口
	AccountBaseService interface {
		Create(ctx context.Context, account *entity.Account) error
	}
	// AccountService 拓展服务接口
	AccountService interface {
		AccountBaseService
		GetAccounts(ctx context.Context, filter *AccountFilter) ([]*entity.Account, error)
		GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error)

		UpdateRoleByEmail(ctx context.Context, email string, roles types.Roles) error
		UpdatePassByEmail(ctx context.Context, email, pass string) error
		Login(ctx context.Context, email, pass string) (*entity.Session, error)
	}
	// AccountRepository 仓库接口
	AccountRepository interface {
		GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error)
		GetOne(ctx context.Context, id uint) (*entity.Account, error)
		Accounts(ctx context.Context, filter *AccountFilter) ([]*entity.Account, error)

		Save(ctx context.Context, account *entity.Account) error
		UpdateRole(ctx context.Context, account *entity.Account) error
		UpdatePass(ctx context.Context, account *entity.Account) error
	}
)
