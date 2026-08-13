import { useState, useMemo } from 'react'
import { Breadcrumb, AutoComplete, Input, Dropdown } from 'antd'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons'
import 'flag-icons/css/flag-icons.min.css'
import { useTranslation } from 'react-i18next'
import { useLocation, useNavigate } from 'react-router-dom'
import { useUserStore } from '@/stores'
import type { MenuItem } from '@/types'

const countryCodeMap: Record<string, string> = {
  'zh-CN': 'cn',
  'en-US': 'us',
}

interface TopBarProps {
  onRefresh: () => void
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
  'my-messages': 'menuMessage',
  'design-system': 'menuDesignSystem',
  'system-settings': 'menuSystemSettings',
}

function getMenuLabel(node: MenuItem, t: (key: string) => string): string {
  const i18nKey = menuKeyMap[node.key]
  return i18nKey ? t(i18nKey) : node.label
}

export default function TopBar({ onRefresh }: TopBarProps) {
  const { t, i18n } = useTranslation()
  const location = useLocation()
  const navigate = useNavigate()
  const menuTree = useUserStore(state => state.menuTree)
  const [searchValue, setSearchValue] = useState('')

  const findBreadcrumb = () => {
    const result: { title: string }[] = [{ title: t('home') }]
    function search(nodes: MenuItem[], parentLabel?: string): boolean {
      for (const node of nodes) {
        const nodeLabel = getMenuLabel(node, t)
        if (node.path === location.pathname) {
          if (parentLabel) result.push({ title: parentLabel })
          result.push({ title: nodeLabel })
          return true
        }
        if (node.children && search(node.children, nodeLabel)) {
          return true
        }
      }
      return false
    }
    if (!search(menuTree)) {
      // 非菜单路径回退：尝试从 menuKeyMap 查找
      const key = location.pathname.replace(/^\/|\/$/g, '')
      const i18nKey = menuKeyMap[key]
      if (i18nKey) {
        result.push({ title: t(i18nKey) })
      }
    }
    return result
  }

  const breadcrumbItems = findBreadcrumb()

  // 将菜单树扁平化为搜索选项
  const searchOptions = useMemo(() => {
    const items: { value: string; label: string; path: string }[] = []
    function flatten(nodes: MenuItem[], parentLabel?: string) {
      for (const node of nodes) {
        const nodeLabel = getMenuLabel(node, t)
        const fullLabel = parentLabel ? `${parentLabel} / ${nodeLabel}` : nodeLabel
        if (node.path) {
          items.push({ value: node.path, label: fullLabel, path: node.path })
        }
        if (node.children) flatten(node.children, fullLabel)
      }
    }
    flatten(menuTree)
    return items
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [menuTree, i18n.language])

  const filteredOptions = useMemo(() => {
    if (!searchValue) return []
    const keyword = searchValue.toLowerCase()
    return searchOptions
      .filter(item => item.label.toLowerCase().includes(keyword) || item.path.toLowerCase().includes(keyword))
      .map(item => ({ value: item.path, label: item.label }))
  }, [searchValue, searchOptions])

  return (
    <div
      style={{
        height: 50,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '16px 28px 0',
        flexShrink: 0,
        background: 'var(--main-bg)',
      }}
    >
      {/* Left: Breadcrumb */}
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <Breadcrumb
          items={breadcrumbItems.map((item, index) => ({
            title: (
              <span
                style={{
                  fontSize: 13,
                  color:
                    index === breadcrumbItems.length - 1
                      ? 'var(--text-secondary)'
                      : 'var(--text-muted)',
                  fontWeight: index === breadcrumbItems.length - 1 ? 500 : 400,
                }}
              >
                {item.title}
              </span>
            ),
          }))}
        />
      </div>

      {/* Right: Global Search + Refresh */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        {/* Global Search */}
        <AutoComplete
          value={searchValue}
          options={filteredOptions}
          onSelect={(path: string) => {
            setSearchValue('')
            navigate(path)
          }}
          onChange={setSearchValue}
          style={{ width: 220 }}
        >
          <Input
            prefix={<SearchOutlined style={{ color: 'var(--text-icon)' }} />}
            placeholder={t('searchMenu')}
            allowClear
            style={{
              height: 34,
              borderRadius: 'var(--radius-md)',
              border: '1px solid var(--border-color)',
              background: 'var(--bg-light)',
              fontSize: 13,
            }}
            onFocus={e => {
              e.target.style.borderColor = 'var(--brand-dark)'
              e.target.style.background = 'var(--bg-white)'
            }}
            onBlur={e => {
              e.target.style.borderColor = 'var(--border-color)'
              e.target.style.background = 'var(--bg-light)'
            }}
          />
        </AutoComplete>

        {/* Refresh Button */}
        <button
          type="button"
          className="icon-btn"
          onClick={onRefresh}
          title={t('refresh')}
        >
          <ReloadOutlined style={{ fontSize: 16, color: 'var(--text-secondary)' }} />
        </button>

        {/* Language Switcher */}
        <Dropdown
          menu={{
            items: [
              { key: 'zh-CN', label: <span style={{ fontSize: 14, display: 'flex', alignItems: 'center', gap: 6 }}><span className="fi fi-cn" style={{ fontSize: 16 }} /> 中文</span> },
              { key: 'en-US', label: <span style={{ fontSize: 14, display: 'flex', alignItems: 'center', gap: 6 }}><span className="fi fi-us" style={{ fontSize: 16 }} /> English</span> },
            ],
            onClick: ({ key }) => i18n.changeLanguage(key),
            selectedKeys: [i18n.language],
          }}
          trigger={['click']}
        >
          <button
            type="button"
            className="icon-btn"
            style={{ fontSize: 18 }}
          >
            <span className={`fi fi-${countryCodeMap[i18n.language] || 'us'}`} style={{ fontSize: 16 }} />
          </button>
        </Dropdown>
      </div>
    </div>
  )
}
