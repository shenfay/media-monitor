import { Button } from 'antd'
import type { ButtonProps } from 'antd'
import { useUserStore } from '@/stores'
import { hasPermission } from '@/config/permission'

interface PermissionButtonProps extends ButtonProps {
  permission?: string
  fallback?: React.ReactNode
}

export default function PermissionButton({
  permission,
  fallback = null,
  children,
  ...props
}: PermissionButtonProps) {
  const { permissions } = useUserStore()

  if (permission && !hasPermission(permissions, permission)) {
    return <>{fallback}</>
  }

  return <Button {...props}>{children}</Button>
}
