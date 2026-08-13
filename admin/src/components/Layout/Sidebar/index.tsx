import { Layout } from 'antd'
import { useUserStore, useLayoutStore } from '@/stores'
import SidebarLogo from './SidebarLogo'
import SidebarMenu from './SidebarMenu'
import SidebarUser from './SidebarUser'

const { Sider } = Layout

export default function Sidebar() {
  const { username, menuTree, logout } = useUserStore()
  const { sidebarCollapsed, toggleSidebar } = useLayoutStore()

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={sidebarCollapsed}
      width={224}
      collapsedWidth={56}
      style={{
        background: 'var(--sidebar-bg)',
        borderRight: '1px solid var(--border-color)',
        transition: 'all var(--transition)',
      }}
    >
      <div
        style={{
          height: '100vh',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <SidebarLogo collapsed={sidebarCollapsed} onToggle={toggleSidebar} />
        <SidebarMenu menuTree={menuTree} collapsed={sidebarCollapsed} />
        <SidebarUser collapsed={sidebarCollapsed} username={username} onLogout={logout} />
      </div>
    </Sider>
  )
}
