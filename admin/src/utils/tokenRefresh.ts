/**
 * Token 自动刷新管理
 *
 * - 被动刷新：401 响应时自动尝试用 refresh_token 续期，成功后重试原请求
 * - 主动刷新：根据 JWT 过期时间设置定时器，提前 60s 刷新
 * - 并发控制：多个请求同时 401 只触发一次刷新，其余排队等待
 */
import axios from 'axios'
import type { AxiosRequestConfig } from 'axios'
import type { LoginResponse } from '@/types'

// ---- 常量 ----
const REFRESH_API = (import.meta.env.VITE_API_BASE_URL || '/api') + '/v1/auth/refresh'
const TOKEN_KEY = 'admin-token'
const REFRESH_TOKEN_KEY = 'admin-refresh-token'
/** 提前刷新的时间窗口（毫秒）：token 过期前 60s 触发刷新 */
const REFRESH_AHEAD_MS = 60_000
/** 最小刷新间隔（毫秒）：防止频繁刷新 */
const MIN_REFRESH_INTERVAL_MS = 10_000

// ---- 刷新状态 ----
let isRefreshing = false
let pendingQueue: Array<{
  resolve: (value: unknown) => void
  reject: (reason?: unknown) => void
  config: AxiosRequestConfig
}> = []

// ---- 主动刷新定时器 ----
let proactiveTimer: ReturnType<typeof setTimeout> | null = null

// 专用 axios 实例（不挂载任何拦截器，避免刷新请求触发 401 重试死循环）
const refreshAxios = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

/**
 * 执行 token 刷新
 *
 * @returns 是否刷新成功
 */
export async function doRefreshToken(): Promise<boolean> {
  if (isRefreshing) return false

  const refreshTokenValue = localStorage.getItem(REFRESH_TOKEN_KEY)
  if (!refreshTokenValue) return false

  isRefreshing = true

  try {
    const res = await refreshAxios.post(REFRESH_API, {
      refresh_token: refreshTokenValue,
    })

    // 后端响应格式: { code, data: { access_token, refresh_token, ... } }
    const tokenData: LoginResponse = res.data?.data ?? res.data
    const newAccessToken = tokenData?.access_token
    const newRefreshToken = tokenData?.refresh_token

    if (!newAccessToken || !newRefreshToken) {
      throw new Error('Invalid refresh response')
    }

    // 更新 localStorage
    localStorage.setItem(TOKEN_KEY, newAccessToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken)

    // 刷新成功 → 重试所有排队请求
    const queue = pendingQueue
    pendingQueue = []
    isRefreshing = false

    // 重新调度主动刷新
    scheduleProactiveRefresh(newAccessToken)

    queue.forEach(item => item.resolve(true))
    return true
  } catch {
    // 刷新失败 → 清理认证状态，排队请求全部拒绝
    isRefreshing = false
    pendingQueue = []
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    return false
  }
}

/**
 * 将请求加入等待队列（刷新完成后自动重试）
 */
export function enqueueRequest(config: AxiosRequestConfig): Promise<unknown> {
  return new Promise((resolve, reject) => {
    pendingQueue.push({ resolve, reject, config })
  })
}

/**
 * 当前是否正在刷新
 */
export function getIsRefreshing(): boolean {
  return isRefreshing
}

// ---- 主动刷新调度 ----

/**
 * 根据 JWT 过期时间，提前 60s 调度一次 token 刷新
 */
export function scheduleProactiveRefresh(token: string): void {
  clearProactiveRefresh()

  try {
    const payload = decodeJWTPayload(token)
    if (!payload?.exp) return

    const exp = payload.exp as number
    const expiresAt = exp * 1000
    const refreshAt = expiresAt - REFRESH_AHEAD_MS
    const delay = refreshAt - Date.now()

    if (delay > MIN_REFRESH_INTERVAL_MS) {
      proactiveTimer = setTimeout(() => {
        doRefreshToken().catch(() => {
          // 静默处理，下次请求触发被动刷新
        })
      }, delay)
    }
  } catch {
    // 解码失败，忽略
  }
}

/**
 * 清除主动刷新定时器
 */
export function clearProactiveRefresh(): void {
  if (proactiveTimer) {
    clearTimeout(proactiveTimer)
    proactiveTimer = null
  }
}

// ---- JWT 解码工具 ----

function decodeJWTPayload(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const payload = parts[1].replace(/-/g, '+').replace(/_/g, '/')
    const padded = payload + '='.repeat((4 - (payload.length % 4)) % 4)
    return JSON.parse(atob(padded))
  } catch {
    return null
  }
}

/**
 * 获取 JWT 过期时间戳（毫秒），用于外部判断
 */
export function getTokenExpMs(token: string): number | null {
  const payload = decodeJWTPayload(token)
  if (!payload?.exp || typeof payload.exp !== 'number') return null
  return (payload.exp as number) * 1000
}
