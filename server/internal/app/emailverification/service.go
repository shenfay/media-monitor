package emailverification

import (
	"context"

	"github.com/shenfay/go-react-admin/internal/app/port"
	appevents "github.com/shenfay/go-react-admin/internal/app/shared/events"
	"github.com/shenfay/go-react-admin/internal/app/shared/operationlog"
	"github.com/shenfay/go-react-admin/internal/app/tokenmanager"
	"github.com/shenfay/go-react-admin/internal/domain/shared/events"
	"github.com/shenfay/go-react-admin/internal/domain/user"
)

// Service 邮箱验证应用服务
type Service struct {
	tokenManager *tokenmanager.TokenManager
	userRepo     user.UserRepository
	sender       port.EmailSender
	eventBus     events.Bus
	recorder     *operationlog.OperationRecorder
}

// NewService 创建邮箱验证服务
func NewService(
	tokenManager *tokenmanager.TokenManager,
	userRepo user.UserRepository,
	sender port.EmailSender,
	eventBus events.Bus,
	recorder *operationlog.OperationRecorder,
) *Service {
	return &Service{
		tokenManager: tokenManager,
		userRepo:     userRepo,
		sender:       sender,
		eventBus:     eventBus,
		recorder:     recorder,
	}
}

// SendVerificationEmail 发送验证邮件
func (s *Service) SendVerificationEmail(ctx context.Context, userID, email string) error {
	// 1. 生成验证令牌
	token, err := s.tokenManager.CreateVerificationToken(ctx, userID)
	if err != nil {
		return err
	}

	// 2. 通过事件总线发布事件 → Bridge → Asynq → Worker 发送邮件
	evt := appevents.NewSendEmailEvent("verification", email, token.Token, userID)
	if err := s.eventBus.Publish(ctx, evt); err != nil {
		return err
	}

	// 3. 记录操作日志
	s.recorder.Record(ctx, "AUTH.SEND_VERIFICATION", "AUTH", "SUCCESS",
		userID, email,
		map[string]interface{}{"email": email},
	)

	return nil
}

// VerifyEmail 验证邮箱
func (s *Service) VerifyEmail(ctx context.Context, token, userID string) error {
	// 1. 验证并消费令牌
	if err := s.tokenManager.VerifyAndConsumeVerificationToken(ctx, token, userID); err != nil {
		return err
	}

	// 2. 获取用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 3. 标记邮箱已验证
	u.VerifyEmail()
	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	// 4. 记录操作日志
	s.recorder.Record(ctx, "AUTH.VERIFY_EMAIL", "AUTH", "SUCCESS",
		userID, u.Email,
		map[string]interface{}{"email": u.Email},
	)

	return nil
}
