import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Select } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function ShopItem() {
  const { t } = useTranslation()
  const [approvalFilter, setApprovalFilter] = useState('')

  const approvalOptions = [
    { label: t('allApproval'), value: '' },
    { label: t('autoPassed'), value: 'auto' },
    { label: t('notifyParent'), value: 'notify' },
    { label: t('needApproval'), value: 'approve' },
  ]

  const columns = [
    { title: t('itemName'), dataIndex: 'name', key: 'name' },
    { title: t('itemDescription'), dataIndex: 'description', key: 'description', ellipsis: true },
    {
      title: t('pointsPrice'),
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => <Tag style={{ background: 'var(--yellow-light)', color: 'var(--yellow-text)' }}>{price} {t('pointsUnit')}</Tag>,
    },
    {
      title: t('approvalLevel'),
      dataIndex: 'approvalLevel',
      key: 'approvalLevel',
      render: (level: string) => {
        const labelMap: Record<string, string> = {
          auto: t('autoPassed'),
          notify: t('notifyParent'),
          approve: t('needApproval'),
        }
        const colorMap: Record<string, { bg: string; color: string }> = {
          auto: { bg: 'var(--green-light)', color: 'var(--green-text)' },
          notify: { bg: 'var(--yellow-light)', color: 'var(--yellow-text)' },
          approve: { bg: 'var(--red-light)', color: 'var(--red-text)' },
        }
        const c = colorMap[level] || { bg: 'var(--gray-light)', color: 'var(--gray-text)' }
        return <Tag style={{ background: c.bg, color: c.color }}>{labelMap[level] || level}</Tag>
      },
    },
    { title: t('stock'), dataIndex: 'stock', key: 'stock' },
    { title: t('updatedAt'), dataIndex: 'updatedAt', key: 'updatedAt' },
  ]

  return (
    <DataPanel
      title={t('shopItemManagement')}
      filters={
        <>
          <FilterSearch placeholder={t('searchItemName')} />
          <Select value={approvalFilter} onChange={setApprovalFilter} style={{ width: 140 }} options={approvalOptions} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
        </>
      }
      toolbarActions={
        <Button type="primary" icon={<PlusOutlined />}>
          {t('addItem')}
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
