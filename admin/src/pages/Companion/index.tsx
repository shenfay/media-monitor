import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Select } from 'antd'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function Companion() {
  const { t } = useTranslation()
  const [personalityFilter, setPersonalityFilter] = useState('')

  const personalityOptions = [
    { label: t('allPersonality'), value: '' },
    { label: t('explorer'), value: 'explorer' },
    { label: t('guardian'), value: 'guardian' },
    { label: t('creator'), value: 'creator' },
    { label: t('thinker'), value: 'thinker' },
  ]

  const columns = [
    { title: t('companionName'), dataIndex: 'name', key: 'name' },
    {
      title: t('personalityType'),
      dataIndex: 'personality',
      key: 'personality',
      render: (personality: string) => {
        const labelMap: Record<string, string> = {
          explorer: t('explorer'),
          guardian: t('guardian'),
          creator: t('creator'),
          thinker: t('thinker'),
        }
        return <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)' }}>{labelMap[personality] || personality}</Tag>
      },
    },
    { title: t('unlockLevel'), dataIndex: 'unlockLevel', key: 'unlockLevel' },
    { title: t('description'), dataIndex: 'description', key: 'description', ellipsis: true },
    { title: t('updatedAt'), dataIndex: 'updatedAt', key: 'updatedAt' },
  ]

  return (
    <DataPanel
      title={t('companionManagement')}
      filters={
        <>
          <FilterSearch placeholder={t('searchCompanionName')} />
          <Select value={personalityFilter} onChange={setPersonalityFilter} style={{ width: 140 }} options={personalityOptions} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
        </>
      }
      toolbarActions={
        <Button type="primary" icon={<PlusOutlined />}>
          {t('addCompanion')}
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
