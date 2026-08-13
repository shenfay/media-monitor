package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-react-admin/internal/domain/notification"
	"github.com/shenfay/go-react-admin/pkg/utils"
	"gorm.io/gorm"
)

// messageRepository GORM 实现
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息仓储
func NewMessageRepository(db *gorm.DB) notification.MessageRepository {
	return &messageRepository{db: db}
}

// Save 保存消息
func (r *messageRepository) Save(ctx context.Context, msg *notification.Message) error {
	if msg.ID == "" {
		msg.ID = utils.GenerateID()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(fromDomainMessage(msg)).Error
}

// FindByRecipient 按接收者查询消息列表（支持筛选 + 分页）
func (r *messageRepository) FindByRecipient(ctx context.Context, params notification.MessageListParams) ([]*notification.Message, int64, error) {
	var pos []*messagePO
	var total int64

	query := r.db.WithContext(ctx).Model(&messagePO{}).Where("recipient_id = ?", params.RecipientID)

	if params.Type != "" {
		query = query.Where("type = ?", string(params.Type))
	}
	if params.Category != "" {
		query = query.Where("category = ?", string(params.Category))
	}
	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return toDomainMessageList(pos), total, nil
}

// CountUnread 统计未读消息数（按类型分组）
func (r *messageRepository) CountUnread(ctx context.Context, recipientID string) ([]notification.UnreadCount, error) {
	type result struct {
		Type  string `gorm:"column:type"`
		Count int64  `gorm:"column:count"`
	}
	var results []result

	err := r.db.WithContext(ctx).
		Model(&messagePO{}).
		Select("type, count(*) as count").
		Where("recipient_id = ? AND is_read = ?", recipientID, false).
		Group("type").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make([]notification.UnreadCount, len(results))
	for i, r := range results {
		counts[i] = notification.UnreadCount{
			Type:  notification.MessageType(r.Type),
			Count: r.Count,
		}
	}
	return counts, nil
}

// FindByID 根据 ID 查找消息
func (r *messageRepository) FindByID(ctx context.Context, id string) (*notification.Message, error) {
	var po messagePO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		return nil, err
	}
	return po.toDomain(), nil
}

// MarkAsRead 标记单条消息已读
func (r *messageRepository) MarkAsRead(ctx context.Context, id, recipientID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&messagePO{}).
		Where("id = ? AND recipient_id = ?", id, recipientID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead 标记全部已读（可按类型）
func (r *messageRepository) MarkAllAsRead(ctx context.Context, recipientID string, msgType notification.MessageType) error {
	now := time.Now()
	query := r.db.WithContext(ctx).
		Model(&messagePO{}).
		Where("recipient_id = ? AND is_read = ?", recipientID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	if msgType != "" {
		query = query.Where("type = ?", string(msgType))
	}

	return query.Error
}

// FindAll 管理员查询所有消息（支持筛选 + 分页）
func (r *messageRepository) FindAll(ctx context.Context, msgType notification.MessageType, category notification.MessageCategory, limit, offset int) ([]*notification.Message, int64, error) {
	var pos []*messagePO
	var total int64

	query := r.db.WithContext(ctx).Model(&messagePO{})

	if msgType != "" {
		query = query.Where("type = ?", string(msgType))
	}
	if category != "" {
		query = query.Where("category = ?", string(category))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return toDomainMessageList(pos), total, nil
}
