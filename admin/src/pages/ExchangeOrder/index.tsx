import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Space, Button, Select } from 'antd'
import { CheckOutlined, CloseOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function ExchangeOrder() {
  const { t } = useTranslation()
  const [statusFilter, setStatusFilter] = useState('')

  const statusOptions = [
    { label: t('allStatus'), value: '' },
    { label: t('pending'), value: 'pending' },
    { label: t('approved'), value: 'approved' },
    { label: t('rejected'), value: 'rejected' },
    { label: t('completed'), value: 'completed' },
  ]

  const columns = [
    { title: t('exchanger'), dataIndex: 'childName', key: 'childName' },
    { title: t('itemName'), dataIndex: 'itemName', key: 'itemName' },
    {
      title: t('points'),
      dataIndex: 'points',
      key: 'points',
      render: (points: number) => <Tag style={{ background: 'var(--yellow-light)', color: 'var(--yellow-text)' }}>{points} {t('pointsUnit')}</Tag>,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, { bg: string; color: string }> = {
          pending: { bg: 'var(--yellow-light)', color: 'var(--yellow-text)' },
          approved: { bg: 'var(--green-light)', color: 'var(--green-text)' },
          rejected: { bg: 'var(--red-light)', color: 'var(--red-text)' },
          completed: { bg: 'var(--blue-light)', color: 'var(--blue-text)' },
        }
        const labelMap: Record<string, string> = {
          pending: t('pending'),
          approved: t('approved'),
          rejected: t('rejected'),
          completed: t('completed'),
        }
        const c = colorMap[status] || { bg: 'var(--gray-light)', color: 'var(--gray-text)' }
        return <Tag style={{ background: c.bg, color: c.color }}>{labelMap[status] || status}</Tag>
      },
    },
    { title: t('applicationTime'), dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: t('actions'),
      key: 'action',
      render: () => (
        <Space size={4}>
          <Button type="link" size="small" icon={<CheckOutlined />}>{t('approveBtn')}</Button>
          <Button type="link" size="small" icon={<CloseOutlined />}>{t('rejectBtn')}</Button>
        </Space>
      ),
    },
  ]

  return (
    <DataPanel
      title={t('exchangeOrder')}
      filters={
        <>
          <FilterSearch placeholder={t('searchExchangerOrItem')} />
          <Select value={statusFilter} onChange={setStatusFilter} style={{ width: 140 }} options={statusOptions} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
        </>
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
