/**
 * 抓取模块 API
 */
import request from '@/utils/request'

// ---- Types ----

export interface Source {
  id: string
  name: string
  platform_type: string
  base_url: string
  list_endpoint: string
  nodes: string[]
  source_filter: string
  months: number
  schedule: string
  tags: string[]
  enabled: boolean
  last_crawl_at: string | null
  created_at: string
  updated_at: string
}

export interface TaskRun {
  id: string
  source_id: string
  task_type: string
  triggered_by: string
  status: string
  total: number
  ingested: number
  failed: number
  error: string
  started_at: string | null
  finished_at: string | null
  created_at: string
  updated_at: string
}

export interface Article {
  id: string
  source_id: string
  platform: string
  external_id: string
  url: string
  title: string
  summary: string
  body: string
  body_format: string
  author: string
  source_name: string
  language: string
  published_at: string | null
  media: unknown
  interactions: unknown
  fetched_at: string
  created_at: string
}

export interface DashboardStats {
  total_sources: number
  enabled_sources: number
  total_articles: number
  today_articles: number
  total_tasks: number
  running_tasks: number
  online_workers: number
}

// ---- Source API ----

export async function getSources(silent = false): Promise<Source[]> {
  const res = await request.get('/v1/admin/crawl/sources', { silent })
  return res || []
}

export async function createSource(data: Partial<Source> & { auth?: unknown }): Promise<Source> {
  const res = await request.post('/v1/admin/crawl/sources', data)
  return res.data
}

export async function updateSource(id: string, data: Partial<Source> & { auth?: unknown }): Promise<Source> {
  const res = await request.put(`/v1/admin/crawl/sources/${id}`, data)
  return res.data
}

export async function deleteSource(id: string): Promise<void> {
  return request.delete(`/v1/admin/crawl/sources/${id}`)
}

export async function runSource(id: string): Promise<{ task_id: string }> {
  const res = await request.post(`/v1/admin/crawl/sources/${id}/run`)
  return res.data
}

// ---- Task API ----

export async function getTasks(params?: { source_id?: string }, silent = false): Promise<TaskRun[]> {
  const res = await request.get('/v1/admin/crawl/tasks', { params, silent })
  return res || []
}

export async function getTask(id: string): Promise<TaskRun> {
  const res = await request.get(`/v1/admin/crawl/tasks/${id}`)
  return res.data
}

export async function createTask(data: { source_id: string; limit?: number; with_body?: boolean; since?: string; mode?: string }): Promise<{ task_id: string }> {
  const res = await request.post('/v1/admin/crawl/tasks', data)
  return res.data
}

// ---- Article API ----

export async function getArticles(params?: {
  source_id?: string
  platform?: string
  language?: string
  keyword?: string
  limit?: number
  offset?: number
}): Promise<{ articles: Article[]; total: number }> {
  const res = await request.get('/v1/admin/crawl/articles', { params })
  return res.data || { articles: [], total: 0 }
}

export async function getArticle(id: string): Promise<Article> {
  const res = await request.get(`/v1/admin/crawl/articles/${id}`)
  return res.data
}

// ---- Dashboard API ----

export async function getDashboardStats(silent = false): Promise<DashboardStats | null> {
  try {
    const res = await request.get('/v1/admin/crawl/dashboard/stats', { silent })
    return res || null
  } catch {
    return null
  }
}
