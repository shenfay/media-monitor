/**
 * 消息管理 API
 */
import request from '@/utils/request'

export interface MessageRecord {
  id: string
  sender_id: string | null
  recipient_id: string
  type: 'system' | 'companion'
  category: string
  title: string
  content: string
  is_read: boolean
  read_at: string | null
  ref_type: string | null
  ref_id: string | null
  metadata: Record<string, unknown>
  created_at: string
}

export interface MessageListResponse {
  messages: MessageRecord[]
  total: number
}

export interface UnreadCountItem {
  type: string
  count: number
}

export interface UnreadCountResponse {
  counts: UnreadCountItem[]
  total: number
}

/** 获取当前用户未读消息数 */
export async function getUnreadCount(): Promise<UnreadCountResponse> {
  return request.get('/v1/messages/unread-count')
}

/** 获取当前用户消息列表 */
export async function getMyMessages(params: {
  type?: string
  category?: string
  is_read?: boolean
  limit?: number
  offset?: number
}): Promise<MessageListResponse> {
  return request.get('/v1/messages', { params })
}

/** 标记消息已读 */
export async function markAsRead(id: string): Promise<void> {
  return request.post(`/v1/messages/${id}/read`)
}

/** 标记全部已读 */
export async function markAllAsRead(type?: string): Promise<void> {
  return request.post('/v1/messages/read-all', { type })
}

/** 管理员获取所有消息列表 */
export async function getMessages(params: {
  type?: string
  category?: string
  limit?: number
  offset?: number
}): Promise<MessageListResponse> {
  return request.get('/v1/admin/messages', { params })
}
