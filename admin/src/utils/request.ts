import axios, { AxiosError, type AxiosRequestConfig, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { message } from 'antd'
import i18n from '@/locales'
import type { ApiResponse } from '@/types'
import { doRefreshToken, enqueueRequest, getIsRefreshing } from './tokenRefresh'

// 扩展 axios config 以支持静默模式（跳过全局错误提示）
declare module 'axios' {
  interface AxiosRequestConfig {
    silent?: boolean
  }
}

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ---- 请求取消：路由切换时取消所有进行中的请求 ----
let abortController: AbortController | null = null

function cancelPendingRequests() {
  if (abortController) {
    abortController.abort('路由切换，取消请求')
    abortController = null
  }
}

/** 在路由切换时调用，取消所有进行中的请求 */
export function cancelAllRequests() {
  cancelPendingRequests()
}

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 从 localStorage 获取 token
    const token = localStorage.getItem('admin-token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    // 为每个请求附加 AbortController
    abortController = new AbortController()
    config.signal = abortController.signal
    return config
  },
  (error: AxiosError) => {
    // 忽略取消请求的错误
    if (axios.isCancel(error)) return Promise.reject(error)
    return Promise.reject(error)
  }
)

// 响应拦截器

/** 统一处理登录跳转：清理状态 + 跳转登录页 */
function redirectToLogin(msg?: string) {
  message.error(msg || i18n.t('sessionExpired'))
  localStorage.removeItem('admin-token')
  localStorage.removeItem('admin-refresh-token')
  // 延迟导入避免循环依赖
  import('@/stores').then(({ useUserStore }) => {
    useUserStore.getState().logout()
  })
  window.location.href = '/login'
}

request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 检查是否为标准 ApiResponse 结构（含 code + data 字段）
    const { data } = response
    if (data && typeof data === 'object' && 'code' in data && 'data' in data) {
      return (data as ApiResponse).data
    }
    // 非标准响应，直接返回
    return data
  },
  (error: AxiosError<{ code?: string; message?: string }>) => {
    // 静默模式：跳过全局错误提示，由调用方自行处理
    if (error.config?.silent) {
      return Promise.reject(error)
    }
    const { response } = error
    if (response) {
      const { status, data } = response
      const msg = data?.message
      if (status === 401) {
        // 登录页的 401 是账号密码错误，不跳转，让登录页自行处理错误提示
        if (window.location.pathname === '/login') {
          return Promise.reject(error)
        }

        const originalConfig = error.config as AxiosRequestConfig & { _retried?: boolean }
        const refreshToken = localStorage.getItem('admin-refresh-token')

        // 有 refresh token 且未重试过 → 尝试自动刷新
        if (refreshToken && originalConfig && !originalConfig._retried) {
          originalConfig._retried = true

          if (getIsRefreshing()) {
            // 已有刷新请求进行中，排队等待
            return enqueueRequest(originalConfig)
          }

          return doRefreshToken().then(success => {
            if (success) {
              // 刷新成功，用新 token 重试原请求
              return request(originalConfig)
            }
            // 刷新失败，跳转登录页
            redirectToLogin(msg)
            return Promise.reject(error)
          })
        }

        // 无 refresh token 或已重试过 → 直接跳转登录
        redirectToLogin(msg)
      } else {
        // 直接使用后端返回的中文消息
        message.error(msg || `${i18n.t('error')} (${status})`)
      }
    } else {
      message.error(i18n.t('networkError'))
    }
    return Promise.reject(error)
  }
)

// ---- 类型化请求接口（响应拦截器已提取 data，返回类型为 T 而非 AxiosResponse） ----
interface TypedRequest {
  get<T>(url: string, config?: AxiosRequestConfig): Promise<T>
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>
  delete<T>(url: string, config?: AxiosRequestConfig): Promise<T>
  patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>
}

export default request as unknown as TypedRequest
