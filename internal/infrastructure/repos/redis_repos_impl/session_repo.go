// Package redis_repos_impl redis 仓库实现
package redis_repos_impl

import (
	"context"
	"demo/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

// SessionRepository 会话仓库
type SessionRepository struct {
	redis *redis.Client
}

// NewSessionRepository 创建会话仓库
func NewSessionRepository(redis *redis.Client) domain.SessionRepository {
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
