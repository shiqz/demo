// Package errs 错误包
package errs

import "github.com/pkg/errors"

var (
	ErrInvalidToken = errors.New("invalid token")
)

// Error 应用错误
type Error struct {
	ec  ErrStatus
	err error
}

// 获取错误消息
func (e Error) Error() string {
	if e.ec == EcInvalidRequest && e.err != nil {
		return e.err.Error()
	}
	return e.ec.Error()
}

// Code 获取错误码
func (e Error) Code() int {
	return e.ec.Code()
}

// HTTPStatus 获取HTTP码
func (e Error) HTTPStatus() int {
	return e.ec.HTTPStatus()
}

// GetError 获取错误
func (e Error) GetError() error {
	return e.err
}

// New 内部错误构造函数
func New(code ErrStatus, err error) Error {
	return Error{ec: code, err: err}
}
