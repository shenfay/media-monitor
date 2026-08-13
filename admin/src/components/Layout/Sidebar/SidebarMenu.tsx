import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Menu } from 'antd'
import { useNavigate, useLocation } from 'react-router-dom'
import { getIcon } from '@/config/menu'
import type { MenuItem } from '@/types'

interface SidebarMenuProps {
  menuTree: MenuItem[]
  collapsed: boolean
}

/** 菜单 key -> i18n key 映射 */
const menuKeyMap: Record<string, string> = {
  overview: 'menuOverview',
  dashboard: 'menuDashboard',
  growth: 'menuGrowth',
  family: 'menuFamily',
  goals: 'menuGoals',
  'card-engine': 'menuCardEngine',
  'card-templates': 'menuCardTemplates',
  'card-instances': 'menuCardInstances',
  companion: 'menuCompanion',
  companions: 'menuCompanions',
  acceptance: 'menuAcceptance',
  'acceptance-pending': 'menuAcceptancePending',
  'points-system': 'menuPointsSystem',
  points: 'menuPoints',
  'shop-items': 'menuShopItems',
  'exchange-orders': 'menuExchangeOrders',
  user: 'menuUser',
  'user-management': 'menuUserManagement',
  'permission-management': 'menuPermissionManagement',
  'menu-management': 'menuMenuManagement',
  profile: 'menuProfile',
  system: 'menuSystem',
  'operation-log': 'menuOperationLog',
  'design-system': 'menuDesignSystem',
  'system-settings': 'menuSystemSettings',
  message: 'menuMessage',
  ws_test: 'menuWsTest',
}

/** 将后端菜单树转换为 Ant Design Menu 项 */
function convertMenuTree(nodes: MenuItem[], showIcons: boolean, isParent: boolean, t: (key: string) => string): Array<{ key: string; label: string; icon?: React.ReactNode; children?: Array<{ key: string; label: string; icon?: React.ReactNode }> }> {
  return nodes
    .filter(node => node.status)
    .map(node => {
      const resolvedIcon = node.icon ? getIcon(node.icon) : undefined
      const hasChildren = node.children && node.children.length > 0
      const i18nKey = menuKeyMap[node.key]
      const label = i18nKey ? t(i18nKey) : node.label
      const renderItem = {
        key: node.key,
        label,
        icon: (isParent && !showIcons) ? undefined : resolvedIcon,
      }

      if (hasChildren) {
        return {
          ...renderItem,
          children: convertMenuTree(node.children!, showIcons, true, t),
        }
      }

      return renderItem
    })
}

/** 从菜单树中查找 key 对应的 path */
function findPathByKey(nodes: MenuItem[], key: string): string | null {
  for (const node of nodes) {
    if (node.key === key && node.path) return node.path
    if (node.children) {
      const found = findPathByKey(node.children, key)
      if (found) return found
    }
  }
  return null
}

/** 收集所有叶子节点 key */
function collectLeafKeys(nodes: MenuItem[]): string[] {
  const keys: string[] = []
  function walk(list: MenuItem[]) {
    list.forEach(node => {
      if (node.children && node.children.length > 0) {
        walk(node.children)
      } else {
        keys.push(node.key)
      }
    })
  }
  walk(nodes)
  return keys
}

/** 收集所有父级 key */
function collectParentKeys(nodes: MenuItem[]): string[] {
  const keys: string[] = []
  function walk(list: MenuItem[]) {
    list.forEach(node => {
      if (node.children && node.children.length > 0) {
        keys.push(node.key)
        walk(node.children)
      }
    })
  }
  walk(nodes)
  return keys
}

/** 查找选中叶子节点所在的父级 key 链路 */
function findAncestorKeys(nodes: MenuItem[], targetPath: string, parents: string[] = []): string[] | null {
  for (const node of nodes) {
    const current = node.path ? [...parents, node.key] : parents
    if (node.path === targetPath) return parents
    if (node.children) {
      const found = findAncestorKeys(node.children, targetPath, current)
      if (found) return found
    }
  }
  return null
}

export default function SidebarMenu({ menuTree, collapsed }: SidebarMenuProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const [openKeys, setOpenKeys] = useState<string[]>([])

  const allLeafKeys = menuTree.length > 0 ? collectLeafKeys(menuTree) : []

  useEffect(() => {
    // 确保当前选中项的父级展开，同时保留用户手动展开的其他菜单
    const ancestors = findAncestorKeys(menuTree, location.pathname)
    if (ancestors && ancestors.length > 0) {
      setOpenKeys(prev => {
        const merged = new Set([...prev, ...ancestors])
        return Array.from(merged)
      })
    }
  }, [menuTree, location.pathname])

  const filteredMenu = menuTree.length > 0
    ? convertMenuTree(menuTree, collapsed, false, t)
    : []

  const handleMenuClick = ({ key }: { key: string }) => {
    const path = findPathByKey(menuTree, key)
    if (path) navigate(path)
  }

  const selectedKey = location.pathname === '/dashboard' ? 'dashboard' : allLeafKeys.find(key => {
    return findPathByKey(menuTree, key) === location.pathname
  }) || ''

  return (
    <div style={{ flex: 1, overflow: 'auto', padding: '0 6px' }}>
      <Menu
        mode="inline"
        inlineCollapsed={collapsed}
        selectedKeys={selectedKey ? [selectedKey] : []}
        openKeys={collapsed ? [] : openKeys}
        onOpenChange={setOpenKeys}
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        items={filteredMenu as any}
        onClick={handleMenuClick}
        style={{
          background: 'transparent',
          borderRight: 'none',
        }}
      />
    </div>
  )
}
