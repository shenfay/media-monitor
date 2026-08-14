package repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
)

// articlePO 文章持久化对象
type articlePO struct {
	ID           string     `gorm:"primaryKey;type:varchar(26)"`
	SourceID     string     `gorm:"column:source_id;type:varchar(26);index"`
	Platform     string     `gorm:"size:50"`
	ExternalID   string     `gorm:"column:external_id;size:100;index"`
	URL          string     `gorm:"column:url;type:text"`
	URLHash      string     `gorm:"column:url_hash;size:64"`
	Title        string     `gorm:"size:512"`
	Subtitle     *string
	Summary      string `gorm:"type:text"`
	Body         string `gorm:"type:text"`
	BodyFormat   string `gorm:"column:body_format;size:20"`
	Channel      string `gorm:"size:100"`
	Author       string `gorm:"size:200"`
	SourceName   string `gorm:"column:source_name;size:200"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	Language     string `gorm:"size:20"`
	Interactions string `gorm:"type:text"`
	Media        string `gorm:"type:text"`
	ThreadID     *string `gorm:"column:thread_id;type:varchar(100)"`
	RawPayload   string `gorm:"column:raw_payload;type:text"`
	FetchedAt    time.Time `gorm:"column:fetched_at"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
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
		Interactions: string(a.Interactions),
		Media:        string(a.Media),
		ThreadID:     a.ThreadID,
		RawPayload:   string(a.RawPayload),
		FetchedAt:    a.FetchedAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    time.Now(),
	}
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

// UpsertBatch 批量 upsert，幂等键为 (source_id, url_hash)；冲突时更新正文/状态等字段
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
					"title", "subtitle", "summary", "body", "body_format",
					"channel", "author", "source_name", "published_at", "language",
					"interactions", "media", "thread_id", "raw_payload", "fetched_at", "updated_at",
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
