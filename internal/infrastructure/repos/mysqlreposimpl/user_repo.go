// Package mysqlreposimpl 仓库MySQL实现
package mysqlreposimpl

import (
	"context"
	"encoding/json"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/infrastructure/depend"
	"example/internal/infrastructure/po"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// UserRepository 用户仓库
type UserRepository struct {
	inject *depend.Injecter
}

// NewUserRepository 实例化User repo
func NewUserRepository(inject *depend.Injecter) domain.UserRepository {
	return &UserRepository{inject}
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
	r.inject.GetDB().Tracef(sql, args...)
	var data po.User
	if err = r.inject.GetDB().GetContext(ctx, &data, sql, args...); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	u := new(po.UserConvertor).ToEntity(data)
	marshal, err := json.Marshal(u)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = r.inject.GetRedis().Set(ctx, r.key(id), string(marshal), -1).Err(); err != nil {
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
	r.inject.GetDB().Tracef(sql, args...)
	var data po.User
	if err = r.inject.GetDB().GetContext(ctx, &data, sql, args...); err != nil {
		return nil, errors.Wrap(err, sql)
	}
	return new(po.UserConvertor).ToEntity(data), nil
}

// GetOne 根据ID获取用户
func (r *UserRepository) GetOne(ctx context.Context, id uint) (*entity.User, error) {
	cacheInfo, err := r.inject.GetRedis().Get(ctx, r.key(id)).Result()
	empty := errors.Is(err, redis.Nil)
	if err != nil && !empty {
		return nil, err
	}
	u := new(entity.User)
	if empty || cacheInfo == "" {
		return r.reload(ctx, id)
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
	r.inject.GetDB().Tracef(sql, args...)
	result, err := r.inject.GetDB().ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, sql)
	}
	id, _ := result.LastInsertId()
	user.UserID = uint(id)
	return nil
}

// Users 查询用户列表
func (r *UserRepository) Users(ctx context.Context, filter *domain.UserFilter) ([]*entity.User, error) {
	var list []po.User
	where := goqu.Ex{}
	if filter.UserID != nil {
		where["user_id"] = *filter.UserID
	}
	if filter.Nickname != nil {
		where["nickname"] = goqu.Op{"like": `%` + *filter.Nickname + `%`}
	}
	if filter.Gender != nil {
		where["gender"] = *filter.Gender
	}
	if filter.Status != nil {
		where["status"] = *filter.Status
	}
	sql, args, err := dialect.From(po.UserTable).Prepared(true).Where(where).
		Offset(filter.Offset).Limit(filter.Limit).ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r.inject.GetDB().Tracef(sql, args...)
	if err = r.inject.GetDB().SelectContext(ctx, &list, sql, args...); err != nil {
		return nil, errors.WithStack(err)
	}
	var result []*entity.User
	for _, user := range list {
		result = append(result, new(po.UserConvertor).ToEntity(user))
	}
	return result, nil
}

// UpdatePass 修改用户密码
func (r *UserRepository) UpdatePass(ctx context.Context, user *entity.User) error {
	update := goqu.Record{
		"passwd": user.Password.String(),
		"salt":   user.Password.GetSalt(),
	}
	return r.update(ctx, user.UserID, update)
}

// UpdateStatus 修改用户状态
func (r *UserRepository) UpdateStatus(ctx context.Context, user *entity.User) error {
	update := goqu.Record{
		"status": user.Status,
	}
	if err := r.update(ctx, user.UserID, update); err != nil {
		return err
	}
	return errors.WithStack(r.inject.GetRedis().Del(ctx, r.key(user.UserID)).Err())
}

func (r *UserRepository) update(ctx context.Context, id uint, update goqu.Record) error {
	sql, args, err := dialect.Update(po.UserTable).Prepared(true).Set(update).
		Where(goqu.C("user_id").Eq(id)).ToSQL()
	if err != nil {
		return errors.Wrap(err, sql)
	}
	r.inject.GetDB().Tracef(sql, args...)
	if _, err = r.inject.GetDB().ExecContext(ctx, sql, args...); err != nil {
		return errors.Wrap(err, sql)
	}
	return nil
}
