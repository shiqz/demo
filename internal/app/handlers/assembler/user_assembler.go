// Package assembler 传输数据与实体转换
package assembler

import (
	"demo/internal/app/handlers/dto"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"time"
)

// User assembler
type User struct{}

func (u *User) ToEntityFromCreateDTO(dto *dto.UserCreateDTO) *entity.User {
	item := &entity.User{
		Username:   dto.Username,
		Nickname:   dto.Nickname,
		CreateTime: time.Now(),
	}
	item.SetPassword(dto.Password)
	return item
}

func (u *User) ToFilterFromQueryDTO(vo *dto.QueryUsersDTO) *domain.UserFilter {
	item := &domain.UserFilter{}
	if vo.UserID != nil && *vo.UserID > 0 {
		item.UserID = vo.UserID
	}
	if vo.Nickname != nil && *vo.Nickname != "" {
		item.Nickname = vo.Nickname
	}
	if vo.Gender != nil {
		tmp := types.UserGender(*vo.Gender)
		item.Gender = &tmp
	}
	if vo.Status != nil && *vo.Status != 0 {
		tmp := types.UserState(*vo.Status)
		item.Status = &tmp
	}
	return item
}

// ToFilterResult 转换查询结果
func (u *User) ToFilterResult(users []*entity.User) dto.ResQueryDTO {
	result := dto.ResQueryDTO{}
	for _, user := range users {
		result.List = append(result.List, dto.ResUserinfoItem{
			UserID:   user.UserID,
			Nickname: user.Nickname,
			Gender:   uint8(user.Gender),
			Status:   uint8(user.Status),
		})
	}
	return result
}
