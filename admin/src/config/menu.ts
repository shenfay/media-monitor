/**
 * 菜单图标映射工具
 * 菜单结构已从后端动态获取，此处仅保留图标名称到 React 组件的映射
 */
import {
  DashboardOutlined,
  TeamOutlined,
  AimOutlined,
  FileTextOutlined,
  SmileOutlined,
  CheckCircleOutlined,
  StarOutlined,
  ShopOutlined,
  SwapOutlined,
  UserOutlined,
  LockOutlined,
  AuditOutlined,
  SettingOutlined,
  ProfileOutlined,
  MenuOutlined,
  BuildOutlined,
  ApiOutlined,
} from '@ant-design/icons'
import React, { type ReactNode } from 'react'

const iconMap: Record<string, ReactNode> = {
  DashboardOutlined: React.createElement(DashboardOutlined),
  TeamOutlined: React.createElement(TeamOutlined),
  AimOutlined: React.createElement(AimOutlined),
  FileTextOutlined: React.createElement(FileTextOutlined),
  SmileOutlined: React.createElement(SmileOutlined),
  CheckCircleOutlined: React.createElement(CheckCircleOutlined),
  StarOutlined: React.createElement(StarOutlined),
  ShopOutlined: React.createElement(ShopOutlined),
  SwapOutlined: React.createElement(SwapOutlined),
  UserOutlined: React.createElement(UserOutlined),
  LockOutlined: React.createElement(LockOutlined),
  AuditOutlined: React.createElement(AuditOutlined),
  SettingOutlined: React.createElement(SettingOutlined),
  ProfileOutlined: React.createElement(ProfileOutlined),
  MenuOutlined: React.createElement(MenuOutlined),
  BuildOutlined: React.createElement(BuildOutlined),
  ApiOutlined: React.createElement(ApiOutlined),
}

export function getIcon(name: string): ReactNode {
  return iconMap[name] || null
}
