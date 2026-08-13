import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Select, Button } from 'antd'
import { ReloadOutlined } from '@ant-design/icons'
import DataPanel from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'
import { useCrudList } from '@/hooks/useCrudList'
import { getMessages, type MessageRecord } from '@/services/message'

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

export default function MessageManagement() {
  const { t, i18n } = useTranslation()
  const [typeFilter, setTypeFilter] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')

  const typeOptions = [
    { label: t('msgAllTypes'), value: '' },
    { label: t('msgTypeSystem'), value: 'system' },
    { label: t('msgTypeCompanion'), value: 'companion' },
  ]

  const categoryOptions = [
    { label: t('msgAllCategories'), value: '' },
    { label: t('msgCatVerification'), value: 'verification' },
    { label: t('msgCatPoints'), value: 'points' },
    { label: t('msgCatGoal'), value: 'goal' },
    { label: t('msgCatCompanionStatus'), value: 'companion_status' },
    { label: t('msgCatExchange'), value: 'exchange' },
    { label: t('msgCatCompanionEncourage'), value: 'companion_encourage' },
    { label: t('msgCatCompanionRemind'), value: 'companion_remind' },
    { label: t('msgCatCompanionCelebrate'), value: 'companion_celebrate' },
  ]

  const categoryLabelMap: Record<string, string> = {
    verification: t('msgCatVerification'),
    points: t('msgCatPoints'),
    goal: t('msgCatGoal'),
    companion_status: t('msgCatCompanionStatus'),
    exchange: t('msgCatExchange'),
    companion_encourage: t('msgCatCompanionEncourage'),
    companion_remind: t('msgCatCompanionRemind'),
    companion_celebrate: t('msgCatCompanionCelebrate'),
  }

  const { loading, dataSource, total, page, pageSize, fetchData, handlePageChange } =
    useCrudList<MessageRecord>(
      async ({ page: p, pageSize: ps }) => {
        const res = await getMessages({
          type: typeFilter || undefined,
          category: categoryFilter || undefined,
          limit: ps,
          offset: (p - 1) * ps,
        })
        const data = res.messages || []
        const inferredTotal = data.length >= ps ? p * ps + 1 : (p - 1) * ps + data.length
        return { data, total: res.total || inferredTotal }
      },
    )

  // 类型标签颜色
  const typeColorMap: Record<string, { bg: string; color: string }> = {
    system: { bg: 'var(--blue-light)', color: 'var(--blue-text)' },
    companion: { bg: 'var(--yellow-light)', color: 'var(--yellow-text)' },
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
      title: t('msgTitle'),
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: t('msgType'),
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (v: string) => {
        const c = typeColorMap[v] || { bg: 'var(--gray-light)', color: 'var(--gray-text)' }
        return (
          <Tag style={{ background: c.bg, color: c.color }}>
            {v === 'system' ? t('msgTypeSystem') : t('msgTypeCompanion')}
          </Tag>
        )
      },
    },
    {
      title: t('msgCategory'),
      dataIndex: 'category',
      key: 'category',
      width: 120,
      render: (v: string) => categoryLabelMap[v] || v,
    },
    {
      title: t('msgRecipient'),
      dataIndex: 'recipient_id',
      key: 'recipient_id',
      width: 120,
      ellipsis: true,
      render: (v: string) => v?.slice(0, 8) + '...' || '-',
    },
    {
      title: t('msgReadStatus'),
      dataIndex: 'is_read',
      key: 'is_read',
      width: 80,
      render: (v: boolean) => (
        <Tag style={{
          background: v ? 'var(--green-light)' : 'var(--gray-light)',
          color: v ? 'var(--green-text)' : 'var(--gray-text)',
        }}>
          {v ? t('msgRead') : t('msgUnread')}
        </Tag>
      ),
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('messageManagement')}
        filters={
          <>
            <Select
              value={typeFilter}
              onChange={(v) => { setTypeFilter(v); setTimeout(fetchData, 0) }}
              style={{ width: 140 }}
              options={typeOptions}
            />
            <Select
              value={categoryFilter}
              onChange={(v) => { setCategoryFilter(v); setTimeout(fetchData, 0) }}
              style={{ width: 160 }}
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
