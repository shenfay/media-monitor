"""MediaMonitor Python 抓取服务。

架构：
  Go 后端 ──XADD crawl:task:dispatch──▶ Python 常驻消费（serve）
  Go 后端 ◀──XADD crawl:article:ingest + crawl:task:event── Python

- run 子命令：直跑某个 source，本地/测试用，可落 JSON。
- serve 子命令：常驻消费 Redis Stream，按 task 调对应 adapter，回写数据。
"""
