// Package entity 相关实体
package entity

import (
	"example/internal/app/errs"
	"example/internal/domain/types"
	"example/internal/pkg/utils"
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

// ValidPassword 验证密码是否正确
func (u *User) ValidPassword(pass string) error {
	if u.Password != utils.EncryptMD5(pass+u.Salt) {
		return errs.EcInvalidUser
	}
	return nil
}

// ValidState 验证状态是否正常
func (u *User) ValidState() error {
	if u.Status == types.UserStateDisabled {
		return errs.EcStatusForbidden
	}
	return nil
}
