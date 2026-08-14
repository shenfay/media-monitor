import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Button, Space, Switch, Form, Input, InputNumber, Select, Popconfirm, message, Drawer, Descriptions, Row, Col, Divider, Typography } from 'antd'
import { PlusOutlined, ReloadOutlined, PlayCircleOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons'
import DataPanel from '@/components/DataPanel'
import { getSources, createSource, updateSource, deleteSource, runSource, type Source } from '@/services/crawl'

const { TextArea } = Input

function formatTime(dateStr: string | null): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
}

export default function SourceManagement() {
  const { t } = useTranslation()
  const [sources, setSources] = useState<Source[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Source | null>(null)
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailSource, setDetailSource] = useState<Source | null>(null)
  const [form] = Form.useForm()

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const data = await getSources()
      setSources(data)
    } catch { /* ignore */ }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { fetchData() }, [fetchData])

  const handleAdd = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ enabled: true, months: 6, platform_type: 'news' })
    setModalOpen(true)
  }

  const handleEdit = (record: Source) => {
    setEditing(record)
    form.setFieldsValue({
      ...record,
      nodes: record.nodes || [],
      tags: record.tags || [],
      auth: record.auth ? JSON.stringify(record.auth, null, 2) : '',
    })
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      let auth = undefined
      if (values.auth) {
        try { auth = JSON.parse(values.auth) } catch { message.error(t('crawlAuthJsonError')); return }
      }
      const payload = { ...values, auth }
      if (editing) {
        await updateSource(editing.id, payload)
        message.success(t('updateSuccess'))
      } else {
        await createSource(payload)
        message.success(t('createSuccess'))
      }
      setModalOpen(false)
      form.resetFields()
      fetchData()
    } catch { /* validation error */ }
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteSource(id)
      message.success(t('deleteSuccess'))
      fetchData()
    } catch { message.error(t('crawlDeleteFailed')) }
  }

  const handleRun = async (id: string) => {
    try {
      const res = await runSource(id)
      message.success(`${t('crawlTaskCreated')}: ${res.task_id}`)
    } catch { message.error(t('crawlRunFailed')) }
  }

  const handleToggleEnabled = async (record: Source, checked: boolean) => {
    try {
      await updateSource(record.id, { enabled: checked })
      message.success(checked ? t('crawlEnabledMsg') : t('crawlDisabledMsg'))
      fetchData()
    } catch { message.error(t('crawlOpFailed')) }
  }

  const columns = [
    { title: t('name'), dataIndex: 'name', key: 'name', width: 160, ellipsis: true },
    { title: t('crawlPlatform'), dataIndex: 'platform_type', key: 'platform_type', width: 80, render: (v: string) => <Tag>{v || t('crawlNews')}</Tag> },
    {
      title: t('crawlNodes'), dataIndex: 'nodes', key: 'nodes', width: 160,
      render: (v: string[]) => v?.length ? v.map(n => <Tag key={n}>{n}</Tag>) : <span style={{ color: 'var(--text-quaternary)' }}>-</span>,
    },
    {
      title: t('crawlTags'), dataIndex: 'tags', key: 'tags', width: 140,
      render: (v: string[]) => v?.length ? v.map(tag => <Tag key={tag} color="blue">{tag}</Tag>) : <span style={{ color: 'var(--text-quaternary)' }}>{t('crawlUniversal')}</span>,
    },
    { title: t('crawlSchedule'), dataIndex: 'schedule', key: 'schedule', width: 100, render: (v: string) => v || <span style={{ color: 'var(--text-quaternary)' }}>{t('crawlManual')}</span> },
    {
      title: t('crawlEnabled'), dataIndex: 'enabled', key: 'enabled', width: 60,
      render: (_: boolean, record: Source) => <Switch size="small" checked={record.enabled} onChange={c => handleToggleEnabled(record, c)} />,
    },
    { title: t('crawlLastCrawl'), dataIndex: 'last_crawl_at', key: 'last_crawl_at', width: 150, render: (v: string | null) => formatTime(v) },
    {
      title: t('actions'), key: 'actions', width: 220,
      render: (_: unknown, record: Source) => (
        <Space size={4}>
          <Button type="link" size="small" icon={<EyeOutlined />} onClick={() => { setDetailSource(record); setDetailOpen(true) }}>
            {t('crawlView')}
          </Button>
          <Button type="link" size="small" icon={<PlayCircleOutlined />} onClick={() => handleRun(record.id)}>
            {t('crawlRun')}
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            {t('edit')}
          </Button>
          <Popconfirm title={t('crawlConfirmDelete')} onConfirm={() => handleDelete(record.id)}>
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              {t('delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('crawlSourceMgmt')}
        extra={
          <Button icon={<PlusOutlined />} type="primary" onClick={handleAdd}>{t('crawlCreateSource')}</Button>
        }
        toolbarActions={
          <Button icon={<ReloadOutlined />} onClick={fetchData}>{t('refresh')}</Button>
        }
      >
        <Table dataSource={sources} columns={columns} rowKey="id" loading={loading} pagination={false} size="small" />
      </DataPanel>

      <Drawer
        title={editing ? t('crawlEditSource') : t('crawlCreateSource')}
        open={modalOpen}
        onClose={() => { setModalOpen(false); form.resetFields() }}
        width={520}
        destroyOnClose
        extra={
          <Space>
            <Button onClick={() => { setModalOpen(false); form.resetFields() }}>{t('cancel')}</Button>
            <Button type="primary" onClick={handleSubmit}>{t('submit')}</Button>
          </Space>
        }
      >
        <Form form={form} layout="vertical" style={{ marginTop: 4 }}>
          <Typography.Text strong style={{ fontSize: 13, marginBottom: 12, display: 'block' }}>{t('crawlBasicInfo')}</Typography.Text>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="name" label={t('name')} rules={[{ required: true }]}><Input /></Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="platform_type" label={t('crawlPlatform')}>
                <Select options={[{ value: 'news', label: t('crawlNews') }, { value: 'social', label: t('crawlSocial') }, { value: 'social_overseas', label: t('crawlSocialOverseas') }]} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="enabled" label={t('crawlEnabled')} valuePropName="checked"><Switch /></Form.Item>

          <Divider style={{ margin: '16px 0 12px' }} />
          <Typography.Text strong style={{ fontSize: 13, marginBottom: 12, display: 'block' }}>{t('crawlFetchConfig')}</Typography.Text>
          <Form.Item name="base_url" label={t('crawlBaseUrl')}><Input placeholder="https://example.com" /></Form.Item>
          <Form.Item name="list_endpoint" label={t('crawlListEndpoint')}><Input placeholder="/api/list" /></Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="source_filter" label={t('crawlSourceFilter')}><Input placeholder={t('crawlSourceFilter')} /></Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="months" label={t('crawlMonths')}><InputNumber min={1} max={120} style={{ width: '100%' }} /></Form.Item>
            </Col>
          </Row>
          <Form.Item name="schedule" label={t('crawlCronSchedule')}><Input placeholder="0 */6 * * *（...）" /></Form.Item>

          <Divider style={{ margin: '16px 0 12px' }} />
          <Typography.Text strong style={{ fontSize: 13, marginBottom: 12, display: 'block' }}>{t('crawlAuth')}</Typography.Text>
          <Form.Item name="nodes" label={t('crawlNodes')}><Select mode="tags" placeholder={t('crawlInputEnter')} /></Form.Item>
          <Form.Item name="tags" label={t('crawlCapabilityTags')}><Select mode="tags" placeholder={t('crawlInputEnter')} /></Form.Item>
          <Form.Item name="auth" label={t('crawlAuth')}><TextArea rows={3} placeholder='{"cookie": "..."}' /></Form.Item>
        </Form>
      </Drawer>

      <Drawer title={t('crawlSourceDetail')} open={detailOpen} onClose={() => setDetailOpen(false)} width={480}>
        {detailSource && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="ID">{detailSource.id}</Descriptions.Item>
            <Descriptions.Item label={t('name')}>{detailSource.name}</Descriptions.Item>
            <Descriptions.Item label={t('crawlPlatform')}>{detailSource.platform_type}</Descriptions.Item>
            <Descriptions.Item label={t('crawlBaseUrl')}>{detailSource.base_url || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('crawlListEndpoint')}>{detailSource.list_endpoint || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('crawlNodes')}>{detailSource.nodes?.join(', ') || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('crawlSourceFilter')}>{detailSource.source_filter || '-'}</Descriptions.Item>
            <Descriptions.Item label={t('crawlMonths')}>{detailSource.months}</Descriptions.Item>
            <Descriptions.Item label={t('crawlSchedule')}>{detailSource.schedule || t('crawlManual')}</Descriptions.Item>
            <Descriptions.Item label={t('crawlTags')}>{detailSource.tags?.join(', ') || t('crawlUniversal')}</Descriptions.Item>
            <Descriptions.Item label={t('crawlEnabled')}>{detailSource.enabled ? t('yes') : t('no')}</Descriptions.Item>
            <Descriptions.Item label={t('crawlLastCrawl')}>{formatTime(detailSource.last_crawl_at)}</Descriptions.Item>
            <Descriptions.Item label={t('createdAt')}>{formatTime(detailSource.created_at)}</Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>
    </div>
  )
}
