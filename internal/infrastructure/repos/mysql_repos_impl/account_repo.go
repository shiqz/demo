// Package mysql_repos_impl 仓库MySQL实现
package mysql_repos_impl

import (
	"context"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/infrastructure/po"
	"example/internal/pkg/db"
	"example/internal/pkg/logger"
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
func (r *AccountRepository) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
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
func (r *AccountRepository) GetOne(ctx context.Context, id uint) (*entity.Account, error) {
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

// Accounts 查询用户列表
func (r *AccountRepository) Accounts(ctx context.Context, filter *domain.AccountFilter) ([]*entity.Account, error) {
	var list []po.Account
	query := dialect.From(po.AccountTable).Prepared(true)
	if filter != nil {
		query = query.Offset(filter.Offset).Limit(filter.Limit)
	}
	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r.db.Tracef(sql, args...)
	if err = r.db.SelectContext(ctx, &list, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*entity.Account
	for _, item := range list {
		result = append(result, new(po.AccountConvertor).ToEntity(item))
	}
	return result, nil
}

func (r *AccountRepository) Save(ctx context.Context, account *entity.Account) error {
	var data = new(po.AccountConvertor).CreatePO(account)
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
	account.AdminID = uint(id)
	return nil
}

// UpdateRole 修改角色
func (r *AccountRepository) UpdateRole(ctx context.Context, account *entity.Account) error {
	update := goqu.Record{
		"roles": account.Roles.String(),
	}
	return r.update(ctx, account.AdminID, update)
}

// UpdatePass 修改密码
func (r *AccountRepository) UpdatePass(ctx context.Context, account *entity.Account) error {
	update := goqu.Record{
		"passwd": account.Password,
	}
	return r.update(ctx, account.AdminID, update)
}

func (r *AccountRepository) update(ctx context.Context, id uint, update goqu.Record) error {
	sql, args, err := dialect.Update(po.AccountTable).Prepared(true).Set(update).
		Where(goqu.C("admin_id").Eq(id)).ToSQL()
	if err != nil {
		return errors.Wrap(err, sql)
	}
	r.db.Tracef(sql, args...)
	if _, err = r.db.ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}
