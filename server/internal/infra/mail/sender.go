package mail

import (
	"context"
	"fmt"

	"github.com/shenfay/go-react-admin/internal/app/port"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

var _ port.EmailSender = (*NoopSender)(nil)

// NoopSender 开发环境邮件发送器
// 仅记录日志，不实际发送邮件
type NoopSender struct{}

// NewNoopSender 创建开发环境邮件发送器
func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

// SendVerificationEmail 记录验证邮件日志
func (s *NoopSender) SendVerificationEmail(ctx context.Context, to, token, userID string) error {
	logger.Info("[NOOP] Send verification email",
		"to", to,
		"user_id", userID,
		"token_prefix", token[:12]+"...",
	)
	return nil
}

// SendPasswordResetEmail 记录密码重置邮件日志
func (s *NoopSender) SendPasswordResetEmail(ctx context.Context, to, token, userID string) error {
	logger.Info("[NOOP] Send password reset email",
		"to", to,
		"user_id", userID,
		"token_prefix", token[:12]+"...",
	)
	return nil
}

// SmtpSender 生产环境邮件发送器
type SmtpSender struct {
	from     string
	host     string
	port     int
	username string
	password string
}

// SmtpConfig SMTP 发件配置
type SmtpConfig struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

// NewSmtpSender 创建 SMTP 邮件发送器
func NewSmtpSender(cfg SmtpConfig) *SmtpSender {
	return &SmtpSender{
		from:     cfg.From,
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
	}
}

// SendVerificationEmail 发送验证邮件
func (s *SmtpSender) SendVerificationEmail(ctx context.Context, to, token, userID string) error {
	// TODO: 实现真实 SMTP 发送
	// 使用 net/smtp 发送 HTML 模板邮件
	logger.Info("[SMTP] Send verification email (not implemented)",
		"to", to,
		"user_id", userID,
	)
	return fmt.Errorf("SMTP sender not yet implemented")
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *SmtpSender) SendPasswordResetEmail(ctx context.Context, to, token, userID string) error {
	// TODO: 实现真实 SMTP 发送
	logger.Info("[SMTP] Send password reset email (not implemented)",
		"to", to,
		"user_id", userID,
	)
	return fmt.Errorf("SMTP sender not yet implemented")
}
