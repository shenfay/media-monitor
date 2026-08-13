package user

import (
	"crypto/rand"
	"encoding/base64"
)

// SecureToken 安全随机令牌值对象
// 使用 crypto/rand 生成 32 字节随机数，Base64URL 编码
type SecureToken string

// NewSecureToken 生成安全随机令牌
func NewSecureToken() (SecureToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return SecureToken(base64.URLEncoding.EncodeToString(b)), nil
}

// String 返回令牌字符串
func (t SecureToken) String() string {
	return string(t)
}
