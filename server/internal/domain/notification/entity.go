package notification

import (
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeSystem    MessageType = "system"    // 系统通知
	MessageTypeCompanion MessageType = "companion" // 伙伴对话
)

// MessageCategory 消息分类
type MessageCategory string

const (
	// 系统通知分类
	CategoryVerification    MessageCategory = "verification"     // 验收
	CategoryPoints          MessageCategory = "points"           // 积分
	CategoryGoal            MessageCategory = "goal"             // 目标
	CategoryCompanionStatus MessageCategory = "companion_status" // 伙伴状态
	CategoryExchange        MessageCategory = "exchange"         // 兑换

	// 伙伴对话分类
	CategoryCompanionEncourage MessageCategory = "companion_encourage" // 鼓励
	CategoryCompanionRemind    MessageCategory = "companion_remind"    // 提醒
	CategoryCompanionCelebrate MessageCategory = "companion_celebrate" // 庆祝
)

// Message 消息实体
type Message struct {
	ID          string                 `json:"id"`
	SenderID    *string                `json:"sender_id,omitempty"` // 发送者（系统通知为 nil）
	RecipientID string                 `json:"recipient_id"`        // 接收者
	Type        MessageType            `json:"type"`                // 消息类型
	Category    MessageCategory        `json:"category"`            // 业务分类
	Title       string                 `json:"title"`               // 标题
	Content     string                 `json:"content"`             // 内容
	IsRead      bool                   `json:"is_read"`             // 是否已读
	ReadAt      *time.Time             `json:"read_at,omitempty"`   // 已读时间
	RefType     string                 `json:"ref_type,omitempty"`  // 关联实体类型
	RefID       string                 `json:"ref_id,omitempty"`    // 关联实体 ID
	Metadata    map[string]interface{} `json:"metadata,omitempty"`  // 扩展元数据
	CreatedAt   time.Time              `json:"created_at"`          // 创建时间
}

// NewSystemMessage 创建系统通知
func NewSystemMessage(recipientID string, category MessageCategory, title, content string) *Message {
	return &Message{
		RecipientID: recipientID,
		Type:        MessageTypeSystem,
		Category:    category,
		Title:       title,
		Content:     content,
		IsRead:      false,
	}
}

// NewCompanionMessage 创建伙伴对话消息
func NewCompanionMessage(senderID, recipientID string, category MessageCategory, title, content string) *Message {
	return &Message{
		SenderID:    &senderID,
		RecipientID: recipientID,
		Type:        MessageTypeCompanion,
		Category:    category,
		Title:       title,
		Content:     content,
		IsRead:      false,
	}
}

// WithRef 设置关联实体
func (m *Message) WithRef(refType, refID string) *Message {
	m.RefType = refType
	m.RefID = refID
	return m
}

// WithMetadata 设置扩展元数据
func (m *Message) WithMetadata(metadata map[string]interface{}) *Message {
	m.Metadata = metadata
	return m
}

// MarkAsRead 标记为已读
func (m *Message) MarkAsRead() {
	now := time.Now()
	m.IsRead = true
	m.ReadAt = &now
}

// MessageListParams 消息列表查询参数
type MessageListParams struct {
	RecipientID string
	Type        MessageType
	Category    MessageCategory
	IsRead      *bool
	Limit       int
	Offset      int
}

// UnreadCount 未读计数（按类型分组）
type UnreadCount struct {
	Type  MessageType `json:"type"`
	Count int64       `json:"count"`
}
