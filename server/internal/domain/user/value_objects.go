package user

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// emailRegex 邮箱格式校验正则（RFC 5322 简化版）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Email 邮箱值对象
// 封装邮箱格式验证逻辑，确保领域层不接受无效邮箱
type Email struct {
	address string
}

// NewEmail 创建邮箱值对象
// 格式无效时返回错误
func NewEmail(address string) (Email, error) {
	if !emailRegex.MatchString(address) {
		return Email{}, fmt.Errorf("invalid email format: %s", address)
	}
	return Email{address: address}, nil
}

// String 返回邮箱字符串
func (e Email) String() string {
	return e.address
}

// Password 密码值对象
// 封装密码哈希和验证逻辑
type Password struct {
	hash string
}

// NewPassword 从明文密码创建密码值对象
// 密码长度不足 8 位时返回错误
func NewPassword(plain string) (Password, error) {
	if len(plain) < 8 {
		return Password{}, fmt.Errorf("password must be at least 8 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, fmt.Errorf("failed to hash password: %w", err)
	}
	return Password{hash: string(hash)}, nil
}

// Hash 返回密码哈希字符串
func (p Password) Hash() string {
	return p.hash
}

// Verify 验证明文密码是否与哈希匹配
func (p Password) Verify(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain)) == nil
}

// HashPassword 哈希密码（保留向后兼容的顶层函数）
func HashPassword(password string) (string, error) {
	p, err := NewPassword(password)
	if err != nil {
		return "", err
	}
	return p.Hash(), nil
}

// verifyPasswordHash 验证明文密码是否与已有哈希匹配
func verifyPasswordHash(hashedPassword, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}
