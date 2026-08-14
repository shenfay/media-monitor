# MediaMonitor Python 抓取层

Go 后端（`server/`）负责调度、存储、UI、权限；本目录是**抓取执行层**，只做"按任务抓数据并回写"。

## 架构边界（与 `integration_design.md` 对齐）

```
Go 后端 ──XADD crawl:tasks──▶ Python 常驻消费（serve）
Go 后端 ◀──POST /ingest/articles + /crawl/tasks/:id/status (X-Scraper-Key)── Python
```

- 任务总线：**Redis Stream** `crawl:tasks`（消费者组，`asynq` 不参与抓取任务）
- 回写鉴权：服务间 **`X-Scraper-Key`**
- Python：**常驻服务**消费队列，复用会话/代理/浏览器实例

## 目录

```
scripts/
├── scrapers/
│   ├── config.py          # 环境变量配置（MM_* 前缀）
│   ├── models.py          # Source / Article / TaskMessage / TaskStatus（契约）
│   ├── fetchers/http.py   # HTTP 获取（urllib，自动优先 httpx）
│   ├── adapters/          # 各数据源适配器
│   │   ├── base.py        #   BaseSourceAdapter 抽象
│   │   └── huanqiu.py     #   环球网（PHASE 1 列表 + PHASE 2 正文）
│   ├── sinks/go_client.py # 回写 Go 的 HTTP 客户端
│   ├── worker.py          # Redis Stream 常驻消费循环
│   └── cli.py             # run / serve 命令
└── requirements.txt
```

## 使用

```bash
cd scripts
export MM_GO_API_BASE=http://localhost:8080
export MM_SCRAPER_API_KEY=xxxx          # 与 Go 侧一致
export MM_REDIS_URL=redis://localhost:6379/0

# 直跑（本地/测试）：拉环球网近半年原创，落 JSON
python -m scrapers.cli run --source huanqiu --limit 100 --out hq.json

# 直跑并抓正文、回写 Go
python -m scrapers.cli run --source huanqiu --with-body --publish

# 常驻消费 Redis Stream（生产形态）
python -m scrapers.cli serve
```

## 加新数据源

1. `adapters/` 下新建 `<site>.py`，继承 `BaseSourceAdapter`，实现 `fetch_list`（必），`fetch_detail`（可选）。
2. 在 `worker.py` 的 `ADAPTERS` 注册 `平台类型/source.name -> 类`。
3. 社媒类（需登录/渲染）后续接入 `crawlee + playwright` 作为 fetch 后端。

## 依赖

- `run` 模式：仅标准库，零依赖。
- `serve` 模式：需 `redis`（消费 Stream）。
- 社媒适配器（未来）：`crawlee[all]` + `playwright`。
