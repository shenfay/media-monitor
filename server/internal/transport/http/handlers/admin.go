package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-react-admin/internal/app/admin"
	"github.com/shenfay/go-react-admin/internal/domain/user"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
	validationErr "github.com/shenfay/go-react-admin/pkg/errors/validation"
)

// AdminHandler 管理员 HTTP 处理器
type AdminHandler struct {
	service *admin.Service
}

// NewAdminHandler 创建管理员处理器实例
func NewAdminHandler(service *admin.Service) *AdminHandler {
	return &AdminHandler{service: service}
}

// RegisterRoutes 注册管理员路由（路由组已由外部应用 auth + perm 中间件）
func (h *AdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	// 用户管理
	rg.GET("/users", h.ListUsers)
	rg.POST("/users", h.CreateUser)
	rg.PUT("/users/:id", h.UpdateUser)
	rg.PATCH("/users/:id/status", h.ToggleUserStatus)

	// 角色管理
	rg.GET("/roles", h.ListRoles)
	rg.POST("/roles", h.CreateRole)
	rg.PUT("/roles/:id", h.UpdateRole)
	rg.DELETE("/roles/:id", h.DeleteRole)
	rg.PATCH("/roles/:id/status", h.ToggleRoleStatus)

	// 权限管理
	rg.GET("/roles/:id/permissions", h.GetRolePermissions)
	rg.PUT("/roles/:id/permissions", h.UpdateRolePermissions)

	// 菜单管理
	rg.GET("/menus", h.ListMenus)
	rg.POST("/menus", h.CreateMenu)
	rg.PUT("/menus/:id", h.UpdateMenu)
	rg.DELETE("/menus/:id", h.DeleteMenu)
	rg.PATCH("/menus/:id/status", h.ToggleMenuStatus)
	rg.PUT("/menus/sort", h.UpdateMenuSort)
}

// RegisterUserMenuRoutes 注册当前用户权限与菜单路由（只需 auth 中间件）
func (h *AdminHandler) RegisterUserMenuRoutes(rg *gin.RouterGroup) {
	rg.GET("/permissions", h.GetCurrentUserPermissions)
	rg.GET("/menus", h.GetUserMenuTree)
}

// ---- 请求 DTO ----

// createUserRequest 创建用户请求
type createUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Name     string   `json:"name" binding:"required"`
	Password string   `json:"password" binding:"required,min=8"`
	RoleIDs  []string `json:"role_ids"`
}

// updateUserRequest 更新用户请求
type updateUserRequest struct {
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	RoleIDs []string `json:"role_ids"`
}

// toggleStatusRequest 切换状态请求
type toggleStatusRequest struct {
	Locked bool `json:"locked"`
}

// createRoleRequest 创建角色请求
type createRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

// updateRoleRequest 更新角色请求
type updateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// updatePermissionsRequest 更新权限请求
type updatePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// createMenuRequest 创建菜单请求
type createMenuRequest struct {
	Key        string `json:"key" binding:"required"`
	Label      string `json:"label" binding:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
	ParentID   string `json:"parent_id"`
	SortOrder  int    `json:"sort_order"`
}

// updateMenuRequest 更新菜单请求
type updateMenuRequest struct {
	Label      string `json:"label" binding:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

// sortMenuRequest 菜单排序请求
type sortMenuRequest struct {
	Items []struct {
		ID        string `json:"id"`
		SortOrder int    `json:"sort_order"`
	} `json:"items" binding:"required"`
}

// ---- 用户管理 ----

// ListUsers 用户列表
// @Summary 获取用户列表（分页）
// @Tags Admin/Users
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Param keyword query string false "搜索关键词"
// @Param role_id query string false "角色ID"
// @Param status query string false "状态筛选" Enums(active, locked)
// @Success 200 {object} response.SuccessResponse "用户列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := user.UserListParams{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		RoleID:   c.Query("role_id"),
	}

	if status := c.Query("status"); status != "" {
		b := status == "active"
		params.Status = &b
	}

	result, err := h.service.ListUsers(c.Request.Context(), params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// CreateUser 创建用户
// @Summary 管理员创建用户
// @Tags Admin/Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{email=string,name=string,password=string,role_ids=[]string} true "用户创建数据"
// @Success 201 {object} response.SuccessResponse "创建成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 409 {object} response.ErrorResponse "邮箱已存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/users [post]
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := admin.CreateUserCmd{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		RoleIDs:  req.RoleIDs,
	}

	dto, err := h.service.CreateUser(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}

// UpdateUser 更新用户
// @Summary 管理员更新用户信息
// @Tags Admin/Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户ID"
// @Param request body object{name=string,email=string,role_ids=[]string} true "用户更新数据"
// @Success 200 {object} response.SuccessResponse "更新成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "用户不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req updateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := admin.UpdateUserCmd{
		UserID:  userID,
		Name:    req.Name,
		Email:   req.Email,
		RoleIDs: req.RoleIDs,
	}

	dto, err := h.service.UpdateUser(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// ToggleUserStatus 启用/禁用用户
// @Summary 切换用户启用/禁用状态
// @Tags Admin/Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户ID"
// @Param request body object{locked=bool} true "状态数据"
// @Success 200 {object} response.SuccessResponse "状态已更新"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "用户不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/users/{id}/status [patch]
func (h *AdminHandler) ToggleUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var req toggleStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	if err := h.service.ToggleUserStatus(c.Request.Context(), userID, req.Locked); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Status updated"})
}

// ---- 角色管理 ----

// ListRoles 角色列表
// @Summary 获取所有角色列表
// @Tags Admin/Roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "角色列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles [get]
func (h *AdminHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, roles)
}

// CreateRole 创建角色
// @Summary 创建新角色
// @Tags Admin/Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{name=string,code=string,description=string} true "角色创建数据"
// @Success 201 {object} response.SuccessResponse "创建成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 409 {object} response.ErrorResponse "角色编码已存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles [post]
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req createRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := admin.CreateRoleCmd{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	dto, err := h.service.CreateRole(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}

// UpdateRole 更新角色
// @Summary 更新角色信息
// @Tags Admin/Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Param request body object{name=string,description=string} true "角色更新数据"
// @Success 200 {object} response.SuccessResponse "更新成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles/{id} [put]
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")

	var req updateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := admin.UpdateRoleCmd{
		RoleID:      roleID,
		Name:        req.Name,
		Description: req.Description,
	}

	dto, err := h.service.UpdateRole(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Tags Admin/Roles
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Success 200 {object} response.SuccessResponse "删除成功"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles/{id} [delete]
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")

	if err := h.service.DeleteRole(c.Request.Context(), roleID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Role deleted"})
}

// ToggleRoleStatus 切换角色状态
// @Summary 切换角色启用/禁用状态
// @Tags Admin/Roles
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Success 200 {object} response.SuccessResponse "状态已切换"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles/{id}/status [patch]
func (h *AdminHandler) ToggleRoleStatus(c *gin.Context) {
	roleID := c.Param("id")

	if err := h.service.ToggleRoleStatus(c.Request.Context(), roleID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Status toggled"})
}

// ---- 权限管理 ----

// GetRolePermissions 获取角色权限
// @Summary 获取指定角色的权限列表
// @Tags Admin/Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Success 200 {object} response.SuccessResponse "权限列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles/{id}/permissions [get]
func (h *AdminHandler) GetRolePermissions(c *gin.Context) {
	roleID := c.Param("id")

	perms, err := h.service.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, perms)
}

// UpdateRolePermissions 更新角色权限
// @Summary 更新角色的权限分配
// @Tags Admin/Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Param request body object{permissions=[]string} true "权限标识列表"
// @Success 200 {object} response.SuccessResponse "权限已更新"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/roles/{id}/permissions [put]
func (h *AdminHandler) UpdateRolePermissions(c *gin.Context) {
	roleID := c.Param("id")

	var req updatePermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	if err := h.service.UpdateRolePermissions(c.Request.Context(), roleID, req.Permissions); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Permissions updated"})
}

// GetCurrentUserPermissions 获取当前用户权限
// @Summary 获取当前登录用户的权限列表
// @Tags Admin/Permissions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "权限列表"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /auth/permissions [get]
func (h *AdminHandler) GetCurrentUserPermissions(c *gin.Context) {
	userID := c.GetString("user_id")

	perms, err := h.service.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, perms)
}

// GetUserMenuTree 获取当前用户可见的菜单树
// @Summary 获取当前登录用户可见的菜单树
// @Tags Admin/Menus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "菜单树"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /auth/menus [get]
func (h *AdminHandler) GetUserMenuTree(c *gin.Context) {
	userID := c.GetString("user_id")

	tree, err := h.service.GetUserMenuTree(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, tree)
}

// ---- 菜单管理 ----

// ListMenus 获取菜单树
// @Summary 获取完整菜单树（管理端）
// @Tags Admin/Menus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "菜单树"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/menus [get]
func (h *AdminHandler) ListMenus(c *gin.Context) {
	tree, err := h.service.ListMenuTree(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, tree)
}

// CreateMenu 创建菜单
// @Summary 创建新菜单
// @Tags Admin/Menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{key=string,label=string,icon=string,path=string,permission=string,parent_id=string,sort_order=int} true "菜单创建数据"
// @Success 201 {object} response.SuccessResponse "创建成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 409 {object} response.ErrorResponse "菜单Key已存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/menus [post]
func (h *AdminHandler) CreateMenu(c *gin.Context) {
	var req createMenuRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := admin.CreateMenuCmd{
		Key:        req.Key,
		Label:      req.Label,
		Icon:       req.Icon,
		Path:       req.Path,
		Permission: req.Permission,
		ParentID:   req.ParentID,
		SortOrder:  req.SortOrder,
	}

	dto, err := h.service.CreateMenu(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}

// UpdateMenu 更新菜单
// @Summary 更新菜单信息
// @Tags Admin/Menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Param request body object{label=string,icon=string,path=string,permission=string} true "菜单更新数据"
// @Success 200 {object} response.SuccessResponse "更新成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "菜单不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/menus/{id} [put]
func (h *AdminHandler) UpdateMenu(c *gin.Context) {
	menuID := c.Param("id")

	var req updateMenuRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := admin.UpdateMenuCmd{
		MenuID:     menuID,
		Label:      req.Label,
		Icon:       req.Icon,
		Path:       req.Path,
		Permission: req.Permission,
	}

	dto, err := h.service.UpdateMenu(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Tags Admin/Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Success 200 {object} response.SuccessResponse "删除成功"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "菜单不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/menus/{id} [delete]
func (h *AdminHandler) DeleteMenu(c *gin.Context) {
	menuID := c.Param("id")

	if err := h.service.DeleteMenu(c.Request.Context(), menuID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Menu deleted"})
}

// ToggleMenuStatus 切换菜单状态
// @Summary 切换菜单启用/禁用状态
// @Tags Admin/Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Success 200 {object} response.SuccessResponse "状态已切换"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "菜单不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/menus/{id}/status [patch]
func (h *AdminHandler) ToggleMenuStatus(c *gin.Context) {
	menuID := c.Param("id")

	if err := h.service.ToggleMenuStatus(c.Request.Context(), menuID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Status toggled"})
}

// UpdateMenuSort 更新菜单排序
// @Summary 批量更新菜单排序顺序
// @Tags Admin/Menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{items=[]object{id=string,sort_order=int}} true "排序数据"
// @Success 200 {object} response.SuccessResponse "排序已更新"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /admin/menus/sort [put]
func (h *AdminHandler) UpdateMenuSort(c *gin.Context) {
	var req sortMenuRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	items := make([]admin.SortMenuItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, admin.SortMenuItem{
			ID:        item.ID,
			SortOrder: item.SortOrder,
		})
	}

	cmd := admin.SortMenuCmd{Items: items}
	if err := h.service.UpdateMenuSort(c.Request.Context(), cmd); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Sort updated"})
}
