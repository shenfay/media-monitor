/**
 * 操作日志 API
 */
import request from '@/utils/request'

export interface OperationLogRecord {
  id: string
  user_id: string
  email: string
  action: string
  category: string
  status: string
  ip: string
  user_agent: string
  device: string
  browser: string
  os: string
  metadata: Record<string, unknown>
  created_at: string
}

export interface OperationLogListResponse {
  data: OperationLogRecord[]
  limit: number
  offset: number
}

/** 获取操作日志列表 */
export async function getOperationLogs(params: {
  category?: string
  action?: string
  limit?: number
  offset?: number
}): Promise<OperationLogListResponse> {
  return request.get('/v1/operation-logs', { params })
}

/** 获取用户操作日志 */
export async function getUserOperationLogs(
  userId: string,
  params?: { limit?: number; offset?: number }
): Promise<OperationLogListResponse> {
  return request.get(`/v1/operation-logs/user/${userId}`, { params })
}

/** 按分类获取操作日志 */
export async function getCategoryOperationLogs(
  category: string,
  params?: { limit?: number; offset?: number }
): Promise<OperationLogListResponse> {
  return request.get(`/v1/operation-logs/category/${category}`, { params })
}
