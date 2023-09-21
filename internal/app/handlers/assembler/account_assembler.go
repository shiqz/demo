// Package assembler 传输数据与实体转换
package assembler

import (
	"example/internal/app/handlers/dto"
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"time"
)

// Account assembler
type Account struct{}

// ToEntityFromCreateDTO 转换为实体
func (u *Account) ToEntityFromCreateDTO(vo *dto.AccountCreateDTO) (*entity.Account, error) {
	item := &entity.Account{
		Email:      vo.Email,
		CreateTime: time.Now(),
	}
	var err error
	item.Roles, err = types.ParseRoles(vo.Role, false)
	if err != nil {
		return nil, err
	}
	if err = item.SetPassword(vo.Password); err != nil {
		return nil, err
	}
	return item, nil
}

// ToAccountInfo 转换为响应列表项
func (u *Account) ToAccountInfo(vo *entity.Account) dto.ResAccountInfo {
	return dto.ResAccountInfo{
		ID:         vo.AdminID,
		Email:      vo.Email,
		Role:       vo.Roles.String(),
		CreateTime: vo.CreateTime.Format(time.DateTime),
	}
}
