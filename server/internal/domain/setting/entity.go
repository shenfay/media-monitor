package setting

import (
	"encoding/json"
	"time"
)

// Setting 系统设置实体
type Setting struct {
	ID          int64           `json:"id"`
	Key         string          `json:"key"`
	Value       json.RawMessage `json:"value"`
	Category    string          `json:"category"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	UpdatedBy   *string         `json:"updated_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// SettingCategory 设置分类常量
const (
	CategoryBasic        = "basic"
	CategoryToggle       = "toggle"
	CategoryBusiness     = "business"
	CategoryNotification = "notification"
)

// SensitiveKeys 需要 AES 加密的敏感设置 key 列表
var SensitiveKeys = map[string]bool{
	"channel_email.password": true,
	"channel_webhook.secret": true,
}

// IsSensitiveKey 判断设置值中是否包含需要加密的敏感字段
// 对于渠道配置类 JSON 对象，内部包含 password/secret 字段
func IsSensitiveKey(key string) bool {
	return key == "channel_email" || key == "channel_webhook"
}
