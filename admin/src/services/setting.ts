/**
 * 系统设置 API
 */
import request from '@/utils/request'

/** 系统设置项 */
export interface SettingItem {
  id: number
  key: string
  value: unknown
  category: string
  label: string
  description: string
  updated_by?: number
  created_at: string
  updated_at: string
}

/** 获取设置列表 */
export async function getSettings(params?: { category?: string }): Promise<SettingItem[]> {
  return request.get('/v1/settings', { params })
}

/** 获取单个设置 */
export async function getSettingByKey(key: string): Promise<SettingItem> {
  return request.get(`/v1/settings/${key}`)
}

/** 批量更新设置项 */
export interface SettingUpdateItem {
  key: string
  value: unknown
}

/** 批量更新设置 */
export async function batchUpdateSettings(settings: SettingUpdateItem[]): Promise<{ message: string }> {
  return request.put('/v1/settings', { settings })
}
