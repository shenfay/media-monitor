import axios, { AxiosError, type InternalAxiosRequestConfig, type AxiosResponse, type CancelTokenSource } from 'axios'
import { message } from 'antd'
import i18n from '@/locales'
import type { ApiResponse } from '@/types'

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ---- 请求取消：路由切换时取消所有进行中的请求 ----
let cancelTokenSource: CancelTokenSource | null = null

function cancelPendingRequests() {
  if (cancelTokenSource) {
    cancelTokenSource.cancel('路由切换，取消请求')
    cancelTokenSource = null
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
    // 为每个请求附加取消令牌
    cancelTokenSource = axios.CancelToken.source()
    config.cancelToken = cancelTokenSource.token
    return config
  },
  (error: AxiosError) => {
    // 忽略取消请求的错误
    if (axios.isCancel(error)) return Promise.reject(error)
    return Promise.reject(error)
  }
)

// 响应拦截器
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
    const { response } = error
    if (response) {
      const { status, data } = response
      const msg = data?.message
      if (status === 401) {
        // 登录页的 401 是账号密码错误，不跳转，让登录页自行处理错误提示
        if (window.location.pathname === '/login') {
          return Promise.reject(error)
        }
        message.error(msg || i18n.t('sessionExpired'))
        // 统一清理：localStorage + Zustand store
        localStorage.removeItem('admin-token')
        localStorage.removeItem('admin-refresh-token')
        // 延迟导入避免循环依赖
        import('@/stores').then(({ useUserStore }) => {
          useUserStore.getState().logout()
        })
        window.location.href = '/login'
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

export default request
