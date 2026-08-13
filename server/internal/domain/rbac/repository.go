package rbac

import "context"

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *Role) error

	// Update 更新角色
	Update(ctx context.Context, role *Role) error

	// Delete 删除角色
	Delete(ctx context.Context, roleID string) error

	// FindByID 根据 ID 查找角色
	FindByID(ctx context.Context, id string) (*Role, error)

	// FindByCode 根据编码查找角色
	FindByCode(ctx context.Context, code string) (*Role, error)

	// FindAll 获取所有角色
	FindAll(ctx context.Context) ([]*Role, error)

	// FindByUserID 查询用户的所有角色（通过 user_roles 关联表）
	FindByUserID(ctx context.Context, userID string) ([]*Role, error)

	// FindByUserIDs 批量查询多个用户的所有角色（解决 N+1 问题）
	FindByUserIDs(ctx context.Context, userIDs []string) (map[string][]*Role, error)

	// AssignRolesToUser 分配角色给用户（先删后插 user_roles）
	AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error

	// HasRole 检查用户是否拥有指定角色
	HasRole(ctx context.Context, userID string, roleCode string) (bool, error)
}

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *Menu) error

	// Update 更新菜单
	Update(ctx context.Context, menu *Menu) error

	// Delete 删除菜单（含子菜单）
	Delete(ctx context.Context, menuID string) error

	// FindByID 根据 ID 查找菜单
	FindByID(ctx context.Context, id string) (*Menu, error)

	// FindByKey 根据 key 查找菜单
	FindByKey(ctx context.Context, key string) (*Menu, error)

	// FindAll 获取所有菜单（扁平列表）
	FindAll(ctx context.Context) ([]*Menu, error)

	// FindChildren 查询子菜单
	FindChildren(ctx context.Context, parentID string) ([]*Menu, error)

	// UpdateSort 批量更新排序
	UpdateSort(ctx context.Context, menuID string, sortOrder int) error
}
