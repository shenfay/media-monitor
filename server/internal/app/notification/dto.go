package notification

import (
	"time"

	domain "github.com/shenfay/go-react-admin/internal/domain/notification"
)

// MessageDTO 消息数据传输对象
type MessageDTO struct {
	ID          string                 `json:"id"`
	SenderID    *string                `json:"sender_id,omitempty"`
	RecipientID string                 `json:"recipient_id"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	IsRead      bool                   `json:"is_read"`
	ReadAt      *time.Time             `json:"read_at,omitempty"`
	RefType     string                 `json:"ref_type,omitempty"`
	RefID       string                 `json:"ref_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// MessageListDTO 消息列表 DTO
type MessageListDTO struct {
	Messages []*MessageDTO `json:"messages"`
	Total    int64         `json:"total"`
}

// UnreadCountDTO 未读计数 DTO
type UnreadCountDTO struct {
	Counts []domain.UnreadCount `json:"counts"`
	Total  int64                `json:"total"`
}

// SendMessageCmd 发送消息命令
type SendMessageCmd struct {
	RecipientID string  `json:"recipient_id" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=system companion"`
	Category    string  `json:"category" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Content     string  `json:"content" binding:"required"`
	SenderID    *string `json:"sender_id,omitempty"`
	RefType     string  `json:"ref_type,omitempty"`
	RefID       string  `json:"ref_id,omitempty"`
}

// ReadAllCmd 标记全部已读命令
type ReadAllCmd struct {
	Type string `json:"type,omitempty"`
}

// toMessageDTO 将领域模型转换为 DTO
func toMessageDTO(msg *domain.Message) *MessageDTO {
	return &MessageDTO{
		ID:          msg.ID,
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Type:        string(msg.Type),
		Category:    string(msg.Category),
		Title:       msg.Title,
		Content:     msg.Content,
		IsRead:      msg.IsRead,
		ReadAt:      msg.ReadAt,
		RefType:     msg.RefType,
		RefID:       msg.RefID,
		Metadata:    msg.Metadata,
		CreatedAt:   msg.CreatedAt,
	}
}

// toMessageDTOList 批量转换
func toMessageDTOList(msgs []*domain.Message) []*MessageDTO {
	dtos := make([]*MessageDTO, len(msgs))
	for i, msg := range msgs {
		dtos[i] = toMessageDTO(msg)
	}
	return dtos
}
