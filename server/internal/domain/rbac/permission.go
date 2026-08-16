package rbac

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
// 使用数据库中菜单的 Permissions 字段进行匹配，fallback 到静态 PermissionMenuMap
// 没有配置任何权限的菜单（permissions 和 permission 均为空）默认可见
func DeriveMenusFromMenus(permissions []string, menus []*Menu) []string {
	// 构建 permission → menu_key 动态映射（使用 Permissions 数组）
	permToMenu := make(map[string]string)
	// 无需权限即可见的菜单集合
	alwaysVisible := make(map[string]bool)

	for _, m := range menus {
		if m.Status {
			// 优先使用 Permissions 数组（包含 view + manage 等所有权限）
			if len(m.Permissions) > 0 {
				for _, p := range m.Permissions {
					if p != "" {
						permToMenu[p] = m.Key
					}
				}
			} else if m.Permission != "" {
				// 兼容：如果 Permissions 为空，使用旧的 Permission 单字段
				permToMenu[m.Permission] = m.Key
			} else {
				// 没有配置任何权限的菜单，默认可见
				alwaysVisible[m.Key] = true
			}
		}
	}

	menuSet := make(map[string]bool)
	// 添加默认可见的菜单
	for key := range alwaysVisible {
		menuSet[key] = true
	}
	// 根据用户权限匹配菜单
	for _, perm := range permissions {
		// 优先从数据库菜单映射查找
		if menu, ok := permToMenu[perm]; ok {
			menuSet[menu] = true
		} else if menu, ok := PermissionMenuMap[perm]; ok {
			// fallback 到静态映射
			menuSet[menu] = true
		}
	}
	result := make([]string, 0, len(menuSet))
	for m := range menuSet {
		result = append(result, m)
	}
	return result
}
