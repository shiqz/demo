package mysql_repos_impl

import (
	"context"
	"demo/internal/domain"
	"demo/internal/infrastructure/po"
	"demo/internal/pkg/db"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// AccountRepository 用户仓库
type AccountRepository struct {
	db    *db.Connector
	redis *redis.Client
}

// NewAccountRepository 实例化Account repo
func NewAccountRepository(dc *db.Connector, redis *redis.Client) domain.AccountRepository {
	return &AccountRepository{db: dc, redis: redis}
}

func (r *AccountRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.AccountAggregate, error) {
	sql, _, err := r.db.PerSQL(r.db.Query(new(po.Account).TableName()).Where(goqu.C("email").Eq(email)).ToSQL())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var account po.Account
	if err = r.db.Driver.GetContext(ctx, &account, sql); err != nil {
		return nil, errors.WithStack(err)
	}
	return new(po.AccountConvertor).CreateEntity(account), nil
}
func (r *AccountRepository) GetOne(ctx context.Context, id uint) (*domain.AccountAggregate, error) {
	return nil, nil
}

func (r *AccountRepository) Save(ctx context.Context, account *domain.AccountAggregate) error {
	return nil
}
func (r *AccountRepository) UpdatePass(ctx context.Context, id uint, pass string) error {
	return nil
}
