package entity

import (
	"encoding/json"
	"example/internal/app/errs"
	"example/internal/domain/types"
	"example/internal/pkg/utils"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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
	return s.state() == types.SessionStateExpired
}

// IsExpireSoon Session 是否快过期
func (s *Session) IsExpireSoon() bool {
	return s.state() == types.SessionStateExpiringSoon
}

// IsNormal Session 是否正常过期
func (s *Session) IsNormal() bool {
	return s.state() == types.SessionStateNormal
}

// Reset 重置
func (s *Session) Reset() {
	s.Token = utils.GetRandomStr(20)
	s.RefreshTime = time.Now().Add(types.SessionExpireTime).Unix()
}

// GetSessionID 获取会话ID
func (s *Session) GetSessionID() uint {
	return s.uid
}

// GetScene 获取会话类型
func (s *Session) GetScene() types.SessionScene {
	return s.scene
}

// Remove 过期token
func (s *Session) Remove() {
	s.RefreshTime = 0
}

// State 获取Session状态
func (s *Session) state() types.SessionState {
	sub := s.RefreshTime - time.Now().Unix()
	if sub <= 0 {
		return types.SessionStateExpired
	} else if sub > 0 && sub < int64(types.SessionExpiringSoonTime.Seconds()) {
		return types.SessionStateExpiringSoon
	}
	return types.SessionStateNormal
}

// NewSession 创建Session
func NewSession(scene types.SessionScene, uid uint) *Session {
	s := &Session{scene: scene, uid: uid}
	s.RefreshTime = time.Now().Add(types.SessionExpireTime).Unix()
	s.Token = utils.GetRandomStr(20)
	return s
}

// LoadSessionByToken 根据 Token 加载会话
func LoadSessionByToken(token string) (*Session, error) {
	tokenInfo := strings.Split(token, types.SessionSplitFlag)
	tl := len(tokenInfo)
	if tl < 2 {
		return nil, errs.ErrInvalidToken
	}
	// parse uid from token
	id, err := strconv.ParseUint(tokenInfo[tl-2], 10, 64)
	if err != nil {
		return nil, errs.ErrInvalidToken
	}
	s := &Session{
		scene: types.UserSession,
		uid:   uint(id),
		Token: tokenInfo[tl-1],
	}
	if strings.HasPrefix(token, types.AdminTokenPrefix) {
		s.scene = types.AdminSession
	}
	return s, nil
}
