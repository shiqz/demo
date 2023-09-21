// Package utils 系统内部常用工具包
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

// TokenPrefix 授权登录token前缀
const TokenPrefix = "Bearer"

var random *rand.Rand
var charsets = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var (
	// ErrInvalidToken token 错误
	ErrInvalidToken = errors.New("invalid token")
)

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// EncryptMD5 MD5加密
func EncryptMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetRandomStr 生成随机字符串
func GetRandomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = charsets[random.Intn(len(charsets))]
	}
	return string(b)
}

// HashPassEncrypt 密码加密处理
func HashPassEncrypt(pass []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pass, 0)
}

// HashPassCheck 密码校验
func HashPassCheck(hashPass, pass []byte) error {
	return bcrypt.CompareHashAndPassword(hashPass, pass)
}

// GetRequestToken 获取请求token
func GetRequestToken(auth string) (string, error) {
	if auth == "" || !strings.HasPrefix(auth, TokenPrefix) {
		return "", ErrInvalidToken
	}
	return strings.TrimSpace(strings.TrimLeft(auth, TokenPrefix)), nil
}
