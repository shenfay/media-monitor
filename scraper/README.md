# MediaMonitor 抓取服务

Go 后端（`server/`）负责调度、存储、UI、权限；本目录是**抓取执行层**，只做"按任务抓数据并回写"。

## 架构

```
Go 后端 ──XADD crawl:task:dispatch──▶ Python 常驻消费（serve）
Go 后端 ◀──XADD crawl:article:ingest + crawl:task:event── Python
```

- 任务总线：**Redis Stream**（消费者组，支持分布式部署）
- Python 只需一个连接：**Redis**，不需要知道 Go 地址或数据库凭证
- 适配器自动注册：新增数据源只需在 `adapters/` 下建子包

## 目录

```
scraper/
├── src/scrapers/
│   ├── contracts/      # 跨服务数据契约（Source / Article / TaskEvent）
│   ├── streams/        # Redis Stream I/O（consumer / publisher / event_emitter）
│   ├── core/           # 业务编排（TaskExecutor / AdapterRegistry）
│   ├── adapters/       # 数据源适配器（自动注册）
│   ├── fetchers/       # HTTP / 浏览器获取层
│   └── cleaning/       # 数据清洗与标准化
├── tests/
├── pyproject.toml
├── Dockerfile
└── Makefile
```

## 使用

```bash
cd scraper

# 安装依赖
make install        # 生产
make dev            # 开发（含 pytest / ruff / mypy）

# 直跑（本地测试）
make run
python -m scrapers run --source huanqiu --limit 10 --out hq.json

# 常驻消费 Redis Stream（生产形态）
make serve
python -m scrapers serve

# Docker 部署
make docker
docker run -d -e MM_REDIS_URL=redis://redis:6379/0 media-monitor-scraper
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MM_REDIS_URL` | `redis://localhost:6379/0` | Redis 连接 |
| `MM_CONSUMER_NAME` | `scraper-1` | 消费者名（分布式部署时每个实例唯一） |
| `MM_BATCH_SIZE` | `20` | 每批回传文章数 |
| `MM_MAX_STREAM_LEN` | `500` | Stream 最大长度（近似） |

## 加新数据源

1. `adapters/` 下新建子目录，如 `adapters/weibo/`
2. 实现 `adapter.py`，继承 `BaseAdapter`
3. 在 `__init__.py` 中调用 `register("weibo", WeiboAdapter)`
4. 无需修改任何其他文件

## 依赖

- 必需：`redis>=5.0.0`、`typer>=0.12.0`、`requests>=2.31.0`
- 可选：`httpx`（更快的 HTTP 后端）、`crawlee + playwright`（浏览器渲染）
