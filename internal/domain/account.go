package domain

import (
	"context"
	"demo/internal/domain/entity"
)

// 管理员账号领域核心业务
type (
	// AccountAggregate 管理员聚合
	AccountAggregate struct {
		Account *entity.Account
	}
	// AccountBaseService 基础服务接口
	AccountBaseService interface {
		Create(ctx context.Context, account *AccountAggregate) error
	}
	// AccountService 拓展服务接口
	AccountService interface {
		AccountBaseService
		UpdatePassByEmail(ctx context.Context, email, pass string) error
		Login(ctx context.Context, email, pass string) (*entity.Session, error)
	}
	// AccountRepository 仓库接口
	AccountRepository interface {
		GetAccountByEmail(ctx context.Context, email string) (*AccountAggregate, error)
		GetOne(ctx context.Context, id uint) (*AccountAggregate, error)

		Save(ctx context.Context, account *AccountAggregate) error
		UpdatePass(ctx context.Context, id uint, pass string) error
	}
)
