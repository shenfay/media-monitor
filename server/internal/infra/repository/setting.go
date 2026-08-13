package repository

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/setting"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// settingPO GORM 持久化对象
type settingPO struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	Key         string          `json:"key" gorm:"uniqueIndex;size:100;not null"`
	Value       json.RawMessage `json:"value" gorm:"type:jsonb;not null"`
	Category    string          `json:"category" gorm:"size:50;not null;index:idx_category"`
	Label       string          `json:"label" gorm:"size:200"`
	Description string          `json:"description" gorm:"type:text"`
	UpdatedBy   *string         `json:"updated_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (settingPO) TableName() string {
	return "system_settings"
}

// toDomain 将持久化对象转换为领域模型
func (po *settingPO) toDomain() *setting.Setting {
	return &setting.Setting{
		ID:          po.ID,
		Key:         po.Key,
		Value:       po.Value,
		Category:    po.Category,
		Label:       po.Label,
		Description: po.Description,
		UpdatedBy:   po.UpdatedBy,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// fromDomain 将领域模型转换为持久化对象
func fromSettingDomain(s *setting.Setting) *settingPO {
	return &settingPO{
		ID:          s.ID,
		Key:         s.Key,
		Value:       s.Value,
		Category:    s.Category,
		Label:       s.Label,
		Description: s.Description,
		UpdatedBy:   s.UpdatedBy,
	}
}

// settingRepository GORM 实现
type settingRepository struct {
	db *gorm.DB
}

// NewSettingRepository 创建系统设置仓储
func NewSettingRepository(db *gorm.DB) setting.Repository {
	return &settingRepository{db: db}
}

// FindAll 获取所有设置
func (r *settingRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	var pos []*settingPO
	err := r.db.WithContext(ctx).Order("category, key").Find(&pos).Error
	if err != nil {
		return nil, err
	}

	result := make([]*setting.Setting, len(pos))
	for i, po := range pos {
		result[i] = po.toDomain()
	}
	return result, nil
}

// FindByCategory 按分类获取设置
func (r *settingRepository) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	var pos []*settingPO
	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("key").
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	result := make([]*setting.Setting, len(pos))
	for i, po := range pos {
		result[i] = po.toDomain()
	}
	return result, nil
}

// FindByKey 根据 key 获取单个设置
func (r *settingRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	var po settingPO
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.toDomain(), nil
}

// BatchUpsert 批量更新/插入设置
// 对于已存在的 key 执行 UPDATE，不存在的执行 INSERT
func (r *settingRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range settings {
			po := fromSettingDomain(s)

			// 对敏感字段进行 AES 加密
			if setting.IsSensitiveKey(s.Key) {
				encrypted, err := encryptSensitiveFields(s.Key, s.Value)
				if err != nil {
					return err
				}
				po.Value = encrypted
			}

			// Upsert: 如果 key 存在则更新，不存在则插入
			err := tx.Where("key = ?", s.Key).
				Assign(map[string]interface{}{
					"value":      string(po.Value),
					"updated_by": po.UpdatedBy,
					"updated_at": time.Now(),
				}).
				FirstOrCreate(&po).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// encryptSensitiveFields 对渠道配置 JSON 中的敏感字段进行加密
func encryptSensitiveFields(key string, value json.RawMessage) (json.RawMessage, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(value, &obj); err != nil {
		return value, nil // 非 JSON 对象，原样返回
	}

	// 需要加密的字段名
	sensitiveFieldNames := []string{"password", "secret"}

	for _, fieldName := range sensitiveFieldNames {
		if v, ok := obj[fieldName]; ok {
			if str, isStr := v.(string); isStr && str != "" {
				// 如果已经是加密后的密文（base64 且较长），跳过
				if isLikelyEncrypted(str) {
					continue
				}
				encrypted, err := utils.Encrypt(str)
				if err != nil {
					return nil, err
				}
				obj[fieldName] = encrypted
			}
		}
	}

	return json.Marshal(obj)
}

// isLikelyEncrypted 简单判断字符串是否可能是已加密的密文
func isLikelyEncrypted(s string) bool {
	// base64 编码的 AES-GCM 密文通常较长（至少 40+ 字符）
	return len(s) > 40
}
