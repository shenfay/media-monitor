import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useUserStore } from '@/stores'
import { hasPermission } from '@/config/permission'

interface PermissionGuardProps {
  permission?: string
  children: ReactNode
  fallback?: ReactNode
}

export default function PermissionGuard({
  permission,
  children,
  fallback,
}: PermissionGuardProps) {
  const { t } = useTranslation()
  const { isLogin, permissions } = useUserStore()

  // 未登录，跳转到登录页
  if (!isLogin) {
    return <Navigate to="/login" replace />
  }

  // 检查权限
  if (permission && !hasPermission(permissions, permission)) {
    if (fallback) {
      return <>{fallback}</>
    }
    return (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          height: '100%',
          color: 'var(--text-muted)',
          fontSize: 16,
        }}
      >
        {t('noPermission')}
      </div>
    )
  }

  return <>{children}</>
}
