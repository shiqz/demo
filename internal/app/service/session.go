package service

import (
	"context"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
)

// SessionService 会话服务
type SessionService struct {
	repo domain.SessionRepository
}

// NewSessionService 创建会话服务
func NewSessionService(repo domain.SessionRepository) domain.SessionService {
	return &SessionService{repo: repo}
}

// Set 设置会话
func (s *SessionService) Set(ctx context.Context, session *entity.Session) error {
	info, err := session.Encode()
	if err != nil {
		return err
	}
	return s.repo.Save(ctx, session.FormatKey(), info, session.GetDuration())
}

// Get 获取Session
func (s *SessionService) Get(ctx context.Context, scene types.SessionScene, id uint) (*entity.Session, error) {
	session := entity.NewSession(scene, id)
	cacheInfo, err := s.repo.Get(ctx, session.FormatKey())
	if err != nil {
		return nil, err
	}
	if err = session.Decode(cacheInfo); err != nil {
		return nil, err
	}
	return session, nil
}

// Refresh 刷新会话
func (s *SessionService) Refresh(ctx context.Context, session *entity.Session) error {
	session.Reset()
	info, err := session.Encode()
	if err != nil {
		return err
	}
	return s.repo.Save(ctx, session.FormatKey(), info, session.GetDuration())
}

// Disconnect 断开会话
func (s *SessionService) Disconnect(ctx context.Context) error {
	session := ctx.Value(types.SessionFlag).(*entity.Session)
	if err := s.repo.Delete(ctx, session.FormatKey()); err != nil {
		return err
	}
	session.Remove()
	return nil
}

// Remove 移除会话
func (s *SessionService) Remove(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
