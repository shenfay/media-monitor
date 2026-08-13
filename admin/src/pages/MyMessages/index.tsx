import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Select, Button, Badge } from 'antd'
import { ReloadOutlined, CheckOutlined } from '@ant-design/icons'
import DataPanel from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'
import { useCrudList } from '@/hooks/useCrudList'
import {
  getMyMessages,
  getUnreadCount,
  markAsRead,
  markAllAsRead,
  type MessageRecord,
} from '@/services/message'
import { message } from 'antd'
import { usePushNotification } from '@/hooks/useWebSocket'

function formatTime(dateStr: string, lang: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString(lang, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

export default function MyMessages() {
  const { t, i18n } = useTranslation()
  const [typeFilter, setTypeFilter] = useState('')
  const [readFilter, setReadFilter] = useState('')
  const [unreadTotal, setUnreadTotal] = useState(0)

  const typeOptions = [
    { label: t('msgAllTypes'), value: '' },
    { label: t('msgTypeSystem'), value: 'system' },
    { label: t('msgTypeCompanion'), value: 'companion' },
  ]

  const readOptions = [
    { label: t('msgAllStatus'), value: '' },
    { label: t('msgUnread'), value: 'false' },
    { label: t('msgRead'), value: 'true' },
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

  const fetchUnread = useCallback(async () => {
    try {
      const res = await getUnreadCount()
      setUnreadTotal(res.total ?? 0)
    } catch {
      // ignore
    }
  }, [])

  // 初始加载未读数 + 切回 tab 刷新
  useEffect(() => {
    fetchUnread()

    const onVisible = () => {
      if (document.visibilityState === 'visible') fetchUnread()
    }
    document.addEventListener('visibilitychange', onVisible)
    return () => document.removeEventListener('visibilitychange', onVisible)
  }, [fetchUnread])

  // WebSocket 实时推送：弹出通知 + 刷新列表和未读数
  usePushNotification(() => {
    fetchUnread()
    fetchData()
  })

  const { loading, dataSource, total, page, pageSize, fetchData, handlePageChange } =
    useCrudList<MessageRecord>(
      async ({ page: p, pageSize: ps }) => {
        const params: Record<string, unknown> = {
          limit: ps,
          offset: (p - 1) * ps,
        }
        if (typeFilter) params.type = typeFilter
        if (readFilter) params.is_read = readFilter === 'true'
        const res = await getMyMessages(params as { type?: string; is_read?: boolean; limit?: number; offset?: number })
        const data = res.messages || []
        const inferredTotal = data.length >= ps ? p * ps + 1 : (p - 1) * ps + data.length
        return { data, total: res.total || inferredTotal }
      },
    )

  const handleMarkAsRead = async (id: string) => {
    try {
      await markAsRead(id)
      message.success(t('markReadSuccess'))
      fetchData()
      fetchUnread()
    } catch {
      message.error(t('operationFailed'))
    }
  }

  const handleMarkAllAsRead = async () => {
    try {
      await markAllAsRead(typeFilter || undefined)
      message.success(t('markAllReadSuccess'))
      fetchData()
      fetchUnread()
    } catch {
      message.error(t('operationFailed'))
    }
  }

  const typeColorMap: Record<string, { bg: string; color: string }> = {
    system: { bg: 'var(--blue-light)', color: 'var(--blue-text)' },
    companion: { bg: 'var(--yellow-light)', color: 'var(--yellow-text)' },
  }

  const columns = [
    {
      title: t('time'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
      render: (v: string) => formatTime(v, i18n.language),
    },
    {
      title: t('msgTitle'),
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
      render: (v: string, record: MessageRecord) => (
        <span style={{ fontWeight: record.is_read ? 400 : 600 }}>
          {!record.is_read && <Badge status="processing" style={{ marginRight: 6 }} />}
          {v}
        </span>
      ),
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
    {
      title: t('actions'),
      key: 'actions',
      width: 80,
      render: (_: unknown, record: MessageRecord) => (
        !record.is_read ? (
          <Button type="link" size="small" onClick={() => handleMarkAsRead(record.id)}>
            {t('markRead')}
          </Button>
        ) : null
      ),
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('myMessages')}
        filters={
          <>
            <Select
              value={typeFilter}
              onChange={(v) => { setTypeFilter(v); setTimeout(fetchData, 0) }}
              style={{ width: 140 }}
              options={typeOptions}
            />
            <Select
              value={readFilter}
              onChange={(v) => { setReadFilter(v); setTimeout(fetchData, 0) }}
              style={{ width: 120 }}
              options={readOptions}
            />
            <Button
              icon={<CheckOutlined />}
              onClick={handleMarkAllAsRead}
              disabled={unreadTotal === 0}
            >
              {t('markAllRead')}
            </Button>
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
