package entity

import (
	"example/internal/app/errs"
	"example/internal/domain/types"
	"time"
)

// Account 管理员账号实体
type Account struct {
	AdminID    uint
	Email      string
	Password   *types.Password
	Roles      types.Roles
	CreateTime time.Time
}

// SetPassword 设置密码
func (u *Account) SetPassword(pass string) error {
	pwd, err := types.NewPassword(types.PassMethodHash, pass)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

// ValidPassword 验证密码是否正确
func (u *Account) ValidPassword(pass string) error {
	if !u.Password.Valid(pass) {
		return errs.EcInvalidUser
	}
	return nil
}

// SetRole 设置角色
func (u *Account) SetRole(roles types.Roles) {
	u.Roles = roles
}

// HasPermission 路由权限校验
func (u *Account) HasPermission(route types.Route) bool {
	return u.Roles.HasPermission(route)
}
