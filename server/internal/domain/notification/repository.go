package notification

import "context"

// MessageRepository 消息仓储接口
type MessageRepository interface {
	// Save 保存消息
	Save(ctx context.Context, msg *Message) error

	// FindByRecipient 按接收者查询消息列表（支持筛选 + 分页）
	FindByRecipient(ctx context.Context, params MessageListParams) ([]*Message, int64, error)

	// CountUnread 统计未读消息数（按类型分组）
	CountUnread(ctx context.Context, recipientID string) ([]UnreadCount, error)

	// FindByID 根据 ID 查找消息
	FindByID(ctx context.Context, id string) (*Message, error)

	// MarkAsRead 标记单条消息已读
	MarkAsRead(ctx context.Context, id, recipientID string) error

	// MarkAllAsRead 标记全部已读（可按类型）
	MarkAllAsRead(ctx context.Context, recipientID string, msgType MessageType) error

	// FindAll 管理员查询所有消息（支持筛选 + 分页）
	FindAll(ctx context.Context, msgType MessageType, category MessageCategory, limit, offset int) ([]*Message, int64, error)
}
