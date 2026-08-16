// Package crawl 抓取领域模型。
//
// 该模块是 go-react-admin 与 Python 抓取服务之间的「集成缝」：
//   - Source：数据源配置（含平台/节点/来源过滤/cron/加密 auth）
//   - Article：抓取到的文章（幂等键 source_id + url_hash）
//   - TaskRun：一次抓取任务的运行记录（状态由 Stream 事件驱动）
//
// 约定：Source.Auth 在内存中为明文 JSON 对象；落库时由仓储层加密（AES-GCM），
// 读取时解密。管理员列表接口对 auth 脱敏。
package crawl

import (
	"encoding/json"
	"time"
)

// Source 数据源配置
type Source struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	PlatformType string          `json:"platform_type"` // news | social | social_overseas
	BaseURL      string          `json:"base_url"`
	ListEndpoint string          `json:"list_endpoint"`
	Nodes        []string        `json:"nodes"`
	SourceFilter string          `json:"source_filter"`
	Months       int             `json:"months"`
	Schedule     string          `json:"schedule"` // cron 表达式；空表示不自动调度
	Auth         json.RawMessage `json:"auth"`     // 明文 JSON 对象（社媒登录态/cookie/token）；落库为密文
	Tags         []string        `json:"tags"`     // 能力标签（如 overseas, youtube），用于 Stream 路由
	Enabled      bool            `json:"enabled"`
	OwnerID      *string         `json:"owner_id,omitempty"`
	ArticleCount int             `json:"article_count"` // 关联文章数（列表查询时填充）
	LastCrawlAt  *time.Time      `json:"last_crawl_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// Article 抓取到的文章
type Article struct {
	ID           string          `json:"id"`
	SourceID     string          `json:"source_id"`
	Platform     string          `json:"platform"`
	ExternalID   string          `json:"external_id"` // 平台侧文章 ID
	URL          string          `json:"url"`
	URLHash      string          `json:"url_hash"` // MD5(url)，用于高效索引
	Title        string          `json:"title"`
	Subtitle     *string         `json:"subtitle,omitempty"`
	Summary      string          `json:"summary"`
	Body         string          `json:"body"`
	BodyFormat   string          `json:"body_format"` // html | markdown
	Channel      string          `json:"channel"`
	Author       string          `json:"author"`
	SourceName   string          `json:"source_name"`
	PublishedAt  *time.Time      `json:"published_at"`
	Language     string          `json:"language"`
	Status       string          `json:"status"`       // pending | completed | failed
	Interactions json.RawMessage `json:"interactions"` // 互动指标（来自 Python extra）
	Media        json.RawMessage `json:"media"`        // 图片/视频等媒体（来自 Python images）
	ThreadID     *string         `json:"thread_id,omitempty"`
	RawPayload   json.RawMessage `json:"raw_payload"` // 原始回写负载，便于排查
	FetchedAt    time.Time       `json:"fetched_at"`
	CreatedAt    time.Time       `json:"created_at"`
}

// TaskRun 抓取任务运行记录
type TaskRun struct {
	ID          string     `json:"id"`
	SourceID    string     `json:"source_id"`
	TaskType    string     `json:"task_type"`    // news | social
	TriggeredBy string     `json:"triggered_by"` // cron | manual
	Status      string     `json:"status"`       // queued | running | success | partial | failed
	Total       int        `json:"total"`
	Ingested    int        `json:"ingested"`
	Failed      int        `json:"failed"`
	Error       string     `json:"error"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
