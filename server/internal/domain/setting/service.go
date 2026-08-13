package setting

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shenfay/go-react-admin/pkg/errors"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// Service 系统设置业务逻辑
type Service struct {
	repo Repository
}

// NewService 创建系统设置服务
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetAllSettings 获取所有设置（可选按分类过滤）
func (s *Service) GetAllSettings(ctx context.Context, category string) ([]*Setting, error) {
	if category != "" {
		return s.repo.FindByCategory(ctx, category)
	}
	return s.repo.FindAll(ctx)
}

// GetSettingByKey 获取单个设置
func (s *Service) GetSettingByKey(ctx context.Context, key string) (*Setting, error) {
	return s.repo.FindByKey(ctx, key)
}

// SettingUpdate 单条设置更新请求
type SettingUpdate struct {
	Key   string          `json:"key" validate:"required"`
	Value json.RawMessage `json:"value" validate:"required"`
}

// BatchUpdate 批量更新设置
func (s *Service) BatchUpdate(ctx context.Context, updates []SettingUpdate, updatedBy *string) error {
	if len(updates) == 0 {
		return errors.NewAppError(
			errors.ErrCodeSystemInvalidRequest,
			"没有要更新的设置项",
			http.StatusBadRequest,
		)
	}

	settings := make([]*Setting, len(updates))
	for i, u := range updates {
		settings[i] = &Setting{
			Key:       u.Key,
			Value:     u.Value,
			UpdatedBy: updatedBy,
		}
	}

	return s.repo.BatchUpsert(ctx, settings)
}

// GetSettingsMap 获取所有设置并返回 map 结构（方便业务代码读取配置）
func (s *Service) GetSettingsMap(ctx context.Context) (map[string]json.RawMessage, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	m := make(map[string]json.RawMessage, len(all))
	for _, setting := range all {
		m[setting.Key] = setting.Value
	}
	return m, nil
}

// GetStringValue 获取字符串类型的设置值
func GetStringValue(raw json.RawMessage) string {
	var v string
	if err := json.Unmarshal(raw, &v); err != nil {
		return ""
	}
	return v
}

// GetIntValue 获取整数类型的设置值
func GetIntValue(raw json.RawMessage) int {
	return utils.ToInt(string(raw))
}

// GetBoolValue 获取布尔类型的设置值
func GetBoolValue(raw json.RawMessage) bool {
	var v bool
	if err := json.Unmarshal(raw, &v); err != nil {
		return false
	}
	return v
}


