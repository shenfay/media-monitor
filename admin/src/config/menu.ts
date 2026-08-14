/**
 * 菜单图标映射工具
 * 菜单结构已从后端动态获取，此处仅保留图标名称到 React 组件的映射
 */
import {
  DashboardOutlined,
  TeamOutlined,
  FileTextOutlined,
  UserOutlined,
  LockOutlined,
  AuditOutlined,
  SettingOutlined,
  ProfileOutlined,
  MenuOutlined,
  BuildOutlined,
  ApiOutlined,
  GlobalOutlined,
  DatabaseOutlined,
  UnorderedListOutlined,
  ReadOutlined,
  MailOutlined,
} from '@ant-design/icons'
import React, { type ReactNode } from 'react'

const iconMap: Record<string, ReactNode> = {
  DashboardOutlined: React.createElement(DashboardOutlined),
  TeamOutlined: React.createElement(TeamOutlined),
  FileTextOutlined: React.createElement(FileTextOutlined),
  UserOutlined: React.createElement(UserOutlined),
  LockOutlined: React.createElement(LockOutlined),
  AuditOutlined: React.createElement(AuditOutlined),
  SettingOutlined: React.createElement(SettingOutlined),
  ProfileOutlined: React.createElement(ProfileOutlined),
  MenuOutlined: React.createElement(MenuOutlined),
  BuildOutlined: React.createElement(BuildOutlined),
  ApiOutlined: React.createElement(ApiOutlined),
  GlobalOutlined: React.createElement(GlobalOutlined),
  DatabaseOutlined: React.createElement(DatabaseOutlined),
  UnorderedListOutlined: React.createElement(UnorderedListOutlined),
  ReadOutlined: React.createElement(ReadOutlined),
  MailOutlined: React.createElement(MailOutlined),
}

export function getIcon(name: string): ReactNode {
  return iconMap[name] || null
}
