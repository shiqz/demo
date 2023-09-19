// Package dto of data transport object
package dto

type (
	// AccountCreateDTO 创建admin结构
	AccountCreateDTO struct {
		Email    string `json:"email" valid:"valid_email~请输入正确的邮箱格式"`
		Password string `json:"password" valid:"valid_pass~密码格式不符合要求"`
		Role     string `json:"role" valid:"-"`
	}

	// AccountLoginDTO admin用户登录结构
	AccountLoginDTO struct {
		Email    string `json:"email" valid:"valid_email~请输入正确的邮箱格式"`
		Password string `json:"password" valid:"valid_pass~密码格式不符合要求"`
	}
)

type (
	// ResAccountLoginDTO 登录响应
	ResAccountLoginDTO struct {
		Token string `json:"authorizeToken"`
	}
)
