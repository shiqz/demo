package repos

import (
	"context"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"demo/internal/infrastructure/convertor"
	"demo/internal/infrastructure/po"
	"encoding/json"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// UserRepository 用户仓库
type UserRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

// NewUserRepository 实例化User repo
func NewUserRepository(db *sqlx.DB, redis *redis.Client) domain.UserRepository {
	return &UserRepository{db: db, redis: redis}
}

// GetUserByUsername 根据用户名获取用户
func (r *UserRepository) GetUserByUsername(ctx context.Context, uname string) (*domain.UserAggregate, error) {
	sql, _, err := goqu.From("users").Where(goqu.Ex{"username": uname}).ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sql = format(sql)
	var data po.User
	if err = r.db.GetContext(ctx, &data, sql); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	return new(convertor.UserConvertor).CreateUserEntity(data), nil
}

// GetOne 根据ID获取用户
func (r *UserRepository) GetOne(ctx context.Context, id uint) (*domain.UserAggregate, error) {
	cacheInfo, err := r.redis.Get(ctx, r.key(id)).Result()
	empty := errors.Is(err, redis.Nil)
	if err != nil && !empty {
		return nil, err
	}
	var ug *domain.UserAggregate
	if empty {
		if ug, err = r.reload(ctx, id); err != nil {
			return nil, err
		}
	} else {
		var user entity.User
		if err = json.Unmarshal([]byte(cacheInfo), &user); err != nil {
			return nil, err
		}
		ug.User = &user
	}
	us := entity.NewSession(types.UserSession, ug.User.UserID)
	usCache, err := r.redis.Get(ctx, us.FormatKey()).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if err = us.Decode(usCache); err == nil {
		ug.Session = us
	}
	return ug, nil
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, ug *domain.UserAggregate) error {
	data := new(convertor.UserConvertor).CreateUserPO(ug)
	sql, _, err := goqu.Insert("users").Rows(data).ToSQL()
	if err != nil {
		return errors.Wrap(err, sql)
	}

	if _, err = r.db.ExecContext(ctx, format(sql)); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}

// SetSession 保存用户会话
func (r *UserRepository) SetSession(ctx context.Context, ug *domain.UserAggregate) error {
	cacheInfo, err := ug.Session.Encode()
	if err != nil {
		return err
	}
	return errors.WithStack(r.redis.Set(ctx, ug.Session.FormatKey(), cacheInfo, ug.Session.GetDuration()).Err())
}

// key format key
func (r *UserRepository) key(id uint) string {
	return "u:info:" + strconv.FormatInt(int64(id), 10)
}

// 根据ID缓存用户
func (r *UserRepository) reload(ctx context.Context, id uint) (*domain.UserAggregate, error) {
	sql, _, err := goqu.From("users").Where(goqu.Ex{"user_id": id}).ToSQL()
	sql = format(sql)
	if err != nil {
		return nil, errors.Wrap(err, sql)
	}
	var data po.User
	if err = r.db.GetContext(ctx, &data, sql); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	ug := new(convertor.UserConvertor).CreateUserEntity(data)
	marshal, err := json.Marshal(ug.User)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = r.redis.Set(ctx, r.key(id), string(marshal), -1).Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return ug, nil
}
