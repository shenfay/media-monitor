package crawl

import (
	"context"
	"time"
)

// SourceRepository 数据源仓储接口
type SourceRepository interface {
	Create(ctx context.Context, s *Source) error
	Update(ctx context.Context, s *Source) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Source, error)
	List(ctx context.Context, enabledOnly bool) ([]*Source, error)
	UpdateLastCrawlAt(ctx context.Context, id string, t time.Time) error
	// CountAll 返回总数
	CountAll(ctx context.Context) (int, error)
	// CountEnabled 返回启用数
	CountEnabled(ctx context.Context) (int, error)
}

// ArticleFilter 文章查询过滤条件
type ArticleFilter struct {
	SourceID string
	Platform string
	Language string
	Keyword  string // 标题模糊搜索
	Status   string // pending | completed | failed
	Limit    int
	Offset   int
}

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	// UpsertBatch 批量 upsert，幂等键为 (source_id, url_hash)
	UpsertBatch(ctx context.Context, articles []*Article) (int, error)
	// UpsertDetailBatch 仅更新正文相关字段（body, status, fetched_at, raw_payload）
	UpsertDetailBatch(ctx context.Context, articles []*Article) (int, error)
	// List 按过滤条件分页查询
	List(ctx context.Context, filter ArticleFilter) ([]*Article, int, error)
	// FindByID 按 ID 查询单篇文章
	FindByID(ctx context.Context, id string) (*Article, error)
	// CountAll 返回总数
	CountAll(ctx context.Context) (int, error)
	// CountSince 返回指定时间之后的数量
	CountSince(ctx context.Context, since time.Time) (int, error)
}

// TaskRunRepository 任务运行仓储接口
type TaskRunRepository interface {
	Create(ctx context.Context, t *TaskRun) error
	UpdateStatus(ctx context.Context, id string, status string, total, ingested, failed int, errMsg string) error
	UpdateTimestamps(ctx context.Context, id string, startedAt, finishedAt *time.Time) error
	FindByID(ctx context.Context, id string) (*TaskRun, error)
	List(ctx context.Context, sourceID string) ([]*TaskRun, error)
	HasActiveTask(ctx context.Context, sourceID string) (bool, error)
	// CountAll 返回总数
	CountAll(ctx context.Context) (int, error)
	// CountByStatus 按状态列表统计数量
	CountByStatus(ctx context.Context, statuses ...string) (int, error)
}
