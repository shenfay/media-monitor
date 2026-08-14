import { useState, useEffect, useCallback, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Space, Popconfirm, message, Card, Row, Col, Tooltip, Badge } from 'antd'
import { ReloadOutlined, PauseCircleOutlined, PlayCircleOutlined, StopOutlined } from '@ant-design/icons'
import DataPanel from '@/components/DataPanel'
import {
  getWorkers,
  getAdapters,
  pauseWorker,
  resumeWorker,
  shutdownWorker,
  type WorkerInfo,
  type AdapterMeta,
} from '@/services/worker'

const REFRESH_INTERVAL = 10_000

const statusColors: Record<string, string> = {
  online: 'success',
  offline: 'default',
  paused: 'warning',
}

const statusKeys: Record<string, string> = {
  online: 'workerOnline',
  offline: 'workerOffline',
  paused: 'workerPaused',
}

function formatTime(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

function timeAgo(dateStr: string, t: (key: string, opts?: Record<string, unknown>) => string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  const diff = Math.floor((Date.now() - d.getTime()) / 1000)
  if (diff < 60) return t('workerTimeAgoSec', { n: diff })
  if (diff < 3600) return t('workerTimeAgoMin', { n: Math.floor(diff / 60) })
  if (diff < 86400) return t('workerTimeAgoHour', { n: Math.floor(diff / 3600) })
  return t('workerTimeAgoDay', { n: Math.floor(diff / 86400) })
}

export default function WorkerManagement() {
  const { t } = useTranslation()
  const [workers, setWorkers] = useState<WorkerInfo[]>([])
  const [adapters, setAdapters] = useState<AdapterMeta[]>([])
  const [loading, setLoading] = useState(false)
  const [actionLoading, setActionLoading] = useState<string | null>(null)
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const [w, a] = await Promise.all([getWorkers(), getAdapters()])
      setWorkers(w)
      setAdapters(a)
    } catch {
      // ignore
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchData()
    timerRef.current = setInterval(fetchData, REFRESH_INTERVAL)
    return () => {
      if (timerRef.current) clearInterval(timerRef.current)
    }
  }, [fetchData])

  const handleAction = async (id: string, action: 'pause' | 'resume' | 'shutdown') => {
    setActionLoading(id)
    try {
      if (action === 'pause') await pauseWorker(id)
      else if (action === 'resume') await resumeWorker(id)
      else await shutdownWorker(id)
      message.success(t('workerCmdSent', { action }))
      setTimeout(fetchData, 1000)
    } catch {
      message.error(t('workerCmdFailed'))
    } finally {
      setActionLoading(null)
    }
  }

  const workerColumns = [
    {
      title: t('workerName'),
      dataIndex: 'name',
      key: 'name',
      width: 140,
      render: (v: string, record: WorkerInfo) => (
        <Tooltip title={`ID: ${record.id}`}>
          <span style={{ fontWeight: 500 }}>{v || record.id.slice(0, 8)}</span>
        </Tooltip>
      ),
    },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (v: string) => {
        const color = statusColors[v] || 'default'
        return <Badge status={color as 'success' | 'default' | 'warning'} text={t(statusKeys[v] || v)} />
      },
    },
    {
      title: t('workerAdapter'),
      dataIndex: 'adapters',
      key: 'adapters',
      render: (v: string[]) =>
        v?.length ? v.map((a) => <Tag key={a}>{a}</Tag>) : <span style={{ color: 'var(--text-quaternary)' }}>-</span>,
    },
    {
      title: t('crawlCapabilityTags'),
      dataIndex: 'capabilities',
      key: 'capabilities',
      render: (v: string[]) =>
        v?.length
          ? v.map((c) => <Tag key={c} color="blue">{c}</Tag>)
          : <span style={{ color: 'var(--text-quaternary)' }}>{t('crawlUniversal')}</span>,
    },
    {
      title: t('workerConcurrency'),
      dataIndex: 'concurrency',
      key: 'concurrency',
      width: 80,
      render: (v: number) => v || '-',
    },
    {
      title: t('workerProcessed'),
      dataIndex: 'processed_count',
      key: 'processed_count',
      width: 80,
      render: (v: number) => v || 0,
    },
    {
      title: t('workerCurrentTask'),
      dataIndex: 'current_task',
      key: 'current_task',
      width: 160,
      ellipsis: true,
      render: (v: string) => v ? <Tooltip title={v}><span>{v}</span></Tooltip> : <span style={{ color: 'var(--text-quaternary)' }}>{t('workerIdle')}</span>,
    },
    {
      title: t('workerLastHeartbeat'),
      dataIndex: 'last_heartbeat',
      key: 'last_heartbeat',
      width: 110,
      render: (v: string) => <Tooltip title={formatTime(v)}>{timeAgo(v, t)}</Tooltip>,
    },
    {
      title: t('workerStartedAt'),
      dataIndex: 'started_at',
      key: 'started_at',
      width: 160,
      render: (v: string) => formatTime(v),
    },
    {
      title: t('actions'),
      key: 'actions',
      width: 200,
      render: (_: unknown, record: WorkerInfo) => {
        const isPaused = record.status === 'paused'
        const isOffline = record.status === 'offline'
        const isLoading = actionLoading === record.id
        return (
          <Space size={4}>
            {!isOffline && !isPaused && (
              <Popconfirm title={t('workerConfirmPause')} onConfirm={() => handleAction(record.id, 'pause')}>
                <Button type="link" size="small" icon={<PauseCircleOutlined />} loading={isLoading}>
                  {t('workerPause')}
                </Button>
              </Popconfirm>
            )}
            {!isOffline && isPaused && (
              <Button
                type="link"
                size="small"
                icon={<PlayCircleOutlined />}
                loading={isLoading}
                onClick={() => handleAction(record.id, 'resume')}
              >
                {t('workerResume')}
              </Button>
            )}
            {!isOffline && (
              <Popconfirm title={t('workerConfirmShutdown')} onConfirm={() => handleAction(record.id, 'shutdown')}>
                <Button type="link" size="small" danger icon={<StopOutlined />} loading={isLoading}>
                  {t('workerShutdown')}
                </Button>
              </Popconfirm>
            )}
          </Space>
        )
      },
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('workerMgmt')}
        toolbarActions={
          <Button
            icon={<ReloadOutlined />}
            onClick={fetchData}
            loading={loading}
          >
            {t('refresh')}
          </Button>
        }
      >
        <Table
          dataSource={workers}
          columns={workerColumns}
          rowKey="id"
          loading={loading}
          pagination={false}
          size="small"
          locale={{ emptyText: t('workerNoWorkers') }}
        />
      </DataPanel>

      {adapters.length > 0 && (
        <DataPanel title={t('workerAdapters')} style={{ marginTop: 16 }}>
          <Row gutter={[16, 16]}>
            {adapters.map((adapter) => (
              <Col key={adapter.name} xs={24} sm={12} md={8} lg={6}>
                <Card size="small" title={adapter.name}>
                  <div style={{ marginBottom: 8 }}>
                    <span style={{ color: 'var(--text-secondary)' }}>{t('workerPlatformType')}：</span>
                    <Tag>{adapter.platform_type || 'news'}</Tag>
                  </div>
                  <div style={{ marginBottom: 8 }}>
                    <span style={{ color: 'var(--text-secondary)' }}>{t('workerRequiredTags')}：</span>
                    {adapter.required_tags?.length
                      ? adapter.required_tags.map((tag) => <Tag key={tag} color="orange">{tag}</Tag>)
                      : <span style={{ color: 'var(--text-quaternary)' }}>{t('workerNoTags')}</span>
                    }
                  </div>
                  <div>
                    <span style={{ color: 'var(--text-secondary)' }}>{t('workerFirstSeen')}：</span>
                    <span style={{ fontSize: 12, color: 'var(--text-tertiary)' }}>
                      {formatTime(adapter.first_seen)}
                    </span>
                  </div>
                </Card>
              </Col>
            ))}
          </Row>
        </DataPanel>
      )}
    </div>
  )
}
