import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Form, Input, InputNumber, Switch, Button, Select, Tabs, Table, Tag,
  Modal, message, Spin, Space,
} from 'antd'
import { EditOutlined, CheckOutlined, FormatPainterOutlined } from '@ant-design/icons'
import type { TabsProps } from 'antd'
import CodeMirror from '@uiw/react-codemirror'
import { json } from '@codemirror/lang-json'
import DataPanel from '@/components/DataPanel'
import {
  getSettings, batchUpdateSettings,
  type SettingItem, type SettingUpdateItem,
} from '@/services/setting'

// ---- 类型 ----
type SettingsMap = Record<string, unknown>

// ---- 通知事件定义 ----
interface NotifyEvent {
  key: string
  label: string
  enabled: boolean
}

// ---- 渠道配置 ----
interface ChannelConfig {
  name: string
  key: string
  configured: boolean
}

export default function SystemSettings() {
  const { t } = useTranslation()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [settingsData, setSettingsData] = useState<SettingsMap>({})
  const [rawSettings, setRawSettings] = useState<SettingItem[]>([])

  // 通知设置状态
  const [notifyEvents, setNotifyEvents] = useState<NotifyEvent[]>([])
  const [channelModalOpen, setChannelModalOpen] = useState(false)
  const [editingChannel, setEditingChannel] = useState<ChannelConfig | null>(null)
  const [channelForm] = Form.useForm()

  // 业务规则 JSON 编辑器状态
  const [streakJson, setStreakJson] = useState('')

  // ---- 加载数据 ----
  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const items = await getSettings()
      setRawSettings(items)

      const map: SettingsMap = {}
      items.forEach((item) => {
        map[item.key] = item.value
      })
      setSettingsData(map)

      // 填充表单
      const formValues: Record<string, unknown> = {}
      items.forEach((item) => {
        formValues[item.key] = item.value
      })
      form.setFieldsValue(formValues)

      // 初始化通知事件
      setNotifyEvents([
        { key: 'notify_acceptance', label: t('notifyAcceptance'), enabled: map['notify_acceptance'] as boolean ?? true },
        { key: 'notify_goal_progress', label: t('notifyGoalProgress'), enabled: map['notify_goal_progress'] as boolean ?? true },
        { key: 'notify_companion_status', label: t('notifyCompanionStatus'), enabled: map['notify_companion_status'] as boolean ?? true },
        { key: 'notify_streak_inactive', label: t('notifyStreakInactive'), enabled: map['notify_streak_inactive'] as boolean ?? true },
      ])

      // 初始化连续加成规则 JSON
      const streakRules = map['streak_bonus_rules']
      if (streakRules && typeof streakRules === 'object') {
        setStreakJson(JSON.stringify(streakRules, null, 2))
      } else if (typeof streakRules === 'string') {
        try {
          setStreakJson(JSON.stringify(JSON.parse(streakRules), null, 2))
        } catch {
          setStreakJson('{}')
        }
      }
    } catch {
      message.error(t('settingsLoadFailed'))
    } finally {
      setLoading(false)
    }
  }, [form])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  // ---- 保存设置 ----
  const handleSave = async () => {
    try {
      const values = await form.validateFields()
      setSaving(true)

      const updates: SettingUpdateItem[] = []

      // 收集所有表单字段
      Object.entries(values).forEach(([key, value]) => {
        if (key === 'streak_bonus_rules') return // JSON 字段单独处理
        updates.push({ key, value })
      })

      // 处理连续加成规则 JSON
      if (streakJson) {
        try {
          const parsed = JSON.parse(streakJson)
          updates.push({ key: 'streak_bonus_rules', value: parsed })
        } catch {
          message.error(t('streakRuleJsonError'))
          setSaving(false)
          return
        }
      }

      // 收集通知事件开关
      notifyEvents.forEach((event) => {
        updates.push({ key: event.key, value: event.enabled })
      })

      await batchUpdateSettings(updates)
      message.success(t('saveSuccess'))
      await fetchData()
    } catch {
      // validation failed
    } finally {
      setSaving(false)
    }
  }

  // ---- 通知事件切换 ----
  const handleNotifyToggle = (key: string, enabled: boolean) => {
    setNotifyEvents((prev) =>
      prev.map((e) => (e.key === key ? { ...e, enabled } : e))
    )
  }

  // ---- 渠道配置 ----
  const channels: ChannelConfig[] = [
    { name: t('emailChannel'), key: 'channel_email', configured: hasChannelConfig('channel_email') },
    { name: t('webhookChannel'), key: 'channel_webhook', configured: hasChannelConfig('channel_webhook') },
  ]

  function hasChannelConfig(key: string): boolean {
    const val = settingsData[key]
    if (!val || typeof val !== 'object') return false
    const obj = val as Record<string, unknown>
    // 检查是否有非空字段
    return Object.values(obj).some((v) => {
      if (typeof v === 'string') return v !== '' && v !== '••••••'
      if (typeof v === 'number') return v !== 0
      return v !== null && v !== undefined
    })
  }

  const handleEditChannel = (channel: ChannelConfig) => {
    setEditingChannel(channel)
    const val = settingsData[channel.key]
    if (val && typeof val === 'object') {
      channelForm.setFieldsValue(val as Record<string, unknown>)
    } else {
      channelForm.resetFields()
    }
    setChannelModalOpen(true)
  }

  const handleChannelSave = async () => {
    if (!editingChannel) return
    try {
      const values = await channelForm.validateFields()
      setSaving(true)

      // 密码/secret 字段：如果是占位符则不提交（保持原值）
      const obj = { ...values }
      if (obj.password === '••••••') delete obj.password
      if (obj.secret === '••••••') delete obj.secret

      await batchUpdateSettings([{ key: editingChannel.key, value: obj }])
      message.success(t('channelConfigSaved'))
      setChannelModalOpen(false)
      fetchData()
    } catch {
      // validation failed
    } finally {
      setSaving(false)
    }
  }

  // ---- JSON 格式化 ----
  const handleFormatJson = () => {
    try {
      const parsed = JSON.parse(streakJson)
      setStreakJson(JSON.stringify(parsed, null, 2))
    } catch {
      message.error(t('jsonFormatError'))
    }
  }

  // ---- Tab 内容 ----
  const tabItems: TabsProps['items'] = [
    {
      key: 'basic',
      label: t('basicConfig'),
      children: (
        <>
          <div style={{ display: 'flex', gap: 24 }}>
            <Form.Item label={t('siteName')} name="site_name" style={{ flex: 1 }}
              rules={[{ required: true, message: t('siteNamePlaceholder') }]}>
              <Input placeholder={t('siteNamePlaceholder')} />
            </Form.Item>
            <Form.Item label={t('siteLogo')} name="logo_url" style={{ flex: 1 }}>
              <Input placeholder={t('logoPlaceholder')} />
            </Form.Item>
          </div>
          <div style={{ display: 'flex', gap: 24 }}>
            <Form.Item label={t('defaultLanguage')} name="default_language" style={{ flex: 1 }}>
              <Select>
                <Select.Option value="zh-CN">{t('zhCN')}</Select.Option>
                <Select.Option value="en-US">English</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item label={t('timezone')} name="timezone" style={{ flex: 1 }}>
              <Select>
                <Select.Option value="Asia/Shanghai">Asia/Shanghai</Select.Option>
                <Select.Option value="UTC">UTC</Select.Option>
              </Select>
            </Form.Item>
          </div>
          <Form.Item label={t('sessionTimeout')} name="session_timeout" style={{ maxWidth: 400 }}
            rules={[{ required: true, message: t('sessionTimeoutPlaceholder') }]}>
            <InputNumber min={5} max={120} addonAfter={t('minutes')} style={{ width: '100%' }} />
          </Form.Item>
        </>
      ),
    },
    {
      key: 'toggle',
      label: t('featureToggles'),
      children: (
        <>
          <Form.Item label={t('openRegistration')} name="enable_register" valuePropName="checked">
            <Switch checkedChildren={t('turnedOn')} unCheckedChildren={t('turnedOff')} />
          </Form.Item>
          <Form.Item label={t('auditLog')} name="enable_audit_log" valuePropName="checked">
            <Switch checkedChildren={t('turnedOn')} unCheckedChildren={t('turnedOff')} />
          </Form.Item>
          <Form.Item label={t('messageNotification')} name="enable_notification" valuePropName="checked">
            <Switch checkedChildren={t('turnedOn')} unCheckedChildren={t('turnedOff')} />
          </Form.Item>
        </>
      ),
    },
    {
      key: 'business',
      label: t('businessRules'),
      children: (
        <>
          <div style={{ display: 'flex', gap: 24 }}>
            <Form.Item label={t('maxChildrenPerFamily')} name="max_children_per_family" style={{ flex: 1 }}
              rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <InputNumber min={1} max={10} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item label={t('maxDailyCards')} name="max_daily_cards" style={{ flex: 1 }}
              rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <InputNumber min={1} max={10} style={{ width: '100%' }} />
            </Form.Item>
          </div>
          <div style={{ display: 'flex', gap: 24 }}>
            <Form.Item label={t('maxGoalCards')} name="max_goal_cards" style={{ flex: 1 }}
              rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <InputNumber min={1} max={10} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item label={t('xpLevelDivisor')} name="xp_level_divisor" style={{ flex: 1 }}
              rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <InputNumber min={1} style={{ width: '100%' }} />
            </Form.Item>
          </div>
          <div style={{ display: 'flex', gap: 24 }}>
            <Form.Item label={t('defaultCompanionName')} name="default_companion_name" style={{ flex: 1 }}
              rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <Input placeholder={t('companionExample')} />
            </Form.Item>
            <Form.Item label={t('defaultGoalReward')} name="default_goal_reward" style={{ flex: 1 }}
              rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <InputNumber min={1} addonAfter={t('pointsUnit2')} style={{ width: '100%' }} />
            </Form.Item>
          </div>

          {/* 连续加成规则 - CodeMirror 编辑器 */}
          <Form.Item label={t('streakBonusRules')}>
            <div style={{ border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', overflow: 'hidden' }}>
              <CodeMirror
                value={streakJson}
                height="150px"
                extensions={[json()]}
                onChange={(val) => setStreakJson(val)}
                basicSetup={{ lineNumbers: true, foldGutter: false }}
              />
            </div>
            <Space style={{ marginTop: 8 }}>
              <Button size="small" icon={<FormatPainterOutlined />} onClick={handleFormatJson}>
                {t('formatJson')}
              </Button>
              <Button size="small" icon={<CheckOutlined />} onClick={() => {
                try { JSON.parse(streakJson); message.success(t('jsonFormatCorrect')) }
                catch { message.error(t('jsonFormatErrorShort')) }
              }}>
                {t('validate')}
              </Button>
            </Space>
          </Form.Item>
        </>
      ),
    },
    {
      key: 'notification',
      label: t('notificationSettings'),
      children: (
        <>
          <div style={{ marginBottom: 16, fontWeight: 500 }}>{t('eventSwitches')}</div>
          <Table<NotifyEvent>
            columns={[
              { title: t('notifyEvent'), dataIndex: 'label', key: 'label' },
              {
                title: t('status'), dataIndex: 'enabled', key: 'enabled',
                render: (enabled: boolean) => (
                  <Tag style={{
                    background: enabled ? 'var(--green-light)' : 'var(--gray-light)',
                    color: enabled ? 'var(--green-text)' : 'var(--gray-text)',
                    border: 'none', borderRadius: 'var(--radius-sm)',
                  }}>{enabled ? t('turnedOn') : t('turnedOff')}</Tag>
                ),
              },
              {
                title: t('actions'), key: 'action', width: 100,
                render: (_: unknown, record: NotifyEvent) => (
                  <Switch
                    size="small"
                    checked={record.enabled}
                    onChange={(checked) => handleNotifyToggle(record.key, checked)}
                  />
                ),
              },
            ]}
            dataSource={notifyEvents}
            rowKey="key"
            pagination={false}
            size="small"
            style={{ marginBottom: 24 }}
          />

          <div style={{ marginBottom: 16, fontWeight: 500 }}>{t('channelConfig')}</div>
          <Table<ChannelConfig>
            columns={[
              { title: t('channelName'), dataIndex: 'name', key: 'name' },
              {
                title: t('configStatus'), dataIndex: 'configured', key: 'configured',
                render: (configured: boolean) => (
                  <Tag style={{
                    background: configured ? 'var(--green-light)' : 'var(--gray-light)',
                    color: configured ? 'var(--green-text)' : 'var(--gray-text)',
                    border: 'none', borderRadius: 'var(--radius-sm)',
                  }}>{configured ? t('configured') : t('notConfigured')}</Tag>
                ),
              },
              {
                title: t('actions'), key: 'action', width: 100,
                render: (_: unknown, record: ChannelConfig) => (
                  <Button type="link" size="small" icon={<EditOutlined />}
                    onClick={() => handleEditChannel(record)}>
                    {t('edit')}
                  </Button>
                ),
              },
            ]}
            dataSource={channels}
            rowKey="key"
            pagination={false}
            size="small"
          />
        </>
      ),
    },
  ]

  // ---- 渠道编辑 Modal ----
  const isEmail = editingChannel?.key === 'channel_email'

  const channelModalNode = (
    <Modal
      title={t('editChannel', { name: editingChannel?.name ?? '' })}
      open={channelModalOpen}
      onOk={handleChannelSave}
      onCancel={() => setChannelModalOpen(false)}
      confirmLoading={saving}
      destroyOnClose
    >
      <Form form={channelForm} layout="vertical" preserve={false}>
        {isEmail ? (
          <>
            <Form.Item label={t('smtpHost')} name="host" rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <Input placeholder={t('smtpHostExample')} />
            </Form.Item>
            <Form.Item label={t('port')} name="port" rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <InputNumber min={1} max={65535} style={{ width: '100%' }} placeholder={t('portExample')} />
            </Form.Item>
            <Form.Item label={t('username')} name="user" rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <Input placeholder={t('smtpUser')} />
            </Form.Item>
            <Form.Item label={t('password')} name="password">
              <Input.Password placeholder="••••••" />
            </Form.Item>
            <Form.Item label={t('sender')} name="from" rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <Input placeholder={t('senderExample')} />
            </Form.Item>
          </>
        ) : (
          <>
            <Form.Item label="Webhook URL" name="url" rules={[{ required: true, message: t('inputPlaceholderShort') }]}>
              <Input placeholder="https://..." />
            </Form.Item>
            <Form.Item label={t('signingKey')} name="secret">
              <Input.Password placeholder="••••••" />
            </Form.Item>
          </>
        )}
      </Form>
    </Modal>
  )

  if (loading) {
    return (
      <DataPanel title={t('systemSettings')}>
        <div style={{ display: 'flex', justifyContent: 'center', padding: 60 }}>
          <Spin size="large" />
        </div>
      </DataPanel>
    )
  }

  return (
    <DataPanel title={t('systemSettings')}>
      <div style={{ padding: '0 28px 20px' }}>
        <Form form={form} layout="vertical">
          <Tabs items={tabItems} />

          <div style={{ display: 'flex', gap: 12, marginTop: 8 }}>
            <Button type="primary" loading={saving} onClick={handleSave}>
              {t('saveSettingsBtn')}
            </Button>
            <Button onClick={() => fetchData()}>
              {t('reset')}
            </Button>
          </div>
        </Form>

        {channelModalNode}
      </div>
    </DataPanel>
  )
}
