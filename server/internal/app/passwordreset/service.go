package passwordreset

import (
	"context"
	"net/http"

	"github.com/shenfay/go-react-admin/internal/app/port"
	appevents "github.com/shenfay/go-react-admin/internal/app/shared/events"
	"github.com/shenfay/go-react-admin/internal/app/shared/operationlog"
	"github.com/shenfay/go-react-admin/internal/app/tokenmanager"
	"github.com/shenfay/go-react-admin/internal/domain/shared/events"
	"github.com/shenfay/go-react-admin/internal/domain/user"
	"github.com/shenfay/go-react-admin/pkg/errors"
)

// Service 密码重置应用服务
type Service struct {
	tokenManager *tokenmanager.TokenManager
	userRepo     user.UserRepository
	sender       port.EmailSender
	eventBus     events.Bus
	recorder     *operationlog.OperationRecorder
}

// NewService 创建密码重置服务
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

// ForgotPassword 请求密码重置
// 始终返回成功（防枚举），仅当邮箱存在时发送重置邮件
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// 邮箱不存在时静默成功（防枚举）
		return nil
	}

	// 生成重置令牌
	token, err := s.tokenManager.CreateResetToken(ctx, u.ID)
	if err != nil {
		return err
	}

	// 发布事件 → Worker 发送邮件
	evt := appevents.NewSendEmailEvent("password_reset", email, token.Token, u.ID)
	if err := s.eventBus.Publish(ctx, evt); err != nil {
		return err
	}

	s.recorder.Record(ctx, "AUTH.FORGOT_PASSWORD", "AUTH", "SUCCESS",
		u.ID, email, nil,
	)

	return nil
}

// ResetPassword 执行密码重置
func (s *Service) ResetPassword(ctx context.Context, token, userID, newPassword string) error {
	// 1. 验证并消费令牌
	if err := s.tokenManager.VerifyAndConsumeResetToken(ctx, token, userID); err != nil {
		return err
	}

	// 2. 获取用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 3. 修改密码
	if err := u.ChangePassword(newPassword); err != nil {
		return errors.NewAppError(
			errors.ErrCodeSystemInvalidRequest,
			"密码不符合要求",
			http.StatusBadRequest,
		)
	}

	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	s.recorder.Record(ctx, "AUTH.RESET_PASSWORD", "AUTH", "SUCCESS",
		userID, u.Email, nil,
	)

	return nil
}
