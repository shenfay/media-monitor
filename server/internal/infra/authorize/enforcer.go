package authorize

import (
	_ "embed"

	casbin "github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

//go:embed model.conf
var modelConf string

// Enforcer Casbin 权限引擎封装
type Enforcer struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
}

// NewEnforcer 创建 Casbin Enforcer 实例
// 初始化 gorm-adapter，加载模型，同步 user_roles → g 规则
func NewEnforcer(db *gorm.DB) (*Enforcer, error) {
	// 创建 gorm adapter（自动创建 casbin_rule 表）
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 从嵌入的模型字符串创建 enforcer
	m, err := model.NewModelFromString(modelConf)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	e := &Enforcer{
		enforcer: enforcer,
		db:       db,
	}

	// 同步 user_roles 表到 casbin g 规则
	if err := e.syncUserRolesFromDB(); err != nil {
		return nil, err
	}

	return e, nil
}

// Enforce 检查用户是否拥有指定权限
func (e *Enforcer) Enforce(userID, permission string) (bool, error) {
	return e.enforcer.Enforce(userID, permission)
}

// GetPermissionsForRole 查询角色拥有的所有 permission_key
func (e *Enforcer) GetPermissionsForRole(roleID string) ([]string, error) {
	permissions, err := e.enforcer.GetPermissionsForUser(roleID)
	if err != nil {
		return nil, err
	}
	// 过滤出 p 规则（g 规则的 len == 1，p 规则的 len == 2）
	result := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if len(p) == 2 && p[0] == roleID {
			result = append(result, p[1])
		}
	}
	return result, nil
}

// GetPermissionsForUser 查询用户通过所有角色获得的 permission_key（去重）
func (e *Enforcer) GetPermissionsForUser(userID string) ([]string, error) {
	// 获取用户的所有角色
	roles, err := e.enforcer.GetRolesForUser(userID)
	if err != nil {
		return nil, err
	}

	permSet := make(map[string]bool)
	for _, role := range roles {
		permissions, err := e.GetPermissionsForRole(role)
		if err != nil {
			return nil, err
		}
		for _, p := range permissions {
			permSet[p] = true
		}
	}

	result := make([]string, 0, len(permSet))
	for p := range permSet {
		result = append(result, p)
	}
	return result, nil
}

// SetRolePermissions 设置角色的权限列表（先删后插）
func (e *Enforcer) SetRolePermissions(roleID string, perms []string) error {
	// 删除角色的所有现有权限
	_, err := e.enforcer.RemoveFilteredPolicy(0, roleID)
	if err != nil {
		return err
	}

	// 添加新权限
	for _, perm := range perms {
		if _, err := e.enforcer.AddPolicy(roleID, perm); err != nil {
			return err
		}
	}

	return nil
}

// SyncUserRoles 同步用户角色（先删所有 g 规则，再添加新的）
func (e *Enforcer) SyncUserRoles(userID string, roleIDs []string) error {
	// 删除用户的所有角色关联
	_, err := e.enforcer.RemoveFilteredGroupingPolicy(0, userID)
	if err != nil {
		return err
	}

	// 添加新角色
	for _, roleID := range roleIDs {
		if _, err := e.enforcer.AddGroupingPolicy(userID, roleID); err != nil {
			return err
		}
	}

	return nil
}

// AddRoleForUser 为用户添加单个角色
func (e *Enforcer) AddRoleForUser(userID, roleID string) error {
	_, err := e.enforcer.AddGroupingPolicy(userID, roleID)
	return err
}

// RemoveRoleForUser 移除用户的单个角色
func (e *Enforcer) RemoveRoleForUser(userID, roleID string) error {
	_, err := e.enforcer.RemoveGroupingPolicy(userID, roleID)
	return err
}

// syncUserRolesFromDB 从 user_roles 表同步所有用户角色到 casbin g 规则
func (e *Enforcer) syncUserRolesFromDB() error {
	type userRole struct {
		UserID string `gorm:"column:user_id"`
		RoleID string `gorm:"column:role_id"`
	}

	var userRoles []userRole
	if err := e.db.Table("user_roles").Find(&userRoles).Error; err != nil {
		// 表可能不存在（首次迁移前），忽略
		return nil
	}

	for _, ur := range userRoles {
		// 检查是否已存在
		has, _ := e.enforcer.HasGroupingPolicy(ur.UserID, ur.RoleID)
		if !has {
			if _, err := e.enforcer.AddGroupingPolicy(ur.UserID, ur.RoleID); err != nil {
				return err
			}
		}
	}

	return nil
}
