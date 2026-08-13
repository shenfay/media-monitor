import { useState, useEffect, useCallback } from 'react'
import { Tag, Space, Button, Modal, Form, Input, Select, message, Switch, Popconfirm, Table } from 'antd'
import type { TableColumnsType } from 'antd'
import { PlusOutlined, EditOutlined, SearchOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { DEFAULT_PAGINATION, getPaginationShowTotal } from '@/config/pagination'
import { useCrudList } from '@/hooks/useCrudList'
import { getUserList, createUser, updateUser, toggleUserStatus } from '@/services/user'
import { getRoleList } from '@/services/role'
import { getEmailRules, getPasswordRules } from '@/utils/formRules'
import type { User, Role } from '@/types'

export default function UserManagement() {
  const { t, i18n } = useTranslation()
  const {
    loading, dataSource: users, total, page, pageSize,
    fetchData: fetchUsers, handlePageChange,
    isModalOpen, editingItem: editingUser, form,
    handleAdd, handleEdit, handleCancel,
  } = useCrudList<User>(
    async ({ page: p, pageSize: ps }) => {
      const res = await getUserList({ page: p, page_size: ps, keyword })
      return { data: res.users || [], total: res.total || 0 }
    },
  )

  const [roles, setRoles] = useState<Role[]>([])
  const [keyword, setKeyword] = useState('')
  const [roleFilter, setRoleFilter] = useState('')

  const fetchRoles = async () => {
    try {
      const res = await getRoleList()
      setRoles(res || [])
    } catch {
      setRoles([])
    }
  }

  useEffect(() => {
    fetchRoles()
  }, [])

  // keyword 变化时重新加载
  const fetchWithKeyword = useCallback(async () => {
    const res = await getUserList({ page, page_size: pageSize, keyword })
    return { data: res.users || [], total: res.total || 0 }
  }, [page, pageSize, keyword])

  useEffect(() => {
    fetchUsers()
  }, [fetchWithKeyword]) // eslint-disable-line react-hooks/exhaustive-deps

  const handleEditUser = (record: User) => {
    handleEdit(record, {
      name: record.name,
      email: record.email,
      role_ids: record.roles?.map(r => r.id) || [],
    })
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingUser) {
        await updateUser(editingUser.id, {
          name: values.name,
          email: values.email,
          role_ids: values.role_ids,
        })
        message.success(t('updateSuccess'))
      } else {
        await createUser({
          name: values.name,
          email: values.email,
          password: values.password,
          role_ids: values.role_ids || [],
        })
        message.success(t('createSuccess'))
      }
      handleCancel()
      fetchUsers()
    } catch {
      // validation or API error
    }
  }

  const handleToggleStatus = async (record: User) => {
    try {
      await toggleUserStatus(record.id, !record.locked)
      message.success(record.locked ? t('enableUser') : t('disableUser'))
      fetchUsers()
    } catch {
      message.error(t('operationFailed'))
    }
  }

  const columns: TableColumnsType<User> = [
    { title: t('name'), dataIndex: 'name', key: 'name', width: 140 },
    { title: t('email'), dataIndex: 'email', key: 'email', ellipsis: true },
    {
      title: t('roles'),
      dataIndex: 'roles',
      key: 'roles',
      width: 160,
      render: (_: unknown, record: User) => (
        <Space wrap size={4}>
          {(record.roles || []).map(r => (
            <Tag key={r.code} style={{ background: 'var(--blue-light)', color: 'var(--blue-text)' }}>{r.name}</Tag>
          ))}
          {(!record.roles || record.roles.length === 0) && <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)' }}>{t('unassigned')}</Tag>}
        </Space>
      ),
    },
    {
      title: t('status'),
      dataIndex: 'locked',
      key: 'locked',
      width: 100,
      render: (_: unknown, record: User) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
          <span style={{
            width: 7, height: 7, borderRadius: '50%',
            background: record.locked ? 'var(--border-hover)' : 'var(--green)',
            display: 'inline-block',
          }} />
          <span style={{ color: 'var(--text-primary)' }}>{record.locked ? t('disabled') : t('active')}</span>
        </div>
      ),
    },
    {
      title: t('createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (_: unknown, record: User) => record.created_at
        ? <span style={{ color: 'var(--text-secondary)' }}>{new Date(record.created_at).toLocaleString(i18n.language)}</span>
        : <span style={{ color: 'var(--text-muted)' }}>-</span>,
    },
    {
      title: t('actions'),
      key: 'action',
      width: 120,
      render: (_: unknown, record: User) => (
        <Space size={12}>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEditUser(record)}>
            {t('edit')}
          </Button>
          <Popconfirm
            title={record.locked ? t('confirmEnable') : t('confirmDisable')}
            onConfirm={() => handleToggleStatus(record)}
          >
            <Switch checked={!record.locked} size="small" />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('userManagement')}
        filters={
          <>
            <FilterSearch value={keyword} onChange={setKeyword} placeholder={t('searchNameOrEmail')} onSearch={() => fetchUsers()} />
            <Select
              value={roleFilter}
              onChange={setRoleFilter}
              style={{ width: 140 }}
              options={[
                { label: t('allRoles'), value: '' },
                ...roles.map(r => ({ label: r.name, value: r.id })),
              ]}
            />
            <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }} onClick={() => fetchUsers()}>{t('query')}</Button>
          </>
        }
        toolbarActions={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            {t('addUser')}
          </Button>
        }
      >
        <Table<User>
          dataSource={users}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            ...DEFAULT_PAGINATION,
            ...getPaginationShowTotal(t),
            onChange: handlePageChange,
          }}
        />
      </DataPanel>

      <Modal
        title={editingUser ? t('editUser') : t('addUser')}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={handleCancel}
        width={520}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item label={t('name')} name="name" rules={[{ required: true, message: t('namePlaceholder') }]}>
            <Input placeholder={t('namePlaceholder')} />
          </Form.Item>
          <Form.Item label={t('email')} name="email" rules={getEmailRules()}>
            <Input placeholder="example@domain.com" />
          </Form.Item>
          {!editingUser && (
            <Form.Item label={t('password')} name="password" rules={getPasswordRules()}>
              <Input.Password placeholder={t('passwordPlaceholder')} />
            </Form.Item>
          )}
          <Form.Item label={t('roles')} name="role_ids">
            <Select mode="multiple" placeholder={t('selectRole')}>
              {roles.map(r => (
                <Select.Option key={r.id} value={r.id}>{r.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
