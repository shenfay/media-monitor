import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Space, Select, Modal, Form, InputNumber, Switch, message, Tooltip, Drawer, Descriptions } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import DataPanel from '@/components/DataPanel'
import { getTasks, getSources, createTask, type TaskRun, type Source } from '@/services/crawl'
import { formatTime, formatDuration } from '@/utils/format'
import { TASK_STATUS_COLORS, TASK_STATUS_KEYS } from '@/constants/status'

export default function TaskManagement() {
  const { t, i18n } = useTranslation()
  const [tasks, setTasks] = useState<TaskRun[]>([])
  const [sources, setSources] = useState<Source[]>([])
  const [loading, setLoading] = useState(false)
  const [sourceFilter, setSourceFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailTask, setDetailTask] = useState<TaskRun | null>(null)
  const [createOpen, setCreateOpen] = useState(false)
  const [form] = Form.useForm()

  const sourceMap = Object.fromEntries(sources.map(s => [s.id, s.name]))

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const data = await getTasks(sourceFilter ? { source_id: sourceFilter } : undefined)
      setTasks(data)
    } catch { /* ignore */ }
    finally { setLoading(false) }
  }, [sourceFilter])

  useEffect(() => {
    getSources().then(setSources).catch(() => {})
  }, [])

  useEffect(() => { fetchData() }, [fetchData])

  const filteredTasks = statusFilter ? tasks.filter(task => task.status === statusFilter) : tasks

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      const res = await createTask(values)
      message.success(`${t('crawlTaskCreated')}: ${res.task_id}`)
      setCreateOpen(false)
      form.resetFields()
      fetchData()
    } catch { /* validation */ }
  }

  const columns = [
    {
      title: t('crawlTaskId'), dataIndex: 'id', key: 'id', width: 100,
      render: (v: string) => <Tooltip title={v}><span style={{ fontFamily: 'monospace', fontSize: 12 }}>{v.slice(0, 10)}…</span></Tooltip>,
    },
    { title: t('crawlDataSource'), dataIndex: 'source_id', key: 'source_id', width: 140, render: (v: string) => sourceMap[v] || v.slice(0, 10) },
    { title: t('crawlTrigger'), dataIndex: 'triggered_by', key: 'triggered_by', width: 70, render: (v: string) => <Tag>{v}</Tag> },
    {
      title: t('status'), dataIndex: 'status', key: 'status', width: 90,
      render: (v: string) => <Tag color={TASK_STATUS_COLORS[v] || 'default'}>{t(TASK_STATUS_KEYS[v] || v)}</Tag>,
    },
    {
      title: t('crawlResult'), key: 'counts', width: 120,
      render: (_: unknown, r: TaskRun) => `${r.ingested}/${r.total}${r.failed ? ` (${r.failed}${t('failed')})` : ''}`,
    },
    { title: t('crawlDuration'), key: 'duration', width: 80, render: (_: unknown, r: TaskRun) => formatDuration(r.started_at, r.finished_at, t) },
    { title: t('createdAt'), dataIndex: 'created_at', key: 'created_at', width: 160, render: (v: string) => formatTime(v, i18n.language) },
  ]

  return (
    <div>
      <DataPanel
        title={t('crawlTaskMgmt')}
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { form.resetFields(); setCreateOpen(true) }}>{t('crawlManualCreate')}</Button>
        }
        filters={
          <>
            <Select
              value={sourceFilter} onChange={setSourceFilter} style={{ width: 160 }}
              placeholder={t('crawlDataSource')} allowClear
              options={[{ value: '', label: t('crawlDataSource') + ' — ' + t('all') }, ...sources.map(s => ({ value: s.id, label: s.name }))]}
            />
            <Select
              value={statusFilter} onChange={setStatusFilter} style={{ width: 120 }}
              options={[
                { value: '', label: t('status') + ' — ' + t('all') },
                { value: 'queued', label: t('crawlQueued') }, { value: 'running', label: t('crawlRunning') },
                { value: 'success', label: t('success') }, { value: 'failed', label: t('failed') },
              ]}
            />
          </>
        }
      >
        <Table
          dataSource={filteredTasks} columns={columns} rowKey="id" loading={loading} pagination={false} size="small"
          onRow={(record) => ({ onClick: () => { setDetailTask(record); setDetailOpen(true) }, style: { cursor: 'pointer' } })}
        />
      </DataPanel>

      <Drawer title={t('crawlTaskDetail')} open={detailOpen} onClose={() => setDetailOpen(false)} width={480}>
        {detailTask && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="ID">{detailTask.id}</Descriptions.Item>
            <Descriptions.Item label={t('crawlDataSource')}>{sourceMap[detailTask.source_id] || detailTask.source_id}</Descriptions.Item>
            <Descriptions.Item label={t('status')}><Tag color={TASK_STATUS_COLORS[detailTask.status]}>{t(TASK_STATUS_KEYS[detailTask.status] || detailTask.status)}</Tag></Descriptions.Item>
            <Descriptions.Item label={t('crawlTotal')}>{detailTask.total}</Descriptions.Item>
            <Descriptions.Item label={t('success')}>{detailTask.ingested}</Descriptions.Item>
            <Descriptions.Item label={t('failed')}>{detailTask.failed}</Descriptions.Item>
            {detailTask.error && <Descriptions.Item label={t('crawlError')}><span style={{ color: 'var(--red)' }}>{detailTask.error}</span></Descriptions.Item>}
            <Descriptions.Item label={t('crawlStartTime')}>{formatTime(detailTask.started_at, i18n.language)}</Descriptions.Item>
            <Descriptions.Item label={t('crawlEndTime')}>{formatTime(detailTask.finished_at, i18n.language)}</Descriptions.Item>
            <Descriptions.Item label={t('crawlDuration')}>{formatDuration(detailTask.started_at, detailTask.finished_at, t)}</Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>

      <Modal title={t('crawlManualCreate')} open={createOpen} onOk={handleCreate} onCancel={() => setCreateOpen(false)} destroyOnClose>
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="source_id" label={t('crawlDataSource')} rules={[{ required: true }]}>
            <Select placeholder={t('crawlSelectSource')} options={sources.map(s => ({ value: s.id, label: s.name }))} />
          </Form.Item>
          <Form.Item name="limit" label={t('crawlLimit')} initialValue={200}><InputNumber min={1} max={10000} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="with_body" label={t('crawlFetchBody')} valuePropName="checked"><Switch /></Form.Item>
          <Form.Item name="mode" label={t('crawlMode')} initialValue="list">
            <Select options={[{ value: 'list', label: t('crawlListMode') }, { value: 'detail', label: t('crawlDetailMode') }]} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
