// Package po 定义相关数据库结构映射
package po

import (
	"demo/internal/domain"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
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

// CreateEntity op 转为 aggregate
func (uc *AccountConvertor) CreateEntity(data Account) *domain.AccountAggregate {
	item := &entity.Account{
		AdminID:    data.AdminID,
		Email:      data.Email,
		Password:   data.Passwd,
		CreateTime: time.Unix(data.CreateTime, 0),
	}
	item.Roles, _ = types.ParseRoles(data.Roles, true)
	return &domain.AccountAggregate{
		Account: item,
	}
}

// CreatePO aggregate -> PO
func (uc *AccountConvertor) CreatePO(vo *domain.AccountAggregate) *Account {
	item := &Account{
		Email:      vo.Account.Email,
		CreateTime: vo.Account.CreateTime.Unix(),
		Roles:      vo.Account.Roles.String(),
		Passwd:     vo.Account.Password,
	}
	return item
}
