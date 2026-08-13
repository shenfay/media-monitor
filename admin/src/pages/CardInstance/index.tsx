import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Select } from 'antd'
import { EyeOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function CardInstance() {
  const { t } = useTranslation()
  const [statusFilter, setStatusFilter] = useState('')

  const statusOptions = [
    { label: t('allStatus'), value: '' },
    { label: t('pendingReview'), value: 'pending' },
    { label: t('approved'), value: 'approved' },
    { label: t('rejected'), value: 'rejected' },
    { label: t('autoPassed'), value: 'auto_passed' },
  ]

  const columns = [
    { title: t('cardContent'), dataIndex: 'content', key: 'content', ellipsis: true },
    { title: t('submitter'), dataIndex: 'childName', key: 'childName' },
    { title: t('belongingGoal'), dataIndex: 'goalName', key: 'goalName' },
    { title: t('cardTemplate'), dataIndex: 'templateName', key: 'templateName' },
    {
      title: t('acceptanceStatus'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, { bg: string; color: string }> = {
          pending: { bg: 'var(--yellow-light)', color: 'var(--yellow-text)' },
          approved: { bg: 'var(--green-light)', color: 'var(--green-text)' },
          rejected: { bg: 'var(--red-light)', color: 'var(--red-text)' },
          auto_passed: { bg: 'var(--blue-light)', color: 'var(--blue-text)' },
        }
        const labelMap: Record<string, string> = {
          pending: t('pendingReview'),
          approved: t('approved'),
          rejected: t('rejected'),
          auto_passed: t('autoPassed'),
        }
        const c = colorMap[status] || { bg: 'var(--gray-light)', color: 'var(--gray-text)' }
        return <Tag style={{ background: c.bg, color: c.color }}>{labelMap[status] || status}</Tag>
      },
    },
    { title: t('submissionTime'), dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: t('actions'),
      key: 'action',
      render: () => (
        <Button type="link" size="small" icon={<EyeOutlined />}>{t('view')}</Button>
      ),
    },
  ]

  return (
    <DataPanel
      title={t('cardInstanceManagement')}
      filters={
        <>
          <FilterSearch placeholder={t('searchCardContent')} />
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
