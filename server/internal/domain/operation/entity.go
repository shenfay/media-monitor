package operation

import "time"

// OperationLog 统一操作日志实体
// 合并原 AuditLog（安全审计）与 ActivityLog（业务活动）
type OperationLog struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Email     string                 `json:"email"`
	Action    string                 `json:"action"`   // AUTH.LOGIN.SUCCESS / USER.PROFILE.UPDATED / ...
	Category  string                 `json:"category"` // AUTH / USER / SYSTEM / BIZ
	Status    string                 `json:"status"`   // SUCCESS / FAILED
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Device    string                 `json:"device"`
	Browser   string                 `json:"browser"`
	OS        string                 `json:"os"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}
