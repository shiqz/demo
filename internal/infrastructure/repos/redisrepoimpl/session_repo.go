// Package redisrepoimpl redis 仓库实现
package redisrepoimpl

import (
	"context"
	"example/internal/domain"
	"example/internal/infrastructure/depend"
	"time"
)

// SessionRepository 会话仓库
type SessionRepository struct {
	inject *depend.Injecter
}

// NewSessionRepository 创建会话仓库
func NewSessionRepository(inject *depend.Injecter) domain.SessionRepository {
	return &SessionRepository{inject}
}

// Save 保持会话
func (s *SessionRepository) Save(ctx context.Context, key string, value string, expire time.Duration) error {
	return s.inject.GetRedis().Set(ctx, key, value, expire).Err()
}

// Get 获取会话
func (s *SessionRepository) Get(ctx context.Context, key string) (string, error) {
	return s.inject.GetRedis().Get(ctx, key).Result()
}

// Delete 删除会话
func (s *SessionRepository) Delete(ctx context.Context, key string) error {
	return s.inject.GetRedis().Del(ctx, key).Err()
}
