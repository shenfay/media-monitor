import { useState, useEffect, useCallback, useRef } from 'react'
import { Form, message } from 'antd'
import { useTranslation } from 'react-i18next'
import type { FormInstance } from 'antd'

/**
 * 列表页公共 CRUD Hook
 *
 * 封装列表页常见的：分页加载、弹窗状态、表单编辑 等模式。
 * 适用于 UserManagement / OperationLog / MenuManagement 等页面。
 */

interface UseCrudListOptions<T> {
  /** 获取列表数据的函数，返回数据数组和总数 */
  fetchFn: (params: { page: number; pageSize: number }) => Promise<{ data: T[]; total: number }>
  /** 每页条数，默认 20 */
  defaultPageSize?: number
}

interface ModalHandlers<T> {
  /** 打开新增弹窗 */
  handleAdd: () => void
  /** 打开编辑弹窗，并填充表单数据 */
  handleEdit: (record: T, formValues?: Record<string, unknown>) => void
  /** 关闭弹窗并重置表单 */
  handleCancel: () => void
}

export function useCrudList<T>(
  fetchFn: UseCrudListOptions<T>['fetchFn'],
  options?: Omit<UseCrudListOptions<T>, 'fetchFn'>,
) {
  const { t } = useTranslation()
  const defaultPageSize = options?.defaultPageSize ?? 20

  // ---- 列表状态 ----
  const [loading, setLoading] = useState(false)
  const [dataSource, setDataSource] = useState<T[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(defaultPageSize)

  // ---- 弹窗 / 表单状态 ----
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingItem, setEditingItem] = useState<T | null>(null)
  const [form] = Form.useForm()

  // ---- 数据加载 ----
  // 用 ref 保存最新的 fetchFn，避免内联函数导致无限循环
  const fetchFnRef = useRef(fetchFn)
  fetchFnRef.current = fetchFn

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetchFnRef.current({ page, pageSize })
      setDataSource(res.data || [])
      setTotal(res.total || 0)
    } catch {
      setDataSource([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }, [page, pageSize])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  // ---- 分页变更 ----
  const handlePageChange = (p: number, ps: number) => {
    setPage(p)
    setPageSize(ps)
  }

  // ---- 弹窗操作 ----
  const handleAdd = () => {
    setEditingItem(null)
    form.resetFields()
    setIsModalOpen(true)
  }

  const handleEdit = (record: T, formValues?: Record<string, unknown>) => {
    setEditingItem(record)
    if (formValues) {
      form.setFieldsValue(formValues)
    }
    setIsModalOpen(true)
  }

  const handleCancel = () => {
    setIsModalOpen(false)
    form.resetFields()
  }

  /**
   * 通用提交处理
   * @param createFn 新增接口
   * @param updateFn 更新接口（接收 editingItem + formValues）
   * @param getFormValues 从 form 提取字段的函数
   */
  const handleSubmit = async (
    createFn: (values: Record<string, unknown>) => Promise<unknown>,
    updateFn: (item: T, values: Record<string, unknown>) => Promise<unknown>,
    getFormValues?: (values: Record<string, unknown>) => Record<string, unknown>,
  ) => {
    try {
      const values = await form.validateFields()
      const payload = getFormValues ? getFormValues(values) : values
      if (editingItem) {
        await updateFn(editingItem, payload)
        message.success(t('updateSuccess'))
      } else {
        await createFn(payload)
        message.success(t('createSuccess'))
      }
      setIsModalOpen(false)
      form.resetFields()
      fetchData()
    } catch {
      // validation or API error
    }
  }

  return {
    // 列表状态
    loading,
    dataSource,
    total,
    page,
    pageSize,
    fetchData,
    handlePageChange,

    // 弹窗 / 表单
    isModalOpen,
    editingItem,
    form: form as FormInstance,
    handleAdd,
    handleEdit,
    handleCancel,
    handleSubmit,
  }
}
