import { useState, useCallback, useEffect } from 'react'
import { Layout } from 'antd'
import { useNavigate, useLocation } from 'react-router-dom'
import Sidebar from './Sidebar'
import TopBar from './TopBar'
import PageContainer from './PageContainer'
import { useUserStore } from '@/stores'
import { getUserMenuTree } from '@/services/auth'
import { cancelAllRequests } from '@/utils/request'
import { useWebSocketInit } from '@/hooks/useWebSocket'
import type { ReactNode } from 'react'

const { Content } = Layout

interface MainLayoutProps {
  children: ReactNode
}

export default function MainLayout({ children }: MainLayoutProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { isLogin, setMenuTree } = useUserStore()
  const [contentKey, setContentKey] = useState(0)

  // 初始化 WebSocket 连接（实时推送）
  useWebSocketInit()

  // 路由切换时取消进行中的 API 请求
  useEffect(() => {
    return () => { cancelAllRequests() }
  }, [location.pathname])

  const handleRefresh = useCallback(() => {
    setContentKey(k => k + 1)
  }, [])

  // 未登录时跳转到登录页
  useEffect(() => {
    if (!isLogin) {
      navigate('/login', { replace: true })
    }
  }, [isLogin, navigate])

  // 每次页面加载时都获取最新菜单树，确保权限变更后菜单实时更新
  useEffect(() => {
    if (isLogin) {
      getUserMenuTree()
        .then(tree => {
          setMenuTree(tree || [])
        })
        .catch(() => {
          // 获取失败静默处理
        })
    }
  }, [isLogin, setMenuTree])

  if (!isLogin) {
    return null
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sidebar />
      <Layout style={{ height: '100vh' }}>
        <TopBar onRefresh={handleRefresh} />
        <Content style={{ overflow: 'hidden', height: 'calc(100vh - 50px)', display: 'flex', flexDirection: 'column' }}>
          <PageContainer key={contentKey}>{children}</PageContainer>
        </Content>
      </Layout>
    </Layout>
  )
}
