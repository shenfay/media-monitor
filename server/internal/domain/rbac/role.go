package rbac

import (
	"time"

	"github.com/shenfay/go-react-admin/pkg/utils"
)

// Role 角色聚合根
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewRole 创建新角色
func NewRole(name, code, description string) *Role {
	now := utils.Now()
	return &Role{
		ID:          utils.GenerateID(),
		Name:        name,
		Code:        code,
		Description: description,
		Status:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update 更新角色信息
func (r *Role) Update(name, description string) {
	r.Name = name
	r.Description = description
	r.UpdatedAt = utils.Now()
}

// ToggleStatus 切换角色状态
func (r *Role) ToggleStatus() {
	r.Status = !r.Status
	r.UpdatedAt = utils.Now()
}
