# Go React Admin

全栈管理后台脚手架，基于 **Go + React** 构建，采用 DDD（领域驱动设计）架构，提供完整的用户认证、RBAC 权限管理、异步任务处理、监控告警等核心功能。

## ✨ 功能特性

### 后端
- **用户认证** — 注册、登录、JWT Token（Access + Refresh）、邮箱验证、密码重置
- **RBAC 权限** — 基于 Casbin 的角色权限控制，动态菜单管理
- **DDD 架构** — 领域层、应用层、基础设施层、传输层清晰分离
- **事件驱动** — 进程内事件总线 + Asynq 异步任务队列
- **操作日志** — 全链路操作审计记录
- **系统设置** — 可配置的系统参数管理
- **消息通知** — 站内消息发送与管理
- **可观测性** — Prometheus 指标采集 + Grafana 可视化仪表盘
- **API 文档** — Swagger 自动生成，在线交互式文档
- **优雅关闭** — 支持信号量处理的优雅停机

### 前端
- **React 19** + **TypeScript** + **Vite** 现代化技术栈
- **Ant Design 6** UI 组件库 + ProComponents 高级组件
- **Zustand** 轻量状态管理
- **i18next** 国际化（中/英双语）
- **动态路由** — 后端驱动菜单，细粒度权限控制
- **页面懒加载** — React.lazy + Suspense 性能优化
- **ErrorBoundary** — 全局错误边界保护
- **请求拦截** — Axios 封装，统一鉴权与错误处理

## 🏗️ 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | React 19, TypeScript, Vite, Ant Design 6, Zustand, i18next, React Router 7 |
| **后端** | Go 1.25, Gin, GORM, Casbin, JWT, Zap |
| **数据库** | PostgreSQL 14+ |
| **缓存** | Redis 7+ |
| **异步任务** | Asynq + Asynqmon |
| **监控** | Prometheus, Grafana |
| **API 文档** | Swagger (swaggo) |
| **数据库迁移** | golang-migrate |

## 📁 项目结构

```
go-react-admin/
├── admin/                    # 前端管理后台
│   ├── src/
│   │   ├── components/       # 通用组件（Layout, PermissionGuard 等）
│   │   ├── config/           # 菜单图标映射、分页、权限配置
│   │   ├── hooks/            # 自定义 Hooks（useCrudList）
│   │   ├── locales/          # 国际化语言文件
│   │   ├── pages/            # 页面组件
│   │   ├── router/           # 路由配置
│   │   ├── services/         # API 请求服务
│   │   ├── stores/           # Zustand 状态管理
│   │   ├── types/            # TypeScript 类型定义
│   │   └── utils/            # 工具函数
│   └── package.json
├── server/                   # 后端服务
│   ├── cmd/
│   │   ├── api/              # API 服务入口
│   │   ├── cli/              # CLI 工具（数据库迁移等）
│   │   ├── docs/             # 文档服务
│   │   └── worker/           # 异步任务 Worker
│   ├── configs/              # 配置文件
│   ├── internal/
│   │   ├── app/              # 应用服务层（认证、管理、通知）
│   │   ├── domain/           # 领域层（用户、RBAC、通知、操作、设置）
│   │   ├── infra/            # 基础设施层（仓储、配置、授权、消息）
│   │   └── transport/        # 传输层（HTTP Handler、路由、中间件）
│   ├── migrations/           # 数据库迁移文件
│   ├── pkg/                  # 公共包（常量、错误、日志、指标、工具）
│   └── Makefile
├── docs/                     # 项目文档
├── scripts/                  # 运维脚本
├── prometheus.yml            # Prometheus 配置
└── grafana-dashboard.json    # Grafana 仪表盘
```

## 🚀 快速开始

### 前置要求

| 软件 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端运行时 |
| Node.js | 24+ | 前端运行时 |
| PostgreSQL | 14+ | 主数据库 |
| Redis | 7+ | 缓存 / 会话存储 |

### 后端启动

```bash
# 1. 进入后端目录
cd server

# 2. 复制并编辑配置文件
cp configs/.env.example configs/.env
# 根据实际环境修改数据库、Redis、JWT 等配置

# 3. 安装依赖
go mod download

# 4. 运行数据库迁移
make migrate up

# 5. 启动 API 服务
make run api
```

服务启动后：
- API 地址：`http://localhost:8080`
- Swagger 文档：`http://localhost:8080/swagger/index.html`
- 健康检查：`curl http://localhost:8080/health`

### 启动 Worker（异步任务处理）

```bash
cd server
make run worker
```

### 前端启动

```bash
# 1. 进入前端目录
cd admin

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev
```

## 🛠️ 常用命令

### 后端（server/）

```bash
# 开发
make run api              # 启动 API 服务
make run worker           # 启动 Worker
make run asynqmon PORT=   # 启动 Asynqmon 任务监控 UI

# 数据库
make migrate up           # 执行数据库迁移
make migrate down         # 回滚数据库迁移
make db-status            # 查看迁移状态

# 构建
make build                # 构建应用
make build-linux          # 交叉编译 Linux 版本
make build-worker         # 构建 Worker

# 测试
make test                 # 运行所有测试
make test-short           # 运行单元测试（跳过集成测试）
make coverage             # 生成测试覆盖率报告

# 代码质量
make fmt                  # 格式化代码
make vet                  # 运行 go vet
make lint                 # 运行代码检查

# 文档
make swagger gen          # 生成 Swagger 文档
make swagger serve        # 生成并预览 Swagger 文档
```

### 前端（admin/）

```bash
npm run dev               # 启动开发服务器
npm run build             # 构建生产版本
npm run lint              # 代码检查
npm run preview           # 预览构建产物
```

## 📊 监控

项目内置 Prometheus + Grafana 监控栈：

```bash
# 启动监控
make start-monitoring

# 或手动启动
prometheus --config.file=prometheus.yml &
# 导入 grafana-dashboard.json 到 Grafana
```

监控面板包含 6 大分组、25 个面板：
- **API Overview** — QPS、错误率、延迟
- **HTTP Metrics** — 状态码分布、QPS 趋势
- **Authentication** — 认证成功率、失败原因
- **Database** — 连接池使用率
- **Business Metrics** — 业务指标
- **Endpoints** — 按路径分解的详细指标

## 📖 文档

| 文档 | 说明 |
|------|------|
| [快速开始指南](docs/backend/development/GETTING_STARTED.md) | 5 分钟运行项目 |
| [DDD 架构设计](docs/backend/architecture/DDD_ARCHITECTURE.md) | 架构理念与分层设计 |
| [领域模型](docs/backend/architecture/DOMAIN_MODEL.md) | 聚合根、实体、值对象、RBAC 权限模型 |
| [数据库设计](docs/backend/database/SCHEMA_DESIGN.md) | 表结构、ER 图、迁移管理 |
| [开发指南](docs/backend/development/DEVELOPMENT_GUIDE.md) | 开发规范与流程 |
| [Docker 部署](docs/backend/deployment/DOCKER_DEPLOYMENT.md) | 容器化部署方案 |
| [监控配置](docs/backend/operations/MONITORING_SETUP.md) | Prometheus + Grafana 配置 |
| [故障排查](docs/backend/operations/TROUBLESHOOTING.md) | 常见问题诊断与解决 |
| [前端架构](docs/admin/ARCHITECTURE.md) | 路由、状态管理、组件体系 |
| [UI 设计规范](docs/UI_DESIGN_SPEC.md) | 布局、色彩、字体、组件规范 |

## 🤝 贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

项目使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范，前端已配置 Commitizen 辅助提交。

## 📄 许可证

[MIT](https://opensource.org/licenses/MIT)
