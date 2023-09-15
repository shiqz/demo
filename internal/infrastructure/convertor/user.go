// Package convertor 数据库数据与实体转换
package convertor

import (
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"demo/internal/infrastructure/po"
	"time"
)

// UserConvertor 用户数据转换
type UserConvertor struct {
}

// CreateUserEntity op 转为 aggregate
func (uc *UserConvertor) CreateUserEntity(u po.User) *domain.UserAggregate {
	et := &entity.User{
		UserID:     u.UserID,
		Username:   u.Username,
		Password:   u.Passwd,
		Salt:       u.Salt,
		Nickname:   u.Nickname,
		CreateTime: time.Unix(u.CreateTime, 0),
	}
	et.Gender = types.UserGender(u.Gender)
	et.Status = types.UserState(u.Status)
	return &domain.UserAggregate{
		User: et,
	}
}

// CreateUserPO aggregate -> PO
func (uc *UserConvertor) CreateUserPO(ug *domain.UserAggregate) *po.User {
	return &po.User{
		Username:   ug.User.Username,
		Salt:       ug.User.Salt,
		Nickname:   ug.User.Nickname,
		Passwd:     ug.User.Password,
		CreateTime: ug.User.CreateTime.Unix(),
	}
}
