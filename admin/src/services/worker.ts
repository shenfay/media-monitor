/**
 * Worker 管理 API
 */
import request from '@/utils/request'

export interface WorkerInfo {
  id: string
  name: string
  adapters: string[]
  capabilities: string[]
  status: string
  concurrency: number
  current_task: string
  processed_count: number
  last_heartbeat: string
  started_at: string
}

export interface AdapterMeta {
  name: string
  required_tags: string[]
  platform_type: string
  class_name: string
  first_seen: string
}

/** 获取所有 Worker 列表 */
export async function getWorkers(): Promise<WorkerInfo[]> {
  const res = await request.get('/v1/admin/crawl/workers')
  return (res as unknown as WorkerInfo[]) || []
}

/** 获取所有适配器列表 */
export async function getAdapters(): Promise<AdapterMeta[]> {
  const res = await request.get('/v1/admin/crawl/workers/adapters')
  return (res as unknown as AdapterMeta[]) || []
}

/** 暂停 Worker */
export async function pauseWorker(id: string): Promise<void> {
  return request.post(`/v1/admin/crawl/workers/${id}/pause`)
}

/** 恢复 Worker */
export async function resumeWorker(id: string): Promise<void> {
  return request.post(`/v1/admin/crawl/workers/${id}/resume`)
}

/** 下线 Worker（优雅停止） */
export async function shutdownWorker(id: string): Promise<void> {
  return request.post(`/v1/admin/crawl/workers/${id}/shutdown`)
}
