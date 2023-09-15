// Package entity 相关实体
package entity

import (
	"demo/internal/domain/types"
	"demo/internal/pkg/utils"
	"time"
)

// User 用户实体
type User struct {
	UserID     uint             `json:"user_id"`
	Username   string           `json:"username"`
	Password   string           `json:"-"`
	Salt       string           `json:"-"`
	Nickname   string           `json:"nickname"`
	Gender     types.UserGender `json:"gender"`
	Status     types.UserState  `json:"status"`
	CreateTime time.Time        `json:"create_time"`
}

// SetPassword 设置密码
func (u *User) SetPassword(pass string) {
	u.Salt = utils.GetRandomStr(8)
	u.Password = utils.EncryptMD5(pass + u.Salt)
}

// IsValidPassword 验证密码是否正确
func (u *User) IsValidPassword(pass string) bool {
	return u.Password == utils.EncryptMD5(pass+u.Salt)
}
