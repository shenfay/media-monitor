-- 抓取任务运行表（状态由 Stream 事件驱动）
CREATE TABLE IF NOT EXISTS crawl_task_runs (
    id           VARCHAR(26) PRIMARY KEY,
    source_id    VARCHAR(26) NOT NULL,
    task_type    VARCHAR(30) DEFAULT '',
    triggered_by VARCHAR(20) DEFAULT '',
    status       VARCHAR(20) NOT NULL DEFAULT 'queued',
    total        INTEGER NOT NULL DEFAULT 0,
    ingested     INTEGER NOT NULL DEFAULT 0,
    failed       INTEGER NOT NULL DEFAULT 0,
    error        TEXT,
    started_at   TIMESTAMP WITH TIME ZONE,
    finished_at  TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_crawl_task_runs_source_id ON crawl_task_runs(source_id);
CREATE INDEX IF NOT EXISTS idx_crawl_task_runs_status ON crawl_task_runs(status);
