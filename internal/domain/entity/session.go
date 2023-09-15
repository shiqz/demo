package entity

import (
	"demo/internal/domain/types"
	"demo/internal/pkg/utils"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// Session 会话实体
type Session struct {
	scene       types.SessionScene
	uid         uint
	Token       string `json:"token"`
	RefreshTime int64  `json:"refreshTime"`
}

// Encode 编码Session
func (s *Session) Encode() (string, error) {
	marshal, err := json.Marshal(s)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(marshal), nil
}

// Decode 解析Session
func (s *Session) Decode(info string) error {
	if info == "" {
		return errors.WithStack(types.ErrInvalidSession)
	}
	return errors.WithStack(json.Unmarshal([]byte(info), s))
}

// FormatKey 格式化会话缓存标识
func (s *Session) FormatKey() string {
	prefix := types.UserSessionPrefix
	if s.scene == types.AdminSession {
		prefix = types.AdminSessionPrefix
	}
	return prefix + strconv.FormatInt(int64(s.uid), 10)
}

// GetExpireTime 获取会话过期时间
func (s *Session) GetExpireTime() time.Time {
	return time.Unix(s.RefreshTime, 0)
}

// GetDuration 获取会话有效时长
func (s *Session) GetDuration() time.Duration {
	sub := time.Unix(s.RefreshTime, 0).Sub(time.Now())
	if sub <= 0 {
		return 0
	}
	return sub
}

// FormatToken 格式化Token标识
func (s *Session) FormatToken() string {
	token := strconv.FormatInt(int64(s.uid), 10) + types.SessionSplitFlag + s.Token
	if s.scene == types.AdminSession {
		token = types.AdminTokenPrefix + token
	}
	return token
}

// IsExpired 已过期
func (s *Session) IsExpired() bool {
	return s.RefreshTime-time.Now().Unix() <= 0
}

// NewSession 创建Session
func NewSession(scene types.SessionScene, uid uint) *Session {
	s := &Session{scene: scene, uid: uid}
	s.RefreshTime = time.Now().Add(types.SessionExpireTime).Unix()
	s.Token = utils.GetRandomStr(20)
	return s
}
