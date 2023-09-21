package domain

import (
	"context"
	"example/internal/domain/types"
)

// PermissionService 权限服务接口
type PermissionService interface {
	CheckPermission(ctx context.Context, route types.Route) error
}
