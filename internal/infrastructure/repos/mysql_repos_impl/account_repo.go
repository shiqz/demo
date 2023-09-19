// Package mysql_repos_impl 仓库MySQL实现
package mysql_repos_impl

import (
	"context"
	"demo/internal/domain"
	"demo/internal/infrastructure/po"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

// AccountRepository 用户仓库
type AccountRepository struct {
	db    *db.Connector
	redis *db.Redis
	lg    *logger.Logger
}

// NewAccountRepository 实例化Account repo
func NewAccountRepository(dc *db.Connector, redis *db.Redis, lg *logger.Logger) domain.AccountRepository {
	return &AccountRepository{db: dc, redis: redis, lg: lg}
}

// GetAccountByEmail 通过邮箱获取管理员账号
func (r *AccountRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.AccountAggregate, error) {
	sql, args, err := dialect.From(po.AccountTable).Prepared(true).
		Where(goqu.C("email").Eq(email)).ToSQL()
	if err != nil {
		return nil, err
	}
	r.db.Tracef(sql, args...)
	var account po.Account
	if err = r.db.GetContext(ctx, &account, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}
	return new(po.AccountConvertor).CreateEntity(account), nil
}

// GetOne 通过ID获取管理员账号
func (r *AccountRepository) GetOne(ctx context.Context, id uint) (*domain.AccountAggregate, error) {
	var data po.Account
	sql, args, err := dialect.From(po.AccountTable).Prepared(true).Where(goqu.C("admin_id").Eq(id)).ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r.db.Tracef(sql, args...)
	if err = r.db.GetContext(ctx, &data, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}
	return new(po.AccountConvertor).CreateEntity(data), nil
}

func (r *AccountRepository) Save(ctx context.Context, ag *domain.AccountAggregate) error {
	var data = new(po.AccountConvertor).CreatePO(ag)
	sql, args, err := dialect.Insert(po.AccountTable).Prepared(true).Rows(data).ToSQL()
	if err != nil {
		return errors.WithStack(err)
	}
	r.db.Tracef(sql, args...)
	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	id, _ := result.LastInsertId()
	ag.Account.AdminID = uint(id)
	return nil
}

func (r *AccountRepository) UpdatePass(ctx context.Context, id uint, pass string) error {
	return nil
}
