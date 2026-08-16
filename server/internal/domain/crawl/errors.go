package crawl

import "errors"

// 领域错误 — 仓储层在 FindByID 等场景返回，应用层与 Handler 层统一使用。
var (
	ErrSourceNotFound  = errors.New("crawl: source not found")
	ErrArticleNotFound = errors.New("crawl: article not found")
	ErrTaskNotFound    = errors.New("crawl: task not found")
)
