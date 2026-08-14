// Package operationlog 提供跨服务共享的操作日志记录器
// 消除 authentication/admin/setting 三个 Service 中 recordOperation 的重复实现
package operationlog

import (
	"context"

	appevents "github.com/shenfay/go-react-admin/internal/app/shared/events"
	"github.com/shenfay/go-react-admin/internal/domain/shared/events"
	"github.com/shenfay/go-react-admin/pkg/logger"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// OperationRecorder 操作日志记录器
// 封装 OperationEvent 构建 + 事件总线发布 + 失败降级逻辑
type OperationRecorder struct {
	eventBus events.Bus
}

// NewOperationRecorder 创建操作日志记录器
// eventBus 为 nil 时静默跳过（兼容测试场景）
func NewOperationRecorder(eventBus events.Bus) *OperationRecorder {
	return &OperationRecorder{eventBus: eventBus}
}

// Record 记录操作日志（显式传入用户信息）
// 适用于认证场景（Login/Register/Logout），此时用户信息尚未注入 context
// 事件发布采用异步方式，避免阻塞主请求链路（如登录流程）
func (r *OperationRecorder) Record(ctx context.Context, action, category, status string, userID, email string, metadata map[string]interface{}) {
	if r.eventBus == nil {
		return
	}

	evt := appevents.NewOperationEvent(action, category, status).
		WithUser(userID, email).
		WithRequestInfo(
			utils.GetRequestIP(ctx),
			utils.GetRequestUserAgent(ctx),
			utils.GetRequestDevice(ctx),
			utils.GetRequestBrowser(ctx),
			utils.GetRequestOS(ctx),
		).
		WithMetadata(metadata)

	// 异步发布事件，避免 Bridge → Asynq 入队阻塞主请求
	go func() {
		// 使用 background context，避免请求 context 取消导致事件丢失
		bgCtx := context.Background()
		if err := r.eventBus.Publish(bgCtx, evt); err != nil {
			logger.Warn("Failed to record operation log",
				"action", action,
				"user_id", userID,
				"error", err,
			)
		}
	}()
}

// RecordFromContext 记录操作日志（从 context 自动提取操作人信息）
// 适用于已认证场景（admin/setting），操作人信息由 JWT 中间件注入
func (r *OperationRecorder) RecordFromContext(ctx context.Context, action, category, status string, metadata map[string]interface{}) {
	r.Record(ctx, action, category, status,
		utils.GetOperatorUserID(ctx),
		utils.GetOperatorEmail(ctx),
		metadata,
	)
}
