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
}

// ArticleFilter 文章查询过滤条件
type ArticleFilter struct {
	SourceID string
	Platform string
	Language string
	Keyword  string // 标题模糊搜索
	Limit    int
	Offset   int
}

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	// UpsertBatch 批量 upsert，幂等键为 (source_id, url_hash)
	UpsertBatch(ctx context.Context, articles []*Article) (int, error)
	// List 按过滤条件分页查询
	List(ctx context.Context, filter ArticleFilter) ([]*Article, int, error)
	// FindByID 按 ID 查询单篇文章
	FindByID(ctx context.Context, id string) (*Article, error)
}

// TaskRunRepository 任务运行仓储接口
type TaskRunRepository interface {
	Create(ctx context.Context, t *TaskRun) error
	UpdateStatus(ctx context.Context, id string, status string, total, ingested, failed int, errMsg string) error
	FindByID(ctx context.Context, id string) (*TaskRun, error)
	List(ctx context.Context, sourceID string) ([]*TaskRun, error)
	HasActiveTask(ctx context.Context, sourceID string) (bool, error)
}
