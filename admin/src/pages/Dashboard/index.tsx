import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Row, Col, Table, Tag } from 'antd'
import DataPanel from '@/components/DataPanel'
import StatCard from '@/components/StatCard'
import { getDashboardStats, getTasks, getSources, type DashboardStats as Stats, type TaskRun, type Source } from '@/services/crawl'
import { formatTimeShort, formatDuration } from '@/utils/format'
import { TASK_STATUS_COLORS, TASK_STATUS_KEYS } from '@/constants/status'


export default function Dashboard() {
  const { t, i18n } = useTranslation()
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
      render: (v: string) => <Tag color={TASK_STATUS_COLORS[v] || 'default'}>{t(TASK_STATUS_KEYS[v] || v)}</Tag>,
    },
    {
      title: t('crawlResult'), key: 'counts', width: 120,
      render: (_: unknown, r: TaskRun) => `${r.ingested}/${r.total}${r.failed ? ` (${t('crawlFailedCount', { n: r.failed })})` : ''}`,
    },
    { title: t('crawlDuration'), key: 'duration', width: 80, render: (_: unknown, r: TaskRun) => formatDuration(r.started_at, r.finished_at, t) },
    { title: t('createdAt'), dataIndex: 'created_at', key: 'created_at', width: 120, render: (v: string) => formatTimeShort(v, i18n.language) },
  ]

  return (
    <DataPanel title={t('crawlDashboard')}>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            label={t('crawlDashSources')}
            value={stats ? `${stats.enabled_sources}/${stats.total_sources}` : '-/-'}
            color="var(--blue)"
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            label={t('crawlDashTodayArticles')}
            value={stats?.today_articles ?? 0}
            color="var(--green)"
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            label={t('crawlDashRunningTasks')}
            value={stats?.running_tasks ?? 0}
            color="var(--yellow)"
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            label={t('crawlDashOnlineWorkers')}
            value={stats?.online_workers ?? 0}
            color="var(--green)"
          />
        </Col>
      </Row>

      <div style={{ marginTop: 16 }}>
        <div
          style={{
            border: '1px solid var(--border-light)',
            borderRadius: 'var(--radius-md)',
            background: 'var(--bg-white)',
            overflow: 'hidden',
          }}
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
        </div>
      </div>
    </DataPanel>
  )
}
