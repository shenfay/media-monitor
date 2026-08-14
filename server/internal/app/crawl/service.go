// Package crawl 抓取应用服务：承载 Go 侧与 Python 抓取服务之间的全部业务逻辑。
//
// 职责边界：
//   - 管理员接口（CRUD / 手动触发 / 任务查询）走用户 JWT + Casbin。
//   - 任务分发通过 Redis Stream crawl:task:dispatch 下发给 Python Worker。
//   - 文章数据通过 crawl:article:ingest Stream 从 Python 回传，Go 消费后写入 PG。
//   - 任务事件通过 crawl:task:event Stream 从 Python 回传，更新 TaskRun 状态。
package crawl

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
	"github.com/shenfay/go-react-admin/internal/infra/config"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// ScrapeParams Redis Stream 任务参数（Python 消费时读取 payload 字段）
type ScrapeParams struct {
	Limit    int    `json:"limit"`
	WithBody bool   `json:"with_body"`
	Since    string `json:"since,omitempty"`
	Mode     string `json:"mode,omitempty"`
}

// 领域错误
var (
	ErrSourceNotFound  = errors.New("crawl: source not found")
	ErrSourceDisabled  = errors.New("crawl: source disabled")
	ErrTaskNotFound    = errors.New("crawl: task not found")
	ErrTaskAlreadyRun  = errors.New("crawl: task already running for this source")
)

// Service 抓取应用服务
type Service struct {
	sources  crawl.SourceRepository
	articles crawl.ArticleRepository
	tasks    crawl.TaskRunRepository
	redis    *redis.Client
	cfg      config.ScraperConfig
}

// NewService 创建抓取服务
func NewService(
	sources crawl.SourceRepository,
	articles crawl.ArticleRepository,
	tasks crawl.TaskRunRepository,
	rdb *redis.Client,
	cfg config.ScraperConfig,
) *Service {
	return &Service{sources: sources, articles: articles, tasks: tasks, redis: rdb, cfg: cfg}
}

// CreateSource 创建数据源（auth 由仓储层加密）
func (s *Service) CreateSource(ctx context.Context, ownerID string, req CreateSourceRequest) (*crawl.Source, error) {
	now := time.Now()
	src := &crawl.Source{
		ID:           utils.GenerateID(),
		Name:         req.Name,
		PlatformType: req.PlatformType,
		BaseURL:      req.BaseURL,
		ListEndpoint: req.ListEndpoint,
		Nodes:        req.Nodes,
		SourceFilter: req.SourceFilter,
		Months:       req.Months,
		Schedule:     req.Schedule,
		Enabled:      req.Enabled,
		OwnerID:      strPtr(ownerID),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if len(req.Auth) > 0 {
		src.Auth = req.Auth
	}
	if err := s.sources.Create(ctx, src); err != nil {
		return nil, err
	}
	return src, nil
}

// ListSources 列出数据源；enabledOnly=true 时仅返回启用项（调度器使用）
func (s *Service) ListSources(ctx context.Context, enabledOnly bool) ([]*crawl.Source, error) {
	return s.sources.List(ctx, enabledOnly)
}

// GetSource 获取数据源（返回明文 auth，供 Python 拉取凭证）
func (s *Service) GetSource(ctx context.Context, id string) (*crawl.Source, error) {
	src, err := s.sources.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSourceNotFound
		}
		return nil, err
	}
	return src, nil
}

// UpdateSource 更新数据源
func (s *Service) UpdateSource(ctx context.Context, id string, req UpdateSourceRequest) (*crawl.Source, error) {
	src, err := s.sources.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSourceNotFound
		}
		return nil, err
	}
	if req.Name != "" {
		src.Name = req.Name
	}
	if req.PlatformType != "" {
		src.PlatformType = req.PlatformType
	}
	if req.BaseURL != "" {
		src.BaseURL = req.BaseURL
	}
	if req.ListEndpoint != "" {
		src.ListEndpoint = req.ListEndpoint
	}
	if req.Nodes != nil {
		src.Nodes = req.Nodes
	}
	if req.SourceFilter != "" {
		src.SourceFilter = req.SourceFilter
	}
	if req.Months != 0 {
		src.Months = req.Months
	}
	if req.Schedule != "" {
		src.Schedule = req.Schedule
	}
	// req.Auth 为 nil 表示调用方未提供，保留现有（已解密）auth；否则覆盖
	if req.Auth != nil {
		src.Auth = req.Auth
	}
	if req.Enabled != nil {
		src.Enabled = *req.Enabled
	}
	src.UpdatedAt = time.Now()
	if err := s.sources.Update(ctx, src); err != nil {
		return nil, err
	}
	return src, nil
}

// DeleteSource 删除数据源
func (s *Service) DeleteSource(ctx context.Context, id string) error {
	return s.sources.Delete(ctx, id)
}

// enqueue 写入 Redis Stream 并创建 TaskRun（状态 queued）
// 消息内嵌 Source 配置（含明文 auth），Python Worker 无需回调 Go 获取配置。
func (s *Service) enqueue(ctx context.Context, src *crawl.Source, trigger string, params ScrapeParams) (string, error) {
	// 任务去重：检查是否有 queued/running 的任务
	hasActive, err := s.tasks.HasActiveTask(ctx, src.ID)
	if err != nil {
		return "", err
	}
	if hasActive {
		return "", ErrTaskAlreadyRun
	}

	taskID := utils.GenerateID()
	now := time.Now()
	tr := &crawl.TaskRun{
		ID:          taskID,
		SourceID:    src.ID,
		TaskType:    src.PlatformType,
		TriggeredBy: trigger,
		Status:      "queued",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.tasks.Create(ctx, tr); err != nil {
		return "", err
	}

	// 构造消息：内嵌 Source 配置 + 参数
	dispatch := TaskDispatchMessage{
		TaskID:   taskID,
		SourceID: src.ID,
		Source:   src,
		Params:   params,
	}
	payload, err := json.Marshal(dispatch)
	if err != nil {
		return "", err
	}

	if err := s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: s.cfg.Stream,
		Values: map[string]interface{}{
			"task_id":   taskID,
			"source_id": src.ID,
			"payload":   string(payload),
		},
	}).Err(); err != nil {
		return "", err
	}
	return taskID, nil
}

// EnqueueForSource 按数据源配置入队一次抓取（cron 或手动触发）
func (s *Service) EnqueueForSource(ctx context.Context, sourceID, trigger string) (string, error) {
	src, err := s.sources.FindByID(ctx, sourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrSourceNotFound
		}
		return "", err
	}
	if !src.Enabled {
		return "", ErrSourceDisabled
	}
	params := ScrapeParams{Limit: 200, WithBody: false, Mode: "list"}
	if src.Months > 0 {
		params.Since = time.Now().AddDate(0, -src.Months, 0).Format(time.RFC3339)
	}
	return s.enqueue(ctx, src, trigger, params)
}

// CreateTask 管理员手动创建任务（可覆盖默认参数）
func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) (string, error) {
	src, err := s.sources.FindByID(ctx, req.SourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrSourceNotFound
		}
		return "", err
	}
	params := ScrapeParams{Limit: req.Limit, WithBody: req.WithBody, Mode: req.Mode}
	if params.Limit == 0 {
		params.Limit = 200
	}
	if req.Since != "" {
		params.Since = req.Since
	} else if src.Months > 0 {
		params.Since = time.Now().AddDate(0, -src.Months, 0).Format(time.RFC3339)
	}
	return s.enqueue(ctx, src, "manual", params)
}

// IngestArticleBatch 从 Stream 消费的文章批次写入 PG
func (s *Service) IngestArticleBatch(ctx context.Context, articles []*crawl.Article) (int, error) {
	return s.articles.UpsertBatch(ctx, articles)
}

// UpdateTaskRunStatus 从 Stream 消费的事件更新 TaskRun 状态
func (s *Service) UpdateTaskRunStatus(ctx context.Context, taskID, status string, total, ingested, failed int, errMsg string) error {
	return s.tasks.UpdateStatus(ctx, taskID, status, total, ingested, failed, errMsg)
}

// UpdateLastCrawlAt 更新数据源最后抓取时间
func (s *Service) UpdateLastCrawlAt(ctx context.Context, sourceID string, t time.Time) error {
	return s.sources.UpdateLastCrawlAt(ctx, sourceID, t)
}

// GetTask 获取任务详情
func (s *Service) GetTask(ctx context.Context, id string) (*crawl.TaskRun, error) {
	t, err := s.tasks.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return t, nil
}

// ListTasks 列出任务（可按 source_id 过滤）
func (s *Service) ListTasks(ctx context.Context, sourceID string) ([]*crawl.TaskRun, error) {
	return s.tasks.List(ctx, sourceID)
}

// mapArticleInput 将 Python 回写字段映射为领域 Article
func mapArticleInput(in ArticleInput) (*crawl.Article, error) {
	now := time.Now()
	a := &crawl.Article{
		ID:         utils.GenerateID(),
		SourceID:   in.SourceID,
		Platform:   in.Platform,
		URL:        in.URL,
		Title:      in.Title,
		Summary:    in.Summary,
		Body:       in.Content,
		BodyFormat: "html",
		Author:     in.Author,
		SourceName: in.SourceName,
		Language:   "",
		CreatedAt:  now,
		FetchedAt:  now,
	}
	// ExternalID 优先取传入值，否则取 URL 最后一段路径
	if in.ExternalID != "" {
		a.ExternalID = in.ExternalID
	} else if in.URL != "" {
		if idx := strings.LastIndex(in.URL, "/"); idx >= 0 && idx < len(in.URL)-1 {
			a.ExternalID = in.URL[idx+1:]
		} else {
			a.ExternalID = in.URL
		}
	}
	// URLHash
	a.URLHash = computeURLHashForArticle(a.URL)
	// 发布时间：ISO 8601 -> time.Time
	if in.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, in.PublishedAt); err == nil {
			a.PublishedAt = &t
		}
	}
	// 媒体：Python images -> media
	if len(in.Images) > 0 {
		if b, err := json.Marshal(in.Images); err == nil {
			a.Media = b
		}
	}
	// 互动指标：Python extra -> interactions
	if len(in.Extra) > 0 {
		if b, err := json.Marshal(in.Extra); err == nil {
			a.Interactions = b
		}
	}
	// 抓取时间
	if in.FetchedAt != "" {
		if t, err := time.Parse(time.RFC3339, in.FetchedAt); err == nil {
			a.FetchedAt = t
		}
	}
	// 原始负载留底
	if raw, err := json.Marshal(in); err == nil {
		a.RawPayload = raw
	}
	return a, nil
}

// computeURLHashForArticle 计算 URL 的 MD5 哈希值
func computeURLHashForArticle(url string) string {
	if url == "" {
		return ""
	}
	h := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", h)
}

func strPtr(s string) *string { return &s }
