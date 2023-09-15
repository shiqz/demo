package types

import (
	"demo/internal/app/errs"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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
	ErrInvalidSession = errors.New("invalid session")
)

// ParseToken 解析Token
func ParseToken(token string) (SessionScene, uint, error) {
	tokenInfo := strings.Split(token, SessionSplitFlag)
	tl := len(tokenInfo)
	if tl < 2 {
		return "", 0, errs.ErrInvalidToken
	}

	// parse uid from token
	id, err := strconv.ParseUint(tokenInfo[tl-2], 10, 64)
	if err != nil {
		return "", 0, errs.ErrInvalidToken
	}
	uid := uint(id)
	scene := UserSession
	if strings.HasPrefix(token, AdminTokenPrefix) {
		scene = AdminSession
	}
	return scene, uid, nil
}
