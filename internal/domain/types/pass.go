package types

import (
	"example/internal/pkg/utils"
	"github.com/pkg/errors"
)

// PassMethod 加密类型
type PassMethod string

const (
	PassMethodHash = "hash" // PassMethodHash 哈希加密
	PassMethodMD5  = "md5"  // PassMethodMD5 MD5加密
)

var (
	// ErrUndefinedMethod 错误的加密方式
	ErrUndefinedMethod = errors.New("invalid pass method")
)

// NewPassword 创建密码
func NewPassword(method PassMethod, pass string) (*Password, error) {
	p := &Password{method: method}
	switch method {
	case PassMethodMD5:
		p.salt = utils.GetRandomStr(8)
		p.pass = utils.EncryptMD5(pass + p.salt)
	case PassMethodHash:
		en, err := utils.HashPassEncrypt([]byte(pass))
		if err != nil {
			return nil, err
		}
		p.pass = string(en)
	default:
		return nil, ErrUndefinedMethod
	}
	return p, nil
}

// ParseMD5Password 解析MD5密码
func ParseMD5Password(pass, salt string) *Password {
	return &Password{method: PassMethodMD5, salt: salt, pass: pass}
}

// ParseHashPassword 解析Hash密码
func ParseHashPassword(pass string) *Password {
	return &Password{method: PassMethodHash, pass: pass}
}

// Password 密码对象
type Password struct {
	method PassMethod
	pass   string
	salt   string // md5 加密时需要
}

// Valid 密码是否匹配
func (p *Password) Valid(pass string) bool {
	switch p.method {
	case PassMethodMD5:
		return p.pass == utils.EncryptMD5(pass+p.salt)
	case PassMethodHash:
		return utils.HashPassCheck([]byte(p.pass), []byte(pass)) == nil
	}
	return false
}

// String 返回密码
func (p *Password) String() string {
	return p.pass
}

// GetSalt 返回密码盐
func (p *Password) GetSalt() string {
	return p.salt
}
