package entity

import (
	"demo/internal/app/errs"
	"demo/internal/pkg/utils"
)

// Account 管理员账号实体
type Account struct {
	AdminID    uint
	Email      string
	Password   string
	Roles      string
	CreateTime int64
}

// SetPassword 设置密码
func (u *Account) SetPassword(pass string) error {
	enPass, err := utils.HashPassEncrypt([]byte(pass))
	if err != nil {
		return err
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
