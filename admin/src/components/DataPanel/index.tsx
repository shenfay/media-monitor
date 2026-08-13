import { Input } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import type { ReactNode } from 'react'

interface DataPanelProps {
  /** 页面主标题 */
  title?: string
  /** 标题右侧操作按钮 */
  extra?: ReactNode
  /** 筛选栏内容（搜索框、筛选按钮等） */
  filters?: ReactNode
  /** 筛选栏右侧工具图标 */
  toolbarActions?: ReactNode
  /** 表格/内容区域 */
  children: ReactNode
  /** 自定义样式 */
  style?: React.CSSProperties
  /** 是否使用紧凑卡片模式（无表格容器包装） */
  compact?: boolean
}

/** 筛选搜索框 */
export function FilterSearch({ value, onChange, placeholder = '搜索...', onSearch }: {
  value?: string
  onChange?: (v: string) => void
  placeholder?: string
  onSearch?: () => void
}) {
  return (
    <div style={{ flex: 1, maxWidth: 320 }}>
      <Input
        prefix={<SearchOutlined style={{ color: 'var(--text-icon)', fontSize: 14 }} />}
        value={value}
        onChange={e => onChange?.(e.target.value)}
        onPressEnter={onSearch}
        placeholder={placeholder}
      />
    </div>
  )
}

export default function DataPanel({
  title,
  extra,
  filters,
  toolbarActions,
  children,
  style,
  compact,
}: DataPanelProps) {
  const hasHeader = title || extra
  const hasFilters = filters || toolbarActions

  return (
    <div style={style}>
      {/* Page Header */}
      {hasHeader && (
        <div
          style={{
            padding: '20px 28px 0',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
          }}
        >
          <div>
            {title && (
              <h2
                style={{
                  fontSize: 20,
                  fontWeight: 600,
                  color: 'var(--text-primary)',
                  lineHeight: 1.3,
                  margin: 0,
                }}
              >
                {title}
              </h2>
            )}
          </div>
          {extra && (
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, flexShrink: 0 }}>
              {extra}
            </div>
          )}
        </div>
      )}

      {/* 标题与下个区域间距 */}
      {hasHeader && <div style={{ height: 16 }} />}

      {/* Filter Bar */}
      {hasFilters && (
        <div
          style={{
            padding: '0 28px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            gap: 10,
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 10, flex: 1 }}>
            {filters}
          </div>
          {toolbarActions && (
            <div style={{ display: 'flex', alignItems: 'center', gap: 6, flexShrink: 0 }}>
              {toolbarActions}
            </div>
          )}
        </div>
      )}

      {/* Filter Bar 与 Table 间距 */}
      {hasFilters && <div style={{ height: 16 }} />}

      {/* Table / Content Area */}
      {compact ? (
        <div style={{ padding: '0 28px 20px' }}>
          {children}
        </div>
      ) : (
        <div style={{ padding: '0 28px 20px' }}>
          <div
            style={{
              border: '1px solid var(--border-light)',
              borderRadius: 'var(--radius-md)',
              background: 'var(--bg-white)',
              overflow: 'hidden',
            }}
          >
            {children}
          </div>
        </div>
      )}
    </div>
  )
}
