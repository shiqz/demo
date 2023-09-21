// Package redisrepoimpl redis 仓库实现
package redisrepoimpl

import (
	"context"
	"example/internal/domain"
	"example/internal/pkg/db"
	"time"
)

// SessionRepository 会话仓库
type SessionRepository struct {
	redis *db.Redis
}

// NewSessionRepository 创建会话仓库
func NewSessionRepository(redis *db.Redis) domain.SessionRepository {
	return &SessionRepository{redis: redis}
}

// Save 保持会话
func (s *SessionRepository) Save(ctx context.Context, key string, value string, expire time.Duration) error {
	return s.redis.Set(ctx, key, value, expire).Err()
}

// Get 获取会话
func (s *SessionRepository) Get(ctx context.Context, key string) (string, error) {
	return s.redis.Get(ctx, key).Result()
}

// Delete 删除会话
func (s *SessionRepository) Delete(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key).Err()
}
