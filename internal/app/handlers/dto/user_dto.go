// Package dto of data transport object
package dto

type (
	// UserCreateDTO 创建用户数据
	UserCreateDTO struct {
		Username string `json:"username" valid:"username_valid~用户名由1-20位字母或数字组成"`
		Password string `json:"password" valid:"valid_pass~密码格式不符合要求"`
		Nickname string `json:"nickname" valid:"runelength(0|15)~昵称不能超过15个字符串"`
	}

	// UserLoginDTO 用户登录请求体
	UserLoginDTO struct {
		Username string `json:"username" valid:"username_valid~用户名由1-20位字母或数字组成"`
		Password string `json:"password" valid:"valid_pass~密码格式不符合要求"`
	}

	// ChangeUserPassDTO 用户修改请求体
	ChangeUserPassDTO struct {
		Password string `json:"password" valid:"valid_pass~密码格式不符合要求"`
	}

	// QueryUsersDTO 查询用户列表结构
	QueryUsersDTO struct {
		BaseFilter `valid:"-"`
		UserID     *uint   `json:"userId" valid:"-"`
		Nickname   *string `json:"nickname" valid:"-"`
		Gender     *uint   `json:"gender" valid:"-"`
		Status     *uint   `json:"status" valid:"-"`
	}

	// ChangeUserStatusDTO 修改用户状态请求
	ChangeUserStatusDTO struct {
		UserID uint   `json:"userId" valid:"required,numeric"`
		Status string `json:"operate" valid:"in(enable|disable)"`
	}

	// ResetUserPassDTO 重置用户密码
	ResetUserPassDTO struct {
		UserID uint `json:"userId" valid:"required,numeric"`
	}
)

type (
	// ResUserLoginDTO 用户登录响应
	ResUserLoginDTO struct {
		Token string `json:"authorizeToken"`
	}

	// ResUserinfoDTO 响应用户信息
	ResUserinfoDTO struct {
		UserID   uint   `json:"userId"`
		Username string `json:"username"`
		Gender   uint8  `json:"gender"`
		Nickname string `json:"nickname"`
	}

	// ResUserinfoItem 响应列表项
	ResUserinfoItem struct {
		UserID   uint   `json:"userId"`
		Gender   uint8  `json:"gender"`
		Nickname string `json:"nickname"`
		Status   uint8  `json:"status"`
	}
	// ResQueryDTO 查询列表响应结果
	ResQueryDTO struct {
		List []ResUserinfoItem `json:"list"`
	}
)
