# 快速开始指南

5 分钟内在本地运行 Kiqi 后端服务。

## 📋 前置要求

### 必需软件

| 软件 | 版本 | 用途 | 安装方式 |
|------|------|------|---------|
| Go | 1.21+ | 后端运行时 | `brew install go` |
| PostgreSQL | 14+ | 主数据库 | `brew install postgresql` |
| Redis | 7+ | 缓存/会话/实时消息总线 | `brew install redis` |

### 可选软件

| 软件 | 用途 | 安装方式 |
|------|------|---------|
| Docker & Docker Compose | 容器化部署 | [Docker Desktop](https://www.docker.com/products/docker-desktop/) |
| Swag | Swagger 文档生成 | `go install github.com/swaggo/swag/cmd/swag@latest` |
| Prometheus | 指标采集和存储 | `brew install prometheus` |
| Grafana | 可视化仪表盘 | `brew install grafana` |
| MailHog | 开发环境邮件捕获 | `brew install mailhog` |

## 🚀 快速启动（推荐方式）

### 方式一：使用 Docker Compose（最简单）

```bash
# 1. 克隆项目
git clone <repository-url>
cd go-react-admin

# 2. 启动所有基础设施（PostgreSQL + Redis）
docker-compose up -d postgres redis

# 3. 等待服务就绪（约 10 秒）
sleep 10

# 4. 进入后端目录
cd server

# 5. 运行数据库迁移
make migrate up

# 6. 启动 API 服务
make run api
```

**验证服务**：
```bash
# 健康检查
curl http://localhost:8080/health

# 访问 Swagger 文档
open http://localhost:8080/swagger/index.html
```

### 方式二：本地安装（开发推荐）

#### 步骤 1：启动基础设施

```bash
# 启动 PostgreSQL
brew services start postgresql

# 启动 Redis（必需：WebSocket 实时推送依赖 Redis Pub/Sub）
brew services start redis

# （可选）启动 MailHog 捕获开发环境邮件
brew services start mailhog

# 验证服务
psql -U postgres -c "SELECT version();"
redis-cli ping  # 应返回 PONG
```

#### 步骤 2：创建数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE admin;"
psql -U postgres -c "CREATE USER admin WITH PASSWORD 'admin';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE admin TO admin;"
```

#### 步骤 3：配置环境变量

```bash
cd server
cp configs/.env.example configs/.env

# 编辑配置文件（根据实际情况修改）
vim configs/.env
```

**关键配置项**：
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=admin

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d

# WebSocket 实时推送开关
WEBSOCKET_ENABLED=true
```

#### 步骤 4：安装依赖并运行

```bash
# 安装 Go 依赖
go mod download

# 运行数据库迁移
make migrate up

# 启动 API 服务
make run api
```

## 📝 测试 API

### 1. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123456!"
  }'
```

**预期响应**：
```json
{
  "code": "SUCCESS",
  "message": "注册成功",
  "data": {
    "user": {
      "id": "01JQMXYZ...",
      "email": "test@example.com",
      "email_verified": false
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

        "email_verified": false
      },
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
      "expires_in": 900
    }
  }
}
```

### 2. 测试 WebSocket 推送

WebSocket 默认启用（配置文件 `configs/development.yaml` 中 `websocket.enabled: true`）。

前端管理后台访问 `/ws-test` 可打开 WebSocket 测试页：

```bash
# 通过 Swagger POST /api/v1/admin/messages/system 发送测试推送
curl -X POST http://localhost:8080/api/v1/admin/messages/system \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "recipient_id": "user_founder",
    "category": "points",
    "title": "测试推送",
    "content": "Hello WebSocket"
  }'
```

### 3. 邮箱验证 & 密码重置

注册后可调用以下接口测试（开发环境使用 NoopSender，验证码/重置链接会打印到日志）：

```bash
# 重新发送验证邮件
curl -X POST http://localhost:8080/api/v1/auth/resend-verification \
  -H "Authorization: Bearer <access_token>"

# 查看日志中的验证链接
grep -i "verify\|reset" /path/to/server/logs
```

详细测试步骤参见 [测试指南](../../server/TESTING_GUIDE.md)。

### 4. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123456!"
  }'
```

### 3. 获取当前用户信息

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

## 🛠️ 常用操作

### 启动 Worker（异步任务处理）

```bash
# 新终端窗口
cd server
make run worker
```

### 生成 Swagger 文档

```bash
make swagger gen

# 查看生成的文档
ls -la api/swagger/
# docs.go
# swagger.json
# swagger.yaml
```

### 查看数据库迁移状态

```bash
make db-status
```

### 回滚数据库迁移

```bash
make migrate down
```

### 运行测试

```bash
# 运行所有测试
make test

# 仅运行单元测试（跳过集成测试）
make test-short

# 生成覆盖率报告
make coverage
open coverage.html
```

## 🔍 监控配置（可选）

项目内置 Prometheus + Grafana 监控栈，用于实时监控服务状态。

### 快速启动监控

```bash
# 1. 安装 Prometheus 和 Grafana
brew install prometheus grafana

# 2. 启动 Prometheus（在项目根目录）
prometheus --config.file=prometheus.yml \
  --storage.tsdb.path=/tmp/prometheus-data \
  --web.enable-lifecycle > /tmp/prometheus.log 2>&1 &

# 3. 启动 Grafana
brew services start grafana

# 4. 验证服务
curl -s http://localhost:9090/-/healthy  # Prometheus
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/login  # Grafana
```

### 配置 Grafana 数据源

1. 访问 http://localhost:3000（账号：admin/admin）
2. **Connections** → **Data Sources** → **Add data source**
3. 选择 **Prometheus**
4. URL: `http://localhost:9090`
5. 点击 **Save & test**

### 导入仪表盘

1. **Dashboards** → **Import**
2. 上传 `grafana-dashboard.json`
3. 选择 Prometheus 数据源
4. 点击 **Import**

### 查看监控面板

导入后可看到 6 个分组，25 个面板：
- 📊 **API Overview** - QPS、错误率、延迟
- 🌐 **HTTP Metrics** - 状态码、QPS 趋势
- 🔐 **Authentication** - 认证成功率、失败原因
- 🗄️ **Database** - 连接池使用率
- 📈 **Business Metrics** - 用户注册数
- 📍 **Endpoints** - 按路径分解的指标

> 📖 详细文档：[监控配置指南](../operations/MONITORING_SETUP.md)

## 🐛 常见问题

### 问题 1：数据库连接失败

**错误信息**：
```
failed to connect database: dial tcp [::1]:5432: connect: connection refused
```

**解决方案**：
```bash
# 检查 PostgreSQL 是否运行
brew services list | grep postgresql

# 启动 PostgreSQL
brew services start postgresql

# 验证连接
psql -U postgres -d admin -c "SELECT 1;"
```

### 问题 2：Redis 连接失败

**错误信息**：
```
redis: connection refused
```

**解决方案**：
```bash
# 检查 Redis 是否运行
brew services list | grep redis

# 启动 Redis
brew services start redis

# 验证连接
redis-cli ping
```

### 问题 3：端口已被占用

**错误信息**：
```
listen tcp :8080: bind: address already in use
```

**解决方案**：
```bash
# 查找占用端口的进程
lsof -i :8080

# 终止进程（替换 PID）
kill -9 <PID>

# 或者修改配置使用其他端口
vim server/configs/.env
# 修改 SERVER_PORT=8081
```

### 问题 4：Swag 命令未找到

**错误信息**：
```
zsh: command not found: swag
```

**解决方案**：
```bash
# 安装 Swag
go install github.com/swaggo/swag/cmd/swag@latest

# 添加到 PATH（如需要）
export PATH=$PATH:$(go env GOPATH)/bin
```

### 问题 5：Prometheus 无法采集指标

**症状**：http://localhost:9090/targets 显示 DOWN

**解决方案**：
```bash
# 1. 检查 API 服务是否运行
curl http://localhost:8080/health

# 2. 检查 /metrics 端点
curl http://localhost:8080/metrics | head -20

# 3. 重启 Prometheus
pkill prometheus
prometheus --config.file=prometheus.yml > /tmp/prometheus.log 2>&1 &
```

### 问题 6：Grafana 显示 No Data

**解决方案**：
1. 检查数据源配置（Connections → Data Sources）
2. 测试连接是否正常
3. 调整时间范围为 "Last 1 hour"
4. 触发一些请求生成数据

> 📖 更多问题排查：[故障排查指南](../operations/TROUBLESHOOTING.md)

## 📚 下一步

- [📖 开发指南](DEVELOPMENT_GUIDE.md) - 了解开发规范和流程
- [🏗️ DDD 架构设计](../architecture/DDD_ARCHITECTURE.md) - 深入理解架构设计
- [📡 API 文档](http://localhost:8080/swagger/index.html) - 查看所有 API 接口
- [🗄️ 数据库设计](../database/SCHEMA_DESIGN.md) - 了解数据库结构
- [🔍 监控配置](../operations/MONITORING_SETUP.md) - 配置 Prometheus + Grafana
- [🐛 故障排查](../operations/TROUBLESHOOTING.md) - 常见问题诊断

## 💡 开发建议

1. **使用 IDE**：推荐 GoLand 或 VS Code（安装 Go 插件）
2. **启用自动保存**：避免忘记保存文件
3. **使用 `make run api`**：而非 `go run`，确保环境变量正确加载
4. **定期运行 `make vet`**：及早发现潜在问题
5. **编写测试**：新功能必须包含单元测试

---

**遇到问题？** 查看 [故障排查指南](../operations/TROUBLESHOOTING.md) 或提交 [GitHub Issue](https://github.com/your-org/go-react-admin/issues)
