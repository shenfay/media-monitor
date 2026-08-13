import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Select, Button } from 'antd'
import { ReloadOutlined } from '@ant-design/icons'
import DataPanel from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'
import { useCrudList } from '@/hooks/useCrudList'
import { getOperationLogs, type OperationLogRecord } from '@/services/operationLog'

function formatTime(dateStr: string, lang: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString(lang, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

export default function OperationLog() {
  const { t, i18n } = useTranslation()
  const [categoryFilter, setCategoryFilter] = useState('')

  const categoryOptions = [
    { label: t('allCategories'), value: '' },
    { label: t('categoryAuth'), value: 'AUTH' },
    { label: t('categoryUser'), value: 'USER' },
    { label: t('categorySystem'), value: 'SYSTEM' },
    { label: t('categoryBiz'), value: 'BIZ' },
  ]

  const actionLabelMap: Record<string, string> = {
    'AUTH.LOGIN.SUCCESS': t('actionLoginSuccess'),
    'AUTH.LOGIN.FAILED': t('actionLoginFailed'),
    'AUTH.LOGOUT': t('actionLogout'),
    'AUTH.TOKEN.REFRESHED': t('actionTokenRefreshed'),
    'AUTH.ACCOUNT.LOCKED': t('actionAccountLocked'),
    'USER.REGISTER': t('actionUserRegister'),
    'USER.PROFILE.UPDATED': t('actionProfileUpdated'),
    'USER.CREATE': t('actionUserCreate', '创建用户'),
    'USER.UPDATE': t('actionUserUpdate', '更新用户'),
    'USER.ENABLE': t('actionUserEnable', '启用用户'),
    'USER.DISABLE': t('actionUserDisable', '禁用用户'),
    'ROLE.CREATE': t('actionRoleCreate', '创建角色'),
    'ROLE.UPDATE': t('actionRoleUpdate', '更新角色'),
    'ROLE.DELETE': t('actionRoleDelete', '删除角色'),
    'ROLE.TOGGLE_STATUS': t('actionRoleToggleStatus', '切换角色状态'),
    'PERMISSION.UPDATE': t('actionPermissionUpdate', '更新权限'),
    'MENU.CREATE': t('actionMenuCreate', '创建菜单'),
    'MENU.UPDATE': t('actionMenuUpdate', '更新菜单'),
    'MENU.DELETE': t('actionMenuDelete', '删除菜单'),
    'MENU.TOGGLE_STATUS': t('actionMenuToggleStatus', '切换菜单状态'),
    'MENU.SORT': t('actionMenuSort', '菜单排序'),
    'SYSTEM.CONFIG.UPDATED': t('actionConfigUpdated'),
    'SYSTEM.PERMISSION.CHANGED': t('actionPermissionChanged'),
  }

  function formatAction(action: string): string {
    return actionLabelMap[action] || action
  }

  const { loading, dataSource, total, page, pageSize, fetchData, handlePageChange } =
    useCrudList<OperationLogRecord>(
      async ({ page: p, pageSize: ps }) => {
        const res = await getOperationLogs({
          category: categoryFilter || undefined,
          limit: ps,
          offset: (p - 1) * ps,
        })
        const data = res.data || []
        const inferredTotal = data.length >= ps ? p * ps + 1 : (p - 1) * ps + data.length
        return { data, total: inferredTotal }
      },
    )

  const handleCategoryChange = (v: string) => {
    setCategoryFilter(v)
  }

  // 分类标签颜色映射
  const categoryColorMap: Record<string, { bg: string; color: string }> = {
    AUTH: { bg: 'var(--blue-light)', color: 'var(--blue-text)' },
    USER: { bg: 'var(--green-light)', color: 'var(--green-text)' },
    SYSTEM: { bg: 'var(--gray-light)', color: 'var(--gray-text)' },
    BIZ: { bg: 'var(--yellow-light)', color: 'var(--yellow-text)' },
  }

  const columns = [
    {
      title: t('time'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      render: (v: string) => formatTime(v, i18n.language),
    },
    {
      title: t('operator'),
      dataIndex: 'email',
      key: 'email',
      render: (v: string, record: OperationLogRecord) => v || record.user_id,
    },
    {
      title: t('category'),
      dataIndex: 'category',
      key: 'category',
      width: 90,
      render: (v: string) => {
        const c = categoryColorMap[v] || { bg: 'var(--gray-light)', color: 'var(--gray-text)' }
        return <Tag style={{ background: c.bg, color: c.color }}>{v}</Tag>
      },
    },
    {
      title: t('actions'),
      dataIndex: 'action',
      key: 'action',
      render: (v: string) => formatAction(v),
    },
    {
      title: t('ipAddress'),
      dataIndex: 'ip',
      key: 'ip',
      width: 140,
      render: (v: string) => v || '-',
    },
    {
      title: t('device'),
      dataIndex: 'device',
      key: 'device',
      width: 100,
      render: (v: string) => v || '-',
    },
    {
      title: t('result'),
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (v: string) => {
        const isSuccess = v === 'SUCCESS'
        return (
          <Tag style={{
            background: isSuccess ? 'var(--green-light)' : 'var(--red-light)',
            color: isSuccess ? 'var(--green-text)' : 'var(--red-text)',
          }}>
            {isSuccess ? t('success') : t('failed')}
          </Tag>
        )
      },
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('operationLog')}
        filters={
          <>
            <Select
              value={categoryFilter}
              onChange={handleCategoryChange}
              style={{ width: 140 }}
              options={categoryOptions}
            />
            <Button
              icon={<ReloadOutlined />}
              onClick={fetchData}
              style={{ color: 'var(--text-primary)' }}
            >
              {t('refresh')}
            </Button>
          </>
        }
      >
        <Table
          dataSource={dataSource}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            ...DEFAULT_PAGINATION,
            ...getPaginationShowTotal(t),
            onChange: handlePageChange,
          }}
        />
      </DataPanel>
    </div>
  )
}
