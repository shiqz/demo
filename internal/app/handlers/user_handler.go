// Package handlers 控制器处理
package handlers

import (
	"example/internal/app/errs"
	"example/internal/app/handlers/assembler"
	"example/internal/app/handlers/dto"
	"example/internal/app/response"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"example/internal/pkg/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
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
	if err := utils.MustParseData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidParams, err))
		return
	}

	user, err := new(assembler.User).ToEntityFromCreateDTO(data)
	if err != nil {
		response.Error(w, err)
		return
	}
	if err = h.userService.Create(r.Context(), user); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "success", nil)
}

// HandleLogin 用户登录
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := new(dto.UserLoginDTO)
	if err := utils.MustParseData(r.Body, data); err != nil {
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
	if err := utils.MustParseData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}
	session := r.Context().Value(types.SessionFlag).(*entity.Session)
	if err := h.userService.UpdatePassword(r.Context(), session.GetSessionID(), data.Password); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", nil)
}

// HandleLogout 注销登录
func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionService.Disconnect(r.Context()); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", nil)
}

// HandleUsers 查询用户列表
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	data := new(dto.QueryUsersDTO)
	if err := utils.ParseData(r.Body, data); err != nil {
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

// ChangeUserStatus 修改用户状态
func (h *UserHandler) ChangeUserStatus(w http.ResponseWriter, r *http.Request) {
	data := new(dto.ChangeUserStatusDTO)
	if err := utils.MustParseData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}
	ctx := r.Context()
	status := types.ParseUserState(data.Status)
	if err := h.userService.UpdateStatus(ctx, data.UserID, status); err != nil {
		response.Error(w, err)
		return
	}
	// 禁用后用户将被强制退出
	if status == types.UserStateDisabled {
		us, err := h.sessionService.Get(ctx, types.UserSession, data.UserID)
		if err != nil && !errors.Is(err, redis.Nil) {
			response.Error(w, err)
			return
		}
		if us != nil {
			if err = h.sessionService.Remove(ctx, us.FormatKey()); err != nil {
				response.Error(w, err)
				return
			}
		}
	}
	response.Success(w, http.StatusOK, "success", nil)
}

// ResetUserPass 重置用户密码
func (h *UserHandler) ResetUserPass(w http.ResponseWriter, r *http.Request) {
	data := new(dto.ResetUserPassDTO)
	if err := utils.MustParseData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}
	if err := h.userService.UpdatePassword(r.Context(), data.UserID, domain.DefaultUserPass); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", nil)
}
