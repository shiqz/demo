// Package errs 内部错误处理包
package errs

import "net/http"

// ErrStatus 内部错误码类型
type ErrStatus int

// 内部错误码
const (
	// EcInvalidRequest 非法请求
	EcInvalidRequest ErrStatus = 1000 + iota
	// EcInternalServerErr 内部服务错误
	EcInternalServerErr
	// EcNotFound 资源不存在
	EcNotFound
	// EcInvalidUser 用户名或密码错误
	EcInvalidUser
	// EcUserHasBeenExist 用户已被注册
	EcUserHasBeenExist
	// EcUnauthorized 未授权登录
	EcUnauthorized
	// EcStatusForbidden 权限不足
	EcStatusForbidden
	// EcStatusForbiddenForPerms 权限不足
	EcStatusForbiddenForPerms
)

var (
	// EcMessages 内部错误码对应错误消息
	EcMessages = map[ErrStatus]string{
		EcInvalidRequest:          "invalid request",
		EcInternalServerErr:       "Internal serve error",
		EcNotFound:                "Not found",
		EcInvalidUser:             "账号或密码错误",
		EcUserHasBeenExist:        "该账号已被注册",
		EcUnauthorized:            "未授权登录或登录会话已过期",
		EcStatusForbidden:         "您的账号已被禁用，如有疑问请联系管理人员",
		EcStatusForbiddenForPerms: "你没有该操作权限",
	}
	// EcHTTPStatus 内部错误对应HTTP状态码
	EcHTTPStatus = map[ErrStatus]int{
		EcInvalidRequest:          http.StatusBadRequest,
		EcInternalServerErr:       http.StatusInternalServerError,
		EcNotFound:                http.StatusNotFound,
		EcInvalidUser:             http.StatusNotAcceptable,
		EcUserHasBeenExist:        http.StatusConflict,
		EcUnauthorized:            http.StatusUnauthorized,
		EcStatusForbidden:         http.StatusForbidden,
		EcStatusForbiddenForPerms: http.StatusForbidden,
	}
)

// Code 获取错误码
func (ec ErrStatus) Code() int {
	return int(ec)
}

// IsInternalErr 判断是否是内部服务错误
func (ec ErrStatus) IsInternalErr() bool {
	return ec == EcInternalServerErr
}

// Error 获取错误消息
func (ec ErrStatus) Error() string {
	vo, ok := EcMessages[ec]
	if ok {
		return vo
	}
	return "unknown error"
}

// HTTPStatus 获取错误消息
func (ec ErrStatus) HTTPStatus() int {
	vo, ok := EcHTTPStatus[ec]
	if ok {
		return vo
	}
	return http.StatusOK
}
