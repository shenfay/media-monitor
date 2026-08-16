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
  auth: Record<string, string> | null
  enabled: boolean
  article_count: number
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
  channel: string
  language: string
  status: string
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
  const res = await request.get<Source[]>('/v1/admin/crawl/sources', { silent })
  return res || []
}

export async function createSource(data: Partial<Source> & { auth?: unknown }): Promise<Source> {
  return await request.post<Source>('/v1/admin/crawl/sources', data)
}

export async function updateSource(id: string, data: Partial<Source> & { auth?: unknown }): Promise<Source> {
  return await request.put<Source>(`/v1/admin/crawl/sources/${id}`, data)
}

export async function deleteSource(id: string): Promise<void> {
  return request.delete(`/v1/admin/crawl/sources/${id}`)
}

export async function runSource(id: string): Promise<{ task_id: string }> {
  return await request.post<{ task_id: string }>(`/v1/admin/crawl/sources/${id}/run`)
}

// ---- Task API ----

export async function getTasks(params?: { source_id?: string }, silent = false): Promise<TaskRun[]> {
  const res = await request.get<TaskRun[]>('/v1/admin/crawl/tasks', { params, silent })
  return res || []
}

export async function getTask(id: string): Promise<TaskRun> {
  return await request.get<TaskRun>(`/v1/admin/crawl/tasks/${id}`)
}

export async function createTask(data: { source_id: string; limit?: number; with_body?: boolean; since?: string; mode?: string }): Promise<{ task_id: string }> {
  return await request.post<{ task_id: string }>('/v1/admin/crawl/tasks', data)
}

// ---- Article API ----

export async function getArticles(params?: {
  source_id?: string
  platform?: string
  language?: string
  keyword?: string
  status?: string
  limit?: number
  offset?: number
}): Promise<{ articles: Article[]; total: number }> {
  const res = await request.get<{ articles: Article[]; total: number }>('/v1/admin/crawl/articles', { params })
  return res || { articles: [], total: 0 }
}

export async function getArticle(id: string): Promise<Article> {
  return await request.get<Article>(`/v1/admin/crawl/articles/${id}`)
}

// ---- Dashboard API ----

export async function getDashboardStats(silent = false): Promise<DashboardStats | null> {
  try {
    const res = await request.get<DashboardStats>('/v1/admin/crawl/dashboard/stats', { silent })
    return res || null
  } catch {
    return null
  }
}
