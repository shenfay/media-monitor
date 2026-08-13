package operation

import "context"

// LogFilter 操作日志查询过滤条件
type LogFilter struct {
	UserID   string // 用户 ID
	Category string // 分类
	Action   string // 操作类型
	Limit    int    // 每页条数
	Offset   int    // 偏移量
}

// LogRepository 操作日志仓储接口
type LogRepository interface {
	// Save 保存操作日志
	Save(ctx context.Context, log *OperationLog) error

	// FindWithFilter 根据过滤条件查询日志（统一查询入口）
	FindWithFilter(ctx context.Context, filter LogFilter) ([]*OperationLog, error)

	// Count 统计日志总数（统一使用 LogFilter 过滤）
	Count(ctx context.Context, filter LogFilter) (int64, error)
}
