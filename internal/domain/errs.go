package domain

import "github.com/pkg/errors"

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在或已被删除")
	// ErrAccountNotFound 账户不存在
	ErrAccountNotFound = errors.New("管理员不存在或已被删除")
)
