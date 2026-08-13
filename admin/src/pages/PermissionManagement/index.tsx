import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Table, Tag, Tree, Switch, Button, Modal, Form, Input, message, Popconfirm, Space } from 'antd'
import { PlusOutlined, SearchOutlined, SettingOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import { getMenuTree } from '@/services/menu'
import type { MenuItem } from '@/types'
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  toggleRoleStatus,
  getRolePermissions,
  updateRolePermissions,
} from '@/services/role'
import type { Role } from '@/types'

const { TextArea } = Input

type PermTreeNode = {
  title: string
  key: string
  permission?: string
  children?: PermTreeNode[]
}

/** 从后端菜单数据生成权限树 */
function buildPermissionTree(menus: MenuItem[]): PermTreeNode[] {
  const tree: PermTreeNode[] = []
  const seenKeys = new Set<string>()

  for (const menu of menus) {
    if (!menu.status || seenKeys.has(menu.key)) continue
    seenKeys.add(menu.key)

    const children = (menu.children || [])
      .filter(c => c.status && !seenKeys.has(c.key))
      .map(item => {
        seenKeys.add(item.key)
        return { title: item.label, key: item.key, permission: item.permission || `${item.key}:view` }
      })

    if (children.length > 0) {
      tree.push({ title: menu.label, key: menu.key, children })
    } else {
      tree.push({ title: menu.label, key: menu.key, permission: menu.permission || `${menu.key}:view` })
    }
  }
  return tree
}

/** 从权限树中提取所有叶子节点的 key */
function getAllLeafKeys(treeData: PermTreeNode[]): string[] {
  const keys: string[] = []
  for (const node of treeData) {
    if (node.children) {
      node.children.forEach(item => keys.push(item.key))
    } else {
      keys.push(node.key)
    }
  }
  return keys
}

/** 构建 key -> permission 映射 */
function buildKeyPermMap(treeData: PermTreeNode[]): Record<string, string> {
  const map: Record<string, string> = {}
  for (const node of treeData) {
    if (node.children) {
      node.children.forEach(item => {
        map[item.key] = item.permission || `${item.key}:view`
      })
    } else {
      map[node.key] = node.permission || `${node.key}:view`
    }
  }
  return map
}

export default function PermissionManagement() {
  const { t } = useTranslation()
  const [roles, setRoles] = useState<Role[]>([])
  const [selectedRole, setSelectedRole] = useState<Role | null>(null)
  const [checkedKeys, setCheckedKeys] = useState<string[]>([])
  const [isRoleModalOpen, setIsRoleModalOpen] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [loading, setLoading] = useState(false)
  const [permLoading, setPermLoading] = useState(false)
  const [menuData, setMenuData] = useState<MenuItem[]>([])
  const [keyword, setKeyword] = useState('')
  const [form] = Form.useForm()

  const permissionTree = buildPermissionTree(menuData)
  const allLeafKeys = getAllLeafKeys(permissionTree)
  const keyPermMap = buildKeyPermMap(permissionTree)

  const fetchRoles = useCallback(async () => {
    setLoading(true)
    try {
      const res = await getRoleList()
      setRoles(res || [])
    } catch {
      setRoles([])
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchMenus = useCallback(async () => {
    try {
      const res = await getMenuTree()
      setMenuData(res || [])
    } catch {
      setMenuData([])
    }
  }, [])

  useEffect(() => {
    fetchRoles()
    fetchMenus()
  }, [fetchRoles, fetchMenus])

  /** 加载角色权限 */
  const loadRolePermissions = async (role: Role) => {
    setPermLoading(true)
    try {
      const perms = await getRolePermissions(role.id)
      const permToKey: Record<string, string> = {}
      Object.entries(keyPermMap).forEach(([key, perm]) => {
        permToKey[perm] = key
      })
      const menuKeys = perms
        .map((p: string) => permToKey[p])
        .filter((k): k is string => !!k)
      setCheckedKeys(menuKeys)
    } catch {
      setCheckedKeys([])
    } finally {
      setPermLoading(false)
    }
  }

  /** 选择角色配置权限 */
  const handleSelectRole = (role: Role) => {
    setSelectedRole(role)
    loadRolePermissions(role)
  }

  /** 保存权限配置 */
  const handleSavePermissions = async () => {
    if (!selectedRole) return
    setPermLoading(true)
    try {
      const permissions: string[] = []
      checkedKeys.forEach(key => {
        const perm = keyPermMap[key]
        if (perm) {
          permissions.push(perm)
        }
      })
      await updateRolePermissions(selectedRole.id, permissions)
      message.success(t('permSaved'))
    } catch {
      message.error(t('permSaveFailed'))
    } finally {
      setPermLoading(false)
    }
  }

  /** 新增角色 */
  const handleAddRole = () => {
    setEditingRole(null)
    form.resetFields()
    setIsRoleModalOpen(true)
  }

  /** 编辑角色 */
  const handleEditRole = (role: Role) => {
    setEditingRole(role)
    form.setFieldsValue({
      name: role.name,
      code: role.code,
      description: role.description,
    })
    setIsRoleModalOpen(true)
  }

  /** 提交角色表单 */
  const handleSubmitRole = async () => {
    try {
      const values = await form.validateFields()
      if (editingRole) {
        await updateRole(editingRole.id, {
          name: values.name,
          description: values.description,
        })
        message.success(t('roleUpdated'))
      } else {
        await createRole({
          name: values.name,
          code: values.code,
          description: values.description,
        })
        message.success(t('roleCreated'))
      }
      setIsRoleModalOpen(false)
      form.resetFields()
      fetchRoles()
    } catch {
      // validation or API error
    }
  }

  /** 删除角色 */
  const handleDeleteRole = async (role: Role) => {
    try {
      await deleteRole(role.id)
      message.success(t('roleDeleted'))
      if (selectedRole?.id === role.id) {
        setSelectedRole(null)
        setCheckedKeys([])
      }
      fetchRoles()
    } catch {
      message.error(t('deleteFailed'))
    }
  }

  /** 切换角色状态 */
  const handleToggleStatus = async (role: Role) => {
    try {
      await toggleRoleStatus(role.id)
      message.success(t('statusUpdated'))
      fetchRoles()
    } catch {
      message.error(t('operationFailed'))
    }
  }

  const columns = [
    { title: t('roleName'), dataIndex: 'name', key: 'name', width: 120 },
    {
      title: t('roleCode'),
      dataIndex: 'code',
      key: 'code',
      width: 120,
      render: (v: string) => <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)' }}>{v}</Tag>,
    },
    { title: t('description'), dataIndex: 'description', key: 'description' },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (v: boolean, record: Role) => (
        <Switch
          checked={v}
          size="small"
          onChange={() => handleToggleStatus(record)}
        />
      ),
    },
    {
      title: t('actions'),
      key: 'action',
      width: 200,
      render: (_: unknown, record: Role) => (
        <Space size={4}>
          <Button type="link" size="small" icon={<SettingOutlined />} onClick={() => handleSelectRole(record)}>
            {t('configPermission')}
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEditRole(record)}>
            {t('edit')}
          </Button>
          <Popconfirm
            title={t('confirmDeleteRole')}
            description={t('roleCannotUndo')}
            onConfirm={() => handleDeleteRole(record)}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>{t('delete')}</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <DataPanel
        title={t('permissionManagement')}
        filters={
          <>
            <FilterSearch value={keyword} onChange={setKeyword} placeholder={t('searchRoleName')} />
            <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
          </>
        }
        toolbarActions={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddRole}>
            {t('addRole')}
          </Button>
        }
        >
        <Table
          dataSource={roles}
          columns={columns}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      </DataPanel>

      {selectedRole && (
        <DataPanel
          title={t('menuPermConfigTitle', { name: selectedRole.name })}
          style={{ marginTop: 0 }}
          extra={
            <Space>
              <Button
                type="primary"
                size="small"
                loading={permLoading}
                onClick={handleSavePermissions}
              >
                {t('savePermission')}
              </Button>
              <Button size="small" onClick={() => setSelectedRole(null)}>
                {t('collapse')}
              </Button>
            </Space>
          }
          compact
        >
          <Tree
            checkable
            checkedKeys={checkedKeys}
            onCheck={(keys) => {
              const leafKeys = (keys as string[]).filter(k => allLeafKeys.includes(k))
              setCheckedKeys(leafKeys)
            }}
            defaultExpandedKeys={permissionTree.map(g => g.key)}
            treeData={permissionTree}
          />
        </DataPanel>
      )}

      <Modal
        title={editingRole ? t('editRole') : t('addRole')}
        open={isRoleModalOpen}
        onOk={handleSubmitRole}
        onCancel={() => { setIsRoleModalOpen(false); form.resetFields() }}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item label={t('roleName')} name="name" rules={[{ required: true, message: t('roleNameRequired') }]}>
            <Input placeholder={t('roleNamePlaceholder')} />
          </Form.Item>
          {!editingRole && (
            <Form.Item label={t('roleCode')} name="code" rules={[{ required: true, message: t('roleCodeRequired') }]}>
              <Input placeholder={t('roleCodePlaceholder')} />
            </Form.Item>
          )}
          <Form.Item label={t('description')} name="description">
            <TextArea rows={3} placeholder={t('roleDescPlaceholder')} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
