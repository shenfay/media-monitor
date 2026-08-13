import { Tag } from 'antd'
import type { ReactNode } from 'react'

export type StatusType =
  | 'success'
  | 'processing'
  | 'pending'
  | 'failed'
  | 'normal'
  | 'warning'
  | 'error'
  | 'active'
  | 'inactive'
  | 'disabled'
  | 'delayed'
  | string

interface StatusTagProps {
  status: StatusType
  children?: ReactNode
}

const statusColorMap: Record<string, { color: string; bg: string }> = {
  success: { color: 'var(--green-text)', bg: 'var(--green-light)' },
  normal: { color: 'var(--green-text)', bg: 'var(--green-light)' },
  active: { color: 'var(--green-text)', bg: 'var(--green-light)' },
  processing: { color: 'var(--blue)', bg: 'var(--blue-light)' },
  pending: { color: 'var(--yellow-text)', bg: 'var(--yellow-light)' },
  warning: { color: 'var(--yellow-text)', bg: 'var(--yellow-light)' },
  delayed: { color: 'var(--yellow-text)', bg: 'var(--yellow-light)' },
  inactive: { color: 'var(--yellow-text)', bg: 'var(--yellow-light)' },
  failed: { color: 'var(--red)', bg: 'var(--red-light)' },
  error: { color: 'var(--red)', bg: 'var(--red-light)' },
  disabled: { color: 'var(--red)', bg: 'var(--red-light)' },
  admin: { color: 'var(--blue)', bg: 'var(--blue-light)' },
  editor: { color: 'var(--text-secondary)', bg: 'var(--hover-bg-light)' },
  member: { color: 'var(--text-secondary)', bg: 'var(--hover-bg-light)' },
  mysql: { color: 'var(--blue)', bg: 'var(--blue-light)' },
  elasticsearch: { color: 'var(--text-secondary)', bg: 'var(--hover-bg-light)' },
  s3: { color: 'var(--text-secondary)', bg: 'var(--hover-bg-light)' },
  kafka: { color: 'var(--text-secondary)', bg: 'var(--hover-bg-light)' },
}

export default function StatusTag({ status, children }: StatusTagProps) {
  const style = statusColorMap[status] || { color: 'var(--text-secondary)', bg: 'var(--hover-bg-light)' }

  return (
    <Tag
      style={{
        color: style.color,
        background: style.bg,
        border: 'none',
        borderRadius: 'var(--radius-sm)',
        fontSize: 12,
        fontWeight: 500,
        padding: '2px 10px',
        lineHeight: '20px',
      }}
    >
      {children || status}
    </Tag>
  )
}
