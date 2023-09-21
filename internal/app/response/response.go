// Package response 响应数据
package response

import (
	"encoding/json"
	"example/internal/app/errs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	// DefaultStatus 默认状态码
	DefaultStatus = http.StatusOK
	// HTTPHeaderStatusFlag 内置响应数据标识
	HTTPHeaderStatusFlag = "_INTERNAL_RESPONSE_DATA_FLAG"
	// HTTPHeaderDataFlag 内置响应状态标识
	HTTPHeaderDataFlag = "_INTERNAL_RESPONSE_STATUS_FLAG"
)

// DefaultData 默认数据
var DefaultData = struct{}{}

// response HTTP响应体
type response struct {
	w       http.ResponseWriter
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// response 实例化响应
func respond(w http.ResponseWriter, code int, msg string, data any) *response {
	if data == nil {
		data = DefaultData
	}
	return &response{w: w, Data: data, Status: code, Message: msg}
}

func (res response) json(status int) {
	res.w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	marshal, err := json.Marshal(res)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "response encode"))
	}
	res.w.Header().Set(HTTPHeaderDataFlag, string(marshal))
	res.w.Header().Set(HTTPHeaderStatusFlag, strconv.FormatInt(int64(status), 10))
}

// Success 响应成功
func Success(w http.ResponseWriter, httpStatus int, message string, data any) {
	respond(w, DefaultStatus, message, data).json(httpStatus)
}

// Error 响应错误
func Error(w http.ResponseWriter, err error) {
	var er errs.Error
	if errors.As(err, &er) {
		respond(w, er.Code(), er.Error(), nil).json(er.HTTPStatus())
		return
	}

	var e errs.ErrStatus
	if !errors.As(err, &e) {
		e = errs.EcInternalServerErr
		log.Errorf("%+v", err)
	}
	respond(w, e.Code(), e.Error(), nil).json(e.HTTPStatus())
}
