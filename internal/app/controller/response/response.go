// Package response 响应数据
package response

import (
	"demo/internal/app/errs"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	// DefaultEmptyData 默认空数据
	DefaultEmptyData = "{}"
	// DefaultStatus 默认状态码
	DefaultStatus = http.StatusOK
	// DefaultMessage 默认返回消息
	DefaultMessage = "success"
)

// Response HTTP响应体
type Response struct {
	w       http.ResponseWriter
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// NewResponse 实例化响应
func NewResponse(w http.ResponseWriter, code int, msg string, data any) *Response {
	return &Response{w: w, Data: data, Status: code, Message: msg}
}

func (res Response) send(status int) {
	if res.Data == nil {
		res.Data = DefaultEmptyData
	}
	res.w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.w.WriteHeader(status)
	err := errors.Wrap(json.NewEncoder(res.w).Encode(res), "response encode")
	if err != nil {
		log.Errorf("%+v", err)
	}
}

// Success 响应成功
func Success(w http.ResponseWriter, httpStatus int, message string, data any) {
	NewResponse(w, http.StatusOK, message, data).send(httpStatus)
}

// Error 响应错误
func Error(w http.ResponseWriter, err error) {
	var er errs.Error
	if errors.As(err, &er) {
		NewResponse(w, er.Code(), er.Error(), nil).send(er.HTTPStatus())
		return
	}

	var e errs.ErrStatus
	if !errors.As(err, &e) {
		e = errs.EcInternalServerErr
		log.Errorf("%+v", err)
	}
	NewResponse(w, e.Code(), e.Error(), nil).send(e.HTTPStatus())
}
