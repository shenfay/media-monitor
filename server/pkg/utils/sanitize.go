package utils

import "strings"

// MaskEmail 邮箱脱敏：john@example.com → j***@example.com
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "****"
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 1 {
		return "***@" + domain
	}

	return string(local[0]) + "***@" + domain
}

// MaskToken Token 脱敏：保留前4位和后4位，中间用 **** 替代
func MaskToken(token string) string {
	if token == "" {
		return ""
	}

	if len(token) <= 8 {
		return "****"
	}

	return token[:4] + "****" + token[len(token)-4:]
}

// MaskSecret 通用敏感信息脱敏：保留前2位，其余用 **** 替代
func MaskSecret(s string) string {
	if s == "" {
		return ""
	}

	if len(s) <= 2 {
		return "****"
	}

	return s[:2] + "****"
}
