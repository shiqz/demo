package domain

import (
	"context"
	"demo/internal/domain/entity"
	"demo/internal/domain/types"
	"time"
)

type (
	// SessionService 会话服务
	SessionService interface {
		Set(ctx context.Context, session *entity.Session) error
		Get(ctx context.Context, scene types.SessionScene, id uint) (*entity.Session, error)
		Refresh(ctx context.Context, session *entity.Session) error
		Remove(ctx context.Context) error
	}

	// SessionRepository 会话仓库
	SessionRepository interface {
		Save(ctx context.Context, key string, value string, expire time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		Delete(ctx context.Context, key string) error
	}
)
