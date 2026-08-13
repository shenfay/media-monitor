# 数据中台管理后台

基于 **React + TypeScript + Vite** 构建的企业级数据中台管理后台模板，采用 Coze 风格 UI 设计，提供开箱即用的菜单体系、权限控制、布局框架和通用组件库。

## 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | React | ^19.2.6 |
| 构建工具 | Vite | ^8.0.12 |
| 类型系统 | TypeScript | ~6.0.2 |
| UI 组件库 | Ant Design | ^5.29.3 |
| 高级组件 | @ant-design/pro-components | ^2.8.10 |
| 图表库 | @ant-design/charts | ^2.6.7 |
| 路由 | react-router-dom | ^7.15.1 |
| 状态管理 | Zustand | ^5.0.13 |
| HTTP 请求 | Axios | ^1.16.1 |
| 代码检查 | ESLint | ^10.3.0 |
| 格式化 | Prettier | ^3.8.3 |
| 规范提交 | Commitizen + Husky | - |

## 快速开始

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览构建产物
npm run preview

# 代码检查
npm run lint
```

## 目录结构

```
src/
├── assets/            # 静态资源（图片等）
├── components/        # 通用业务组件
│   ├── DataPanel/     #   数据面板卡片
│   ├── DataTable/     #   数据表格
│   ├── Layout/        #   主布局
│   │   ├── Sidebar/   #     侧边栏
│   │   ├── TopBar/    #     顶部导航栏
│   │   └── PageContainer/ # 页面容器
│   ├── PermissionButton/ # 权限按钮
│   ├── PermissionGuard/  # 权限守卫
│   ├── ProgressBar/   #   进度条
│   ├── SearchForm/    #   搜索表单
│   ├── StatCard/      #   统计卡片
│   └── StatusTag/     #   状态标签
├── config/            # 配置文件
│   ├── menu.ts        #   菜单配置
│   └── permission.ts  #   权限配置
├── hooks/             # 自定义 Hooks
├── pages/             # 页面（按功能模块划分）
│   ├── Dashboard/     #   工作台
│   ├── DataOverview/  #   数据概览
│   ├── DataReport/    #   数据报表
│   ├── DataQuality/   #   数据质量
│   ├── DataAssets/    #   数据资产
│   ├── ETLTasks/      #   ETL 任务管理
│   ├── DataModeling/  #   数据建模
│   ├── ScriptDev/     #   脚本开发
│   ├── DataInterface/ #   数据接口
│   ├── DataSubscription/ # 数据订阅
│   ├── MetadataMgmt/  #   元数据管理
│   ├── DataStandard/  #   数据标准
│   ├── DataSecurity/  #   数据安全
│   ├── AlertRules/    #   告警规则
│   ├── MonitorDashboard/ # 监控大盘
│   ├── UserManagement/ #  用户管理
│   ├── PermissionManagement/ # 权限管理
│   ├── Organization/  #   组织架构
│   ├── Profile/       #   个人中心
│   ├── ContentManagement/ # 内容管理
│   ├── Workflow/      #   工作流
│   ├── APIManagement/ #   API 管理
│   ├── OperationLog/  #   操作日志
│   ├── AuditLog/      #   审计日志
│   └── SystemSettings/ #  系统设置
├── router/            # 路由配置
├── services/          # 数据服务（Mock）
├── stores/            # 全局状态管理
├── styles/            # 全局样式
│   ├── global.css     #   全局样式 & Ant Design 覆盖
│   └── theme.css      #   CSS 变量主题系统
├── types/             # 全局类型定义
├── utils/             # 工具函数
│   └── request.ts     #   Axios HTTP 请求封装
└── main.tsx           # 应用入口
```

## 架构说明

### 1. 布局系统

整体采用 Ant Design Layout 实现**侧边栏 + 顶部栏 + 内容区**三栏布局：

- **Sidebar**（侧边栏）：展示多级菜单，支持折叠/展开，自动按权限过滤菜单项；底部展示用户头像和通知
- **TopBar**（顶部栏）：显示面包屑导航（根据当前路由自动匹配菜单层级）、搜索入口、帮助按钮
- **PageContainer**（页面容器）：包裹内容区域，提供滚动容器和内边距

### 2. 路由系统

使用 `react-router-dom` v7 的 `createBrowserRouter` 创建扁平路由表，所有页面统一挂在 `/` 路径下：

- 根路径自动重定向到 `/dashboard`
- 每个路由通过 `PermissionGuard` 组件包裹，实现路由级权限控制
- 未知路径重定向到 `/dashboard`

### 3. 权限体系

权限系统基于 `角色 → 权限项 → 菜单项` 三层模型：

- **权限配置**：[src/config/permission.ts](/Users/shenfay/Projects/react-admin-scaffold/src/config/permission.ts) — 定义角色的权限列表和可见菜单
- **路由级守卫**：[PermissionGuard](/Users/shenfay/Projects/react-admin-scaffold/src/components/PermissionGuard/index.tsx) — 页面级别的权限拦截
- **组件级权限**：[PermissionButton](/Users/shenfay/Projects/react-admin-scaffold/src/components/PermissionButton/index.tsx) — 按钮级别的权限控制
- **菜单过滤**：[Sidebar](/Users/shenfay/Projects/react-admin-scaffold/src/components/Layout/Sidebar/index.tsx) — 根据用户 `menus` 列表过滤菜单项
- **状态管理**：[Zustand store](/Users/shenfay/Projects/react-admin-scaffold/src/stores/index.ts) — 持久化用户信息和权限状态到 localStorage

内置角色：
| 角色 | 标识 | 权限范围 |
|------|------|----------|
| 管理员 | `admin` | 全部功能 |
| 编辑 | `editor` | 工作台、数据概览/报表、内容管理、工作流 |
| 成员 | `member` | 工作台、数据概览 |

### 4. 菜单配置

菜单定义在 [src/config/menu.ts](/Users/shenfay/Projects/react-admin-scaffold/src/config/menu.ts)，采用树形结构：

```typescript
interface MenuItem {
  key: string        // 菜单唯一标识
  label: string      // 显示名称
  icon?: ReactNode   // 图标（字符串形式，通过 getIcon 函数解析）
  path?: string      // 路由路径（叶子节点必填）
  children?: MenuItem[]  // 子菜单
  badge?: number     // 徽标数字
  permission?: string    // 所需权限标识
}
```

菜单按功能域分为以下模块：
| 模块 | 包含页面 |
|------|----------|
| 概览 | 工作台 |
| 数据中台 | 数据概览、数据报表、数据质量、数据资产 |
| 数据开发 | ETL 任务管理、数据建模、脚本开发 |
| 数据服务 | 数据接口、数据订阅 |
| 数据治理 | 元数据管理、数据标准、数据安全 |
| 监控告警 | 告警规则、监控大盘 |
| 用户中心 | 用户管理、权限管理、组织架构、个人中心 |
| 业务中台 | 内容管理、工作流、API 管理 |
| 系统 | 操作日志、审计日志、系统设置 |

新增页面时需同步配置：菜单项 → 路由 → 权限项。

### 5. 状态管理

基于 Zustand 实现，支持 `persist` 中间件持久化到 localStorage：

- **用户状态**：`userId`, `username`, `avatar`, `role`, `roleName`, `permissions`, `menus`, `isLogin`
- **布局状态**：`sidebarCollapsed`（侧边栏折叠状态）
- **权限状态**：`permissionConfig`（权限配置），支持从后端动态加载

首次加载自动以管理员身份登录（开发模式），实际项目接入时需替换为真实登录流程。

### 6. HTTP 请求

[src/utils/request.ts](/Users/shenfay/Projects/react-admin-scaffold/src/utils/request.ts) 基于 Axios 封装：

- **请求拦截器**：自动携带 `Bearer token`（从 localStorage 读取）
- **响应拦截器**：统一处理响应码，401 自动跳转登录，403/404/500 给出错误提示
- **基础配置**：`baseURL` 通过 `VITE_API_BASE_URL` 环境变量配置，超时 30s

### 7. Mock 数据

[src/services/mock.ts](/Users/shenfay/Projects/react-admin-scaffold/src/services/mock.ts) 提供开发环境 Mock 服务：

- `mockLogin()` — 模拟登录（admin/editor/member）
- `mockGetPermissions()` — 获取权限配置
- `mockGetStats()` — 获取统计数据
- `mockGetActivities()` — 获取最近活动
- `mockGetServices()` — 获取服务状态
- `mockGetDataSources()` — 获取数据源
- `mockGetUsers()` — 获取用户列表

## 组件库

### 通用组件

| 组件 | 用途 | 路径 |
|------|------|------|
| [StatCard](/Users/shenfay/Projects/react-admin-scaffold/src/components/StatCard/index.tsx) | 统计卡片 — 带趋势变化指示，支持正/负值 |
| [StatusTag](/Users/shenfay/Projects/react-admin-scaffold/src/components/StatusTag/index.tsx) | 状态标签 — 根据状态值自动匹配颜色主题 |
| [DataPanel](/Users/shenfay/Projects/react-admin-scaffold/src/components/DataPanel/index.tsx) | 数据面板卡片 — 封装 Ant Design Card，统一标题/内边距 |
| [ProgressBar](/Users/shenfay/Projects/react-admin-scaffold/src/components/ProgressBar/index.tsx) | 进度条 — 支持 normal/warning/exception 状态 |
| [PermissionButton](/Users/shenfay/Projects/react-admin-scaffold/src/components/PermissionButton/index.tsx) | 权限按钮 — 无权限时隐藏或显示 fallback |
| [PermissionGuard](/Users/shenfay/Projects/react-admin-scaffold/src/components/PermissionGuard/index.tsx) | 权限守卫 — 页面级权限拦截，无权限时显示提示或 fallback |

> `DataTable` 和 `SearchForm` 为预留组件目录，可根据业务需求封装通用的表格和搜索表单组件。

## 主题系统

CSS 变量集中定义在 [src/styles/theme.css](/Users/shenfay/Projects/react-admin-scaffold/src/styles/theme.css)，覆盖以下维度：

- **颜色**：侧边栏、主背景、文字层级、边框、强调色、状态色（绿/红/黄/蓝/灰）
- **圆角**：sm(8px) / md(12px) / lg(16px)
- **阴影**：sm / md / lg
- **字体**：系统字体栈
- **过渡**：统一缓动曲线

[src/styles/global.css](/Users/shenfay/Projects/react-admin-scaffold/src/styles/global.css) 包含全局样式重置和 Ant Design 组件样式覆盖。

## 构建与部署

```bash
# 生产构建
npm run build

# 构建产物输出到 dist/ 目录
# 可通过 npm run preview 本地预览
```

环境变量：
| 变量 | 说明 | 默认值 |
|------|------|--------|
| `VITE_API_BASE_URL` | API 基础路径 | `/api` |

## 开发规范

- **路径别名**：使用 `@/` 替代相对路径（如 `@/components`, `@/pages`）
- **目录结构**：严格遵循模块化划分
- **UI 组件**：采用薄封装 wrapper 模式
- **状态管理**：全局状态统一使用 Zustand store
- **提交规范**：使用 `npm run commit`（Commitizen）规范提交信息
