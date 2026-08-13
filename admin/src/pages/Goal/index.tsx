import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Select } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function Goal() {
  const { t } = useTranslation()
  const [statusFilter, setStatusFilter] = useState('')

  const statusOptions = [
    { label: t('allStatus'), value: '' },
    { label: t('inProgress'), value: 'active' },
    { label: t('completed'), value: 'completed' },
    { label: t('archived'), value: 'archived' },
  ]

  const columns = [
    { title: t('goalName'), dataIndex: 'name', key: 'name' },
    { title: t('familyName'), dataIndex: 'familyName', key: 'familyName' },
    { title: t('goalType'), dataIndex: 'type', key: 'type' },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, { bg: string; color: string }> = {
          active: { bg: 'var(--green-light)', color: 'var(--green-text)' },
          completed: { bg: 'var(--blue-light)', color: 'var(--blue-text)' },
          archived: { bg: 'var(--gray-light)', color: 'var(--gray-text)' },
        }
        const labelMap: Record<string, string> = {
          active: t('inProgress'),
          completed: t('completed'),
          archived: t('archived'),
        }
        const c = colorMap[status] || { bg: 'var(--gray-light)', color: 'var(--gray-text)' }
        return <Tag style={{ background: c.bg, color: c.color }}>{labelMap[status] || status}</Tag>
      },
    },
    { title: t('createdAt'), dataIndex: 'createdAt', key: 'createdAt' },
  ]

  return (
    <DataPanel
      title={t('goalManagement')}
      filters={
        <>
          <FilterSearch placeholder={t('searchGoalName')} />
          <Select value={statusFilter} onChange={setStatusFilter} style={{ width: 140 }} options={statusOptions} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
        </>
      }
      toolbarActions={
        <Button type="primary" icon={<PlusOutlined />}>
          {t('addGoal')}
        </Button>
      }
    >
      <Table
        columns={columns}
        dataSource={[]}
        rowKey="id"
        locale={{ emptyText: t('noData') }}
        pagination={{ ...DEFAULT_PAGINATION, ...getPaginationShowTotal(t) }}
      />
    </DataPanel>
  )
}
