// Package response 响应数据
package response

import (
	"bytes"
	"encoding/json"
	"example/internal/app/errs"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	// DefaultStatus 默认状态码
	DefaultStatus = http.StatusOK
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

// 响应JSON
func (res response) json(status int) {
	res.w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	marshal, err := json.Marshal(res)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "response encode"))
	}
	// Unwrap 到自定义响应器
	r := res.w.(middleware.WrapResponseWriter).Unwrap().(*Responser)
	r.WriteHeader(status)
	_, err = r.Write(marshal)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "response write"))
	}
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

// Responser 自定义响应器
type Responser struct {
	w     http.ResponseWriter
	bytes int
	tee   bytes.Buffer
	code  int
}

// NewResponser 实例化响应器
func NewResponser(w http.ResponseWriter) http.ResponseWriter {
	return &Responser{w: w}
}

// Header 实现Header
func (r *Responser) Header() http.Header {
	return r.w.Header()
}

// Write 覆盖写入
func (r *Responser) Write(buf []byte) (int, error) {
	n, err := r.tee.Write(buf)
	r.bytes += n
	return n, err
}

// WriteHeader 覆盖状态码写入
func (r *Responser) WriteHeader(statusCode int) {
	r.code = statusCode
}

// Render 最终写入
func (r *Responser) Render() (n int64, err error) {
	r.w.WriteHeader(r.code)
	if n, err = r.tee.WriteTo(r.w); err != nil {
		err = errors.WithStack(err)
	}
	return n, err
}

// Status 获取状态码
func (r *Responser) Status() int {
	return r.code
}

// BytesWritten 获取写入字节数
func (r *Responser) BytesWritten() int {
	return r.bytes
}
