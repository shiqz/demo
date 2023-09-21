package types

import (
	"github.com/pkg/errors"
	"time"
)

// SessionScene 会话场景
type SessionScene string

const (
	// AdminSession 管理员会话
	AdminSession SessionScene = "admin"
	// UserSession 普通用户会话
	UserSession SessionScene = "user"
)

// SessionState 会话状态
type SessionState uint

const (
	// SessionStateNormal 正常状态
	SessionStateNormal SessionState = iota + 1
	// SessionStateExpiringSoon 即将过期
	SessionStateExpiringSoon
	// SessionStateExpired 已过期
	SessionStateExpired
)

const (
	// SessionExpiringSoonTime 会话即将过期下限
	SessionExpiringSoonTime = 7 * 24 * time.Hour
	// SessionExpireTime 会话有效期
	SessionExpireTime = 30 * 24 * time.Hour
)

const (
	// SessionFlag 会话标识
	SessionFlag = "session"
	// SessionSplitFlag 会话ID格式分隔标识符
	SessionSplitFlag = "::"
	// AdminTokenPrefix 管理员会话Token ID前缀标识
	AdminTokenPrefix = "admin" + SessionSplitFlag
	// UserSessionPrefix 用户会话缓存标识前缀
	UserSessionPrefix = "u:login:"
	// AdminSessionPrefix 管理员会话缓存标识前缀
	AdminSessionPrefix = "a:login:"
)

var (
	// ErrInvalidSession session 错误
	ErrInvalidSession = errors.New("invalid session")
)
