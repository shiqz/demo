// Package dto of data transport object
package dto

type (
	// UserCreateDTO 创建用户数据
	UserCreateDTO struct {
		Username string `json:"username" valid:"username_valid~用户名由1-20位字母或数字组成"`
		Password string `json:"password" valid:"pass_valid~密码格式不符合要求"`
		Nickname string `json:"nickname" valid:"runelength(0|15)~昵称不能超过15个字符串"`
	}

	// UserLoginDTO 用户登录请求体
	UserLoginDTO struct {
		Username string `json:"username" valid:"username_valid~用户名由1-20位字母或数字组成"`
		Password string `json:"password" valid:"pass_valid~密码格式不符合要求"`
	}

	// ChangeUserPassDTO 用户修改请求体
	ChangeUserPassDTO struct {
		Password string `json:"password" valid:"pass_valid~密码格式不符合要求"`
	}
)

type (
	// ResUserLoginDTO 用户登录响应
	ResUserLoginDTO struct {
		Token string `json:"authorizeToken"`
	}
)
