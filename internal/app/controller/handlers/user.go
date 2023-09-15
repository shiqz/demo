// Package handlers 控制器处理
package handlers

import (
	"demo/internal/app/controller/assembler"
	"demo/internal/app/controller/dto"
	"demo/internal/app/controller/response"
	"demo/internal/app/errs"
	"demo/internal/domain"
	"demo/internal/pkg/utils"
	"net/http"
)

// UserHandler 用户接口控制器
type UserHandler struct {
	srv domain.UserService
}

// NewUserHandler 实例化用户控制器
func NewUserHandler(srv domain.UserService) domain.UserHandler {
	return &UserHandler{srv: srv}
}

// HandleRegister 用户注册
func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	data := new(dto.UserCreateDTO)
	if err := utils.ParseRequestData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}

	ua := new(assembler.User).ToEntityFromCreateDTO(data)
	if err := h.srv.Create(r.Context(), ua); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "success", nil)
}

// HandleLogin 用户登录
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := new(dto.UserLoginDTO)
	if err := utils.ParseRequestData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}

	ug, err := h.srv.Login(r.Context(), data.Username, data.Password)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", dto.ResUserLoginDTO{
		Token: ug.Session.FormatToken(),
	})
}
