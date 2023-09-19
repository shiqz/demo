// Package mysql_repos_impl 仓库MySQL实现
package mysql_repos_impl

import (
	"context"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/infrastructure/po"
	"demo/internal/pkg/db"
	"demo/internal/pkg/logger"
	"encoding/json"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// UserRepository 用户仓库
type UserRepository struct {
	db    *db.Connector
	redis *db.Redis
	lg    *logger.Logger
}

// NewUserRepository 实例化User repo
func NewUserRepository(dc *db.Connector, redis *db.Redis, lg *logger.Logger) domain.UserRepository {
	return &UserRepository{db: dc, redis: redis, lg: lg}
}

// key format key
func (r *UserRepository) key(id uint) string {
	return "u:info:" + strconv.FormatInt(int64(id), 10)
}

// 根据ID缓存用户
func (r *UserRepository) reload(ctx context.Context, id uint) (*entity.User, error) {
	sql, args, err := dialect.From(po.UserTable).Prepared(true).Where(goqu.Ex{"user_id": id}).ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, sql)
	}
	r.db.Tracef(sql, args...)
	var data po.User
	if err = r.db.GetContext(ctx, &data, sql, args...); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	u := new(po.UserConvertor).CreateUserEntity(data)
	marshal, err := json.Marshal(u)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = r.redis.Set(ctx, r.key(id), string(marshal), -1).Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}

// GetUserByUsername 根据用户名获取用户
func (r *UserRepository) GetUserByUsername(ctx context.Context, uname string) (*entity.User, error) {
	sql, args, err := dialect.From(po.UserTable).Prepared(true).Where(goqu.Ex{"username": uname}).ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r.db.Tracef(sql, args...)
	var data po.User
	if err = r.db.GetContext(ctx, &data, sql, args...); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	return new(po.UserConvertor).CreateUserEntity(data), nil
}

// GetOne 根据ID获取用户
func (r *UserRepository) GetOne(ctx context.Context, id uint) (*entity.User, error) {
	cacheInfo, err := r.redis.Get(ctx, r.key(id)).Result()
	empty := errors.Is(err, redis.Nil)
	if err != nil && !empty {
		return nil, err
	}
	u := new(entity.User)
	if empty || cacheInfo == "" {
		if u, err = r.reload(ctx, id); err != nil {
			return nil, err
		}
	}
	u = new(entity.User)
	if err = json.Unmarshal([]byte(cacheInfo), u); err != nil {
		return nil, err
	}
	return u, nil
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	data := new(po.UserConvertor).CreateUserPO(user)
	sql, args, err := dialect.Insert(po.UserTable).Prepared(true).Rows(data).ToSQL()
	if err != nil {
		return errors.Wrap(err, sql)
	}
	r.db.Tracef(sql, args...)
	if _, err = r.db.ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}

// UpdatePass 修改用户密码
func (r *UserRepository) UpdatePass(ctx context.Context, uid uint, pass string) error {
	user, err := r.GetOne(ctx, uid)
	if err != nil {
		return err
	}
	user.SetPassword(pass)
	userPO := new(po.UserConvertor).CreateUserPO(user)
	sql, args, err := dialect.Update(po.UserTable).Prepared(true).Set(userPO).
		Where(goqu.C("user_id").Eq(uid)).ToSQL()
	if err != nil {
		return errors.Wrap(err, sql)
	}
	r.db.Tracef(sql, args...)
	if _, err = r.db.ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}
