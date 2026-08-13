/**
 * 全局类型定义
 */

// 通用响应结构
export interface ApiResponse<T = unknown> {
  code: string
  message: string
  data: T
}

// 分页参数
export interface PaginationParams {
  current: number
  pageSize: number
}

// 分页结果
export interface PaginationResult<T> {
  list: T[]
  total: number
  current: number
  pageSize: number
}

// 角色简要信息
export interface RoleBrief {
  id: string
  name: string
  code: string
}

// 用户权限信息
export interface UserPermission {
  roles: RoleBrief[]
  permissions: string[]
  menus: string[]
}

// 用户
export interface User {
  id: string
  email: string
  name: string
  email_verified: boolean
  locked: boolean
  roles: RoleBrief[]
  last_login_at?: string
  created_at: string
  updated_at: string
}

// 用户列表响应
export interface UserListResponse {
  users: User[]
  total: number
}

// 角色
export interface Role {
  id: string
  name: string
  code: string
  description: string
  status: boolean
  created_at: string
  updated_at: string
}

// 角色权限
export interface RolePermission {
  id?: number
  role_id: string
  permission_key: string
  menu_key: string
}

// 登录请求
export interface LoginRequest {
  email: string
  password: string
}

// 登录响应
export interface LoginResponse {
  user: {
    id: string
    email: string
    name: string
    email_verified: boolean
    created_at: string
  }
  access_token: string
  refresh_token: string
  expires_in: number
  permissions?: UserPermission
}

// 菜单项（后端返回的完整结构）
export interface MenuItem {
  id: string
  key: string
  label: string
  icon: string
  path: string
  permission: string
  parent_id: string
  sort_order: number
  status: boolean
  children?: MenuItem[]
  created_at?: string
  updated_at?: string
}

// 菜单树展示节点（用于 Ant Design Tree 等组件）
export interface MenuTreeDisplayNode {
  key: string
  title: string
  children?: MenuTreeDisplayNode[]
}
