// Package mysqlreposimpl 仓库MySQL实现
package mysqlreposimpl

import (
	"context"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/infrastructure/depend"
	"example/internal/infrastructure/po"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

// AccountRepository 用户仓库
type AccountRepository struct {
	inject *depend.Injecter
}

// NewAccountRepository 实例化Account repo
func NewAccountRepository(inject *depend.Injecter) domain.AccountRepository {
	return &AccountRepository{inject}
}

// GetAccountByEmail 通过邮箱获取管理员账号
func (r *AccountRepository) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	sql, args, err := dialect.From(po.AccountTable).Prepared(true).
		Where(goqu.C("email").Eq(email)).ToSQL()
	if err != nil {
		return nil, err
	}
	r.inject.GetDB().Tracef(sql, args...)
	var account po.Account
	if err = r.inject.GetDB().GetContext(ctx, &account, sql, args...); err != nil {
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
	r.inject.GetDB().Tracef(sql, args...)
	if err = r.inject.GetDB().GetContext(ctx, &data, sql, args...); err != nil {
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
	r.inject.GetDB().Tracef(sql, args...)
	if err = r.inject.GetDB().SelectContext(ctx, &list, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*entity.Account
	for _, item := range list {
		result = append(result, new(po.AccountConvertor).CreateEntity(item))
	}
	return result, nil
}

// Save 保存
func (r *AccountRepository) Save(ctx context.Context, account *entity.Account) error {
	var data = new(po.AccountConvertor).CreatePO(account)
	sql, args, err := dialect.Insert(po.AccountTable).Prepared(true).Rows(data).ToSQL()
	if err != nil {
		return errors.WithStack(err)
	}
	r.inject.GetDB().Tracef(sql, args...)
	result, err := r.inject.GetDB().ExecContext(ctx, sql, args...)
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
		"passwd": account.Password.String(),
	}
	return r.update(ctx, account.AdminID, update)
}

func (r *AccountRepository) update(ctx context.Context, id uint, update goqu.Record) error {
	sql, args, err := dialect.Update(po.AccountTable).Prepared(true).Set(update).
		Where(goqu.C("admin_id").Eq(id)).ToSQL()
	if err != nil {
		return errors.Wrap(err, sql)
	}
	r.inject.GetDB().Tracef(sql, args...)
	if _, err = r.inject.GetDB().ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}
