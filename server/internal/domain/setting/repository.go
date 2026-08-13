package setting

import "context"

// Repository 系统设置仓储接口
type Repository interface {
	// FindAll 获取所有设置
	FindAll(ctx context.Context) ([]*Setting, error)

	// FindByCategory 按分类获取设置
	FindByCategory(ctx context.Context, category string) ([]*Setting, error)

	// FindByKey 根据 key 获取单个设置
	FindByKey(ctx context.Context, key string) (*Setting, error)

	// BatchUpsert 批量更新/插入设置
	BatchUpsert(ctx context.Context, settings []*Setting) error
}
