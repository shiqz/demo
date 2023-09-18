// Package mysql_repos_impl 仓库MySQL实现
package mysql_repos_impl

import (
	"context"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/infrastructure/po"
	"demo/internal/pkg/db"
	"encoding/json"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// UserRepository 用户仓库
type UserRepository struct {
	db    *db.Connector
	redis *redis.Client
}

// NewUserRepository 实例化User repo
func NewUserRepository(dc *db.Connector, redis *redis.Client) domain.UserRepository {
	return &UserRepository{db: dc, redis: redis}
}

// key format key
func (r *UserRepository) key(id uint) string {
	return "u:info:" + strconv.FormatInt(int64(id), 10)
}

// 根据ID缓存用户
func (r *UserRepository) reload(ctx context.Context, id uint) (*domain.UserAggregate, error) {
	sql, _, err := r.db.PerSQL(r.db.Query(new(po.User).TableName()).Where(goqu.Ex{"user_id": id}).ToSQL())
	if err != nil {
		return nil, errors.Wrap(err, sql)
	}
	var data po.User
	if err = r.db.Driver.GetContext(ctx, &data, sql); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	ug := new(po.UserConvertor).CreateUserEntity(data)
	marshal, err := json.Marshal(ug.User)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = r.redis.Set(ctx, r.key(id), string(marshal), -1).Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return ug, nil
}

// GetUserByUsername 根据用户名获取用户
func (r *UserRepository) GetUserByUsername(ctx context.Context, uname string) (*domain.UserAggregate, error) {
	sql, _, err := r.db.PerSQL(r.db.Query(new(po.User).TableName()).Where(goqu.Ex{"username": uname}).ToSQL())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var data po.User
	if err = r.db.Driver.GetContext(ctx, &data, sql); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	return new(po.UserConvertor).CreateUserEntity(data), nil
}

// GetOne 根据ID获取用户
func (r *UserRepository) GetOne(ctx context.Context, id uint) (*domain.UserAggregate, error) {
	cacheInfo, err := r.redis.Get(ctx, r.key(id)).Result()
	empty := errors.Is(err, redis.Nil)
	if err != nil && !empty {
		return nil, err
	}
	ug := new(domain.UserAggregate)
	if empty || cacheInfo == "" {
		if ug, err = r.reload(ctx, id); err != nil {
			return nil, err
		}
	}
	ug.User = new(entity.User)
	if err = json.Unmarshal([]byte(cacheInfo), ug.User); err != nil {
		return nil, err
	}
	return ug, nil
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, ug *domain.UserAggregate) error {
	data := new(po.UserConvertor).CreateUserPO(ug)
	sql, _, err := r.db.PerSQL(r.db.Insert("users").Rows(data).ToSQL())
	if err != nil {
		return errors.Wrap(err, sql)
	}

	if _, err = r.db.Driver.ExecContext(ctx, sql); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}

// UpdatePass 修改用户密码
func (r *UserRepository) UpdatePass(ctx context.Context, uid uint, pass string) error {
	ug, err := r.GetOne(ctx, uid)
	if err != nil {
		return err
	}
	ug.User.SetPassword(pass)
	userPO := new(po.UserConvertor).CreateUserPO(ug)
	sql, _, err := r.db.PerSQL(r.db.Update(userPO.TableName()).Set(userPO).Where(goqu.C("user_id").Eq(uid)).ToSQL())
	if err != nil {
		return errors.Wrap(err, sql)
	}
	if _, err = r.db.Driver.ExecContext(ctx, sql); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}
