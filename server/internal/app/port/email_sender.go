package port

import "context"

// EmailSender 邮件发送出站端口
// 由 app 层定义，infra 层实现（适配器模式）
type EmailSender interface {
	// SendVerificationEmail 发送邮箱验证邮件
	SendVerificationEmail(ctx context.Context, to, token, userID string) error

	// SendPasswordResetEmail 发送密码重置邮件
	SendPasswordResetEmail(ctx context.Context, to, token, userID string) error
}
