package commands

import (
	"context"
	"example/internal/domain"
	"example/internal/domain/entity"
	"example/internal/domain/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// UserHandler 用户命令行操作
type UserHandler struct {
	sessionService domain.SessionService
}

// NewUserHandler 实例化
func NewUserHandler(s domain.SessionService) *UserHandler {
	return &UserHandler{sessionService: s}
}

// UpdateSession 更新会话过期时间
func (c *UserHandler) UpdateSession(uid, t string) error {
	id, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("请输入操作用户ID")
	}
	expire, err := strconv.ParseUint(t, 10, 64)
	if expire == 0 {
		return errors.New("请输入操作会话过期时间戳")
	}
	ctx := context.Background()
	session, err := c.sessionService.Get(ctx, types.UserSession, uint(id))
	if err != nil {
		return err
	}
	session.RefreshTime = int64(expire)
	if err = c.sessionService.Set(ctx, session); err != nil {
		return err
	}
	log.Infof("会话刷新成功，过期时间为：%s", session.GetExpireTime().Format(time.DateTime))
	return nil
}

// RemoveSession 删除会话
func (c *UserHandler) RemoveSession(uid string) error {
	id, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.New("请输入操作用户ID")
	}
	s := entity.NewSession(types.UserSession, uint(id))
	if err = c.sessionService.Remove(context.Background(), s.FormatKey()); err != nil {
		return err
	}
	log.Info("done")
	return nil
}
