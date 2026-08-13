package setting

import (
	"context"
	"encoding/json"

	appshared "github.com/shenfay/go-react-admin/internal/app/shared/operationlog"
	domain "github.com/shenfay/go-react-admin/internal/domain/setting"
)

// Service 系统设置应用服务
// 包装 domain.Service，在其基础上附加操作日志记录等应用层职责
type Service struct {
	domainSvc *domain.Service
	recorder  *appshared.OperationRecorder
}

// NewService 创建系统设置应用服务
func NewService(domainSvc *domain.Service, recorder *appshared.OperationRecorder) *Service {
	return &Service{
		domainSvc: domainSvc,
		recorder:  recorder,
	}
}

// GetAllSettings 获取所有设置
func (s *Service) GetAllSettings(ctx context.Context, category string) ([]*domain.Setting, error) {
	return s.domainSvc.GetAllSettings(ctx, category)
}

// GetSettingByKey 获取单个设置
func (s *Service) GetSettingByKey(ctx context.Context, key string) (*domain.Setting, error) {
	return s.domainSvc.GetSettingByKey(ctx, key)
}

// BatchUpdate 批量更新设置
func (s *Service) BatchUpdate(ctx context.Context, updates []domain.SettingUpdate, updatedBy *string) error {
	// 1. 委托领域服务执行校验和持久化
	if err := s.domainSvc.BatchUpdate(ctx, updates, updatedBy); err != nil {
		return err
	}

	// 2. 在应用层记录操作日志
	keys := make([]string, len(updates))
	for i, u := range updates {
		keys[i] = u.Key
	}
	s.recorder.RecordFromContext(ctx, "SYSTEM.CONFIG.UPDATED", "SYSTEM", "SUCCESS",
		map[string]interface{}{"updated_keys": keys},
	)
	return nil
}

// GetSettingsMap 获取所有设置并返回 map
func (s *Service) GetSettingsMap(ctx context.Context) (map[string]json.RawMessage, error) {
	return s.domainSvc.GetSettingsMap(ctx)
}
