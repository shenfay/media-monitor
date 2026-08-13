/**
 * WebSocket 实时推送 Hook
 *
 * 单例连接管理，自动重连，消息分发
 */
import { useEffect, useRef, useState } from 'react'
import { notification } from 'antd'
import { BellOutlined } from '@ant-design/icons'

// ---- WebSocket 连接管理器（单例） ----

type MessageHandler = (data: WSPushMessage) => void

export interface WSPushMessage {
  type: string
  category: string
  title: string
  content: string
}

class WebSocketManager {
  private ws: WebSocket | null = null
  private handlers: Set<MessageHandler> = new Set()
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 10
  private _connected = false
  private listeners: Set<(connected: boolean) => void> = new Set()

  get connected() {
    return this._connected
  }

  connect() {
    if (this.ws) return

    const token = localStorage.getItem('admin-token')
    if (!token) return

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.hostname
    const port = window.location.port || (protocol === 'wss:' ? '443' : '80')
    const url = `${protocol}//${host}:${port}/api/ws?token=${token}`

    try {
      this.ws = new WebSocket(url)

      this.ws.onopen = () => {
        this._connected = true
        this.reconnectAttempts = 0
        this.notifyConnectionChange()
        console.log('[WS] Connected')
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as WSPushMessage
          this.handlers.forEach((handler) => handler(data))
        } catch {
          console.warn('[WS] Failed to parse message:', event.data)
        }
      }

      this.ws.onclose = () => {
        this._connected = false
        this.ws = null
        this.notifyConnectionChange()
        console.log('[WS] Disconnected')
        this.scheduleReconnect()
      }

      this.ws.onerror = () => {
        // onclose will be called after onerror
      }
    } catch {
      this.scheduleReconnect()
    }
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    this.reconnectAttempts = this.maxReconnectAttempts // prevent reconnect
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this._connected = false
    this.notifyConnectionChange()
  }

  subscribe(handler: MessageHandler) {
    this.handlers.add(handler)
    return () => {
      this.handlers.delete(handler)
    }
  }

  onConnectionChange(listener: (connected: boolean) => void) {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) return
    if (this.reconnectTimer) return

    const delay = Math.min(1000 * 2 ** this.reconnectAttempts, 30000)
    this.reconnectAttempts++

    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      this.connect()
    }, delay)
  }

  private notifyConnectionChange() {
    this.listeners.forEach((listener) => listener(this._connected))
  }
}

// 全局单例
const wsManager = new WebSocketManager()

// ---- React Hooks ----

/**
 * 初始化 WebSocket 连接（在 Layout 层调用一次）
 */
export function useWebSocketInit() {
  useEffect(() => {
    wsManager.connect()
    return () => {
      wsManager.disconnect()
    }
  }, [])
}

/**
 * 订阅 WebSocket 推送消息
 */
export function useWebSocketPush(onMessage: MessageHandler) {
  const handlerRef = useRef(onMessage)
  handlerRef.current = onMessage

  useEffect(() => {
    return wsManager.subscribe((data) => handlerRef.current(data))
  }, [])
}

/**
 * WebSocket 连接状态
 */
export function useWebSocketStatus() {
  const [connected, setConnected] = useState(wsManager.connected)

  useEffect(() => {
    return wsManager.onConnectionChange(setConnected)
  }, [])

  return connected
}

/**
 * 消息推送通知 Hook（自动弹出 Toast + 更新未读数）
 * 只对真实消息（system/companion）弹通知，跳过内部事件
 */
export function usePushNotification(onNewMessage?: () => void) {
  const [api] = notification.useNotification()

  useWebSocketPush((data) => {
    if (data.type === 'unread_update') return
    api.open({
      message: data.title || '新消息',
      description: data.content,
      icon: <BellOutlined style={{ color: 'var(--primary-color)' }} />,
      duration: 4,
    })

    onNewMessage?.()
  })
}
