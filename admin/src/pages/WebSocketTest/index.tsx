import { useState, useCallback, useRef, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Card, Tag, Button, Space, Typography, Empty, Badge, Switch, Form, Input, Select, Row, Col, Statistic } from 'antd'
import {
  ClearOutlined,
  WifiOutlined,
  DisconnectOutlined,
  SendOutlined,
} from '@ant-design/icons'
import {
  useWebSocketStatus,
  useWebSocketPush,
  type WSPushMessage,
} from '@/hooks/useWebSocket'
import request from '@/utils/request'
import { message } from 'antd'

const { Text } = Typography

interface LogEntry {
  id: number
  timestamp: Date
  message: WSPushMessage
}

interface TestFormValues {
  recipient_id: string
  category: string
  title: string
  content: string
}

export default function WebSocketTest() {
  const { t } = useTranslation()
  const connected = useWebSocketStatus()
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [autoScroll, setAutoScroll] = useState(true)
  const [sending, setSending] = useState(false)
  const [form] = Form.useForm<TestFormValues>()
  const logIdRef = useRef(0)
  const logContainerRef = useRef<HTMLDivElement>(null)

  // 分类选项（使用 i18n）
  const categoryOptions = [
    { label: t('wsCatPoints'), value: 'points' },
    { label: t('wsCatGoal'), value: 'goal' },
    { label: t('wsCatVerification'), value: 'verification' },
    { label: t('wsCatCompanionStatus'), value: 'companion_status' },
    { label: t('wsCatExchange'), value: 'exchange' },
    { label: t('wsCatCompanionEncourage'), value: 'companion_encourage' },
    { label: t('wsCatCompanionRemind'), value: 'companion_remind' },
    { label: t('wsCatCompanionCelebrate'), value: 'companion_celebrate' },
  ]

  const categoryMap = Object.fromEntries(categoryOptions.map(o => [o.value, o.label]))
  const getCategoryLabel = (value: string) => categoryMap[value] || value

  // 收到推送时记录日志
  useWebSocketPush((data) => {
    logIdRef.current += 1
    setLogs((prev) => [
      ...prev,
      { id: logIdRef.current, timestamp: new Date(), message: data },
    ])
  })

  // 自动滚动到底部
  useEffect(() => {
    if (autoScroll && logContainerRef.current) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight
    }
  }, [logs, autoScroll])

  const handleClear = useCallback(() => {
    setLogs([])
  }, [])

  const handleSendTest = useCallback(async (values: TestFormValues) => {
    setSending(true)
    try {
      await request.post('/v1/admin/messages/system', values)
      message.success(t('wsSendSuccess'))
    } catch {
      message.error(t('wsSendFailed'))
    } finally {
      setSending(false)
    }
  }, [t])

  return (
    <div>
      {/* ===== 页面 Header ===== */}
      <div
        style={{
          padding: '20px 28px 0',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <h2
          style={{
            fontSize: 20,
            fontWeight: 600,
            color: 'var(--text-primary)',
            lineHeight: 1.3,
            margin: 0,
            display: 'flex',
            alignItems: 'center',
            gap: 10,
          }}
        >
          <WifiOutlined style={{ fontSize: 18 }} />
          <span>{t('wsTestTitle')}</span>
          <Badge
            status={connected ? 'success' : 'error'}
            text={
              <Text type={connected ? 'success' : 'danger'} style={{ fontSize: 14 }}>
                {connected ? t('wsConnected') : t('wsDisconnected')}
              </Text>
            }
          />
        </h2>
        <Space>
          <Switch
            checkedChildren={t('wsAutoScroll')}
            unCheckedChildren={t('wsAutoScroll')}
            checked={autoScroll}
            onChange={setAutoScroll}
          />
          <Button icon={<ClearOutlined />} onClick={handleClear}>
            {t('wsClearLog')}
          </Button>
        </Space>
      </div>

      {/* ===== 标题与内容间距 ===== */}
      <div style={{ height: 20 }} />

      {/* ===== 内容区域 ===== */}
      <div style={{ padding: '0 28px 20px' }}>
        <div
          style={{
            border: '1px solid var(--border-light)',
            borderRadius: 'var(--radius-md)',
            background: 'var(--bg-white)',
            overflow: 'hidden',
            padding: 20,
          }}
        >
          <Row gutter={[20, 20]}>
            {/* 左侧：发送表单 + 统计 */}
            <Col xs={24} lg={8}>
              <Card
                title={
                  <Space>
                    <SendOutlined />
                    <span>{t('wsSendTest')}</span>
                  </Space>
                }
                styles={{ header: { padding: '16px 24px' }, body: { padding: '16px 24px' } }}
              >
                <Form
                  form={form}
                  layout="vertical"
                  onFinish={handleSendTest}
                  initialValues={{
                    recipient_id: 'user_founder',
                    category: 'points',
                    title: '',
                    content: '',
                  }}
                >
                  <Form.Item
                    name="recipient_id"
                    label={t('wsRecipient')}
                    rules={[{ required: true, message: t('required') }]}
                  >
                    <Input placeholder="user_founder" />
                  </Form.Item>
                  <Form.Item
                    name="category"
                    label={t('wsCategory')}
                    rules={[{ required: true }]}
                  >
                    <Select options={categoryOptions} />
                  </Form.Item>
                  <Form.Item
                    name="title"
                    label={t('wsTitle')}
                    rules={[{ required: true, message: t('required') }]}
                  >
                    <Input placeholder={t('wsTitle')} />
                  </Form.Item>
                  <Form.Item
                    name="content"
                    label={t('wsContent')}
                    rules={[{ required: true, message: t('required') }]}
                  >
                    <Input.TextArea rows={3} placeholder={t('wsContent')} />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      icon={<SendOutlined />}
                      loading={sending}
                      block
                    >
                      {t('wsSend')}
                    </Button>
                  </Form.Item>
                </Form>

                {/* 统计 */}
                <div style={{ display: 'flex', gap: 16, marginTop: 8 }}>
                  <Statistic
                    title={t('wsReceivedCount')}
                    value={logs.length}
                    suffix={t('wsCountUnit')}
                    valueStyle={{ fontSize: 20 }}
                  />
                  <Statistic
                    title={t('wsConnectionStatus')}
                    value={connected ? t('wsStatusNormal') : t('wsStatusBroken')}
                    valueStyle={{
                      fontSize: 20,
                      color: connected ? 'var(--success-color, #52c41a)' : 'var(--error-color, #ff4d4f)',
                    }}
                  />
                </div>
              </Card>
            </Col>

            {/* 右侧：日志区域 */}
            <Col xs={24} lg={16}>
              <div
                ref={logContainerRef}
                style={{
                  background: 'var(--bg-color-secondary, #f5f5f5)',
                  borderRadius: 8,
                  padding: 16,
                  minHeight: 420,
                  maxHeight: 620,
                  overflow: 'auto',
                  fontFamily: 'JetBrains Mono, monospace',
                  fontSize: 13,
                }}
              >
                {logs.length === 0 ? (
                  <Empty
                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                    description={
                      connected
                        ? t('wsWaiting')
                        : t('wsNotConnected')
                    }
                    style={{ padding: '100px 0' }}
                  >
                    {!connected && (
                      <Space>
                        <DisconnectOutlined style={{ color: 'var(--error-color, #ff4d4f)', fontSize: 24 }} />
                        <Text type="danger">{t('wsConnectionBroken')}</Text>
                      </Space>
                    )}
                  </Empty>
                ) : (
                  logs.map((log) => (
                    <div
                      key={log.id}
                      style={{
                        padding: '12px 16px',
                        marginBottom: 12,
                        background: 'var(--bg-color, #fff)',
                        borderRadius: 8,
                        borderLeft: '3px solid var(--primary-color, #1677ff)',
                        boxShadow: '0 1px 2px rgba(0,0,0,0.04)',
                      }}
                    >
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                        <Space size={4}>
                          <Tag color="blue">{log.message.type}</Tag>
                          <Tag>{getCategoryLabel(log.message.category)}</Tag>
                        </Space>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          {log.timestamp.toLocaleTimeString()}
                        </Text>
                      </div>
                      <div style={{ fontWeight: 600, fontSize: 14, marginBottom: 4 }}>
                        {log.message.title}
                      </div>
                      <div style={{ color: 'var(--text-secondary, #666)', fontSize: 13 }}>
                        {log.message.content}
                      </div>
                      <details style={{ marginTop: 12 }}>
                        <summary style={{ cursor: 'pointer', fontSize: 12, color: 'var(--text-secondary, #999)', userSelect: 'none' }}>
                          {t('wsRawJson')}
                        </summary>
                        <pre
                          style={{
                            marginTop: 8,
                            padding: 12,
                            background: 'var(--bg-color-secondary, #fafafa)',
                            borderRadius: 6,
                            fontSize: 11,
                            overflow: 'auto',
                            border: '1px solid var(--border-color, #e8e8e8)',
                          }}
                        >
                          {JSON.stringify(log.message, null, 2)}
                        </pre>
                      </details>
                    </div>
                  ))
                )}
              </div>
            </Col>
          </Row>
        </div>
      </div>
    </div>
  )
}
