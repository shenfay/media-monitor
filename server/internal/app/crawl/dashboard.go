package crawl

import (
	"context"
	"time"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
)

// DashboardStats Dashboard 统计数据
type DashboardStats struct {
	TotalSources   int `json:"total_sources"`
	EnabledSources int `json:"enabled_sources"`
	TotalArticles  int `json:"total_articles"`
	TodayArticles  int `json:"today_articles"`
	TotalTasks     int `json:"total_tasks"`
	RunningTasks   int `json:"running_tasks"`
	OnlineWorkers  int `json:"online_workers"`
}

// GetDashboardStats 获取 Dashboard 统计数据
func (s *Service) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 数据源统计
	totalSources, _ := s.sources.CountAll(ctx)
	enabledSources, _ := s.sources.CountEnabled(ctx)
	stats.TotalSources = totalSources
	stats.EnabledSources = enabledSources

	// 文章统计
	totalArticles, _ := s.articles.CountAll(ctx)
	stats.TotalArticles = totalArticles

	// 今日文章
	today := time.Now().Truncate(24 * time.Hour)
	todayArticles, _ := s.articles.CountSince(ctx, today)
	stats.TodayArticles = todayArticles

	// 任务统计
	totalTasks, _ := s.tasks.CountAll(ctx)
	stats.TotalTasks = totalTasks

	runningTasks, _ := s.tasks.CountByStatus(ctx, "queued", "running")
	stats.RunningTasks = runningTasks

	// 在线 Worker（从 Redis 读取）
	ids, err := s.redis.SMembers(ctx, workerIDsKey).Result()
	if err == nil {
		for _, id := range ids {
			ttl := s.redis.TTL(ctx, workerKeyPrefix+id).Val()
			if ttl > 0 {
				stats.OnlineWorkers++
			}
		}
	}

	return stats, nil
}

// ListArticles 列出文章（分页）
func (s *Service) ListArticles(ctx context.Context, filter crawl.ArticleFilter) ([]*crawl.Article, int, error) {
	return s.articles.List(ctx, filter)
}

// GetArticle 获取单篇文章
func (s *Service) GetArticle(ctx context.Context, id string) (*crawl.Article, error) {
	return s.articles.FindByID(ctx, id)
}
