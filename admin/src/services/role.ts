/**
 * 角色权限管理 API
 */
import request from '@/utils/request'
import type { Role } from '@/types'

/** 获取角色列表 */
export async function getRoleList(): Promise<Role[]> {
  return request.get('/v1/admin/roles')
}

/** 创建角色 */
export async function createRole(data: {
  name: string
  code: string
  description?: string
}): Promise<Role> {
  return request.post('/v1/admin/roles', data)
}

/** 更新角色 */
export async function updateRole(
  id: string,
  data: { name: string; description?: string }
): Promise<Role> {
  return request.put(`/v1/admin/roles/${id}`, data)
}

/** 删除角色 */
export async function deleteRole(id: string) {
  return request.delete(`/v1/admin/roles/${id}`)
}

/** 切换角色状态 */
export async function toggleRoleStatus(id: string) {
  return request.patch(`/v1/admin/roles/${id}/status`)
}

/** 获取角色权限（返回权限标识字符串数组，如 ["dashboard:view", "family:manage"]） */
export async function getRolePermissions(roleId: string): Promise<string[]> {
  return request.get(`/v1/admin/roles/${roleId}/permissions`)
}

/** 更新角色权限 */
export async function updateRolePermissions(
  roleId: string,
  permissions: string[]
) {
  return request.put(`/v1/admin/roles/${roleId}/permissions`, { permissions })
}
