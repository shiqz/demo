package domain

import (
	"context"
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"time"
)

type (
	// SessionService 会话服务
	SessionService interface {
		Set(ctx context.Context, session *entity.Session) error
		Get(ctx context.Context, scene types.SessionScene, id uint) (*entity.Session, error)

		Refresh(ctx context.Context, session *entity.Session) error
		Disconnect(ctx context.Context) error
		Remove(ctx context.Context, id string) error
	}

	// SessionRepository 会话仓库
	SessionRepository interface {
		Get(ctx context.Context, key string) (string, error)

		Save(ctx context.Context, key string, value string, expire time.Duration) error
		Delete(ctx context.Context, key string) error
	}
)
