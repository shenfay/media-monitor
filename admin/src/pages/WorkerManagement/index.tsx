import { useState, useEffect, useCallback, useRef } from 'react'
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

const statusMap: Record<string, { color: string; label: string }> = {
  online: { color: 'success', label: '在线' },
  offline: { color: 'default', label: '离线' },
  paused: { color: 'warning', label: '暂停' },
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

function timeAgo(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  const diff = Math.floor((Date.now() - d.getTime()) / 1000)
  if (diff < 60) return `${diff}秒前`
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

export default function WorkerManagement() {
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
      message.success(`已发送 ${action} 指令`)
      setTimeout(fetchData, 1000)
    } catch {
      message.error('指令发送失败')
    } finally {
      setActionLoading(null)
    }
  }

  const workerColumns = [
    {
      title: '名称',
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
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (v: string) => {
        const s = statusMap[v] || { color: 'default', label: v }
        return <Badge status={s.color as 'success' | 'default' | 'warning'} text={s.label} />
      },
    },
    {
      title: '适配器',
      dataIndex: 'adapters',
      key: 'adapters',
      render: (v: string[]) =>
        v?.length ? v.map((a) => <Tag key={a}>{a}</Tag>) : <span style={{ color: 'var(--text-quaternary)' }}>-</span>,
    },
    {
      title: '能力标签',
      dataIndex: 'capabilities',
      key: 'capabilities',
      render: (v: string[]) =>
        v?.length
          ? v.map((c) => <Tag key={c} color="blue">{c}</Tag>)
          : <span style={{ color: 'var(--text-quaternary)' }}>通用</span>,
    },
    {
      title: '并发',
      dataIndex: 'concurrency',
      key: 'concurrency',
      width: 60,
      render: (v: number) => v || '-',
    },
    {
      title: '已处理',
      dataIndex: 'processed_count',
      key: 'processed_count',
      width: 80,
      render: (v: number) => v || 0,
    },
    {
      title: '当前任务',
      dataIndex: 'current_task',
      key: 'current_task',
      width: 160,
      ellipsis: true,
      render: (v: string) => v ? <Tooltip title={v}><span>{v}</span></Tooltip> : <span style={{ color: 'var(--text-quaternary)' }}>空闲</span>,
    },
    {
      title: '最后心跳',
      dataIndex: 'last_heartbeat',
      key: 'last_heartbeat',
      width: 110,
      render: (v: string) => <Tooltip title={formatTime(v)}>{timeAgo(v)}</Tooltip>,
    },
    {
      title: '启动时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 160,
      render: (v: string) => formatTime(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 160,
      render: (_: unknown, record: WorkerInfo) => {
        const isPaused = record.status === 'paused'
        const isOffline = record.status === 'offline'
        const isLoading = actionLoading === record.id
        return (
          <Space size="small">
            {!isOffline && !isPaused && (
              <Popconfirm title="确认暂停该 Worker？" onConfirm={() => handleAction(record.id, 'pause')}>
                <Tooltip title="暂停">
                  <Button type="text" size="small" icon={<PauseCircleOutlined />} loading={isLoading} />
                </Tooltip>
              </Popconfirm>
            )}
            {!isOffline && isPaused && (
              <Tooltip title="恢复">
                <Button
                  type="text"
                  size="small"
                  icon={<PlayCircleOutlined />}
                  loading={isLoading}
                  onClick={() => handleAction(record.id, 'resume')}
                />
              </Tooltip>
            )}
            {!isOffline && (
              <Popconfirm title="确认下线该 Worker？（处理完当前任务后停止）" onConfirm={() => handleAction(record.id, 'shutdown')}>
                <Tooltip title="下线">
                  <Button type="text" size="small" danger icon={<StopOutlined />} loading={isLoading} />
                </Tooltip>
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
        title="Worker 管理"
        filters={
          <Button
            icon={<ReloadOutlined />}
            onClick={fetchData}
            loading={loading}
            style={{ color: 'var(--text-primary)' }}
          >
            刷新
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
          locale={{ emptyText: '暂无在线 Worker' }}
        />
      </DataPanel>

      {adapters.length > 0 && (
        <DataPanel title="适配器" style={{ marginTop: 16 }}>
          <Row gutter={[16, 16]}>
            {adapters.map((adapter) => (
              <Col key={adapter.name} xs={24} sm={12} md={8} lg={6}>
                <Card size="small" title={adapter.name}>
                  <div style={{ marginBottom: 8 }}>
                    <span style={{ color: 'var(--text-secondary)' }}>类型：</span>
                    <Tag>{adapter.platform_type || 'news'}</Tag>
                  </div>
                  <div style={{ marginBottom: 8 }}>
                    <span style={{ color: 'var(--text-secondary)' }}>所需标签：</span>
                    {adapter.required_tags?.length
                      ? adapter.required_tags.map((t) => <Tag key={t} color="orange">{t}</Tag>)
                      : <span style={{ color: 'var(--text-quaternary)' }}>无（通用）</span>
                    }
                  </div>
                  <div>
                    <span style={{ color: 'var(--text-secondary)' }}>首次发现：</span>
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
