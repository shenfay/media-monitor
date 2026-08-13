# 管理后台前端架构文档

本文档详细说明管理后台（admin）的前端架构设计，包括技术栈、项目结构、路由权限、状态管理、国际化、组件体系和样式规范。

## 目录

- [技术栈概览](#技术栈概览)
- [项目结构](#项目结构)
- [路由与权限](#路由与权限)
- [状态管理](#状态管理)
- [国际化 i18n](#国际化-i18n)
- [组件体系](#组件体系)
- [请求层](#请求层)
- [样式体系](#样式体系)

---

## 技术栈概览

| 类别 | 技术 | 版本 | 说明 |
|------|------|------|------|
| 构建工具 | Vite | 6+ | 开发服务器 + 打包 |
| UI 框架 | React | 19 | 函数组件 + Hooks |
| 语言 | TypeScript | 5+ | 严格模式 |
| 组件库 | Ant Design | 5.x | 基础 UI 组件 |
| 高级表格 | @ant-design/pro-components | - | ProTable 列表页模式 |
| 路由 | react-router-dom | 7.x | BrowserRouter |
| 状态管理 | Zustand | 5.x | 轻量级 Store |
| 国际化 | react-i18next + i18next | - | 多语言支持 |
| HTTP 请求 | axios | - | 拦截器 + 请求取消 |
| CSS 方案 | CSS Variables + 全局覆盖 | - | 主题变量体系 |

---

## 项目结构

```
admin/src/
├── assets/              # 静态资源（图片、SVG）
├── components/          # 通用组件
│   ├── DataPanel/       # 数据面板（页面容器，含标题区+筛选栏+内容区）
│   ├── ErrorBoundary/   # 错误边界
│   ├── Layout/          # 布局组件（Sidebar + TopBar + PageContainer）
│   ├── PermissionButton/# 权限按钮（根据权限标识控制按钮可见性）
│   ├── PermissionGuard/ # 路由权限守卫
│   ├── ProgressBar/     # 进度条
│   ├── StatCard/        # 统计卡片
│   └── StatusTag/       # 状态标签
├── config/              # 配置文件
│   ├── menu.ts          # 菜单图标映射（菜单结构从后端动态获取）
│   ├── pagination.ts    # 表格分页默认配置
│   └── permission.ts    # 权限工具函数
├── hooks/               # 自定义 Hooks
│   ├── useCrudList.ts   # 列表页 CRUD 通用 Hook
│   └── useWebSocket.tsx  # WebSocket 连接管理 + 推送通知 Hook
├── locales/             # 国际化语言包
│   ├── index.ts         # i18n 初始化配置
│   ├── zh-CN.json       # 中文语言包
│   └── en-US.json       # 英文语言包
├── pages/               # 页面组件（按业务模块分目录）
├── router/              # 路由配置
│   └── index.tsx        # 路由定义 + 权限守卫
├── services/            # API 服务层
├── stores/              # Zustand Store
│   ├── index.ts         # Store 统一导出
│   ├── useUserStore.ts  # 用户信息、Token、权限
│   └── useLayoutStore.ts# 侧边栏折叠、布局状态
├── styles/              # 全局样式
│   ├── global.css       # 全局样式覆盖
│   └── theme.css        # CSS 变量主题定义
├── types/               # TypeScript 类型定义
│   └── index.ts
├── utils/               # 工具函数
│   ├── formRules.ts     # 表单验证规则
│   └── request.ts       # axios 封装
├── App.tsx              # 根组件（ConfigProvider + Router）
└── main.tsx             # 入口文件
```

---

## 路由与权限

### 路由结构

使用 `react-router-dom` 的 `createBrowserRouter`，所有业务页面嵌套在 `MainLayout` 中，通过 `Suspense` + `lazy` 实现懒加载。

```
/login                    → 登录页（同步加载）
/                         → MainLayout
  /dashboard              → 工作台
  /family                 → 家庭管理
  /goals                  → 目标管理
  /card-templates         → 卡片模板
  /card-instances         → 卡片实例
  /companions             → 伙伴管理
  /acceptance             → 验收管理
  /points                 → 积分记录
  /shop-items             → 商品管理
  /exchange-orders        → 兑换订单
  /users                  → 用户管理
  /permissions            → 权限管理
  /menus                  → 菜单管理
  /profile                → 个人资料
  /operation-log          → 操作日志
  /design-system          → 设计规范展示
  /settings               → 系统设置
  /messages               → 我的消息
  /ws-test                → WebSocket 测试（开发调试）
```

### 权限守卫

每个路由页面使用 `PermissionGuard` 组件包裹，传入权限标识：

```tsx
<PermissionGuard permission="user:manage">
  <UserManagement />
</PermissionGuard>
```

**权限检查流程**：
1. 用户登录后，后端返回权限列表（`permissions: string[]`）
2. 权限存入 `useUserStore`
3. `PermissionGuard` 从 Store 读取权限，与路由要求的权限比对
4. 无权限时渲染空内容或跳转

### 动态菜单

菜单结构**不再硬编码在前端**，而是从后端 API 动态获取：

```
登录 → 获取用户权限+菜单 → 存入 useUserStore
                            ↓
              SidebarMenu 组件渲染菜单树
                            ↓
              菜单图标通过 config/menu.ts 的 iconMap 映射
```

**图标映射**：`config/menu.ts` 维护 Ant Design 图标名称到 React 组件的映射表，后端返回的菜单数据中 `icon` 字段对应图标名称。

---

## 状态管理

采用 Zustand 拆分模式，按职责分为两个独立 Store：

### useUserStore

**文件**：`stores/useUserStore.ts`

管理用户认证和权限状态，通过 `persist` 中间件持久化到 `localStorage`。

| 状态 | 类型 | 说明 |
|------|------|------|
| `userId` | string \| null | 用户 ID |
| `username` | string \| null | 用户名 |
| `email` | string \| null | 邮箱 |
| `roles` | RoleBrief[] | 角色列表 |
| `permissions` | string[] | 权限标识列表 |
| `menus` | string[] | 菜单标识列表 |
| `menuTree` | MenuItem[] | 菜单树结构 |
| `isLogin` | boolean | 登录状态 |

**关键 Action**：
- `login()` — 登录，存储 Token 到 localStorage
- `logout()` — 登出，清理 localStorage + 重置状态
- `updatePermissions()` — 更新权限（菜单刷新后调用）
- `setMenuTree()` — 设置动态菜单树

**持久化配置**：
- Storage Key: `admin-user-storage`
- 持久化字段: userId, username, email, avatar, roles, permissions, menus, menuTree, isLogin

### useLayoutStore

**文件**：`stores/useLayoutStore.ts`

管理 UI 布局状态，同样持久化到 `localStorage`。

| 状态 | 类型 | 说明 |
|------|------|------|
| `sidebarCollapsed` | boolean | 侧边栏是否折叠 |

**持久化配置**：
- Storage Key: `admin-layout-storage`
- 持久化字段: sidebarCollapsed

---

## 国际化 i18n

### 技术栈

- `i18next` — 核心国际化框架
- `react-i18next` — React Hooks 集成（`useTranslation`）
- `i18next-browser-languagedetector` — 自动语言检测

### 配置

**文件**：`locales/index.ts`

| 配置项 | 值 | 说明 |
|--------|-----|------|
| `fallbackLng` | `zh-CN` | 默认语言 |
| `supportedLngs` | `['zh-CN', 'en-US']` | 支持的语言 |
| `detection.order` | `['localStorage', 'navigator']` | 语言检测优先级 |
| `detection.lookupLocalStorage` | `admin-lang` | localStorage 中的语言 Key |

### 语言包结构

国际化文件位于 `locales/` 目录，按业务模块组织 Key：

```json
// locales/zh-CN.json（约 250+ key）
{
  "dashboard": "工作台",
  "userManagement": "用户管理",
  "totalRecords": "共 {{total}} 条记录",
  "sessionExpired": "登录已过期，请重新登录",
  // ...

  // 消息模块
  "menuMessage": "消息管理",
  "menuWsTest": "WebSocket 测试",
  "myMessages": "我的消息",
  "markAsRead": "标记已读",
  "markAllAsRead": "全部已读",
  "unreadCount": "未读消息",
  "noMessages": "暂无消息",

  // WebSocket 测试页
  "wsTest": "WebSocket 测试",
  "wsCatPoints": "积分变动",
  "wsCatGoal": "目标相关",
  "wsCatRemind": "提醒通知",
  "wsCatExchange": "兑换通知",
  "wsCatCompanion": "伙伴相关",
  "wsCatReview": "验收通知"
}
```

### 使用方式

```tsx
import { useTranslation } from 'react-i18next'

function MyComponent() {
  const { t } = useTranslation()
  return <h1>{t('dashboard')}</h1>
}
```

### Ant Design 集成

在 `App.tsx` 中通过 `ConfigProvider` 注入 Ant Design 的 locale：

```tsx
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'

<ConfigProvider locale={currentLang === 'zh-CN' ? zhCN : enUS}>
  ...
</ConfigProvider>
```

---

## 组件体系

### 通用组件

#### DataPanel（数据面板）

**文件**：`components/DataPanel/index.tsx`

页面级容器组件，提供统一的页面结构：

```
+------------------------------------------+
| 标题区（title + extra）                    |
+------------------------------------------+
| 筛选栏（filters + toolbarActions）         |
+------------------------------------------+
| 内容区（表格 / 卡片）                      |
+------------------------------------------+
```

**Props**（见组件源码）

#### WebSocket Hooks

**文件**：`hooks/useWebSocket.tsx`

提供 WebSocket 连接管理和实时推送通知能力：

| Hook | 用途 |
|------|------|
| `useWebSocket()` | 建立 WebSocket 连接（自动重连 + 指数退避） |
| `useWebSocketPush(callback)` | 收到推送时执行回调（用于刷新未读数等） |
| `usePushNotification(callback)` | 收到推送时弹出 Ant Design 通知 Toast + 执行回调 |

**连接管理**：
- 自动连接：`useWebSocket()` 在 `MainLayout` 中调用，页面加载即连接
- 自动重连：断线后间隔递增重试（1s → 2s → 4s → 8s → 16s）
- 连接销毁：页面关闭 / 组件卸载时自动关闭

**WebSocket URL**：通过 Vite 代理 `/api/ws` → 后端 `GET /api/ws`

**签名机制**：连接 URL 携带 `token` 参数，后端验证 JWT 后建立连接

#### PermissionGuard（权限守卫）

**文件**：`components/PermissionGuard/index.tsx`

根据用户权限标识控制子组件的渲染：

```tsx
<PermissionGuard permission="user:manage">
  <Button>管理用户</Button>
</PermissionGuard>
```

#### PermissionButton（权限按钮）

**文件**：`components/PermissionButton/index.tsx`

带权限控制的按钮，无权限时不渲染。

### 列表页模式：ProTable + useCrudList

所有管理列表页采用统一模式：

**ProTable**：来自 `@ant-design/pro-components`，提供搜索、筛选、表格、分页一体化方案。

**useCrudList Hook**：封装列表页通用逻辑。

```tsx
const {
  loading, dataSource, total, page, pageSize,
  handlePageChange, fetchData,
  isModalOpen, editingItem, form,
  handleAdd, handleEdit, handleCancel, handleSubmit,
} = useCrudList(fetchApi, { defaultPageSize: 20 })
```

**返回值**：
| 字段 | 说明 |
|------|------|
| `loading` | 加载状态 |
| `dataSource` | 列表数据 |
| `total` | 总记录数 |
| `page` / `pageSize` | 当前页 / 每页条数 |
| `handlePageChange` | 分页变更处理 |
| `fetchData` | 手动刷新 |
| `isModalOpen` | 弹窗开关 |
| `editingItem` | 当前编辑项（null 表示新增） |
| `form` | Ant Design Form 实例 |
| `handleAdd` | 打开新增弹窗 |
| `handleEdit` | 打开编辑弹窗 |
| `handleCancel` | 关闭弹窗并重置表单 |
| `handleSubmit` | 通用提交（自动区分新增/编辑） |

### Layout 组件

```
components/Layout/
├── index.tsx           # MainLayout（侧边栏 + 顶栏 + 内容区）
├── Sidebar/
│   ├── index.tsx       # 侧边栏容器
│   ├── SidebarLogo.tsx # Logo 区域
│   ├── SidebarMenu.tsx # 动态菜单（从 useUserStore 读取 menuTree）
│   └── SidebarUser.tsx # 用户信息区域
├── TopBar/
│   └── index.tsx       # 顶部栏（面包屑、语言切换、刷新按钮）
└── PageContainer/
    └── index.tsx       # 页面容器（统一内边距）
```

---

## 请求层

### axios 封装

**文件**：`utils/request.ts`

**核心特性**：
- 统一 baseURL：`/api`（通过 `VITE_API_BASE_URL` 环境变量配置）
- 请求超时：30s
- Token 自动注入：从 `localStorage` 读取 `admin-token`
- 请求取消：路由切换时取消所有进行中的请求
- 响应拦截：自动解包 `ApiResponse` 结构（提取 `data` 字段）
- 401 处理：自动清理 Token + 跳转登录页

### 请求取消机制

```ts
// 路由切换时调用
export function cancelAllRequests() {
  if (cancelTokenSource) {
    cancelTokenSource.cancel('路由切换，取消请求')
    cancelTokenSource = null
  }
}
```

在 `App.tsx` 的路由监听中调用，避免页面切换后旧请求的结果覆盖新页面数据。

### 错误处理

| HTTP 状态 | 处理方式 |
|-----------|---------|
| 401（非登录页） | 清理 Token → 跳转登录页 |
| 401（登录页） | 返回错误，由登录页显示提示 |
| 其他错误 | 显示后端返回的 `message` 字段 |
| 网络错误 | 显示「网络错误」提示 |

---

## 样式体系

### CSS 变量主题

**文件**：`styles/theme.css`

定义全局 CSS 变量，覆盖布局、色彩、圆角、阴影等设计 Token：

```css
:root {
  /* 布局 */
  --sidebar-bg: #F5F3EF;
  --sidebar-expanded: 224px;
  --sidebar-collapsed: 56px;
  --main-bg: #FFFFFF;

  /* 品牌色 */
  --brand-dark: #2b2b2b;
  --brand-dark-hover: #4d4d4d;

  /* 文字色 */
  --text-primary: #2b2b2b;
  --text-secondary: #6b6258;
  --text-muted: #b0a89a;

  /* 边框 */
  --border-color: #e8e2d8;
  --border-hover: #d4cdc0;
  --border-light: #efeae2;

  /* 圆角 */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;

  /* ...更多变量 */
}
```

### 全局样式覆盖

**文件**：`styles/global.css`

覆盖 Ant Design 默认样式，使其符合设计规范：
- 按钮高度统一为 34px
- 深色按钮移除阴影
- 表格操作按钮统一为 link 模式
- 分页激活项样式定制
- 表格行分隔线样式

### 设计规范参考

完整的 UI 设计规范详见 [UI_DESIGN_SPEC.md](../UI_DESIGN_SPEC.md)，涵盖：
- 布局规范（TopBar、Sidebar、页面结构）
- 字体与排版
- 色彩体系（中性色、品牌色、状态色）
- 按钮体系（Primary、Default、Text、表格操作按钮）
- 表格规范（表头、行、分页）
- 标签体系
- 表单规范
- 筛选栏与工具栏
- 交互规范（hover 状态、过渡时间）

---

## 延伸阅读

- [UI 设计规范](../UI_DESIGN_SPEC.md) - 完整视觉设计规范
- [后端领域模型](../backend/architecture/DOMAIN_MODEL.md) - RBAC 权限模型
- [数据库设计](../backend/database/SCHEMA_DESIGN.md) - 菜单表与权限表结构
