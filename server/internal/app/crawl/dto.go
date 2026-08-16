package crawl

import (
	"encoding/json"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
)

// CreateSourceRequest 创建数据源请求
type CreateSourceRequest struct {
	Name         string          `json:"name" binding:"required"`
	PlatformType string          `json:"platform_type"`
	BaseURL      string          `json:"base_url"`
	ListEndpoint string          `json:"list_endpoint"`
	Nodes        []string        `json:"nodes"`
	SourceFilter string          `json:"source_filter"`
	Months       int             `json:"months"`
	Schedule     string          `json:"schedule"`
	Auth         json.RawMessage `json:"auth"`
	Tags         []string        `json:"tags"`
	Enabled      bool            `json:"enabled"`
}

// UpdateSourceRequest 更新数据源请求
type UpdateSourceRequest struct {
	Name         string          `json:"name"`
	PlatformType string          `json:"platform_type"`
	BaseURL      string          `json:"base_url"`
	ListEndpoint string          `json:"list_endpoint"`
	Nodes        []string        `json:"nodes"`
	SourceFilter string          `json:"source_filter"`
	Months       int             `json:"months"`
	Schedule     string          `json:"schedule"`
	Auth         json.RawMessage `json:"auth"`
	Tags         []string        `json:"tags"`
	Enabled      *bool           `json:"enabled"`
}

// CreateTaskRequest 创建抓取任务请求（管理员手动）
type CreateTaskRequest struct {
	SourceID string `json:"source_id" binding:"required"`
	Limit    int    `json:"limit"`
	WithBody bool   `json:"with_body"`
	Since    string `json:"since"`
	Mode     string `json:"mode"`
}

// ArticleInput 单篇文章（来自 Python Stream，字段名与 scrapers/contracts/article.py 对齐）
type ArticleInput struct {
	SourceID    string                 `json:"source_id" binding:"required"`
	Platform    string                 `json:"platform"`
	Title       string                 `json:"title"`
	URL         string                 `json:"url" binding:"required"`
	ExternalID  string                 `json:"external_id"` // 平台侧文章 ID
	Author      string                 `json:"author"`
	SourceName  string                 `json:"source_name"`
	Summary     string                 `json:"summary"`
	Content     string                 `json:"content"`
	PublishedAt string                 `json:"published_at"` // ISO 8601
	Images      []interface{}          `json:"images"`
	Extra       map[string]interface{} `json:"extra"`      // 互动指标等
	FetchedAt   string                 `json:"fetched_at"` // ISO 8601
}

// TaskDispatchMessage 写入 crawl:task:dispatch 的消息体（内嵌 Source 配置）
type TaskDispatchMessage struct {
	TaskID   string        `json:"task_id"`
	SourceID string        `json:"source_id"`
	Source   *crawl.Source `json:"source"`
	Params   ScrapeParams  `json:"params"`
}

// ArticleIngestMessage 从 crawl:article:ingest 消费到的消息体
type ArticleIngestMessage struct {
	TaskID   string         `json:"task_id"`
	SourceID string         `json:"source_id"`
	Phase    string         `json:"phase"` // list | detail
	BatchSeq int            `json:"batch_seq"`
	Articles []ArticleInput `json:"articles"`
}

// TaskEventMessage 从 crawl:task:event 消费到的消息体
type TaskEventMessage struct {
	TaskID       string `json:"task_id"`
	Type         string `json:"type"` // status | phase_start | phase_done | list_synced | detail_progress | task_done | task_failed
	Phase        string `json:"phase"`
	Status       string `json:"status"`
	Total        int    `json:"total"`
	ListCount    int    `json:"list_count"`
	DetailCount  int    `json:"detail_count"`
	DetailFailed int    `json:"detail_failed"`
	Error        string `json:"error"`
	Timestamp    string `json:"timestamp"`
}
