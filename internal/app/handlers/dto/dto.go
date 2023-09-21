// Package dto 请求结构体
package dto

import "example/internal/domain"

const (
	defaultPage  uint = 1
	defaultLimit uint = 30
)

// Filter 基础过滤接口
type Filter interface {
	Convert()
}

// BaseFilter 基础过滤
type BaseFilter struct {
	Page  *uint `json:"page" valid:"-"`
	Limit *uint `json:"limit" valid:"-"`
}

// GetBaseFilter 转换
func (f *BaseFilter) GetBaseFilter() domain.Filter {
	var page, limit = defaultPage, defaultLimit
	if f.Page != nil && *f.Page > 0 {
		page = *f.Page
	}
	if f.Limit != nil && *f.Limit > 0 {
		limit = *f.Limit
	}
	return domain.Filter{
		Offset: (page - 1) * limit,
		Limit:  limit,
	}
}
