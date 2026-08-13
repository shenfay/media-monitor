/**
 * 认证相关 API
 */
import request from '@/utils/request'
import type { LoginRequest, LoginResponse, UserPermission, MenuItem } from '@/types'

/** 登录 */
export async function login(data: LoginRequest): Promise<LoginResponse> {
  return request.post('/v1/auth/login', data)
}

/** 获取当前用户信息 */
export async function getCurrentUser() {
  return request.get('/v1/auth/me')
}

/** 获取当前用户权限 */
export async function getPermissions(): Promise<UserPermission> {
  return request.get('/v1/auth/permissions')
}

/** 获取当前用户菜单树 */
export async function getUserMenuTree(): Promise<MenuItem[]> {
  return request.get('/v1/auth/menus')
}

/** 登出 */
export async function logout() {
  return request.post('/v1/auth/logout')
}
