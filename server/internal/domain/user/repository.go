package user

import "context"

// UserListParams 用户列表查询参数
type UserListParams struct {
	Page     int
	PageSize int
	Keyword  string
	RoleID   string
	Status   *bool
}

// UserListResult 用户列表结果
type UserListResult struct {
	Users []*User
	Total int64
}

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// Save 保存用户（创建或更新）
	Save(ctx context.Context, user *User) error

	// FindByID 根据 ID 查找用户
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) bool

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// List 分页查询用户列表
	List(ctx context.Context, params UserListParams) (*UserListResult, error)
}
