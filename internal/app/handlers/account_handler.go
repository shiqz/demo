package handlers

import (
	"demo/internal/app/errs"
	"demo/internal/app/handlers/dto"
	"demo/internal/app/response"
	"demo/internal/domain"
	"demo/internal/pkg/utils"
	"net/http"
)

type AccountHandler struct {
	srv            domain.AccountService
	sessionService domain.SessionService
}

func NewAccountHandler(srv domain.AccountService, session domain.SessionService) *AccountHandler {
	return &AccountHandler{srv: srv, sessionService: session}
}

// HandleLogin 登录
func (h *AccountHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := new(dto.AccountLoginDTO)
	if err := utils.ParseRequestData(r.Body, data); err != nil {
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
