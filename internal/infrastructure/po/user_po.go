// Package po 定义相关数据库结构映射
package po

import (
	"example/internal/domain/entity"
	"example/internal/domain/types"
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

// UserConvertor 用户数据转换
type UserConvertor struct{}

// ToEntity op 转为 aggregate
func (uc *UserConvertor) ToEntity(vo User) *entity.User {
	item := &entity.User{
		UserID:     vo.UserID,
		Username:   vo.Username,
		Password:   vo.Passwd,
		Salt:       vo.Salt,
		Nickname:   vo.Nickname,
		CreateTime: time.Unix(vo.CreateTime, 0),
	}
	item.Gender = types.UserGender(vo.Gender)
	item.Status = types.UserState(vo.Status)
	return item
}

// CreateUserPO aggregate -> PO
func (uc *UserConvertor) CreateUserPO(vo *entity.User) *User {
	return &User{
		UserID:     vo.UserID,
		Username:   vo.Username,
		Salt:       vo.Salt,
		Nickname:   vo.Nickname,
		Passwd:     vo.Password,
		CreateTime: vo.CreateTime.Unix(),
		Gender:     uint(vo.Gender),
		Status:     uint(vo.Status),
	}
}
