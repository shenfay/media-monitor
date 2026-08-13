/**
 * 表格分页默认配置
 * 所有列表页统一使用此配置，避免重复定义
 */
export const DEFAULT_PAGINATION = {
  showSizeChanger: true,
  showQuickJumper: true,
} as const

/**
 * 获取带国际化 showTotal 的分页配置
 * 在组件内调用：pagination={{ ...DEFAULT_PAGINATION, ...getPaginationShowTotal(t) }}
 */
export function getPaginationShowTotal(t: (key: string, opts?: Record<string, unknown>) => string) {
  return {
    showTotal: (total: number) => t('totalRecords', { total }),
  }
}
