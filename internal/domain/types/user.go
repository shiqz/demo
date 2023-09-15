// Package types 值对象定义
package types

// UserGender 性别类型
type UserGender uint8

// UserState 用户状态类型
type UserState uint8

// 用户性别枚举
const (
	UserGenderUnknown UserGender = 0 // 未知
	UserGenderMale    UserGender = 1 // 男
	UserGenderFemale  UserGender = 2 // 女
)

// 用户状态枚举
const (
	UserStateNormal   UserState = 1 // 启用
	UserStateDisabled UserState = 2 // 已禁用
)

// String 用户状态类型转换为字符串
func (s UserState) String() string {
	if s == UserStateNormal {
		return "enable"
	}
	return "disable"
}

// ParseUserState 解析字符串格式用户状态
func ParseUserState(state string) UserState {
	if state == UserStateNormal.String() {
		return UserStateNormal
	}
	return UserStateDisabled
}
