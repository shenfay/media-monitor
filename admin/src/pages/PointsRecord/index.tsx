import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Statistic, Row, Col, Card, Button, Select } from 'antd'
import { ArrowUpOutlined, ArrowDownOutlined, SearchOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'

export default function PointsRecord() {
  const { t } = useTranslation()
  const [typeFilter, setTypeFilter] = useState('')

  const typeOptions = [
    { label: t('allTypes'), value: '' },
    { label: t('earn'), value: 'earn' },
    { label: t('spend'), value: 'spend' },
  ]

  const columns = [
    { title: t('user'), dataIndex: 'userName', key: 'userName' },
    {
      title: t('type'),
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag style={{
          background: type === 'earn' ? 'var(--green-light)' : 'var(--red-light)',
          color: type === 'earn' ? 'var(--green-text)' : 'var(--red-text)',
        }}>
          {type === 'earn' ? t('earn') : t('spend')}
        </Tag>
      ),
    },
    {
      title: t('points'),
      dataIndex: 'points',
      key: 'points',
      render: (points: number, record: { type: string }) => (
        <span style={{ color: record.type === 'earn' ? 'var(--green)' : 'var(--red)', fontSize: 13 }}>
          {record.type === 'earn' ? '+' : '-'}{points}
        </span>
      ),
    },
    { title: t('source'), dataIndex: 'source', key: 'source' },
    { title: t('description'), dataIndex: 'description', key: 'description', ellipsis: true },
    { title: t('time'), dataIndex: 'createdAt', key: 'createdAt' },
  ]

  return (
    <div>
      <div style={{ padding: '20px 28px 0' }}>
        <Row gutter={16}>
          <Col span={8}>
            <Card style={{ borderRadius: 'var(--radius-md)', borderColor: 'var(--border-light)' }}>
              <Statistic title={t('todayPointsIssued')} value={0} prefix={<ArrowUpOutlined />} valueStyle={{ color: 'var(--green)' }} />
            </Card>
          </Col>
          <Col span={8}>
            <Card style={{ borderRadius: 'var(--radius-md)', borderColor: 'var(--border-light)' }}>
              <Statistic title={t('todayPointsSpent')} value={0} prefix={<ArrowDownOutlined />} valueStyle={{ color: 'var(--red)' }} />
            </Card>
          </Col>
          <Col span={8}>
            <Card style={{ borderRadius: 'var(--radius-md)', borderColor: 'var(--border-light)' }}>
              <Statistic title={t('totalPointsBalance')} value={0} />
            </Card>
          </Col>
        </Row>
      </div>
      <DataPanel
        title={t('pointsRecord')}
        style={{ marginTop: 16 }}
        filters={
          <>
            <FilterSearch placeholder={t('searchUserOrSource')} />
            <Select value={typeFilter} onChange={setTypeFilter} style={{ width: 140 }} options={typeOptions} />
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
    </div>
  )
}
