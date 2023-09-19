// Package handlers 控制器处理
package handlers

import (
	"demo/internal/app/errs"
	"demo/internal/app/handlers/assembler"
	"demo/internal/app/handlers/dto"
	"demo/internal/app/response"
	"demo/internal/domain"
	"demo/internal/pkg/utils"
	"net/http"
)

// UserHandler 用户接口控制器
type UserHandler struct {
	userService    domain.UserService
	sessionService domain.SessionService
}

// NewUserHandler 实例化用户控制器
func NewUserHandler(s1 domain.UserService, s2 domain.SessionService) *UserHandler {
	return &UserHandler{userService: s1, sessionService: s2}
}

// HandleRegister 用户注册
func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	data := new(dto.UserCreateDTO)
	if err := utils.ParseRequestData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}

	ua := new(assembler.User).ToEntityFromCreateDTO(data)
	if err := h.userService.Create(r.Context(), ua); err != nil {
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
	session, err := h.userService.Login(r.Context(), data.Username, data.Password)
	if err != nil {
		response.Error(w, err)
		return
	}
	if err = h.sessionService.Set(r.Context(), session); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", dto.ResUserLoginDTO{
		Token: session.FormatToken(),
	})
}

// HandleIdentity 获取用户信息
func (h *UserHandler) HandleIdentity(w http.ResponseWriter, r *http.Request) {
	u, err := h.userService.GetUserinfo(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", dto.ResUserinfoDTO{
		UserID:   u.UserID,
		Username: u.Username,
		Gender:   uint8(u.Gender),
		Nickname: u.Nickname,
	})
}

// HandleChangePass 用户修改密码接口
func (h *UserHandler) HandleChangePass(w http.ResponseWriter, r *http.Request) {
	data := new(dto.ChangeUserPassDTO)
	if err := utils.ParseRequestData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}
	if err := h.userService.UpdatePassword(r.Context(), data.Password); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", nil)
}

// HandleLogout 注销登录
func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionService.Remove(r.Context()); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", nil)
}

// HandleUsers 查询用户列表
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	data := new(dto.QueryUsersDTO)
	if err := utils.ParseRequestData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}
	users, err := h.userService.Users(r.Context(), new(assembler.User).ToFilterFromQueryDTO(data))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", new(assembler.User).ToFilterResult(users))
}
