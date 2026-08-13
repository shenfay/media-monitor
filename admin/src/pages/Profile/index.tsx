import { useState } from 'react'
import { Form, Input, Button, Avatar, Descriptions, Divider, message, Tabs, Tag } from 'antd'
import { UserOutlined, LockOutlined, BellOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useUserStore } from '@/stores'
import DataPanel from '@/components/DataPanel'

export default function Profile() {
  const { t } = useTranslation()
  const [profileForm] = Form.useForm()
  const [passwordForm] = Form.useForm()
  const [activeTab, setActiveTab] = useState('profile')
  const { username, email, roles } = useUserStore()

  const roleLabels = roles.map(r => r.name).join('、') || t('noRole')

  const handleSaveProfile = () => {
    message.success(t('updateSuccess'))
  }

  const handleChangePassword = () => {
    message.success(t('updateSuccess'))
  }

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      {/* User Info Card */}
      <DataPanel title="" compact>
        <div style={{ padding: '20px 28px', display: 'flex', alignItems: 'center', gap: 24 }}>
          <Avatar size={80} icon={<UserOutlined />} style={{ background: 'var(--brand-dark)', flexShrink: 0 }} />
          <div>
            <div style={{ fontSize: 20, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 4 }}>
              {username || t('name')}
            </div>
            <div style={{ color: 'var(--text-secondary)', marginBottom: 8 }}>{email || ''}</div>
            <div style={{ display: 'flex', gap: 8 }}>
              {roles.map(r => <Tag key={r.code} style={{ background: 'var(--blue-light)', color: 'var(--blue-text)' }}>{r.name}</Tag>)}
              {roles.length === 0 && <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)' }}>{t('noRole')}</Tag>}
            </div>
          </div>
        </div>
      </DataPanel>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        style={{ padding: '0 28px' }}
        items={[
          {
            key: 'profile',
            label: (
              <span>
                <UserOutlined />
                {t('basicInfo')}
              </span>
            ),
            children: (
              <div style={{ padding: '16px 0' }}>
                <Descriptions bordered column={2}>
                  <Descriptions.Item label={t('name')} span={1}>{username || '-'}</Descriptions.Item>
                  <Descriptions.Item label={t('email')} span={1}>{email || '-'}</Descriptions.Item>
                  <Descriptions.Item label={t('roles')} span={1}>{roleLabels}</Descriptions.Item>
                  <Descriptions.Item label={t('department')} span={1}>-</Descriptions.Item>
                  <Descriptions.Item label={t('phone')} span={1}>-</Descriptions.Item>
                  <Descriptions.Item label={t('registerDate')} span={1}>-</Descriptions.Item>
                </Descriptions>

                <Divider />

                <Form
                  form={profileForm}
                  layout="vertical"
                  style={{ maxWidth: 500 }}
                >
                  <Form.Item label={t('nickname')} name="nickname" initialValue={username || ''}>
                    <Input />
                  </Form.Item>
                  <Form.Item label={t('email')} name="email" initialValue={email || ''} rules={[{ type: 'email', message: t('emailInvalid') }]}>
                    <Input />
                  </Form.Item>
                  <Form.Item label={t('phone')} name="phone" initialValue="">
                    <Input />
                  </Form.Item>
                  <Form.Item label={t('bio')} name="bio">
                    <Input.TextArea rows={3} placeholder={t('bioPlaceholder')} />
                  </Form.Item>
                  <Button type="primary" onClick={handleSaveProfile}>{t('saveChanges')}</Button>
                </Form>
              </div>
            ),
          },
          {
            key: 'password',
            label: (
              <span>
                <LockOutlined />
                {t('changePassword')}
              </span>
            ),
            children: (
              <div style={{ padding: '16px 0' }}>
                <Form
                  form={passwordForm}
                  layout="vertical"
                  style={{ maxWidth: 400 }}
                >
                  <Form.Item label={t('currentPassword')} name="oldPassword" rules={[{ required: true, message: t('pleaseEnter', { field: t('currentPassword') }) }]}>
                    <Input.Password />
                  </Form.Item>
                  <Form.Item label={t('newPassword')} name="newPassword" rules={[
                    { required: true, message: t('pleaseEnter', { field: t('newPassword') }) },
                    { min: 8, message: t('passwordMinLength') },
                  ]}>
                    <Input.Password />
                  </Form.Item>
                  <Form.Item label={t('confirmPassword')} name="confirmPassword" rules={[
                    { required: true, message: t('pleaseEnter', { field: t('confirmPassword') }) },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue('newPassword') === value) return Promise.resolve()
                        return Promise.reject(new Error(t('passwordMismatch')))
                      },
                    }),
                  ]}>
                    <Input.Password />
                  </Form.Item>
                  <Button type="primary" onClick={handleChangePassword}>{t('changePassword')}</Button>
                </Form>
              </div>
            ),
          },
          {
            key: 'notifications',
            label: (
              <span>
                <BellOutlined />
                {t('notificationSettings')}
              </span>
            ),
            children: (
              <div style={{ padding: '16px 0' }}>
                <div style={{ color: 'var(--text-secondary)', marginBottom: 16 }}>
                  {t('notificationDesc')}
                </div>
                <Form layout="vertical" style={{ maxWidth: 500 }}>
                  <Form.Item label={t('emailNotify')} name="emailNotify">
                    <div style={{ color: 'var(--text-secondary)', fontSize: 13, marginBottom: 8 }}>
                      {t('emailNotifyDesc')}
                    </div>
                  </Form.Item>
                  <Form.Item label={t('smsNotify')} name="smsNotify">
                    <div style={{ color: 'var(--text-secondary)', fontSize: 13, marginBottom: 8 }}>
                      {t('smsNotifyDesc')}
                    </div>
                  </Form.Item>
                  <Button type="primary">{t('saveSettings')}</Button>
                </Form>
              </div>
            ),
          },
        ]}
      />
    </div>
  )
}
