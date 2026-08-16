package rbac

import (
	"time"

	"github.com/shenfay/go-react-admin/pkg/utils"
)

// Menu 菜单聚合根
type Menu struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Icon        string    `json:"icon"`
	Path        string    `json:"path"`
	Permission  string    `json:"permission"`  // 主权限标识（用于菜单可见性）
	Permissions []string  `json:"permissions"` // 关联的所有 API 权限
	ParentID    string    `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewMenu 创建新菜单
func NewMenu(key, label, icon, path, permission, parentID string, sortOrder int) *Menu {
	now := utils.Now()
	return &Menu{
		ID:         utils.GenerateID(),
		Key:        key,
		Label:      label,
		Icon:       icon,
		Path:       path,
		Permission: permission,
		ParentID:   parentID,
		SortOrder:  sortOrder,
		Status:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Update 更新菜单信息
func (m *Menu) Update(label, icon, path, permission string) {
	m.Label = label
	m.Icon = icon
	m.Path = path
	m.Permission = permission
	m.UpdatedAt = utils.Now()
}

// ToggleStatus 切换菜单状态
func (m *Menu) ToggleStatus() {
	m.Status = !m.Status
	m.UpdatedAt = utils.Now()
}

// MenuTreeNode 菜单树节点（用于返回树形结构）
type MenuTreeNode struct {
	ID          string          `json:"id"`
	Key         string          `json:"key"`
	Label       string          `json:"label"`
	Icon        string          `json:"icon"`
	Path        string          `json:"path"`
	Permission  string          `json:"permission"`
	Permissions []string        `json:"permissions"`
	ParentID    string          `json:"parent_id"`
	SortOrder   int             `json:"sort_order"`
	Status      bool            `json:"status"`
	Children    []*MenuTreeNode `json:"children,omitempty"`
}
