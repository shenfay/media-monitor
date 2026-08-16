package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/shenfay/go-react-admin/internal/app/shared/operationlog"
	"github.com/shenfay/go-react-admin/internal/domain/rbac"
	"github.com/shenfay/go-react-admin/internal/domain/shared/events"
	"github.com/shenfay/go-react-admin/internal/domain/user"
	"github.com/shenfay/go-react-admin/internal/infra/authorize"
	"github.com/shenfay/go-react-admin/pkg/errors"
)

// Service 管理员应用服务
type Service struct {
	userRepo user.UserRepository
	roleRepo rbac.RoleRepository
	menuRepo rbac.MenuRepository
	enforcer *authorize.Enforcer
	recorder *operationlog.OperationRecorder
}

// NewService 创建管理员应用服务
func NewService(userRepo user.UserRepository, roleRepo rbac.RoleRepository, menuRepo rbac.MenuRepository, enforcer *authorize.Enforcer, eventBus events.Bus) *Service {
	return &Service{
		userRepo: userRepo,
		roleRepo: roleRepo,
		menuRepo: menuRepo,
		enforcer: enforcer,
		recorder: operationlog.NewOperationRecorder(eventBus),
	}
}

// ---- 用户管理 ----

// ListUsers 分页查询用户列表
func (s *Service) ListUsers(ctx context.Context, params user.UserListParams) (*UserListDTO, error) {
	result, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	// 批量查询所有用户的角色（解决 N+1 问题）
	userIDs := make([]string, 0, len(result.Users))
	for _, u := range result.Users {
		userIDs = append(userIDs, u.ID)
	}
	rolesMap, _ := s.roleRepo.FindByUserIDs(ctx, userIDs)

	userDTOs := make([]*UserDTO, 0, len(result.Users))
	for _, u := range result.Users {
		dto := toUserDTO(u)
		if roles, ok := rolesMap[u.ID]; ok {
			dto.Roles = rolesToBriefs(roles)
		}
		userDTOs = append(userDTOs, dto)
	}

	return &UserListDTO{
		Users: userDTOs,
		Total: result.Total,
	}, nil
}

// CreateUser 创建用户并分配角色
func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCmd) (*UserDTO, error) {
	u, err := user.NewUser(cmd.Email, cmd.Name, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	// 分配角色（DB + Casbin 同步）
	if len(cmd.RoleIDs) > 0 {
		if err := s.roleRepo.AssignRolesToUser(ctx, u.ID, cmd.RoleIDs); err != nil {
			return nil, err
		}
		if err := s.enforcer.SyncUserRoles(u.ID, cmd.RoleIDs); err != nil {
			return nil, err
		}
	}

	dto := s.buildUserDTO(ctx, u)

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "USER.CREATE", "USER", "SUCCESS",
		map[string]interface{}{"target_user_id": u.ID, "email": u.Email},
	)

	return dto, nil
}

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(ctx context.Context, cmd UpdateUserCmd) (*UserDTO, error) {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if cmd.Name != "" {
		_ = u.UpdateName(cmd.Name)
	}
	if cmd.Email != "" {
		_ = u.UpdateEmail(cmd.Email)
	}

	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 更新角色（DB + Casbin 同步）
	if cmd.RoleIDs != nil {
		if err := s.roleRepo.AssignRolesToUser(ctx, u.ID, cmd.RoleIDs); err != nil {
			return nil, err
		}
		if err := s.enforcer.SyncUserRoles(u.ID, cmd.RoleIDs); err != nil {
			return nil, err
		}
	}

	dto := s.buildUserDTO(ctx, u)

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "USER.UPDATE", "USER", "SUCCESS",
		map[string]interface{}{"target_user_id": u.ID, "name": cmd.Name, "email": cmd.Email},
	)

	return dto, nil
}

// ToggleUserStatus 启用/禁用用户
func (s *Service) ToggleUserStatus(ctx context.Context, userID string, locked bool) error {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	u.SetLocked(locked)

	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	// 记录操作日志
	action := "USER.ENABLE"
	if locked {
		action = "USER.DISABLE"
	}
	s.recorder.RecordFromContext(ctx, action, "USER", "SUCCESS",
		map[string]interface{}{"target_user_id": userID},
	)

	return nil
}

// ---- 角色管理 ----

// ListRoles 获取所有角色
func (s *Service) ListRoles(ctx context.Context) ([]*RoleDTO, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*RoleDTO, 0, len(roles))
	for _, r := range roles {
		dtos = append(dtos, toRoleDTO(r))
	}
	return dtos, nil
}

// CreateRole 创建角色
func (s *Service) CreateRole(ctx context.Context, cmd CreateRoleCmd) (*RoleDTO, error) {
	r := rbac.NewRole(cmd.Name, cmd.Code, cmd.Description)
	if err := s.roleRepo.Create(ctx, r); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "ROLE.CREATE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"role_id": r.ID, "role_name": cmd.Name, "role_code": cmd.Code},
	)

	return toRoleDTO(r), nil
}

// UpdateRole 更新角色
func (s *Service) UpdateRole(ctx context.Context, cmd UpdateRoleCmd) (*RoleDTO, error) {
	r, err := s.roleRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return nil, err
	}

	r.Update(cmd.Name, cmd.Description)
	if err := s.roleRepo.Update(ctx, r); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "ROLE.UPDATE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"role_id": cmd.RoleID, "role_name": cmd.Name},
	)

	return toRoleDTO(r), nil
}

// DeleteRole 删除角色
func (s *Service) DeleteRole(ctx context.Context, roleID string) error {
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "ROLE.DELETE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"role_id": roleID},
	)

	return nil
}

// ToggleRoleStatus 切换角色状态
func (s *Service) ToggleRoleStatus(ctx context.Context, roleID string) error {
	r, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	r.ToggleStatus()
	if err := s.roleRepo.Update(ctx, r); err != nil {
		return err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "ROLE.TOGGLE_STATUS", "SYSTEM", "SUCCESS",
		map[string]interface{}{"role_id": roleID, "status": r.Status},
	)

	return nil
}

// ---- 权限管理（通过 Casbin）----

// GetRolePermissions 获取角色权限列表（从 Casbin 查询）
func (s *Service) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	return s.enforcer.GetPermissionsForRole(roleID)
}

// UpdateRolePermissions 更新角色权限（通过 Casbin）
func (s *Service) UpdateRolePermissions(ctx context.Context, roleID string, permissions []string) error {
	if err := s.enforcer.SetRolePermissions(roleID, permissions); err != nil {
		return err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "PERMISSION.UPDATE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"role_id": roleID, "permissions_count": len(permissions)},
	)

	return nil
}

// ---- 菜单管理 ----

// ListMenuTree 获取菜单树
func (s *Service) ListMenuTree(ctx context.Context) ([]*rbac.MenuTreeNode, error) {
	menus, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, ""), nil
}

// CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx context.Context, cmd CreateMenuCmd) (*MenuDTO, error) {
	// 检查 key 是否已存在
	existing, _ := s.menuRepo.FindByKey(ctx, cmd.Key)
	if existing != nil {
		return nil, errors.NewAppError(
			errors.ErrCodeMenuKeyAlreadyExists,
			"菜单标识已存在",
			http.StatusConflict,
		)
	}

	m := rbac.NewMenu(cmd.Key, cmd.Label, cmd.Icon, cmd.Path, cmd.Permission, cmd.ParentID, cmd.SortOrder)
	if err := s.menuRepo.Create(ctx, m); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "MENU.CREATE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"menu_id": m.ID, "menu_key": cmd.Key, "label": cmd.Label},
	)

	return toMenuDTO(m), nil
}

// UpdateMenu 更新菜单
func (s *Service) UpdateMenu(ctx context.Context, cmd UpdateMenuCmd) (*MenuDTO, error) {
	m, err := s.menuRepo.FindByID(ctx, cmd.MenuID)
	if err != nil {
		return nil, err
	}

	m.Update(cmd.Label, cmd.Icon, cmd.Path, cmd.Permission)
	if err := s.menuRepo.Update(ctx, m); err != nil {
		return nil, err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "MENU.UPDATE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"menu_id": cmd.MenuID, "label": cmd.Label},
	)

	return toMenuDTO(m), nil
}

// DeleteMenu 删除菜单
func (s *Service) DeleteMenu(ctx context.Context, menuID string) error {
	if err := s.menuRepo.Delete(ctx, menuID); err != nil {
		return err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "MENU.DELETE", "SYSTEM", "SUCCESS",
		map[string]interface{}{"menu_id": menuID},
	)

	return nil
}

// ToggleMenuStatus 切换菜单状态
func (s *Service) ToggleMenuStatus(ctx context.Context, menuID string) error {
	m, err := s.menuRepo.FindByID(ctx, menuID)
	if err != nil {
		return err
	}
	m.ToggleStatus()
	if err := s.menuRepo.Update(ctx, m); err != nil {
		return err
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "MENU.TOGGLE_STATUS", "SYSTEM", "SUCCESS",
		map[string]interface{}{"menu_id": menuID, "status": m.Status},
	)

	return nil
}

// UpdateMenuSort 更新菜单排序
func (s *Service) UpdateMenuSort(ctx context.Context, cmd SortMenuCmd) error {
	for _, item := range cmd.Items {
		if err := s.menuRepo.UpdateSort(ctx, item.ID, item.SortOrder); err != nil {
			return err
		}
	}

	// 记录操作日志
	s.recorder.RecordFromContext(ctx, "MENU.SORT", "SYSTEM", "SUCCESS",
		map[string]interface{}{"items_count": len(cmd.Items)},
	)

	return nil
}

// GetUserPermissions 获取用户权限（登录时使用）
func (s *Service) GetUserPermissions(ctx context.Context, userID string) (*rbac.UserPermission, error) {
	perm, _, err := s.getUserPermissionsWithMenus(ctx, userID)
	return perm, err
}

// getUserPermissionsWithMenus 内部方法：获取用户权限并返回菜单列表（供 GetUserMenuTree 复用）
func (s *Service) getUserPermissionsWithMenus(ctx context.Context, userID string) (*rbac.UserPermission, []*rbac.Menu, error) {
	// 1. 查用户角色
	roles, err := s.roleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	roleBriefs := make([]rbac.RoleBrief, 0, len(roles))
	for _, role := range roles {
		roleBriefs = append(roleBriefs, rbac.RoleBrief{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	// 2. 从 Casbin 查权限
	permissions, err := s.enforcer.GetPermissionsForUser(userID)
	if err != nil {
		return nil, nil, err
	}

	// 3. 从数据库查询所有菜单，动态推导菜单 key
	allMenus, _ := s.menuRepo.FindAll(ctx)
	menus := rbac.DeriveMenusFromMenus(permissions, allMenus)

	perm := &rbac.UserPermission{
		Roles:       roleBriefs,
		Permissions: permissions,
		Menus:       menus,
	}
	return perm, allMenus, nil
}

// buildUserDTO 构建用户 DTO（包含角色信息）
func (s *Service) buildUserDTO(ctx context.Context, u *user.User) *UserDTO {
	dto := toUserDTO(u)
	if roles, _ := s.roleRepo.FindByUserID(ctx, u.ID); roles != nil {
		dto.Roles = rolesToBriefs(roles)
	}
	return dto
}

// ---- 命令对象 ----

// CreateMenuCmd 创建菜单命令
type CreateMenuCmd struct {
	Key        string
	Label      string
	Icon       string
	Path       string
	Permission string
	ParentID   string
	SortOrder  int
}

// UpdateMenuCmd 更新菜单命令
type UpdateMenuCmd struct {
	MenuID     string
	Label      string
	Icon       string
	Path       string
	Permission string
}

// SortMenuCmd 菜单排序命令
type SortMenuCmd struct {
	Items []SortMenuItem
}

// SortMenuItem 排序项
type SortMenuItem struct {
	ID        string
	SortOrder int
}

// CreateUserCmd 创建用户命令
type CreateUserCmd struct {
	Email    string
	Name     string
	Password string
	RoleIDs  []string
}

// UpdateUserCmd 更新用户命令
type UpdateUserCmd struct {
	UserID  string
	Name    string
	Email   string
	RoleIDs []string // nil 表示不更新角色
}

// CreateRoleCmd 创建角色命令
type CreateRoleCmd struct {
	Name        string
	Code        string
	Description string
}

// UpdateRoleCmd 更新角色命令
type UpdateRoleCmd struct {
	RoleID      string
	Name        string
	Description string
}

// ---- DTO ----

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID            string           `json:"id"`
	Email         string           `json:"email"`
	Name          string           `json:"name"`
	EmailVerified bool             `json:"email_verified"`
	Locked        bool             `json:"locked"`
	Roles         []rbac.RoleBrief `json:"roles"`
	LastLoginAt   *time.Time       `json:"last_login_at,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// UserListDTO 用户列表 DTO
type UserListDTO struct {
	Users []*UserDTO `json:"users"`
	Total int64      `json:"total"`
}

// RoleDTO 角色数据传输对象
type RoleDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MenuDTO 菜单数据传输对象
type MenuDTO struct {
	ID          string     `json:"id"`
	Key         string     `json:"key"`
	Label       string     `json:"label"`
	Icon        string     `json:"icon"`
	Path        string     `json:"path"`
	Permission  string     `json:"permission"`
	Permissions []string   `json:"permissions"`
	ParentID    string     `json:"parent_id"`
	SortOrder   int        `json:"sort_order"`
	Status      bool       `json:"status"`
	Children    []*MenuDTO `json:"children,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ---- 转换函数 ----

func toUserDTO(u *user.User) *UserDTO {
	return &UserDTO{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name,
		EmailVerified: u.EmailVerified,
		Locked:        u.Locked,
		Roles:         []rbac.RoleBrief{},
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func toRoleDTO(r *rbac.Role) *RoleDTO {
	return &RoleDTO{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Description: r.Description,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func rolesToBriefs(roles []*rbac.Role) []rbac.RoleBrief {
	briefs := make([]rbac.RoleBrief, 0, len(roles))
	for _, r := range roles {
		briefs = append(briefs, rbac.RoleBrief{
			ID:   r.ID,
			Name: r.Name,
			Code: r.Code,
		})
	}
	return briefs
}

func toMenuDTO(m *rbac.Menu) *MenuDTO {
	return &MenuDTO{
		ID:          m.ID,
		Key:         m.Key,
		Label:       m.Label,
		Icon:        m.Icon,
		Path:        m.Path,
		Permission:  m.Permission,
		Permissions: m.Permissions,
		ParentID:    m.ParentID,
		SortOrder:   m.SortOrder,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func buildMenuTree(menus []*rbac.Menu, parentID string) []*rbac.MenuTreeNode {
	// 一次遍历建立 parentID → children 索引，O(n) 构建树
	childrenMap := make(map[string][]*rbac.Menu, len(menus))
	for _, m := range menus {
		childrenMap[m.ParentID] = append(childrenMap[m.ParentID], m)
	}
	var build func(pid string) []*rbac.MenuTreeNode
	build = func(pid string) []*rbac.MenuTreeNode {
		children := childrenMap[pid]
		if len(children) == 0 {
			return nil
		}
		nodes := make([]*rbac.MenuTreeNode, 0, len(children))
		for _, m := range children {
			nodes = append(nodes, &rbac.MenuTreeNode{
				ID:          m.ID,
				Key:         m.Key,
				Label:       m.Label,
				Icon:        m.Icon,
				Path:        m.Path,
				Permission:  m.Permission,
				Permissions: m.Permissions,
				ParentID:    m.ParentID,
				SortOrder:   m.SortOrder,
				Status:      m.Status,
				Children:    build(m.ID),
			})
		}
		return nodes
	}
	return build(parentID)
}

// GetUserMenuTree 获取当前用户可见的菜单树（根据权限过滤）
func (s *Service) GetUserMenuTree(ctx context.Context, userID string) ([]*rbac.MenuTreeNode, error) {
	// 1. 获取用户权限并复用菜单查询结果（只查一次 FindAll）
	userPerm, allMenus, err := s.getUserPermissionsWithMenus(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 构建用户可见菜单 key 集合
	menuKeySet := make(map[string]bool, len(userPerm.Menus))
	for _, k := range userPerm.Menus {
		menuKeySet[k] = true
	}

	// 3. 过滤出用户可见的菜单（包含父级菜单）
	visibleMenus := make([]*rbac.Menu, 0)
	parentIDs := make(map[string]bool)
	includedIDs := make(map[string]bool)
	for _, m := range allMenus {
		if menuKeySet[m.Key] && m.Status {
			visibleMenus = append(visibleMenus, m)
			includedIDs[m.ID] = true
			if m.ParentID != "" {
				parentIDs[m.ParentID] = true
			}
		}
	}
	// 添加父级菜单（分组菜单），跳过已包含的
	for _, m := range allMenus {
		if parentIDs[m.ID] && m.ParentID == "" && !includedIDs[m.ID] {
			visibleMenus = append(visibleMenus, m)
			includedIDs[m.ID] = true
		}
	}

	// 4. 构建树
	return buildMenuTree(visibleMenus, ""), nil
}
