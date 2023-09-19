// Package assembler 传输数据与实体转换
package assembler

import (
	"demo/internal/app/handlers/dto"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"time"
)

// Account assembler
type Account struct{}

func (u *Account) ToEntityFromCreateDTO(vo *dto.AccountCreateDTO) (*domain.AccountAggregate, error) {
	item := new(domain.AccountAggregate)
	item.Account = &entity.Account{
		Email:      vo.Email,
		CreateTime: time.Now(),
	}
	var err error
	item.Account.Roles, err = types.ParseRoles(vo.Role, false)
	if err != nil {
		return nil, err
	}
	if err = item.Account.SetPassword(vo.Password); err != nil {
		return nil, err
	}
	return item, nil
}
