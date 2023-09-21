package handlers

import (
	"example/internal/app/errs"
	"example/internal/app/handlers/dto"
	"example/internal/app/response"
	"example/internal/domain"
	"example/internal/pkg/utils"
	"net/http"
)

// AccountHandler 账户控制器
type AccountHandler struct {
	srv            domain.AccountService
	sessionService domain.SessionService
}

// NewAccountHandler 实例化账户控制器
func NewAccountHandler(srv domain.AccountService, session domain.SessionService) *AccountHandler {
	return &AccountHandler{srv: srv, sessionService: session}
}

// HandleLogin 登录
func (h *AccountHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := new(dto.AccountLoginDTO)
	if err := utils.MustParseData(r.Body, data); err != nil {
		response.Error(w, errs.New(errs.EcInvalidRequest, err))
		return
	}
	session, err := h.srv.Login(r.Context(), data.Email, data.Password)
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

// HandleLogout 注销登录
func (h *AccountHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionService.Disconnect(r.Context()); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, http.StatusOK, "success", nil)
}
