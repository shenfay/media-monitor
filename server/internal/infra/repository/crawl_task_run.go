package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
)

// taskRunPO 任务运行持久化对象
type taskRunPO struct {
	ID          string     `gorm:"primaryKey;type:varchar(26)"`
	SourceID    string     `gorm:"column:source_id;type:varchar(26);index"`
	TaskType    string     `gorm:"column:task_type;size:30"`
	TriggeredBy string     `gorm:"column:triggered_by;size:20"`
	Status      string     `gorm:"size:20;index"`
	Total       int        `gorm:"default:0"`
	Ingested    int        `gorm:"default:0"`
	Failed      int        `gorm:"column:failed;default:0"`
	Error       string     `gorm:"type:text"`
	StartedAt   *time.Time `gorm:"column:started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (taskRunPO) TableName() string { return "crawl_task_runs" }

func taskRunFromDomain(t *crawl.TaskRun) *taskRunPO {
	return &taskRunPO{
		ID:          t.ID,
		SourceID:    t.SourceID,
		TaskType:    t.TaskType,
		TriggeredBy: t.TriggeredBy,
		Status:      t.Status,
		Total:       t.Total,
		Ingested:    t.Ingested,
		Failed:      t.Failed,
		Error:       t.Error,
		StartedAt:   t.StartedAt,
		FinishedAt:  t.FinishedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func (po *taskRunPO) toDomain() *crawl.TaskRun {
	return &crawl.TaskRun{
		ID:          po.ID,
		SourceID:    po.SourceID,
		TaskType:    po.TaskType,
		TriggeredBy: po.TriggeredBy,
		Status:      po.Status,
		Total:       po.Total,
		Ingested:    po.Ingested,
		Failed:      po.Failed,
		Error:       po.Error,
		StartedAt:   po.StartedAt,
		FinishedAt:  po.FinishedAt,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

type taskRunRepository struct {
	db *gorm.DB
}

// NewCrawlTaskRunRepository 创建任务运行仓储
func NewCrawlTaskRunRepository(db *gorm.DB) crawl.TaskRunRepository {
	return &taskRunRepository{db: db}
}

func (r *taskRunRepository) Create(ctx context.Context, t *crawl.TaskRun) error {
	return r.db.WithContext(ctx).Create(taskRunFromDomain(t)).Error
}

func (r *taskRunRepository) UpdateStatus(ctx context.Context, id string, status string, total, ingested, failed int, errMsg string) error {
	return r.db.WithContext(ctx).Model(&taskRunPO{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"total":      total,
		"ingested":   ingested,
		"failed":     failed,
		"error":      errMsg,
		"updated_at": time.Now(),
	}).Error
}

func (r *taskRunRepository) UpdateTimestamps(ctx context.Context, id string, startedAt, finishedAt *time.Time) error {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if startedAt != nil {
		updates["started_at"] = startedAt
	}
	if finishedAt != nil {
		updates["finished_at"] = finishedAt
	}
	return r.db.WithContext(ctx).Model(&taskRunPO{}).Where("id = ?", id).Updates(updates).Error
}

func (r *taskRunRepository) FindByID(ctx context.Context, id string) (*crawl.TaskRun, error) {
	var po taskRunPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, crawl.ErrTaskNotFound
		}
		return nil, err
	}
	return po.toDomain(), nil
}

func (r *taskRunRepository) List(ctx context.Context, sourceID string) ([]*crawl.TaskRun, error) {
	var pos []taskRunPO
	q := r.db.WithContext(ctx).Order("created_at DESC")
	if sourceID != "" {
		q = q.Where("source_id = ?", sourceID)
	}
	if err := q.Find(&pos).Error; err != nil {
		return nil, err
	}
	out := make([]*crawl.TaskRun, 0, len(pos))
	for i := range pos {
		out = append(out, pos[i].toDomain())
	}
	return out, nil
}

// HasActiveTask 检查是否有 queued/running 状态的任务（任务去重）
func (r *taskRunRepository) HasActiveTask(ctx context.Context, sourceID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&taskRunPO{}).
		Where("source_id = ? AND status IN ?", sourceID, []string{"queued", "running"}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *taskRunRepository) CountAll(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&taskRunPO{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *taskRunRepository) CountByStatus(ctx context.Context, statuses ...string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&taskRunPO{}).Where("status IN ?", statuses).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
