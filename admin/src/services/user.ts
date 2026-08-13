/**
 * 用户管理 API
 */
import request from '@/utils/request'
import type { User, UserListResponse } from '@/types'

/** 获取用户列表 */
export async function getUserList(params: {
  page?: number
  page_size?: number
  keyword?: string
  role_id?: string
  status?: string
}): Promise<UserListResponse> {
  return request.get('/v1/admin/users', { params })
}

/** 创建用户 */
export async function createUser(data: {
  email: string
  name: string
  password: string
  role_ids: string[]
}): Promise<User> {
  return request.post('/v1/admin/users', data)
}

/** 更新用户 */
export async function updateUser(
  id: string,
  data: {
    name?: string
    email?: string
    role_ids?: string[]
  }
): Promise<User> {
  return request.put(`/v1/admin/users/${id}`, data)
}

/** 切换用户状态 */
export async function toggleUserStatus(id: string, locked: boolean) {
  return request.patch(`/v1/admin/users/${id}/status`, { locked })
}
