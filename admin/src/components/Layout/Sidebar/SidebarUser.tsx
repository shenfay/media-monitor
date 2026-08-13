import { useState, useEffect, useCallback } from 'react'
import { Avatar, Dropdown, Badge } from 'antd'
import {
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  BellOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { message } from 'antd'
import { useTranslation } from 'react-i18next'
import { getUnreadCount } from '@/services/message'
import { useWebSocketPush } from '@/hooks/useWebSocket'

interface SidebarUserProps {
  collapsed: boolean
  username: string
  onLogout: () => void
}

export default function SidebarUser({ collapsed, username, onLogout }: SidebarUserProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [unreadCount, setUnreadCount] = useState(0)

  const fetchUnread = useCallback(async () => {
    try {
      const res = await getUnreadCount()
      setUnreadCount(res.total ?? 0)
    } catch {
      // 静默失败，不影响主流程
    }
  }, [])

  // 初始加载 + WebSocket 实时推送 + 切回 tab 刷新
  useEffect(() => {
    fetchUnread()

    const onVisible = () => {
      if (document.visibilityState === 'visible') fetchUnread()
    }
    document.addEventListener('visibilitychange', onVisible)
    return () => document.removeEventListener('visibilitychange', onVisible)
  }, [fetchUnread])

  useWebSocketPush(() => {
    fetchUnread()
  })

  const handleUserMenuClick = ({ key }: { key: string }) => {
    switch (key) {
      case 'profile':
        navigate('/profile')
        break
      case 'settings':
        navigate('/settings')
        break
      case 'logout':
        onLogout()
        message.success(t('logoutSuccess'))
        navigate('/login', { replace: true })
        break
    }
  }

  const userMenuItems = [
    { key: 'profile', icon: <UserOutlined />, label: t('personalCenter') },
    { key: 'settings', icon: <SettingOutlined />, label: t('systemSettings') },
    { type: 'divider' as const },
    { key: 'logout', icon: <LogoutOutlined />, label: t('logout'), danger: true },
  ]

  return (
    <div
      style={{
        padding: collapsed ? '8px 0 12px' : '12px 12px 16px',
        borderTop: '1px solid var(--border-color)',
        flexShrink: 0,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
      }}
    >
      {collapsed ? (
        <Dropdown menu={{ items: userMenuItems, onClick: handleUserMenuClick }} placement="topLeft">
          <Avatar
            size={32}
            style={{
              background: 'linear-gradient(135deg, var(--avatar-gradient-start), var(--avatar-gradient-end))',
              fontSize: 13,
              fontWeight: 600,
              cursor: 'pointer',
            }}
          >
            {username?.charAt(0) || t('user').charAt(0)}
          </Avatar>
        </Dropdown>
      ) : (
        <>
          <Dropdown menu={{ items: userMenuItems, onClick: handleUserMenuClick }} placement="topLeft">
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                padding: 4,
                borderRadius: 'var(--radius-sm)',
                cursor: 'pointer',
                transition: 'background 0.15s',
                flex: 1,
                minWidth: 0,
                overflow: 'hidden',
              }}
              onMouseEnter={e => {
                e.currentTarget.style.background = 'var(--hover-bg)'
              }}
              onMouseLeave={e => {
                e.currentTarget.style.background = 'transparent'
              }}
            >
              <Avatar
                size={28}
                style={{
                  background: 'linear-gradient(135deg, var(--avatar-gradient-start), var(--avatar-gradient-end))',
                  fontSize: 11,
                  fontWeight: 600,
                  minWidth: 28,
                }}
              >
                {username?.charAt(0) || t('user').charAt(0)}
              </Avatar>
              <span
                style={{
                  fontSize: 13,
                  fontWeight: 500,
                  color: 'var(--text-primary)',
                  whiteSpace: 'nowrap',
                  overflow: 'hidden',
                }}
              >
                {username || t('user')}
              </span>
            </div>
          </Dropdown>
          <div
            style={{
              width: 28,
              height: 28,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              flexShrink: 0,
              cursor: 'pointer',
            }}
            onClick={() => navigate('/my-messages')}
          >
            <Badge count={unreadCount} size="small" offset={[-2, 2]}>
              <BellOutlined
                style={{
                  fontSize: 16,
                  color: 'var(--text-secondary)',
                }}
              />
            </Badge>
          </div>
        </>
      )}
    </div>
  )
}
