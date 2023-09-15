// Package assembler 传输数据与实体转换
package assembler

import (
	"demo/internal/app/controller/dto"
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"time"
)

// User assembler
type User struct{}

func (u *User) ToEntityFromCreateDTO(dto *dto.UserCreateDTO) *domain.UserAggregate {
	item := new(domain.UserAggregate)
	item.User = &entity.User{
		Username:   dto.Username,
		Nickname:   dto.Nickname,
		CreateTime: time.Now(),
	}
	item.User.SetPassword(dto.Password)
	return item
}
