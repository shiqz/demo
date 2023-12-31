// Package po 定义相关数据库结构映射
package po

import (
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"time"
)

// Account 用户表映射
type Account struct {
	AdminID    uint   `db:"admin_id"`
	Email      string `db:"email"`
	Passwd     string `db:"passwd"`
	Roles      string `db:"roles"`
	CreateTime int64  `db:"create_time"`
}

// AccountConvertor 用户数据转换
type AccountConvertor struct{}

// CreateEntity 转化为实体
func (c *AccountConvertor) CreateEntity(data Account) *entity.Account {
	item := &entity.Account{
		AdminID:    data.AdminID,
		Email:      data.Email,
		Password:   types.ParseHashPassword(data.Passwd),
		CreateTime: time.Unix(data.CreateTime, 0),
	}
	item.Roles, _ = types.ParseRoles(data.Roles, true)
	return item
}

// CreatePO aggregate -> PO
func (c *AccountConvertor) CreatePO(vo *entity.Account) *Account {
	item := &Account{
		Email:      vo.Email,
		CreateTime: vo.CreateTime.Unix(),
		Roles:      vo.Roles.String(),
		Passwd:     vo.Password.String(),
	}
	return item
}
