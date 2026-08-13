import { useState, useEffect, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Switch,
  Tag,
  Space,
  Popconfirm,
  message,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
} from '@ant-design/icons'
import DataPanel, { FilterSearch } from '@/components/DataPanel'
import {
  getMenuTree,
  createMenu,
  updateMenu,
  deleteMenu,
  toggleMenuStatus,
} from '@/services/menu'
import type { MenuItem } from '@/types'


export default function MenuManagement() {
  const { t } = useTranslation()
  const [menuTree, setMenuTree] = useState<MenuItem[]>([])
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState<MenuItem | null>(null)
  const [expandedKeys, setExpandedKeys] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [keyword, setKeyword] = useState('')
  const [form] = Form.useForm()

  /** 加载菜单树 */
  const fetchMenus = useCallback(async () => {
    setLoading(true)
    try {
      const tree = await getMenuTree()
      setMenuTree(tree || [])
      if (expandedKeys.length === 0 && tree) {
        setExpandedKeys(tree.map(m => m.key))
      }
    } catch {
      setMenuTree([])
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchMenus()
  }, [fetchMenus])

  /** 扁平化菜单列表 */
  const flatMenus: MenuItem[] = []
  function flatten(nodes: MenuItem[]) {
    for (const node of nodes) {
      flatMenus.push(node)
      if (node.children) flatten(node.children)
    }
  }
  flatten(menuTree)

  /** 新增菜单 */
  const handleAddRoot = () => {
    setEditingItem(null)
    form.resetFields()
    form.setFieldsValue({ parent_id: '' })
    setIsModalOpen(true)
  }

  /** 新增子菜单 */
  const handleAddChild = (parentKey: string) => {
    const parent = flatMenus.find(m => m.key === parentKey)
    setEditingItem(null)
    form.resetFields()
    form.setFieldsValue({ parent_id: parent?.id || '' })
    setIsModalOpen(true)
  }

  /** 编辑 */
  const handleEdit = (record: MenuItem) => {
    setEditingItem(record)
    form.setFieldsValue({
      parent_id: record.parent_id || '',
      label: record.label,
      key: record.key,
      icon: record.icon || undefined,
      path: record.path || undefined,
      permission: record.permission || undefined,
    })
    setIsModalOpen(true)
  }

  /** 删除 */
  const handleDelete = async (record: MenuItem) => {
    try {
      await deleteMenu(record.id)
      message.success(t('menuDeleted'))
      fetchMenus()
    } catch {
      message.error(t('deleteFailed'))
    }
  }

  /** 切换状态 */
  const handleToggleStatus = async (record: MenuItem) => {
    try {
      await toggleMenuStatus(record.id)
      message.success(t('statusUpdated'))
      fetchMenus()
    } catch {
      message.error(t('operationFailed'))
    }
  }

  /** 提交表单 */
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingItem) {
        await updateMenu(editingItem.id, {
          label: values.label,
          icon: values.icon || '',
          path: values.path || '',
          permission: values.permission || '',
        })
        message.success(t('menuUpdated'))
      } else {
        const parent = flatMenus.find(m => m.id === values.parent_id)
        await createMenu({
          key: values.key,
          label: values.label,
          icon: values.icon || '',
          path: values.path || '',
          permission: values.permission || '',
          parent_id: values.parent_id || '',
          sort_order: 0,
        })
        message.success(t('menuAdded'))
        if (parent && !expandedKeys.includes(parent.key)) {
          setExpandedKeys(prev => [...prev, parent.key])
        }
      }
      setIsModalOpen(false)
      form.resetFields()
      fetchMenus()
    } catch {
      // validation or API error
    }
  }

  const columns = [
    {
      title: t('menuName'),
      dataIndex: 'label',
      key: 'label',
      width: 200,
    },
    {
      title: t('menuKey'),
      dataIndex: 'key',
      key: 'key',
      width: 160,
      render: (v: string) => <Tag style={{ background: 'var(--blue-light)', color: 'var(--blue-text)' }}>{v}</Tag>,
    },
    {
      title: t('icon'),
      dataIndex: 'icon',
      key: 'icon',
      width: 160,
      render: (v: string) => v ? <Tag style={{ background: 'var(--gray-light)', color: 'var(--gray-text)' }}>{v}</Tag> : <span style={{ color: 'var(--text-muted)' }}>-</span>,
    },
    {
      title: t('routePath'),
      dataIndex: 'path',
      key: 'path',
      width: 160,
      render: (v: string) => v || <span style={{ color: 'var(--text-muted)' }}>—</span>,
    },
    {
      title: t('permission'),
      dataIndex: 'permission',
      key: 'permission',
      width: 180,
      render: (v: string) => v ? <Tag style={{ background: 'var(--green-light)', color: 'var(--green-text)' }}>{v}</Tag> : <span style={{ color: 'var(--text-muted)' }}>-</span>,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (v: boolean, record: MenuItem) => (
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
      width: 140,
      render: (_: unknown, record: MenuItem) => (
        <Space size={4}>
          <Button type="link" size="small" icon={<PlusOutlined />} onClick={() => handleAddChild(record.key)}>
            {t('addSubMenu')}
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            {t('edit')}
          </Button>
          <Popconfirm
            title={t('confirmDeleteMenu')}
            description={t('subMenuWillBeDeleted')}
            onConfirm={() => handleDelete(record)}
          >
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
        title={t('menuManagement')}
        filters={
          <>
            <FilterSearch value={keyword} onChange={setKeyword} placeholder={t('searchMenuName')} />
            <Button icon={<SearchOutlined />} style={{ color: 'var(--text-primary)' }}>{t('query')}</Button>
          </>
        }
        toolbarActions={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddRoot}>
            {t('addMenu')}
          </Button>
        }
        >
        <Table
          dataSource={menuTree}
          columns={columns}
          rowKey="key"
          loading={loading}
          pagination={false}
          expandable={{
            expandedRowKeys: expandedKeys,
            onExpandedRowsChange: (keys) => setExpandedKeys(keys as string[]),
          }}
          size="middle"
        />
      </DataPanel>

      <Modal
        title={editingItem ? t('editMenu') : t('addMenu')}
        open={isModalOpen}
        onOk={handleSubmit}
        onCancel={() => { setIsModalOpen(false); form.resetFields() }}
        width={520}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item label={t('parentMenu')} name="parent_id">
            <Select
              placeholder={t('noParentMenu')}
              allowClear
              options={[
                { value: '', label: t('noParentMenu') },
                ...flatMenus
                  .filter(f => !editingItem || f.id !== editingItem.id)
                  .map(f => ({ value: f.id, label: f.label })),
              ]}
            />
          </Form.Item>
          <div style={{ display: 'flex', gap: 16 }}>
            <Form.Item
              label={t('menuName')}
              name="label"
              rules={[{ required: true, message: t('pleaseEnter', { field: t('menuName') }) }]}
              style={{ flex: 1 }}
            >
              <Input placeholder={t('menuNamePlaceholder')} />
            </Form.Item>
            <Form.Item
              label={t('menuKey')}
              name="key"
              rules={[
                { required: true, message: t('pleaseEnter', { field: t('menuKey') }) },
                { pattern: /^[a-z0-9-]+$/, message: t('menuKeyRule') },
              ]}
              style={{ flex: 1 }}
            >
              <Input placeholder={t('menuKeyPlaceholder')} disabled={!!editingItem} />
            </Form.Item>
          </div>
          <div style={{ display: 'flex', gap: 16 }}>
            <Form.Item label={t('icon')} name="icon" style={{ flex: 1 }}>
              <Input placeholder={t('iconPlaceholder')} />
            </Form.Item>
            <Form.Item label={t('routePath')} name="path" style={{ flex: 1 }}>
              <Input placeholder={t('routePathPlaceholder')} />
            </Form.Item>
          </div>
          <Form.Item label={t('permission')} name="permission">
            <Input placeholder={t('permissionPlaceholder')} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
