package events

import "time"

// OperationEvent 统一操作日志事件
// 用于 Service 层记录操作日志，替代直接发布特定领域事件
type OperationEvent struct {
	Action    string                 `json:"action"`     // AUTH.LOGIN.SUCCESS / USER.PROFILE.UPDATED / ...
	Category  string                 `json:"category"`   // AUTH / USER / SYSTEM / BIZ
	Status    string                 `json:"status"`     // SUCCESS / FAILED
	UserID    string                 `json:"user_id"`    // 用户 ID
	Email     string                 `json:"email"`      // 用户邮箱
	IP        string                 `json:"ip"`         // IP 地址
	UserAgent string                 `json:"user_agent"` // User-Agent
	Device    string                 `json:"device"`     // 设备类型
	Browser   string                 `json:"browser"`    // 浏览器
	OS        string                 `json:"os"`         // 操作系统
	Metadata  map[string]interface{} `json:"metadata"`   // 额外元数据
	Timestamp time.Time              `json:"timestamp"`  // 事件发生时间
}

// NewOperationEvent 创建操作日志事件
func NewOperationEvent(action, category, status string) *OperationEvent {
	return &OperationEvent{
		Action:    action,
		Category:  category,
		Status:    status,
		Timestamp: time.Now(),
	}
}

// WithUser 设置用户信息
func (e *OperationEvent) WithUser(userID, email string) *OperationEvent {
	e.UserID = userID
	e.Email = email
	return e
}

// WithRequestInfo 设置请求信息
func (e *OperationEvent) WithRequestInfo(ip, userAgent, device, browser, os string) *OperationEvent {
	e.IP = ip
	e.UserAgent = userAgent
	e.Device = device
	e.Browser = browser
	e.OS = os
	return e
}

// WithMetadata 设置额外元数据
func (e *OperationEvent) WithMetadata(metadata map[string]interface{}) *OperationEvent {
	e.Metadata = metadata
	return e
}

// EventName 返回事件名称（实现 DomainEvent 接口）
func (e *OperationEvent) EventName() string {
	return "operation.log"
}

// OccurredAt 返回事件发生时间（实现 DomainEvent 接口）
func (e *OperationEvent) OccurredAt() time.Time {
	return e.Timestamp
}
