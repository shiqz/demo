// Package errs 错误包
package errs

import "github.com/pkg/errors"

var (
	// ErrInvalidToken token 错误
	ErrInvalidToken = errors.New("invalid token")
	// ErrAPINotFound 访问不存在接口错误
	ErrAPINotFound = errors.New("API not found")
)

// Error 应用错误
type Error struct {
	ec  ErrStatus
	err error
}

// 获取错误消息
func (e Error) Error() string {
	if e.err != nil {
		switch e.ec {
		case EcInvalidRequest:
			return e.err.Error()
		case EcNotFound:
			return e.err.Error()
		}
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
