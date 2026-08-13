import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Button, Select } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function Family() {
  const { t } = useTranslation()
  const [statusFilter, setStatusFilter] = useState('')

  const statusOptions = [
    { label: t('allStatus'), value: '' },
    { label: t('active'), value: 'active' },
    { label: t('inactive'), value: 'inactive' },
  ]

  const columns = [
    { title: t('name'), dataIndex: 'name', key: 'name' },
    { title: t('parent'), dataIndex: 'parentName', key: 'parentName' },
    { title: t('childrenCount'), dataIndex: 'childrenCount', key: 'childrenCount' },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
          <span style={{
            width: 7, height: 7, borderRadius: '50%',
            background: status === 'active' ? 'var(--green)' : 'var(--border-hover)',
            display: 'inline-block',
          }} />
          <span style={{ color: 'var(--text-primary)' }}>{status === 'active' ? t('active') : t('inactive')}</span>
        </div>
      ),
    },
    { title: t('createdAt'), dataIndex: 'createdAt', key: 'createdAt' },
  ]

  return (
    <DataPanel
      title={t('familyManagement')}
      filters={
        <>
          <FilterSearch placeholder={t('searchFamilyName')} />
          <Select value={statusFilter} onChange={setStatusFilter} style={{ width: 140 }} options={statusOptions} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
        </>
      }
      toolbarActions={
        <Button type="primary" icon={<PlusOutlined />}>
          {t('addFamily')}
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
