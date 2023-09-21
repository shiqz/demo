package entity

import (
	"example/internal/app/errs"
	"example/internal/domain/types"
	"example/internal/pkg/utils"
	"github.com/pkg/errors"
	"time"
)

// Account 管理员账号实体
type Account struct {
	AdminID    uint
	Email      string
	Password   string
	Roles      types.Roles
	CreateTime time.Time
}

// SetPassword 设置密码
func (u *Account) SetPassword(pass string) error {
	enPass, err := utils.HashPassEncrypt([]byte(pass))
	if err != nil {
		return errors.WithStack(err)
	}
	u.Password = string(enPass)
	return nil
}

// ValidPassword 验证密码是否正确
func (u *Account) ValidPassword(pass string) error {
	if utils.HashPassCheck([]byte(u.Password), []byte(pass)) != nil {
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
