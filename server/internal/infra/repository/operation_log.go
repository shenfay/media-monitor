package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-react-admin/internal/domain/operation"
	"github.com/shenfay/go-react-admin/pkg/utils"
	"gorm.io/gorm"
)

// operationLogPO GORM 持久化对象
type operationLogPO struct {
	ID        string                 `json:"id" gorm:"primaryKey"`
	UserID    string                 `json:"user_id" gorm:"not null;index:idx_user_id"`
	Email     string                 `json:"email" gorm:"index:idx_email"`
	Action    string                 `json:"action" gorm:"not null;index:idx_action"`
	Category  string                 `json:"category" gorm:"not null;index:idx_category"`
	Status    string                 `json:"status" gorm:"not null;index:idx_status"`
	IP        string                 `json:"ip" gorm:"size:45"`
	UserAgent string                 `json:"user_agent" gorm:"size:500"`
	Device    string                 `json:"device" gorm:"size:100"`
	Browser   string                 `json:"browser" gorm:"size:50"`
	OS        string                 `json:"os" gorm:"size:50"`
	Metadata  map[string]interface{} `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb;serializer:json"`
	CreatedAt time.Time              `json:"created_at" gorm:"not null;index:idx_created_at"`
}

// TableName 指定表名
func (operationLogPO) TableName() string {
	return "operation_logs"
}

// toDomain 将持久化对象转换为领域模型
func (po *operationLogPO) toDomain() *operation.OperationLog {
	return &operation.OperationLog{
		ID:        po.ID,
		UserID:    po.UserID,
		Email:     po.Email,
		Action:    po.Action,
		Category:  po.Category,
		Status:    po.Status,
		IP:        po.IP,
		UserAgent: po.UserAgent,
		Device:    po.Device,
		Browser:   po.Browser,
		OS:        po.OS,
		Metadata:  po.Metadata,
		CreatedAt: po.CreatedAt,
	}
}

// fromDomain 将领域模型转换为持久化对象
func fromDomain(log *operation.OperationLog) *operationLogPO {
	return &operationLogPO{
		ID:        log.ID,
		UserID:    log.UserID,
		Email:     log.Email,
		Action:    log.Action,
		Category:  log.Category,
		Status:    log.Status,
		IP:        log.IP,
		UserAgent: log.UserAgent,
		Device:    log.Device,
		Browser:   log.Browser,
		OS:        log.OS,
		Metadata:  log.Metadata,
		CreatedAt: log.CreatedAt,
	}
}

// operationLogRepository GORM 实现
type operationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository 创建操作日志仓储
func NewOperationLogRepository(db *gorm.DB) operation.LogRepository {
	return &operationLogRepository{db: db}
}

// Save 保存操作日志
func (r *operationLogRepository) Save(ctx context.Context, log *operation.OperationLog) error {
	if log.ID == "" {
		log.ID = utils.GenerateID()
	}

	return r.db.WithContext(ctx).Create(fromDomain(log)).Error
}

// FindWithFilter 根据过滤条件查询日志（统一查询入口）
// 替代原有的 FindByUserID/FindByCategory/FindByAction/FindAll 四个方法
func (r *operationLogRepository) FindWithFilter(ctx context.Context, filter operation.LogFilter) ([]*operation.OperationLog, error) {
	var pos []*operationLogPO

	query := r.db.WithContext(ctx).Model(&operationLogPO{})
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}

	err := query.Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&pos).Error

	return toDomainList(pos), err
}

// Count 统计日志总数（统一使用 LogFilter 过滤）
func (r *operationLogRepository) Count(ctx context.Context, filter operation.LogFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&operationLogPO{})
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// toDomainList 批量转换持久化对象为领域模型
func toDomainList(pos []*operationLogPO) []*operation.OperationLog {
	result := make([]*operation.OperationLog, len(pos))
	for i, po := range pos {
		result[i] = po.toDomain()
	}
	return result
}
