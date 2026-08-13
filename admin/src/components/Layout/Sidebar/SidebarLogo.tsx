import { Button } from 'antd'
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons'

interface SidebarLogoProps {
  collapsed: boolean
  onToggle: () => void
}

export default function SidebarLogo({ collapsed, onToggle }: SidebarLogoProps) {
  return (
    <div
      style={{
        padding: collapsed ? '16px 8px' : '16px 12px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: collapsed ? 'center' : 'flex-start',
        gap: 10,
        height: 56,
        flexShrink: 0,
      }}
    >
      {collapsed ? (
        <Button
          type="text"
          className="sidebar-toggle-btn"
          icon={<MenuUnfoldOutlined style={{ fontSize: 16 }} />}
          onClick={onToggle}
          style={{
            width: 40,
            height: 40,
            padding: 0,
            color: 'var(--text-muted)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderRadius: 'var(--radius-md)',
            transition: 'all 0.15s',
          }}
        />
      ) : (
        <>
          <div
            style={{
              width: 32,
              height: 32,
              minWidth: 32,
              background: 'var(--brand-dark)',
              borderRadius: 'var(--radius-sm)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: 'var(--bg-white)',
              fontWeight: 700,
              fontSize: 15,
            }}
          >
            A
          </div>
          <span
            style={{
              fontSize: 15,
              fontWeight: 600,
              color: 'var(--text-primary)',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              flex: 1,
            }}
          >
            管理系统
          </span>
          <Button
            type="text"
            className="sidebar-toggle-btn"
            icon={<MenuFoldOutlined style={{ fontSize: 16 }} />}
            onClick={onToggle}
            style={{
              width: 28,
              height: 28,
              minWidth: 28,
              padding: 0,
              color: 'var(--text-muted)',
              borderRadius: 'var(--radius-md)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              transition: 'all 0.15s',
            }}
          />
        </>
      )}
    </div>
  )
}
