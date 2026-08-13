import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Select } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function CardTemplate() {
  const { t } = useTranslation()
  const [typeFilter, setTypeFilter] = useState('')
  const [acceptanceFilter, setAcceptanceFilter] = useState('')

  const typeOptions = [
    { label: t('allTypes2'), value: '' },
    { label: t('flashCard'), value: 'flash' },
    { label: t('causeCard'), value: 'cause' },
    { label: t('digestCard'), value: 'digest' },
    { label: t('feynmanCard'), value: 'feynman' },
    { label: t('daily3Card'), value: 'daily3' },
  ]

  const acceptanceOptions = [
    { label: t('allAcceptance'), value: '' },
    { label: t('autoPassed'), value: 'auto' },
    { label: t('parentAcceptance'), value: 'manual' },
  ]

  const columns = [
    { title: t('templateName'), dataIndex: 'name', key: 'name' },
    {
      title: t('cardType'),
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        const labelMap: Record<string, string> = {
          flash: t('flashCard'),
          cause: t('causeCard'),
          digest: t('digestCard'),
          feynman: t('feynmanCard'),
          daily3: t('daily3Card'),
        }
        return <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)' }}>{labelMap[type] || type}</Tag>
      },
    },
    { title: t('difficulty'), dataIndex: 'difficulty', key: 'difficulty' },
    { title: t('basePoints'), dataIndex: 'basePoints', key: 'basePoints' },
    {
      title: t('acceptanceType'),
      dataIndex: 'acceptanceType',
      key: 'acceptanceType',
      render: (type: string) => (
        <Tag style={{
          background: type === 'auto' ? 'var(--green-light)' : 'var(--yellow-light)',
          color: type === 'auto' ? 'var(--green-text)' : 'var(--yellow-text)',
        }}>
          {type === 'auto' ? t('autoPassed') : t('parentAcceptance')}
        </Tag>
      ),
    },
    { title: t('updatedAt'), dataIndex: 'updatedAt', key: 'updatedAt' },
  ]

  return (
    <DataPanel
      title={t('cardTemplateManagement')}
      filters={
        <>
          <FilterSearch placeholder={t('searchTemplateName')} />
          <Select value={typeFilter} onChange={setTypeFilter} style={{ width: 140 }} options={typeOptions} />
          <Select value={acceptanceFilter} onChange={setAcceptanceFilter} style={{ width: 140 }} options={acceptanceOptions} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
        </>
      }
      toolbarActions={
        <Button type="primary" icon={<PlusOutlined />}>
          {t('addTemplate')}
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
