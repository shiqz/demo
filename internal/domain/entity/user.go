// Package entity 相关实体
package entity

import (
	"example/internal/app/errs"
	"example/internal/domain/types"
	"time"
)

// User 用户实体
type User struct {
	UserID     uint
	Username   string
	Password   *types.Password
	Nickname   string
	Gender     types.UserGender
	Status     types.UserState
	CreateTime time.Time
}

// SetPassword 设置密码
func (u *User) SetPassword(pass string) error {
	pwd, err := types.NewPassword(types.PassMethodMD5, pass)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

// ValidState 验证状态是否正常
func (u *User) ValidState() error {
	if u.Status == types.UserStateDisabled {
		return errs.EcStatusForbidden
	}
	return nil
}
