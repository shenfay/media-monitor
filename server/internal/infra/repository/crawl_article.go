package repository

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
)

// articlePO 文章持久化对象
type articlePO struct {
	ID           string `gorm:"primaryKey;type:varchar(26)"`
	SourceID     string `gorm:"column:source_id;type:varchar(26);index"`
	Platform     string `gorm:"size:50"`
	ExternalID   string `gorm:"column:external_id;size:100;index"`
	URL          string `gorm:"column:url;type:text"`
	URLHash      string `gorm:"column:url_hash;size:64"`
	Title        string `gorm:"size:512"`
	Subtitle     *string
	Summary      string     `gorm:"type:text"`
	Body         string     `gorm:"type:text"`
	BodyFormat   string     `gorm:"column:body_format;size:20"`
	Channel      string     `gorm:"size:100"`
	Author       string     `gorm:"size:200"`
	SourceName   string     `gorm:"column:source_name;size:200"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	Language     string     `gorm:"size:20"`
	Status       string     `gorm:"size:20"`
	Interactions string     `gorm:"type:text"`
	Media        string     `gorm:"type:text"`
	ThreadID     *string    `gorm:"column:thread_id;type:varchar(100)"`
	RawPayload   string     `gorm:"column:raw_payload;type:text"`
	FetchedAt    time.Time  `gorm:"column:fetched_at"`
	CreatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (articlePO) TableName() string { return "crawl_articles" }

func articleFromDomain(a *crawl.Article) *articlePO {
	return &articlePO{
		ID:           a.ID,
		SourceID:     a.SourceID,
		Platform:     a.Platform,
		ExternalID:   a.ExternalID,
		URL:          a.URL,
		URLHash:      computeURLHash(a.URL),
		Title:        a.Title,
		Subtitle:     a.Subtitle,
		Summary:      a.Summary,
		Body:         a.Body,
		BodyFormat:   a.BodyFormat,
		Channel:      a.Channel,
		Author:       a.Author,
		SourceName:   a.SourceName,
		PublishedAt:  a.PublishedAt,
		Language:     a.Language,
		Status:       a.Status,
		Interactions: string(a.Interactions),
		Media:        string(a.Media),
		ThreadID:     a.ThreadID,
		RawPayload:   string(a.RawPayload),
		FetchedAt:    a.FetchedAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    time.Now(),
	}
}

func articleToDomain(po *articlePO) *crawl.Article {
	return &crawl.Article{
		ID:           po.ID,
		SourceID:     po.SourceID,
		Platform:     po.Platform,
		ExternalID:   po.ExternalID,
		URL:          po.URL,
		URLHash:      po.URLHash,
		Title:        po.Title,
		Subtitle:     po.Subtitle,
		Summary:      po.Summary,
		Body:         po.Body,
		BodyFormat:   po.BodyFormat,
		Channel:      po.Channel,
		Author:       po.Author,
		SourceName:   po.SourceName,
		PublishedAt:  po.PublishedAt,
		Language:     po.Language,
		Status:       po.Status,
		Interactions: safeJSON(po.Interactions),
		Media:        safeJSON(po.Media),
		ThreadID:     po.ThreadID,
		RawPayload:   safeJSON(po.RawPayload),
		FetchedAt:    po.FetchedAt,
		CreatedAt:    po.CreatedAt,
	}
}

// safeJSON 将空字符串转为 nil（避免 json.RawMessage 序列化失败）
func safeJSON(s string) json.RawMessage {
	if s == "" {
		return nil
	}
	return json.RawMessage(s)
}

// computeURLHash 计算 URL 的 MD5 哈希值
func computeURLHash(url string) string {
	if url == "" {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))
}

type articleRepository struct {
	db *gorm.DB
}

// NewCrawlArticleRepository 创建文章仓储
func NewCrawlArticleRepository(db *gorm.DB) crawl.ArticleRepository {
	return &articleRepository{db: db}
}

// UpsertBatch 批量 upsert（列表阶段），幂等键为 (source_id, url_hash)；
// 冲突时更新元数据字段，但不覆盖 body（正文由 detail 阶段单独写入）
func (r *articleRepository) UpsertBatch(ctx context.Context, articles []*crawl.Article) (int, error) {
	if len(articles) == 0 {
		return 0, nil
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, a := range articles {
			po := articleFromDomain(a)
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "source_id"}, {Name: "url_hash"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"title", "subtitle", "summary", "body_format",
					"channel", "author", "source_name", "published_at", "language",
					"status", "interactions", "media", "thread_id", "raw_payload", "fetched_at", "updated_at",
				}),
			}).Create(&po).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return len(articles), nil
}

// UpsertDetailBatch 详情阶段 upsert：更新全部字段（详情页数据更权威）
func (r *articleRepository) UpsertDetailBatch(ctx context.Context, articles []*crawl.Article) (int, error) {
	if len(articles) == 0 {
		return 0, nil
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, a := range articles {
			po := articleFromDomain(a)
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "source_id"}, {Name: "url_hash"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"title", "subtitle", "body", "body_format",
					"author", "source_name", "published_at",
					"status", "media", "raw_payload", "fetched_at", "updated_at",
				}),
			}).Create(&po).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return len(articles), nil
}

// List 按过滤条件分页查询文章
func (r *articleRepository) List(ctx context.Context, filter crawl.ArticleFilter) ([]*crawl.Article, int, error) {
	query := r.db.WithContext(ctx).Model(&articlePO{})
	if filter.SourceID != "" {
		query = query.Where("source_id = ?", filter.SourceID)
	}
	if filter.Platform != "" {
		query = query.Where("platform = ?", filter.Platform)
	}
	if filter.Language != "" {
		query = query.Where("language = ?", filter.Language)
	}
	if filter.Keyword != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Keyword+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var pos []*articlePO
	if err := query.Order("created_at DESC").Offset(filter.Offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	articles := make([]*crawl.Article, 0, len(pos))
	for _, po := range pos {
		articles = append(articles, articleToDomain(po))
	}
	return articles, int(total), nil
}

// FindByID 按 ID 查询单篇文章
func (r *articleRepository) FindByID(ctx context.Context, id string) (*crawl.Article, error) {
	var po articlePO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, crawl.ErrArticleNotFound
		}
		return nil, err
	}
	return articleToDomain(&po), nil
}

func (r *articleRepository) CountAll(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&articlePO{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *articleRepository) CountSince(ctx context.Context, since time.Time) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&articlePO{}).Where("created_at >= ?", since).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
