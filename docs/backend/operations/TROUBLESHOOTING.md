# 故障排查指南

本文档提供常见问题诊断和解决方案。

## 📋 目录

- [服务启动问题](#服务启动问题)
- [数据库问题](#数据库问题)
- [Redis 问题](#redis-问题)
- [API 问题](#api-问题)
- [监控问题](#监控问题)
- [性能问题](#性能问题)
- [日志查看](#日志查看)

## 服务启动问题

### API 服务无法启动

#### 症状 1：端口被占用

```
listen tcp :8080: bind: address already in use
```

**解决方案**：
```bash
# 查找占用端口的进程
lsof -ti:8080

# 杀掉进程
kill -9 $(lsof -ti:8080)

# 或者使用其他端口
export SERVER_PORT=8081
./bin/api
```

#### 症状 2：数据库连接失败

```
failed to connect to database
```

**排查步骤**：
```bash
# 1. 检查 PostgreSQL 是否运行
brew services list | grep postgresql

# 2. 检查数据库配置
cat backend/configs/.env | grep DB_

# 3. 测试数据库连接
psql -h localhost -U admin -d admin -c "SELECT 1"

# 4. 检查数据库是否存在
psql -h localhost -U postgres -c "\l" | grep admin
```

**解决方案**：
```bash
# 启动 PostgreSQL
brew services start postgresql

# 创建数据库
createdb -h localhost -U postgres admin

# 运行迁移
cd server
make migrate up
```

#### 症状 3：Redis 连接失败

```
failed to connect to Redis
```

**排查步骤**：
```bash
# 1. 检查 Redis 是否运行
brew services list | grep redis

# 2. 测试 Redis 连接
redis-cli ping
# 预期输出：PONG

# 3. 检查 Redis 配置
cat backend/configs/.env | grep REDIS_
```

**解决方案**：
```bash
# 启动 Redis
brew services start redis

# 如果设置了密码，配置密码
redis-cli CONFIG SET requirepass "your_password"
```

### Prometheus 无法启动

#### 症状：配置文件错误

```
error loading config from "prometheus.yml"
```

**解决方案**：
```bash
# 验证配置文件
promtool check config prometheus.yml

# 检查 YAML 格式
python3 -c "import yaml; yaml.safe_load(open('prometheus.yml'))"
```

### Grafana 无法启动

#### 症状：端口冲突

```
listen tcp 0.0.0.0:3000: bind: address already in use
```

**解决方案**：
```bash
# 查找占用进程
lsof -ti:3000

# 杀掉进程或修改 Grafana 端口
# 编辑 /opt/homebrew/etc/grafana/grafana.ini
# 修改 http_port = 3001
```

## 数据库问题

### 迁移失败

#### 症状：表已存在

```
ERROR: relation "users" already exists
```

**解决方案**：
```bash
# 查看迁移状态
cd server
go run ./cmd/cli migrate version

# 强制设置版本号
go run ./cmd/cli migrate force 5

# 重新运行
go run ./cmd/cli migrate up
```

> 也可以使用 `make db-status` 查看迁移状态，使用 `make migrate up/down` 执行迁移。

### 查询慢

#### 诊断步骤

```bash
# 1. 查看慢查询日志
psql -h localhost -U admin -d admin -c "
SELECT query, calls, total_time, mean_time 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;
"

# 2. 查看当前运行的查询
psql -h localhost -U admin -d admin -c "
SELECT pid, now() - pg_stat_activity.query_start AS duration, query, state
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;
"

# 3. 查看索引使用情况
psql -h localhost -U admin -d admin -c "
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;
"
```

**优化方案**：
```sql
-- 添加缺失的索引
CREATE INDEX idx_users_email ON users(email);

-- 查看查询计划
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';
```

### 连接池耗尽

#### 症状

```
too many clients already
```

**解决方案**：
```bash
# 1. 查看当前连接数
psql -h localhost -U admin -d admin -c "
SELECT count(*) FROM pg_stat_activity;
"

# 2. 查看最大连接数
psql -h localhost -U admin -d admin -c "
SHOW max_connections;
"

# 3. 调整应用连接池配置
# backend/configs/.env
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
```

## Redis 问题

### 内存使用过高

```bash
# 查看内存使用
redis-cli INFO memory | grep used_memory_human

# 查看键数量
redis-cli DBSIZE

# 查找大键
redis-cli --bigkeys

# 清理过期键
redis-cli INFO stats | grep expired_keys
```

### 缓存命中率低

```bash
# 查看命中率
redis-cli INFO stats | grep -E "keyspace_hits|keyspace_misses"

# 计算命中率
# hit_rate = keyspace_hits / (keyspace_hits + keyspace_misses)
```

## API 问题

### 请求返回 500 错误

**排查步骤**：
```bash
# 1. 查看 API 日志
tail -f /tmp/api.log

# 2. 检查 Trace ID
curl -v http://localhost:8080/api/v1/users/1
# 响应头中包含 X-Trace-ID

# 3. 根据 Trace ID 搜索日志
grep "TRACE_ID_HERE" /tmp/api.log
```

### 认证失败

#### 症状：Token 无效

```json
{
  "code": "AUTH.INVALID_TOKEN",
  "message": "无效的访问令牌"
}
```

**排查**：
```bash
# 1. 检查 JWT Secret 配置
cat backend/configs/.env | grep JWT_SECRET

# 2. 检查 Token 是否过期
# 访问 https://jwt.io 解码 Token

# 3. 刷新 Token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_REFRESH_TOKEN"}'
```

### 请求限流

#### 症状：429 Too Many Requests

```json
{
  "code": "AUTH.RATE_LIMIT_EXCEEDED",
  "message": "请求过于频繁"
}
```

**解决方案**：
```bash
# 调整限流配置
# backend/internal/transport/http/middleware/ratelimit.go
# 修改 rate 和 burst 参数
```

## 监控问题

### Prometheus 无法采集指标

**诊断**：
```bash
# 1. 检查目标状态
curl http://localhost:9090/api/v1/targets | python3 -m json.tool

# 2. 检查 /metrics 端点
curl http://localhost:8080/metrics | head -20

# 3. 查看 Prometheus 日志
tail -f /tmp/prometheus.log
```

### Grafana 显示 No Data

**排查步骤**：
```bash
# 1. 验证数据源
curl -u admin:admin http://localhost:3000/api/datasources | python3 -m json.tool

# 2. 测试查询
curl -G http://localhost:9090/api/v1/query \
  --data-urlencode 'query=up' | python3 -m json.tool

# 3. 检查时间范围
# Grafana 右上角选择 "Last 1 hour"
```

### 指标数据为 0

**原因**：
- API 服务重启后 Counter 归零（正常行为）
- 未触发相关业务操作

**验证**：
```bash
# 触发一些请求
for i in {1..10}; do
  curl -s http://localhost:8080/health > /dev/null
done

# 等待 1 分钟后查看
curl http://localhost:8080/metrics | grep http_requests_total
```

## 性能问题

### CPU 使用率高

```bash
# 1. 查看进程 CPU 使用
top -pid $(lsof -ti:8080)

# 2. 查看 Goroutine 数量
curl http://localhost:8080/metrics | grep go_goroutines

# 3. 生成 CPU profile
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# 4. 分析 profile
go tool pprof -http=:8081 /path/to/profile
```

### 内存泄漏

```bash
# 1. 查看内存使用
curl http://localhost:8080/metrics | grep go_memstats

# 2. 生成 Heap profile
go tool pprof http://localhost:8080/debug/pprof/heap

# 3. 分析内存
go tool pprof -http=:8081 /path/to/heap
```

### 响应慢

```bash
# 1. 使用 ab 压测
ab -n 1000 -c 10 http://localhost:8080/health

# 2. 查看 P99 延迟
# Grafana → P99 Latency 面板

# 3. 分析慢请求
grep "latency_ms" /tmp/api.log | sort -t' ' -k2 -rn | head -20
```

## 日志查看

### API 日志

```bash
# 查看最近日志
tail -f /tmp/api.log

# 搜索错误日志
grep "ERROR" /tmp/api.log

# 搜索特定 Trace ID
grep "TRACE_ID" /tmp/api.log

# 查看慢请求（>1s）
grep "latency_ms" /tmp/api.log | awk '{if ($NF > 1000) print}'
```

### Prometheus 日志

```bash
tail -f /tmp/prometheus.log
```

### 数据库日志

```bash
# macOS (Homebrew)
tail -f /opt/homebrew/var/log/postgresql*.log

# 查看慢查询日志
grep "duration:" /opt/homebrew/var/log/postgresql*.log
```

### Redis 日志

```bash
# macOS (Homebrew)
tail -f /opt/homebrew/var/log/redis*.log
```

## 快速诊断脚本

创建 `scripts/diagnose.sh`：

```bash
#!/bin/bash

echo "🔍 DDD-Scaffold 诊断工具"
echo "━━━━━━━━━━━━━━━━━━━━━━"

# 检查服务
echo -n "API (8080): "
curl -s http://localhost:8080/health > /dev/null && echo "✅" || echo "❌"

echo -n "Prometheus (9090): "
curl -s http://localhost:9090/-/healthy > /dev/null && echo "✅" || echo "❌"

echo -n "Grafana (3000): "
curl -s -o /dev/null -w "" http://localhost:3000/login && echo "✅" || echo "❌"

# 检查指标
echo -n "Metrics endpoint: "
curl -s http://localhost:8080/metrics > /dev/null && echo "✅" || echo "❌"

# 检查数据库
echo -n "Database: "
psql -h localhost -U admin -d admin -c "SELECT 1" > /dev/null 2>&1 && echo "✅" || echo "❌"

# 检查 Redis
echo -n "Redis: "
redis-cli ping > /dev/null 2>&1 && echo "✅" || echo "❌"

echo ""
echo "📊 指标统计:"
echo "  HTTP Requests: $(curl -s http://localhost:8080/metrics | grep '^http_requests_total{' | awk -F' ' '{sum+=$2} END {print sum}')"
echo "  Auth Success: $(curl -s http://localhost:8080/metrics | grep '^auth_success_total' | awk -F' ' '{print $2}')"
echo "  Auth Failure: $(curl -s http://localhost:8080/metrics | grep '^auth_failure_total' | awk -F' ' '{print $2}')"
```

使用：
```bash
chmod +x scripts/diagnose.sh
./scripts/diagnose.sh
```

## 📚 延伸阅读

- [监控配置指南](MONITORING_SETUP.md)
- [Docker 部署指南](../deployment/DOCKER_DEPLOYMENT.md)
- [快速开始指南](../development/GETTING_STARTED.md)
