package commands

import (
	"context"
	"demo/internal/app/handlers/assembler"
	"demo/internal/app/handlers/dto"
	"demo/internal/domain"
	"demo/internal/pkg/utils"
)

type AccountHandler struct {
	srv domain.AccountService
}

func NewAccountHandler(srv domain.AccountService) *AccountHandler {
	return &AccountHandler{srv: srv}
}

func (c *AccountHandler) Create(ctx context.Context, email, pass, role string) error {
	data := &dto.AccountCreateDTO{
		Email:    email,
		Password: pass,
		Role:     role,
	}
	if err := utils.Validator(data); err != nil {
		return err
	}
	ag, err := new(assembler.Account).ToEntityFromCreateDTO(data)
	if err != nil {
		return err
	}
	return c.srv.Create(ctx, ag)
}
