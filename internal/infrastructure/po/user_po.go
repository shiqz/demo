// Package po 定义相关数据库结构映射
package po

import (
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"time"
)

// User 用户表映射
type User struct {
	UserID     uint   `db:"user_id"`
	Username   string `db:"username"`
	Passwd     string `db:"passwd"`
	Salt       string `db:"salt"`
	Nickname   string `db:"nickname"`
	Gender     uint   `db:"gender"`
	Status     uint   `db:"status"`
	CreateTime int64  `db:"create_time"`
}

func (*User) TableName() string {
	return "users"
}

// UserConvertor 用户数据转换
type UserConvertor struct {
}

// CreateUserEntity op 转为 aggregate
func (uc *UserConvertor) CreateUserEntity(u User) *domain.UserAggregate {
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
func (uc *UserConvertor) CreateUserPO(ug *domain.UserAggregate) *User {
	return &User{
		UserID:     ug.User.UserID,
		Username:   ug.User.Username,
		Salt:       ug.User.Salt,
		Nickname:   ug.User.Nickname,
		Passwd:     ug.User.Password,
		CreateTime: ug.User.CreateTime.Unix(),
		Gender:     uint(ug.User.Gender),
		Status:     uint(ug.User.Status),
	}
}
