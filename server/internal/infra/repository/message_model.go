package repository

import (
	"time"

	"github.com/shenfay/go-react-admin/internal/domain/notification"
)

// messagePO GORM 持久化对象
type messagePO struct {
	ID          string                 `gorm:"primaryKey;type:varchar(50)"`
	SenderID    *string                `gorm:"type:varchar(50)"`
	RecipientID string                 `gorm:"type:varchar(50);not null;index:idx_messages_recipient"`
	Type        string                 `gorm:"type:varchar(20);not null;index:idx_messages_type_category"`
	Category    string                 `gorm:"type:varchar(30);not null;index:idx_messages_type_category"`
	Title       string                 `gorm:"type:varchar(200);not null"`
	Content     string                 `gorm:"type:text;not null"`
	IsRead      bool                   `gorm:"default:false;index:idx_messages_recipient"`
	ReadAt      *time.Time             `gorm:"type:timestamptz"`
	RefType     string                 `gorm:"type:varchar(30)"`
	RefID       string                 `gorm:"type:varchar(50)"`
	Metadata    map[string]interface{} `gorm:"type:jsonb;default:'{}'::jsonb;serializer:json"`
	CreatedAt   time.Time              `gorm:"not null;index:idx_messages_recipient;default:now()"`
}

// TableName 指定表名
func (messagePO) TableName() string {
	return "messages"
}

// toDomain 将持久化对象转换为领域模型
func (po *messagePO) toDomain() *notification.Message {
	msg := &notification.Message{
		ID:          po.ID,
		SenderID:    po.SenderID,
		RecipientID: po.RecipientID,
		Type:        notification.MessageType(po.Type),
		Category:    notification.MessageCategory(po.Category),
		Title:       po.Title,
		Content:     po.Content,
		IsRead:      po.IsRead,
		ReadAt:      po.ReadAt,
		RefType:     po.RefType,
		RefID:       po.RefID,
		Metadata:    po.Metadata,
		CreatedAt:   po.CreatedAt,
	}
	return msg
}

// fromDomain 将领域模型转换为持久化对象
func fromDomainMessage(msg *notification.Message) *messagePO {
	return &messagePO{
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

// toDomainList 批量转换持久化对象为领域模型
func toDomainMessageList(pos []*messagePO) []*notification.Message {
	result := make([]*notification.Message, len(pos))
	for i, po := range pos {
		result[i] = po.toDomain()
	}
	return result
}
