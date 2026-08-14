import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Card, Row, Col, Statistic, Table, Tag } from 'antd'
import {
  DatabaseOutlined,
  FileTextOutlined,
  SyncOutlined,
  CloudServerOutlined,
} from '@ant-design/icons'
import { getDashboardStats, getTasks, getSources, type DashboardStats as Stats, type TaskRun, type Source } from '@/services/crawl'

const statCardStyle = { borderRadius: 'var(--radius-md)', borderColor: 'var(--border-light)' } as React.CSSProperties

const statusColors: Record<string, string> = {
  queued: 'default',
  running: 'processing',
  success: 'success',
  partial: 'warning',
  failed: 'error',
}

const statusKeys: Record<string, string> = {
  queued: 'crawlQueued',
  running: 'crawlRunning',
  success: 'success',
  partial: 'crawlPartial',
  failed: 'failed',
}

function formatTime(dateStr: string | null): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
}

function formatDuration(start: string | null, end: string | null, t: (key: string, opts?: Record<string, unknown>) => string): string {
  if (!start) return '-'
  const s = new Date(start).getTime()
  const e = end ? new Date(end).getTime() : Date.now()
  const diff = Math.floor((e - s) / 1000)
  if (diff < 60) return t('crawlDurSec', { n: diff })
  if (diff < 3600) return t('crawlDurMinSec', { m: Math.floor(diff / 60), s: diff % 60 })
  return t('crawlDurHourMin', { h: Math.floor(diff / 3600), m: Math.floor((diff % 3600) / 60) })
}

export default function Dashboard() {
  const { t } = useTranslation()
  const [stats, setStats] = useState<Stats | null>(null)
  const [tasks, setTasks] = useState<TaskRun[]>([])
  const [sources, setSources] = useState<Source[]>([])
  const [loading, setLoading] = useState(true)

  const sourceMap = Object.fromEntries(sources.map(s => [s.id, s.name]))

  useEffect(() => {
    Promise.all([
      getDashboardStats(true).then(s => { if (s) setStats(s) }).catch(() => {}),
      getTasks(undefined, true).then(setTasks).catch(() => {}),
      getSources(true).then(setSources).catch(() => {}),
    ]).finally(() => setLoading(false))
  }, [])

  const recentTasks = tasks.slice(0, 8)

  const columns = [
    {
      title: t('crawlTaskId'), dataIndex: 'id', key: 'id', width: 100,
      render: (v: string) => <span style={{ fontFamily: 'monospace', fontSize: 12 }}>{v.slice(0, 10)}…</span>,
    },
    { title: t('crawlDataSource'), dataIndex: 'source_id', key: 'source_id', width: 140, render: (v: string) => sourceMap[v] || v.slice(0, 10) },
    { title: t('crawlTrigger'), dataIndex: 'triggered_by', key: 'triggered_by', width: 70, render: (v: string) => <Tag>{v}</Tag> },
    {
      title: t('status'), dataIndex: 'status', key: 'status', width: 90,
      render: (v: string) => <Tag color={statusColors[v] || 'default'}>{t(statusKeys[v] || v)}</Tag>,
    },
    {
      title: t('crawlResult'), key: 'counts', width: 120,
      render: (_: unknown, r: TaskRun) => `${r.ingested}/${r.total}${r.failed ? ` (${t('crawlFailedCount', { n: r.failed })})` : ''}`,
    },
    { title: t('crawlDuration'), key: 'duration', width: 80, render: (_: unknown, r: TaskRun) => formatDuration(r.started_at, r.finished_at, t) },
    { title: t('createdAt'), dataIndex: 'created_at', key: 'created_at', width: 120, render: (v: string) => formatTime(v) },
  ]

  return (
    <div style={{ padding: '20px 28px' }}>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card style={statCardStyle} loading={loading}>
            <Statistic
              title={t('crawlDashSources')}
              value={stats ? `${stats.enabled_sources}/${stats.total_sources}` : '-/-'}
              prefix={<DatabaseOutlined />}
              valueStyle={{ color: 'var(--blue)' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card style={statCardStyle} loading={loading}>
            <Statistic
              title={t('crawlDashTodayArticles')}
              value={stats?.today_articles ?? 0}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: 'var(--green)' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card style={statCardStyle} loading={loading}>
            <Statistic
              title={t('crawlDashRunningTasks')}
              value={stats?.running_tasks ?? 0}
              prefix={<SyncOutlined />}
              valueStyle={{ color: 'var(--yellow)' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card style={statCardStyle} loading={loading}>
            <Statistic
              title={t('crawlDashOnlineWorkers')}
              value={stats?.online_workers ?? 0}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: 'var(--green)' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card
            title={t('crawlDashRecentTasks')}
            style={{ borderRadius: 'var(--radius-md)', borderColor: 'var(--border-light)' }}
            styles={{ body: { padding: 0 } }}
          >
            <Table
              dataSource={recentTasks}
              columns={columns}
              rowKey="id"
              loading={loading}
              pagination={false}
              size="small"
              locale={{ emptyText: t('crawlDashNoTasks') }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
