import { lazy, Suspense } from 'react'
import { createBrowserRouter, Navigate, Outlet } from 'react-router-dom'
import { Spin } from 'antd'
import MainLayout from '@/components/Layout'
import PermissionGuard from '@/components/PermissionGuard'

// 登录页同步加载（首屏关键路径）
import Login from '@/pages/Login'

// 业务页面懒加载
const Dashboard = lazy(() => import('@/pages/Dashboard'))
const UserManagement = lazy(() => import('@/pages/UserManagement'))
const PermissionManagement = lazy(() => import('@/pages/PermissionManagement'))
const MenuManagement = lazy(() => import('@/pages/MenuManagement'))
const Profile = lazy(() => import('@/pages/Profile'))
const OperationLog = lazy(() => import('@/pages/OperationLog'))
const MyMessages = lazy(() => import('@/pages/MyMessages'))
const MessageManagement = lazy(() => import('@/pages/MessageManagement'))
const SystemSettings = lazy(() => import('@/pages/SystemSettings'))
const DesignSystem = lazy(() => import('@/pages/DesignSystem'))
const WebSocketTest = lazy(() => import('@/pages/WebSocketTest'))
const WorkerManagement = lazy(() => import('@/pages/WorkerManagement'))
const SourceManagement = lazy(() => import('@/pages/SourceManagement'))
const TaskManagement = lazy(() => import('@/pages/TaskManagement'))
const ArticleManagement = lazy(() => import('@/pages/ArticleManagement'))

/** 懒加载 fallback 加载指示器 */
const PageLoading = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', minHeight: 200 }}>
    <Spin size="large" />
  </div>
)

const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/',
    element: (
      <MainLayout>
        <Suspense fallback={<PageLoading />}>
          <Outlet />
        </Suspense>
      </MainLayout>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      {
        path: 'dashboard',
        element: <PermissionGuard permission="dashboard:view"><Dashboard /></PermissionGuard>,
      },
      {
        path: 'users',
        element: <PermissionGuard permission="user:manage"><UserManagement /></PermissionGuard>,
      },
      {
        path: 'permissions',
        element: <PermissionGuard permission="permission:manage"><PermissionManagement /></PermissionGuard>,
      },
      {
        path: 'menus',
        element: <PermissionGuard permission="menu:manage"><MenuManagement /></PermissionGuard>,
      },
      {
        path: 'profile',
        element: <PermissionGuard permission="profile:view"><Profile /></PermissionGuard>,
      },
      {
        path: 'operation-log',
        element: <PermissionGuard permission="operation:log"><OperationLog /></PermissionGuard>,
      },
      {
        path: 'my-messages',
        element: <MyMessages />,
      },
      {
        path: 'messages',
        element: <PermissionGuard permission="message:view"><MessageManagement /></PermissionGuard>,
      },
      {
        path: 'design-system',
        element: <PermissionGuard permission="design:view"><DesignSystem /></PermissionGuard>,
      },
      {
        path: 'ws-test',
        element: <WebSocketTest />,
      },
      {
        path: 'settings',
        element: <PermissionGuard permission="setting:manage"><SystemSettings /></PermissionGuard>,
      },
      {
        path: 'crawl/sources',
        element: <PermissionGuard permission="source:view"><SourceManagement /></PermissionGuard>,
      },
      {
        path: 'crawl/tasks',
        element: <PermissionGuard permission="task:view"><TaskManagement /></PermissionGuard>,
      },
      {
        path: 'crawl/articles',
        element: <PermissionGuard permission="article:view"><ArticleManagement /></PermissionGuard>,
      },
      {
        path: 'crawl/workers',
        element: <PermissionGuard permission="worker:view"><WorkerManagement /></PermissionGuard>,
      },
    ],
  },
  { path: '*', element: <Navigate to="/dashboard" replace /> },
])

export default router
