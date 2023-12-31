// Package assembler 传输数据与实体转换
package assembler

import (
	"example/internal/app/handlers/dto"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"time"
)

// User assembler
type User struct{}

// ToEntityFromCreateDTO 用户创建实体
func (u *User) ToEntityFromCreateDTO(vo *dto.UserCreateDTO) (*entity.User, error) {
	item := &entity.User{
		Username:   vo.Username,
		Nickname:   vo.Nickname,
		CreateTime: time.Now(),
	}
	item.Gender = types.UserGenderUnknown
	item.Status = types.UserStateNormal
	if err := item.SetPassword(vo.Password); err != nil {
		return nil, err
	}
	return item, nil
}

// ToFilterFromQueryDTO 转换为查询
func (u *User) ToFilterFromQueryDTO(vo *dto.QueryUsersDTO) *domain.UserFilter {
	item := &domain.UserFilter{
		Filter: vo.GetBaseFilter(),
	}
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
	result := dto.ResQueryDTO{
		List: make([]dto.ResUserinfoItem, 0),
	}
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
