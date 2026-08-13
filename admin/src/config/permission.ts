/**
 * 权限工具函数
 * 权限配置已迁移至后端动态管理，此处仅保留工具函数
 */

/** 权限检查 */
export function hasPermission(
  userPermissions: string[],
  requiredPermission: string
): boolean {
  if (!requiredPermission) return true
  return userPermissions.includes(requiredPermission)
}

/** 检查是否有任意一个权限 */
export function hasAnyPermission(
  userPermissions: string[],
  requiredPermissions: string[]
): boolean {
  if (!requiredPermissions || requiredPermissions.length === 0) return true
  return requiredPermissions.some(p => userPermissions.includes(p))
}
