package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/rbac"
	"github.com/shenfay/go-react-admin/internal/domain/shared"
)

// MenuPO 菜单持久化对象
type MenuPO struct {
	ID         string         `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Key        string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"key"`
	Label      string         `gorm:"type:varchar(100);not null" json:"label"`
	Icon       string         `gorm:"type:varchar(100);default:''" json:"icon"`
	Path       string         `gorm:"type:varchar(200);default:''" json:"path"`
	Permission string         `gorm:"type:varchar(100);default:''" json:"permission"`
	ParentID   sql.NullString `gorm:"type:varchar(50);index" json:"parent_id"`
	SortOrder  int            `gorm:"default:0;index" json:"sort_order"`
	Status     bool           `gorm:"default:true" json:"status"`
	CreatedAt  TimeNull       `json:"created_at"`
	UpdatedAt  TimeNull       `json:"updated_at"`
}

func (MenuPO) TableName() string { return "menus" }

// ToDomain 转换为领域模型
func (po *MenuPO) ToDomain() *rbac.Menu {
	if po == nil {
		return nil
	}
	createdAt := time.Time{}
	updatedAt := time.Time{}
	if po.CreatedAt.Valid {
		createdAt = po.CreatedAt.Time
	}
	if po.UpdatedAt.Valid {
		updatedAt = po.UpdatedAt.Time
	}
	return &rbac.Menu{
		ID:         po.ID,
		Key:        po.Key,
		Label:      po.Label,
		Icon:       po.Icon,
		Path:       po.Path,
		Permission: po.Permission,
		ParentID:   po.ParentID.String,
		SortOrder:  po.SortOrder,
		Status:     po.Status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// MenuPOFromDomain 从领域模型转换
func MenuPOFromDomain(m *rbac.Menu) *MenuPO {
	po := &MenuPO{
		ID:         m.ID,
		Key:        m.Key,
		Label:      m.Label,
		Icon:       m.Icon,
		Path:       m.Path,
		Permission: m.Permission,
		SortOrder:  m.SortOrder,
		Status:     m.Status,
		CreatedAt:  TimeNull{Time: m.CreatedAt, Valid: true},
		UpdatedAt:  TimeNull{Time: m.UpdatedAt, Valid: true},
	}
	if m.ParentID != "" {
		po.ParentID = sql.NullString{String: m.ParentID, Valid: true}
	}
	return po
}

// menuRepository 菜单仓储 GORM 实现
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储实例
func NewMenuRepository(db *gorm.DB) rbac.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(ctx context.Context, menu *rbac.Menu) error {
	po := MenuPOFromDomain(menu)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *menuRepository) Update(ctx context.Context, menu *rbac.Menu) error {
	po := MenuPOFromDomain(menu)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *menuRepository) Delete(ctx context.Context, menuID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 递归删除子菜单
		if err := r.deleteChildren(ctx, tx, menuID); err != nil {
			return err
		}
		return tx.Delete(&MenuPO{}, "id = ?", menuID).Error
	})
}

func (r *menuRepository) deleteChildren(ctx context.Context, tx *gorm.DB, parentID string) error {
	var children []MenuPO
	if err := tx.WithContext(ctx).Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}
	for _, child := range children {
		if err := r.deleteChildren(ctx, tx, child.ID); err != nil {
			return err
		}
		if err := tx.Delete(&MenuPO{}, "id = ?", child.ID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *menuRepository) FindByID(ctx context.Context, id string) (*rbac.Menu, error) {
	var po MenuPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *menuRepository) FindByKey(ctx context.Context, key string) (*rbac.Menu, error) {
	var po MenuPO
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *menuRepository) FindAll(ctx context.Context) ([]*rbac.Menu, error) {
	var pos []MenuPO
	if err := r.db.WithContext(ctx).Order("sort_order ASC, created_at ASC").Find(&pos).Error; err != nil {
		return nil, err
	}
	menus := make([]*rbac.Menu, 0, len(pos))
	for i := range pos {
		menus = append(menus, pos[i].ToDomain())
	}
	return menus, nil
}

func (r *menuRepository) FindChildren(ctx context.Context, parentID string) ([]*rbac.Menu, error) {
	var pos []MenuPO
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort_order ASC").Find(&pos).Error; err != nil {
		return nil, err
	}
	menus := make([]*rbac.Menu, 0, len(pos))
	for i := range pos {
		menus = append(menus, pos[i].ToDomain())
	}
	return menus, nil
}

func (r *menuRepository) UpdateSort(ctx context.Context, menuID string, sortOrder int) error {
	return r.db.WithContext(ctx).Model(&MenuPO{}).Where("id = ?", menuID).Update("sort_order", sortOrder).Error
}
