package rbac

import "strings"

// RolePermission 角色权限值对象（权限由 Casbin 管理，此处仅保留领域模型）
type RolePermission struct {
	RoleID        string `json:"role_id"`
	PermissionKey string `json:"permission_key"`
}

// UserPermission 用户权限聚合（登录时返回）
type UserPermission struct {
	Roles       []RoleBrief `json:"roles"`
	Permissions []string    `json:"permissions"`
	Menus       []string    `json:"menus"`
}

// RoleBrief 角色简要信息
type RoleBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// MenuItem 菜单权限树节点
type MenuItem struct {
	Key      string     `json:"key"`
	Title    string     `json:"title"`
	Children []MenuItem `json:"children,omitempty"`
}

// PermissionMenuMap 权限标识 → 菜单 key 静态映射
var PermissionMenuMap = map[string]string{
	"dashboard:view":        "dashboard",
	"family:manage":         "family",
	"goal:manage":           "goals",
	"card_template:manage":  "card-templates",
	"card_instance:view":    "card-instances",
	"companion:manage":      "companions",
	"acceptance:manage":     "acceptance",
	"points:view":           "points",
	"shop_item:manage":      "shop-items",
	"exchange_order:manage": "exchange-orders",
	"user:manage":           "user-management",
	"user:list":             "user-management",
	"user:create":           "user-management",
	"user:update":           "user-management",
	"permission:manage":     "permission-management",
	"permission:view":       "permission-management",
	"menu:manage":           "menu-management",
	"profile:view":          "profile",
	"operation:log":         "operation-log",
	"setting:manage":        "system-settings",
}

// DeriveMenus 根据权限列表推导菜单 key 列表（去重）
// 保留作为 fallback，优先使用 DeriveMenusFromMenus
func DeriveMenus(permissions []string) []string {
	menuSet := make(map[string]bool)
	for _, perm := range permissions {
		if menu, ok := PermissionMenuMap[perm]; ok {
			menuSet[menu] = true
		}
	}
	menus := make([]string, 0, len(menuSet))
	for m := range menuSet {
		menus = append(menus, m)
	}
	return menus
}

// DeriveMenusFromMenus 根据权限列表和数据库菜单推导菜单 key 列表（去重）
// 优先使用数据库中菜单的 Permission 字段进行匹配，fallback 到静态 PermissionMenuMap
func DeriveMenusFromMenus(permissions []string, menus []*Menu) []string {
	// 构建 permission → menu_key 动态映射
	permToMenu := make(map[string]string)
	// 所有有效的菜单 key 集合，用于兜底匹配
	menuKeySet := make(map[string]bool)
	for _, m := range menus {
		if m.Status {
			if m.Permission != "" {
				permToMenu[m.Permission] = m.Key
			}
			menuKeySet[m.Key] = true
		}
	}

	menuSet := make(map[string]bool)
	for _, perm := range permissions {
		// 优先从数据库菜单映射查找
		if menu, ok := permToMenu[perm]; ok {
			menuSet[menu] = true
		} else if menu, ok := PermissionMenuMap[perm]; ok {
			// fallback 到静态映射
			menuSet[menu] = true
		} else if idx := strings.LastIndex(perm, ":"); idx > 0 {
			// 兜底：从 permission 中提取 key 前缀匹配菜单
			// 如 "design-system:view" → "design-system"
			if key := perm[:idx]; menuKeySet[key] {
				menuSet[key] = true
			}
		}
	}
	result := make([]string, 0, len(menuSet))
	for m := range menuSet {
		result = append(result, m)
	}
	return result
}
