package service

import (
	"context"
	"example/internal/app/errs"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
)

// PermissionService 权限服务
type PermissionService struct {
	repo domain.AccountRepository
}

// NewPermissionService 实例化权限服务
func NewPermissionService(repo domain.AccountRepository) domain.PermissionService {
	return &PermissionService{repo: repo}
}

// CheckPermission 校验权限
func (s *PermissionService) CheckPermission(ctx context.Context, route types.Route) error {
	session := ctx.Value(types.SessionFlag).(*entity.Session)
	account, err := s.repo.GetOne(ctx, session.GetSessionID())
	if err != nil {
		return err
	}
	if !account.HasPermission(route) {
		return errs.EcStatusForbiddenForPerms
	}
	return nil
}
