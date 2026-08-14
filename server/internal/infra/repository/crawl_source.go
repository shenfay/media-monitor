package repository

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// sourcePO 数据源持久化对象（auth 落库为加密字符串）
type sourcePO struct {
	ID           string     `gorm:"primaryKey;type:varchar(26)"`
	Name         string     `gorm:"size:200;not null"`
	PlatformType string     `gorm:"column:platform_type;size:30;not null;default:'news'"`
	BaseURL      string     `gorm:"column:base_url;size:512"`
	ListEndpoint string     `gorm:"column:list_endpoint;size:512"`
	Nodes        string     `gorm:"type:text"` // JSON 数组字符串
	SourceFilter string     `gorm:"column:source_filter;size:200"`
	Months       int        `gorm:"default:6"`
	Schedule     string     `gorm:"size:100"`
	Auth         string     `gorm:"type:text"` // 加密后的 JSON 字符串
	Tags         string     `gorm:"type:text;default:'[]'"` // JSON 数组字符串
	Enabled      bool       `gorm:"default:true"`
	OwnerID      *string    `gorm:"column:owner_id;type:varchar(50)"`
	LastCrawlAt  *time.Time `gorm:"column:last_crawl_at"`
	CreatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (sourcePO) TableName() string { return "crawl_sources" }

// toDomain 转换为领域模型（解密 auth）
func (po *sourcePO) toDomain() *crawl.Source {
	s := &crawl.Source{
		ID:           po.ID,
		Name:         po.Name,
		PlatformType: po.PlatformType,
		BaseURL:      po.BaseURL,
		ListEndpoint: po.ListEndpoint,
		SourceFilter: po.SourceFilter,
		Months:       po.Months,
		Schedule:     po.Schedule,
		Enabled:      po.Enabled,
		OwnerID:      po.OwnerID,
		LastCrawlAt:  po.LastCrawlAt,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
	if po.Nodes != "" {
		_ = json.Unmarshal([]byte(po.Nodes), &s.Nodes)
	}
	if po.Tags != "" {
		_ = json.Unmarshal([]byte(po.Tags), &s.Tags)
	}
	if po.Auth != "" {
		if dec, err := utils.Decrypt(po.Auth); err == nil {
			s.Auth = json.RawMessage(dec)
		}
	}
	return s
}

// sourceFromDomain 从领域模型构造 PO（加密 auth）
func sourceFromDomain(s *crawl.Source) *sourcePO {
	po := &sourcePO{
		ID:           s.ID,
		Name:         s.Name,
		PlatformType: s.PlatformType,
		BaseURL:      s.BaseURL,
		ListEndpoint: s.ListEndpoint,
		SourceFilter: s.SourceFilter,
		Months:       s.Months,
		Schedule:     s.Schedule,
		Enabled:      s.Enabled,
		OwnerID:      s.OwnerID,
		LastCrawlAt:  s.LastCrawlAt,
	}
	if len(s.Nodes) > 0 {
		if b, err := json.Marshal(s.Nodes); err == nil {
			po.Nodes = string(b)
		}
	}
	if len(s.Tags) > 0 {
		if b, err := json.Marshal(s.Tags); err == nil {
			po.Tags = string(b)
		}
	}
	if len(s.Auth) > 0 {
		if enc, err := utils.Encrypt(string(s.Auth)); err == nil {
			po.Auth = enc
		}
	}
	return po
}

type sourceRepository struct {
	db *gorm.DB
}

// NewCrawlSourceRepository 创建数据源仓储
func NewCrawlSourceRepository(db *gorm.DB) crawl.SourceRepository {
	return &sourceRepository{db: db}
}

func (r *sourceRepository) Create(ctx context.Context, s *crawl.Source) error {
	return r.db.WithContext(ctx).Create(sourceFromDomain(s)).Error
}

func (r *sourceRepository) Update(ctx context.Context, s *crawl.Source) error {
	po := sourceFromDomain(s)
	return r.db.WithContext(ctx).Model(&sourcePO{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
		"name":          po.Name,
		"platform_type": po.PlatformType,
		"base_url":      po.BaseURL,
		"list_endpoint": po.ListEndpoint,
		"nodes":         po.Nodes,
		"source_filter": po.SourceFilter,
		"months":        po.Months,
		"schedule":      po.Schedule,
		"auth":          po.Auth,
		"tags":          po.Tags,
		"enabled":       po.Enabled,
		"updated_at":    time.Now(),
	}).Error
}

func (r *sourceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&sourcePO{}).Error
}

func (r *sourceRepository) FindByID(ctx context.Context, id string) (*crawl.Source, error) {
	var po sourcePO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		return nil, err
	}
	return po.toDomain(), nil
}

func (r *sourceRepository) List(ctx context.Context, enabledOnly bool) ([]*crawl.Source, error) {
	var pos []sourcePO
	q := r.db.WithContext(ctx).Order("created_at DESC")
	if enabledOnly {
		q = q.Where("enabled = ?", true)
	}
	if err := q.Find(&pos).Error; err != nil {
		return nil, err
	}
	out := make([]*crawl.Source, 0, len(pos))
	for i := range pos {
		out = append(out, pos[i].toDomain())
	}
	return out, nil
}

// UpdateLastCrawlAt 更新数据源最后抓取时间
func (r *sourceRepository) UpdateLastCrawlAt(ctx context.Context, id string, t time.Time) error {
	return r.db.WithContext(ctx).Model(&sourcePO{}).Where("id = ?", id).
		Update("last_crawl_at", t).Error
}
