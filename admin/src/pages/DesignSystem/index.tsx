import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { SearchOutlined, ExportOutlined, EditOutlined, SettingOutlined, DeleteOutlined } from '@ant-design/icons'
import {
  Tag, Button, Input, Select, Switch, Table, Typography, Space, Tabs,
  Progress, Alert, Badge, Avatar, Card, Descriptions, Empty, Skeleton,
  Tooltip, Popover, Breadcrumb, Radio, Checkbox, InputNumber, DatePicker,
  Rate, Steps, Divider, Row, Col, Modal, Statistic, Timeline, message,
} from 'antd'
import DataPanel, { FilterSearch } from '@/components/DataPanel'

const { Title, Text } = Typography

// ─── Shared ───────────────────────────────────

function SectionTitle({ children }: { children: string }) {
  return <Title level={4} style={{ margin: 0, marginBottom: 20, fontSize: 16 }}>{children}</Title>
}

function Swatch({ color, name, variable }: { color: string; name: string; variable: string }) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
      <div style={{
        width: 36, height: 36, borderRadius: 'var(--radius-md)',
        backgroundColor: color,
        border: color === '#ffffff' || color.startsWith('rgba') ? '1px solid var(--border-color)' : 'none',
        flexShrink: 0,
      }} />
      <div>
        <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-primary)' }}>{name}</div>
        <Text style={{ fontSize: 12, color: 'var(--text-muted)', fontFamily: 'monospace' }}>{variable} — {color}</Text>
      </div>
    </div>
  )
}

// ─── Tab 1: 设计令牌 ──────────────────────────

function DesignTokens() {
  const { t } = useTranslation()

  const colorGroups = [
    { title: t('brandColor'), items: [
      { name: t('brandDark'), variable: '--brand-dark', color: '#2b2b2b' },
      { name: t('brandDarkHover'), variable: '--brand-dark-hover', color: '#4d4d4d' },
    ]},
    { title: t('textColor'), items: [
      { name: t('textPrimary'), variable: '--text-primary', color: '#2b2b2b' },
      { name: t('textSecondary'), variable: '--text-secondary', color: '#6b6258' },
      { name: t('textMuted'), variable: '--text-muted', color: '#b0a89a' },
      { name: t('textIcon'), variable: '--text-icon', color: '#c4bdb0' },
    ]},
    { title: t('border'), items: [
      { name: t('borderDefault'), variable: '--border-color', color: '#e8e2d8' },
      { name: t('borderHover'), variable: '--border-hover', color: '#d4cdc0' },
      { name: t('borderLight'), variable: '--border-light', color: '#efeae2' },
      { name: t('divider'), variable: '--divider', color: '#f5f2ed' },
    ]},
    { title: t('bgColor'), items: [
      { name: t('bgWhite'), variable: '--bg-white', color: '#ffffff' },
      { name: t('bgSidebar'), variable: '--sidebar-bg', color: '#F5F3EF' },
      { name: t('bgHoverDark'), variable: '--hover-bg', color: '#f0ece6' },
      { name: t('bgHoverLight'), variable: '--hover-bg-light', color: '#f5f2ed' },
      { name: t('bgActive'), variable: '--active-bg', color: '#E4E0D8' },
    ]},
    { title: t('statusColor'), items: [
      { name: t('successActive'), variable: '--green', color: '#22c55e' },
      { name: t('failAbnormal'), variable: '--red', color: '#e74c3c' },
      { name: t('pendingWarn'), variable: '--yellow', color: '#f59e0b' },
      { name: t('info'), variable: '--blue', color: '#3b6fdf' },
      { name: t('defaultDisabled'), variable: '--gray', color: '#6b6258' },
    ]},
  ]

  const varCategories = [
    { title: t('layout'), items: [
      { name: '--sidebar-bg', value: '#F5F3EF', desc: t('sidebarBg') },
      { name: '--main-bg', value: '#FFFFFF', desc: t('mainBg') },
    ]},
    { title: t('font'), items: [
      { name: '--text-primary', value: '#2b2b2b', desc: t('textPrimary') },
      { name: '--text-secondary', value: '#6b6258', desc: t('textSecondary') },
      { name: '--text-muted', value: '#b0a89a', desc: t('textMuted') },
      { name: '--text-icon', value: '#c4bdb0', desc: t('textIcon') },
    ]},
    { title: t('borderAndBg'), items: [
      { name: '--border-color', value: '#e8e2d8', desc: t('borderDefault') },
      { name: '--border-hover', value: '#d4cdc0', desc: t('borderHover') },
      { name: '--divider', value: '#f5f2ed', desc: t('divider') },
      { name: '--hover-bg', value: '#f0ece6', desc: t('bgHoverDark') },
      { name: '--active-bg', value: '#E4E0D8', desc: t('bgActive') },
    ]},
    { title: t('brand'), items: [
      { name: '--brand-dark', value: '#2b2b2b', desc: 'Primary button / Active page' },
      { name: '--brand-dark-hover', value: '#4d4d4d', desc: 'Primary button hover' },
    ]},
    { title: t('statusColor'), items: [
      { name: '--green / --green-light / --green-text', value: '#22c55e / #dcfce7 / #166534', desc: t('successActive') },
      { name: '--red / --red-light / --red-text', value: '#e74c3c / #fef2f2 / #e74c3c', desc: t('failAbnormal') },
      { name: '--yellow / --yellow-light / --yellow-text', value: '#f59e0b / #fef3c7 / #92400e', desc: t('pendingWarn') },
      { name: '--blue / --blue-light / --blue-text', value: '#3b6fdf / #edf2ff / #3b6fdf', desc: t('info') },
      { name: '--gray / --gray-light / --gray-text', value: '#6b6258 / #f5f2ed / #b0a89a', desc: t('defaultDisabled') },
    ]},
    { title: t('table'), items: [
      { name: '--table-header-bg / --table-header-text', value: '#faf8f5 / #8a8276', desc: 'Header bg / text' },
      { name: '--table-border / --table-row-divider', value: '#efeae2 / #f5f2ed', desc: 'Border / row divider' },
    ]},
    { title: t('radiusShadowTransition'), items: [
      { name: '--radius-sm / --radius-md / --radius-lg', value: '8px / 12px / 16px', desc: t('smallMedLarge') },
      { name: '--shadow-sm / --shadow-md / --shadow-lg', value: '3 levels', desc: t('lightMedHeavy') },
      { name: '--transition', value: '0.25s cubic-bezier(0.4,0,0.2,1)', desc: t('defaultTransition') },
    ]},
  ]

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
      <div>
        <SectionTitle>{t('colorSystem')}</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
          {colorGroups.map(g => (
            <div key={g.title}>
              <Text style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-secondary)', marginBottom: 10, display: 'block' }}>{g.title}</Text>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: 20 }}>
                {g.items.map(i => <Swatch key={i.variable} {...i} />)}
              </div>
            </div>
          ))}
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('fontSpec')}</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          {[
            { label: t('enFont'), family: 'var(--font-family-en)', sample: 'The Quick Brown Fox Jumps Over The Lazy Dog', stack: 'Inter, -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica, Arial, sans-serif' },
            { label: t('zhFont'), family: 'var(--font-family-zh)', sample: t('zhFontSample'), stack: 'PingFang SC, Microsoft YaHei, Noto Sans SC, sans-serif' },
            { label: t('numFont'), family: 'var(--font-family-num)', sample: '0123456789 .,%', stack: 'Inter, JetBrains Mono, SF Mono, Roboto Mono, monospace' },
            { label: t('monoFont'), family: 'var(--font-family-mono)', sample: 'const hello = "Qoder" // 2024-07-09', stack: 'JetBrains Mono, SF Mono, Menlo, Roboto Mono, monospace' },
          ].map(f => (
            <div key={f.label} style={{ background: 'var(--bg-light)', borderRadius: 'var(--radius-md)', padding: '16px 20px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                <Text style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-secondary)' }}>{f.label}</Text>
                <Text style={{ fontSize: 11, color: 'var(--text-muted)', fontFamily: 'monospace' }}>{f.stack}</Text>
              </div>
              <div style={{ fontSize: 24, fontFamily: f.family, color: 'var(--text-primary)', lineHeight: 1.6 }}>{f.sample}</div>
            </div>
          ))}
          <div style={{ background: 'var(--bg-light)', borderRadius: 'var(--radius-md)', padding: '16px 20px' }}>
            <Text style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-secondary)', display: 'block', marginBottom: 8 }}>{t('mainFontStack')}</Text>
            <Text style={{ fontSize: 12, color: 'var(--text-muted)', fontFamily: 'monospace' }}>Inter, -apple-system, BlinkMacSystemFont, Segoe UI, PingFang SC, Microsoft YaHei, Noto Sans SC, system-ui, sans-serif, Apple Color Emoji, Segoe UI Emoji</Text>
            <div style={{ marginTop: 10, fontSize: 15, fontFamily: 'var(--font-family)', color: 'var(--text-primary)', lineHeight: 1.8 }}>
              {t('mixedSample')}
            </div>
          </div>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('typography')}</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          <Title style={{ fontSize: 20, fontWeight: 600, color: 'var(--text-primary)', margin: 0 }}>{t('pageTitleDesc')}</Title>
          <Title style={{ fontSize: 14, fontWeight: 600, margin: 0 }}>{t('cardTitleDesc')}</Title>
          <Text style={{ fontSize: 13, color: 'var(--text-primary)' }}>{t('tableBodyDesc')}</Text>
          <Text style={{ fontSize: 13, fontWeight: 500, color: 'var(--text-secondary)' }}>{t('formLabelDesc')}</Text>
          <Text style={{ fontSize: 12, fontWeight: 600, textTransform: 'uppercase', color: 'var(--table-header-text)' }}>{t('tableHeadDesc')}</Text>
          <Text style={{ fontSize: 12, fontWeight: 500, color: 'var(--text-muted)' }}>{t('auxTextDesc')}</Text>
          <Text style={{ fontSize: 11, color: 'var(--text-icon)' }}>{t('miniTextDesc')}</Text>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('radiusShadow')}</SectionTitle>
        <Row gutter={24}>
          <Col span={12}>
            <Card size="small" title={t('radius')} bordered={false} style={{ background: 'var(--bg-light)' }}>
              {[['--radius-sm', '8px', t('btnInputTag')], ['--radius-md', '12px', t('cardDataPanel')], ['--radius-lg', '16px', t('reserved')]].map(([v, val, desc]) => (
                <div key={v} style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
                  <div style={{ width: 40, height: 40, background: 'var(--brand-dark)', borderRadius: v === '--radius-sm' ? 'var(--radius-sm)' : v === '--radius-md' ? 'var(--radius-md)' : 'var(--radius-lg)', flexShrink: 0 }} />
                  <div><Text style={{ fontSize: 13, fontWeight: 500 }}>{v}: {val}</Text><br /><Text style={{ fontSize: 12, color: 'var(--text-muted)' }}>{desc}</Text></div>
                </div>
              ))}
            </Card>
          </Col>
          <Col span={12}>
            <Card size="small" title={t('shadow')} bordered={false} style={{ background: 'var(--bg-light)' }}>
              {[
                ['--shadow-sm', '0 1px 3px rgba(0,0,0,0.06)', t('standby')],
                ['--shadow-md', '0 4px 12px rgba(0,0,0,0.08)', 'StatCard hover'],
                ['--shadow-lg', '0 8px 24px rgba(0,0,0,0.12)', t('reserved')],
              ].map(([v, val, desc]) => (
                <div key={v} style={{ padding: '10px 14px', background: 'var(--bg-white)', borderRadius: 'var(--radius-md)', boxShadow: v === '--shadow-sm' ? 'var(--shadow-sm)' : v === '--shadow-md' ? 'var(--shadow-md)' : 'var(--shadow-lg)', marginBottom: 8 }}>
                  <Text style={{ fontSize: 13, fontWeight: 500 }}>{v}</Text><br />
                  <Text style={{ fontSize: 12, color: 'var(--text-muted)' }}>{val} — {desc}</Text>
                </div>
              ))}
            </Card>
          </Col>
        </Row>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>CSS Variables</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
          {varCategories.map(cat => (
            <div key={cat.title}>
              <Text style={{ fontSize: 13, fontWeight: 600, color: 'var(--text-secondary)', marginBottom: 8, display: 'block' }}>{cat.title}</Text>
              <Table
                columns={[
                  { title: t('varName'), dataIndex: 'name', key: 'name', width: 300 },
                  { title: t('value'), dataIndex: 'value', key: 'value', width: 280 },
                  { title: t('desc'), dataIndex: 'desc', key: 'desc' },
                ]}
                dataSource={cat.items.map((i, idx) => ({ ...i, key: idx }))}
                pagination={false} size="small" style={{ background: 'var(--bg-white)' }}
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

// ─── Tab 2: 组件库 ────────────────────────────

function Components() {
  const { t } = useTranslation()

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
      <div>
        <SectionTitle>{t('button')}</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <Space wrap>
            <Button type="primary">{t('primaryBtn')}</Button>
            <Button>{t('defaultBtn')}</Button>
            <Button type="text">{t('textBtn')}</Button>
            <Button type="dashed">{t('dashedBtn')}</Button>
            <Button type="link">{t('linkBtn')}</Button>
            <Button type="primary" danger>{t('dangerBtn')}</Button>
          </Space>
          <Space wrap>
            <Button type="primary" size="small">{t('smallPrimaryBtn')}</Button>
            <Button size="small">{t('smallDefaultBtn')}</Button>
            <Button type="primary" size="large">{t('largePrimaryBtn')}</Button>
            <Button size="large">{t('largeDefaultBtn')}</Button>
            <Button type="primary" shape="round">{t('roundPrimaryBtn')}</Button>
            <Button shape="round">{t('roundDefaultBtn')}</Button>
          </Space>
          <Space wrap>
            <Button type="primary" disabled>{t('disabledPrimaryBtn')}</Button>
            <Button disabled>{t('disabledDefaultBtn')}</Button>
            <Button type="primary" loading>{t('loadingBtn')}</Button>
          </Space>
        </div>
      </div>

      <div>
        <SectionTitle>{t('tableActionBtn')}</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          <div style={{ background: 'var(--bg-light)', padding: '14px 20px', borderRadius: 'var(--radius-md)' }}>
            <Space size={0} className="action-btn-group">
              <Button type="link" size="small" icon={<EditOutlined />}>{t('edit')}</Button>
              <Button type="link" size="small" icon={<SettingOutlined />}>{t('configPerm')}</Button>
              <Button type="link" size="small" danger icon={<DeleteOutlined />}>{t('delete')}</Button>
            </Space>
          </div>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('tag')}</SectionTitle>
        <Space wrap>
          {[
            { bg: 'var(--green-light)', color: 'var(--green-text)', label: t('pass') },
            { bg: 'var(--green-light)', color: 'var(--green-text)', label: t('active') },
            { bg: 'var(--red-light)', color: 'var(--red-text)', label: t('failed') },
            { bg: 'var(--red-light)', color: 'var(--red-text)', label: t('abnormal') },
            { bg: 'var(--yellow-light)', color: 'var(--yellow-text)', label: t('pending') },
            { bg: 'var(--blue-light)', color: 'var(--blue-text)', label: 'Admin' },
            { bg: 'var(--blue-light)', color: 'var(--blue-text)', label: t('edit') },
            { bg: 'var(--gray-light)', color: 'var(--gray-text)', label: t('unassigned') },
            { bg: 'var(--gray-light)', color: 'var(--gray-text)', label: 'Viewer' },
          ].map(item => (
            <Tag key={item.label} style={{ background: item.bg, color: item.color, border: 'none', borderRadius: 'var(--radius-sm)' }}>{item.label}</Tag>
          ))}
        </Space>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('formElement')}</SectionTitle>
        <div style={{ display: 'flex', gap: 24, flexWrap: 'wrap', alignItems: 'flex-start' }}>
          <div style={{ width: 200 }}>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('input')}</Text>
            <Input placeholder={t('inputPlaceholder')} />
          </div>
          <div style={{ width: 160 }}>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('select')}</Text>
            <Select placeholder={t('selectPlaceholder')} style={{ width: '100%' }} options={[{ value: '1', label: t('option1') }, { value: '2', label: t('option2') }]} />
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('numberInput')}</Text>
            <InputNumber placeholder={t('numberPlaceholder')} />
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('datePicker')}</Text>
            <DatePicker />
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('switch')}</Text>
            <Switch checkedChildren={t('on')} unCheckedChildren={t('off')} defaultChecked />
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('radio')}</Text>
            <Radio.Group defaultValue="a">
              <Radio value="a">{t('optionA')}</Radio>
              <Radio value="b">{t('optionB')}</Radio>
            </Radio.Group>
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('checkbox')}</Text>
            <Checkbox.Group options={['A', 'B', 'C']} defaultValue={['A']} />
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 6 }}>{t('rate')}</Text>
            <Rate defaultValue={3} />
          </div>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('progressBar')}</SectionTitle>
        <div style={{ display: 'flex', gap: 32, alignItems: 'center', flexWrap: 'wrap' }}>
          <div style={{ flex: 1, minWidth: 200 }}>
            <Progress percent={30} size="small" strokeColor="var(--brand-dark)" trailColor="var(--progress-bg)" />
            <Progress percent={60} size="small" strokeColor="var(--green)" trailColor="var(--progress-bg)" />
            <Progress percent={90} size="small" strokeColor="var(--blue)" trailColor="var(--progress-bg)" />
          </div>
          <Space size={16}>
            <Progress type="circle" percent={45} size={60} strokeColor="var(--brand-dark)" trailColor="var(--progress-bg)" />
            <Progress type="circle" percent={75} size={60} strokeColor="var(--green)" trailColor="var(--progress-bg)" />
            <Progress type="circle" percent={100} size={60} strokeColor="var(--blue)" trailColor="var(--progress-bg)" />
          </Space>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('badgeAvatar')}</SectionTitle>
        <Space size={24} align="start">
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 10 }}>Badge</Text>
            <Space size={16}>
              <Badge count={5}><Avatar shape="square" size={40} style={{ background: 'var(--gray-light)', color: 'var(--text-muted)' }}>U</Avatar></Badge>
              <Badge dot><Avatar shape="square" size={40} style={{ background: 'var(--gray-light)', color: 'var(--text-muted)' }}>U</Avatar></Badge>
              <Badge count={99} overflowCount={99}><Avatar shape="square" size={40} style={{ background: 'var(--gray-light)', color: 'var(--text-muted)' }}>U</Avatar></Badge>
            </Space>
          </div>
          <div>
            <Text style={{ fontSize: 13, fontWeight: 500, display: 'block', marginBottom: 10 }}>Avatar</Text>
            <Space size={12}>
              <Avatar size={32} style={{ background: 'linear-gradient(135deg, #4ECDC4, #44B09E)' }}>K</Avatar>
              <Avatar size={32} style={{ background: 'var(--blue)' }}>A</Avatar>
              <Avatar size={32} style={{ background: 'var(--green)' }}>B</Avatar>
              <Avatar size={32} style={{ background: 'var(--gray-light)', color: 'var(--text-muted)' }}>G</Avatar>
            </Space>
          </div>
        </Space>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('alert')}</SectionTitle>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <Alert message={t('successAlert')} type="success" showIcon />
          <Alert message={t('infoAlert')} type="info" showIcon />
          <Alert message={t('warningAlert')} type="warning" showIcon />
          <Alert message={t('errorAlert')} type="error" showIcon />
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('statusIndicator')}</SectionTitle>
        <Space size={24}>
          <Space>
            <span style={{ width: 7, height: 7, borderRadius: '50%', background: 'var(--green)', display: 'inline-block' }} />
            <Text style={{ fontSize: 13, color: 'var(--text-primary)' }}>{t('normalCertified')}</Text>
          </Space>
          <Space>
            <span style={{ width: 7, height: 7, borderRadius: '50%', background: 'var(--border-hover)', display: 'inline-block' }} />
            <Text style={{ fontSize: 13, color: 'var(--text-primary)' }}>{t('abnormalUnverified')}</Text>
          </Space>
        </Space>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('breadcrumb')}</SectionTitle>
        <Breadcrumb items={[
          { title: t('home') },
          { title: t('systemMgmt') },
          { title: t('userManagement'), href: '' },
        ]} />
      </div>
    </div>
  )
}

// ─── Tab 3: 数据展示 ──────────────────────────

function DataDisplay() {
  const { t } = useTranslation()

  const demoColumns = [
    { title: t('user'), dataIndex: 'name', key: 'name', width: 160 },
    { title: t('email'), dataIndex: 'email', key: 'email', width: 240 },
    { title: t('role'), dataIndex: 'role', key: 'role', width: 120 },
    { title: t('status'), dataIndex: 'status', key: 'status', width: 100 },
  ]
  const demoData = [
    { key: '1', name: 'Zhang San', email: 'zhangsan@example.com', role: <Tag style={{ background: 'var(--blue-light)', color: 'var(--blue-text)', border: 'none', borderRadius: 'var(--radius-sm)' }}>Admin</Tag>, status: <span style={{ color: 'var(--green)' }}>{t('active')}</span> },
    { key: '2', name: 'Li Si', email: 'lisi@example.com', role: <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)', border: 'none', borderRadius: 'var(--radius-sm)' }}>Viewer</Tag>, status: <span style={{ color: 'var(--green)' }}>{t('active')}</span> },
    { key: '3', name: 'Wang Wu', email: 'wangwu@example.com', role: <Tag style={{ background: 'var(--yellow-light)', color: 'var(--yellow-text)', border: 'none', borderRadius: 'var(--radius-sm)' }}>Pending</Tag>, status: <span style={{ color: 'var(--yellow)' }}>{t('pending')}</span> },
    { key: '4', name: 'Zhao Liu', email: 'zhaoliu@example.com', role: <Tag style={{ background: 'var(--red-light)', color: 'var(--red-text)', border: 'none', borderRadius: 'var(--radius-sm)' }}>{t('disabled')}</Tag>, status: <span style={{ color: 'var(--red)' }}>{t('abnormal')}</span> },
  ]

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
      <div>
        <SectionTitle>{t('table2')}</SectionTitle>
        <Table columns={demoColumns} dataSource={demoData} pagination={{ pageSize: 4, showSizeChanger: true, showQuickJumper: true, showTotal: (total: number) => t('total', { total }) }} size="middle" />
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('card2')}</SectionTitle>
        <Row gutter={16}>
          <Col span={8}>
            <Card title={t('defaultCard')} size="small"><Text>{t('cardBasicContent')}</Text></Card>
          </Col>
          <Col span={8}>
            <Card title={t('actionCard')} size="small" extra={<Button type="link" size="small" style={{ color: 'var(--text-secondary)' }}>{t('more')}</Button>}>
              <Text>{t('cardWithExtra')}</Text>
            </Card>
          </Col>
          <Col span={8}>
            <Card size="small" style={{ background: 'var(--bg-light)' }} styles={{ body: { padding: 20 } }}>
              <Statistic title={t('todayNewUsers')} value={128} suffix={t('userUnit')} />
              <div style={{ marginTop: 8 }}><Text style={{ fontSize: 12, color: 'var(--green)' }}>↑ 12.5%</Text><Text style={{ fontSize: 12, color: 'var(--text-muted)', marginLeft: 6 }}>{t('comparedYesterday')}</Text></div>
            </Card>
          </Col>
        </Row>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('descList')}</SectionTitle>
        <Descriptions column={2} size="small" bordered style={{ background: 'var(--bg-white)' }}>
          <Descriptions.Item label={t('username')}>admin</Descriptions.Item>
          <Descriptions.Item label={t('email')}>admin@example.com</Descriptions.Item>
          <Descriptions.Item label={t('role')}>{t('founder')}</Descriptions.Item>
          <Descriptions.Item label={t('status')}><Tag style={{ background: 'var(--green-light)', color: 'var(--green-text)', border: 'none', borderRadius: 'var(--radius-sm)' }}>{t('normal')}</Tag></Descriptions.Item>
          <Descriptions.Item label={t('registerTime')}>2024-01-15</Descriptions.Item>
          <Descriptions.Item label={t('lastLogin')}>2026-07-08</Descriptions.Item>
        </Descriptions>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('emptyState')}</SectionTitle>
        <Space size={24}>
          <Empty description={t('noData')} image={Empty.PRESENTED_IMAGE_SIMPLE} />
          <Empty description={t('noSearchResult')} />
        </Space>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('timeline')}</SectionTitle>
        <Timeline items={[
          { color: 'var(--green)', children: 'Project created 2024-01-15' },
          { color: 'var(--blue)', children: 'User module completed 2024-03-20' },
          { color: 'var(--yellow)', children: 'Permission system launched 2024-06-01' },
          { color: 'var(--red)', children: 'Security fix 2024-08-10' },
          { color: 'gray', children: 'v2.0 released 2025-01-01' },
        ]} />
      </div>
    </div>
  )
}

// ─── Tab 4: 布局与页面模式 ────────────────────

function LayoutAndPatterns() {
  const { t } = useTranslation()
  const [modalOpen, setModalOpen] = useState(false)

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
      <div>
        <SectionTitle>{t('dataPanelStructure')}</SectionTitle>
        <div style={{ border: '1px solid var(--border-color)', borderRadius: 'var(--radius-md)', overflow: 'hidden' }}>
          <div style={{ padding: '16px 20px', borderBottom: '1px solid var(--divider)', background: 'var(--bg-white)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Text style={{ fontSize: 14, fontWeight: 600 }}>title</Text>
              <div>
                <Button type="primary" size="small" style={{ background: 'var(--brand-dark)' }}>extra — New</Button>
              </div>
            </div>
          </div>
          <div style={{ padding: '12px 20px', borderBottom: '1px solid var(--divider)', background: '#faf8f5' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Space>
                <FilterSearch placeholder="filters" />
                <Button size="small" style={{ color: 'var(--text-primary)' }} icon={<SearchOutlined />}>{t('query')}</Button>
              </Space>
              <Button size="small" icon={<ExportOutlined />} style={{ color: 'var(--text-secondary)' }}>toolbarActions</Button>
            </div>
          </div>
          <div style={{ padding: '40px 20px', background: 'var(--bg-white)', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <Text style={{ color: 'var(--text-muted)', fontSize: 13 }}>children</Text>
          </div>
        </div>
        <div style={{ marginTop: 8, display: 'flex', gap: 16 }}>
          {[
            { label: 'title + extra', color: '#2b2b2b' },
            { label: 'filters + toolbarActions', color: '#3b6fdf' },
            { label: 'children', color: '#b0a89a' },
          ].map(s => (
            <Space key={s.label} size={4}>
              <span style={{ width: 10, height: 10, borderRadius: '50%', background: s.color, display: 'inline-block' }} />
              <Text style={{ fontSize: 12, color: 'var(--text-muted)' }}>{s.label}</Text>
            </Space>
          ))}
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('pageTitleArea')}</SectionTitle>
        <div style={{ background: 'var(--bg-light)', padding: 20, borderRadius: 'var(--radius-md)' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
            <div>
              <Text style={{ fontSize: 20, fontWeight: 600, color: 'var(--text-primary)', display: 'block' }}>{t('userManagement')}</Text>
              <Text style={{ fontSize: 13, color: 'var(--text-muted)', marginTop: 4, display: 'block' }}>Manage all users, roles and permissions</Text>
            </div>
            <Space>
              <Button style={{ color: 'var(--text-secondary)' }}>Secondary</Button>
              <Button type="primary" style={{ background: 'var(--brand-dark)' }}>Primary</Button>
            </Space>
          </div>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('filterBarLayout')}</SectionTitle>
        <div style={{ display: 'flex', gap: 12, alignItems: 'center', flexWrap: 'wrap' }}>
          <FilterSearch placeholder={t('search')} />
          <Select placeholder={t('category')} style={{ width: 140 }} options={[
            { value: 'all', label: t('all') },
            { value: 'active', label: t('active') },
            { value: 'pending', label: t('pending') },
          ]} />
          <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
          <Button style={{ color: 'var(--text-secondary)' }}>{t('reset')}</Button>
          <div style={{ flex: 1 }} />
          <Button icon={<ExportOutlined />} style={{ color: 'var(--text-secondary)' }}>Export</Button>
        </div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('stepsGuide')}</SectionTitle>
        <Steps current={1} size="small" items={[
          { title: 'Create', description: 'Submit' },
          { title: 'Review', description: 'Processing' },
          { title: t('completed'), description: 'Done' },
        ]} />
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('hoverMatrix')}</SectionTitle>
        <Table
          columns={[
            { title: 'Element', dataIndex: 'element', key: 'element', width: 200 },
            { title: 'Normal', dataIndex: 'normal', key: 'normal', width: 160 },
            { title: 'Hover', dataIndex: 'hover', key: 'hover', width: 260 },
            { title: 'Direction', dataIndex: 'direction', key: 'direction' },
          ]}
          dataSource={[
            { element: 'Primary Button', normal: '#2b2b2b', hover: '#4d4d4d', direction: 'Lighter' },
            { element: 'Default Button', normal: 'transparent', hover: '#f5f2ed', direction: 'Fill in' },
            { element: 'Pagination (inactive)', normal: 'transparent', hover: '#f5f2ed / #d4cdc0', direction: 'Fill in' },
            { element: 'Pagination (active) hover', normal: '#2b2b2b', hover: 'opacity 0.85', direction: 'Darker' },
            { element: 'Table Row', normal: '—', hover: '#faf8f5', direction: 'Light bg' },
            { element: 'Menu / Text Button', normal: 'transparent', hover: '#E4E0D8', direction: 'Darker' },
            { element: 'Sidebar Toggle', normal: 'transparent', hover: '#E4E0D8', direction: 'Fill in' },
          ]}
          pagination={false} size="small"
        />
      </div>
    </div>
  )
}

// ─── Tab 5: 交互与反馈 ────────────────────────

function InteractionFeedback() {
  const { t } = useTranslation()
  const [modalOpen, setModalOpen] = useState(false)

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
      <div>
        <SectionTitle>{t('tooltipPopover')}</SectionTitle>
        <Space size={16}>
          <Tooltip title="Tooltip text"><Button>Tooltip Hover</Button></Tooltip>
          <Popover content="Popover content" title={t('name')}><Button>Popover Click</Button></Popover>
          <Tooltip title={t('editUser')} placement="bottom"><Button type="primary" size="small">{t('edit')}</Button></Tooltip>
        </Space>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('modal')}</SectionTitle>
        <Button onClick={() => setModalOpen(true)} style={{ color: 'var(--text-secondary)' }}>{t('openModal')}</Button>
        <Modal title={t('modalExample')} open={modalOpen} onOk={() => setModalOpen(false)} onCancel={() => setModalOpen(false)} width={420}>
          <Text>{t('modalContent')}</Text>
        </Modal>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('skeleton')}</SectionTitle>
        <div style={{ maxWidth: 400 }}><Skeleton active paragraph={{ rows: 3 }} /></div>
      </div>

      <Divider style={{ margin: 0 }} />

      <div>
        <SectionTitle>{t('transitionTime')}</SectionTitle>
        <Table
          columns={[
            { title: 'Property', dataIndex: 'prop', key: 'prop', width: 200 },
            { title: t('value'), dataIndex: 'value', key: 'value' },
          ]}
          dataSource={[
            { prop: 'Default', value: 'all 0.15s ease' },
            { prop: 'CSS Variable', value: '0.25s cubic-bezier(0.4, 0, 0.2, 1)' },
          ]}
          pagination={false} size="small"
        />
      </div>
    </div>
  )
}

// ─── Main ─────────────────────────────────────

export default function DesignSystem() {
  const { t } = useTranslation()

  return (
    <DataPanel title={t('designSystem')}>
      <div style={{ padding: '0 28px 20px' }}>
        <Tabs
          defaultActiveKey="tokens"
          items={[
            { key: 'tokens', label: t('designTokens'), children: <DesignTokens /> },
            { key: 'components', label: t('componentLib'), children: <Components /> },
            { key: 'data', label: t('dataDisplay'), children: <DataDisplay /> },
            { key: 'patterns', label: t('layoutPatterns'), children: <LayoutAndPatterns /> },
            { key: 'interaction', label: t('interactionFeedback'), children: <InteractionFeedback /> },
          ]}
        />
      </div>
    </DataPanel>
  )
}
