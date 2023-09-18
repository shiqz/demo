// Package po 定义相关数据库结构映射
package po

import (
	"demo/internal/domain"
)

// Account 用户表映射
type Account struct {
	UserID     uint   `db:"user_id"`
	Username   string `db:"username"`
	Passwd     string `db:"passwd"`
	Salt       string `db:"salt"`
	Nickname   string `db:"nickname"`
	Gender     uint   `db:"gender"`
	Status     uint   `db:"status"`
	CreateTime int64  `db:"create_time"`
}

func (*Account) TableName() string {
	return "admins"
}

// AccountConvertor 用户数据转换
type AccountConvertor struct {
}

// CreateEntity op 转为 aggregate
func (uc *AccountConvertor) CreateEntity(u Account) *domain.AccountAggregate {
	return &domain.AccountAggregate{}
}

// CreatePO aggregate -> PO
func (uc *AccountConvertor) CreatePO(ug *domain.AccountAggregate) *Account {
	return &Account{}
}
