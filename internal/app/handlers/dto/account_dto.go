// Package dto of data transport object
package dto

type (
	// AccountCreateDTO 创建admin结构
	AccountCreateDTO struct {
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"pass_valid~密码格式不符合要求"`
	}

	// AccountLoginDTO admin用户登录结构
	AccountLoginDTO struct {
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"pass_valid~密码格式不符合要求"`
	}
)

type (
	// ResAccountLoginDTO 登录响应
	ResAccountLoginDTO struct {
		Token string `json:"authorizeToken"`
	}
)
