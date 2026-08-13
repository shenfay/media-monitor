# Go DDD Scaffold - 代码注释规范

## 核心原则

采用**混合分层策略**，不同架构层使用不同详细程度的注释：
- **Handler 层**：Swagger 注解 + 完整中文注释
- **Service 层**：完整中文注释（不含 Swagger）
- **Domain 层**：简洁中文注释
- **Infrastructure 层**：简洁中文注释 + 技术细节

---

## 分层注释规范

### 1. Handler 层（Transport Layer）

**位置**：`internal/transport/http/handlers/`

**规范**：
- ✅ 必须包含 Swagger 注解
- ✅ 完整的函数注释（参数、返回值、业务描述）- 使用中文
- ✅ 错误码说明
- ✅ Swagger 类型引用使用 `middleware.SuccessResponse` 和 `middleware.ErrorResponse`

**示例**：
```go
// Register 处理用户注册
//
// 创建新用户账户并返回认证令牌。
// 邮箱必须在系统中唯一。
//
// @Summary 用户注册并返回令牌
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "用户注册数据"
// @Success 201 {object} middleware.SuccessResponse{data=authentication.AuthResponse} "注册成功"
// @Failure 400 {object} middleware.ErrorResponse "请求参数错误"
// @Failure 409 {object} middleware.ErrorResponse "邮箱已存在"
// @Failure 500 {object} middleware.ErrorResponse "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // 1. 解析并验证请求
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, validationErr.FromGinError(err))
        return
    }

    // 2. 执行注册
    cmd := authentication.RegisterCommand{
        Email:    req.Email,
        Password: req.Password,
    }
    
    resp, err := h.service.Register(c.Request.Context(), cmd)
    if err != nil {
        response.Error(c, err)
        return
    }

    // 3. 返回响应
    response.Created(c, authentication.ToAuthResponse(resp))
}
```

---

### 2. Service 层（Application Layer）

**位置**：`internal/app/*/service.go`

**规范**：
- ✅ 完整函数注释（参数、返回值、业务流程、错误场景）- 使用中文
- ❌ 不使用 Swagger 注解
- ✅ 复杂逻辑使用步骤注释（中文）

**示例**：
```go
// Register 创建用户账户并返回认证令牌
//
// 注册流程：
// 1. 验证邮箱唯一性
// 2. 创建用户实体（密码已加密）
// 3. 生成访问令牌和刷新令牌
// 4. 发布 UserRegistered 领域事件
//
// 参数：
//   - ctx: 请求上下文
//   - cmd: 注册命令（包含邮箱和密码）
//
// 返回：
//   - *ServiceAuthResponse: 用户数据和认证令牌
//   - error: 注册失败时返回错误（邮箱已存在、验证错误等）
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*ServiceAuthResponse, error) {
    // 1. 检查邮箱是否已存在
    if s.userRepo.ExistsByEmail(ctx, cmd.Email) {
        return nil, userErr.ErrEmailAlreadyExists
    }

    // 2. 创建用户实体
    u, err := user.NewUser(cmd.Email, cmd.Password)
    if err != nil {
        return nil, err
    }

    // 3. 持久化用户到数据库
    if err := s.userRepo.Create(ctx, u); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 4. 生成认证令牌
    tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to generate tokens: %w", err)
    }

    // 5. 发布领域事件（异步，非阻塞）
    s.publisher.Publish(ctx, user.NewUserRegisteredEvent(u.ID, u.Email))

    return &ServiceAuthResponse{
        User:         u,
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
        ExpiresIn:    tokens.ExpiresIn,
    }, nil
}
```

---

### 3. Domain 层

**位置**：`internal/domain/*/`

**规范**：
- ✅ 简洁注释（一行说明）
- ✅ 聚合根、值对象、领域事件必须注释
- ❌ 不需要详细参数说明

**示例**：
```go
// User represents a registered user aggregate root.
// It encapsulates user identity, credentials, and lifecycle state.
type User struct {
    ID            string
    Email         string
    PasswordHash  string
    EmailVerified bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// NewUser creates a new user with validated email and password.
// Returns error if email format is invalid or password doesn't meet requirements.
func NewUser(email, password string) (*User, error) {
    // Implementation...
}

// IsLocked checks if the user account is locked due to too many failed attempts.
func (u *User) IsLocked() bool {
    return u.LockedUntil.After(time.Now())
}
```

**领域事件**：
```go
// UserRegistered is raised when a new user account is created.
type UserRegistered struct {
    UserID    string
    Email     string
    Timestamp time.Time
}

// AggregateID returns the user ID that triggered this event.
func (e *UserRegistered) AggregateID() string {
    return e.UserID
}
```

---

### 4. Infrastructure 层

**位置**：`internal/infra/*/`

**规范**：
- ✅ 简洁注释
- ✅ 技术实现细节说明
- ✅ 特殊配置说明

**示例**：
```go
// UserRepository implements user.Repository interface using GORM.
// It handles persistence operations for the User aggregate.
type UserRepository struct {
    db *gorm.DB
}

// NewUserRepository creates a new user repository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create persists a new user to the database.
// Converts domain entity to persistence object before saving.
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
    po := toUserPO(u)
    return r.db.WithContext(ctx).Create(&po).Error
}
```

---

## 通用注释规则

### 1. 结构体字段注释

```go
// Config contains application configuration settings.
type Config struct {
    Port     int    // HTTP server port
    LogLevel string // logging level (debug, info, warn, error)
    Timeout  time.Duration // request timeout duration
}
```

### 2. 接口注释

```go
// Repository defines the interface for user persistence operations.
// All methods accept context for request lifecycle management.
type Repository interface {
    // FindByID retrieves a user by their unique identifier.
    // Returns ErrUserNotFound if the user doesn't exist.
    FindByID(ctx context.Context, id string) (*User, error)

    // Create persists a new user to the database.
    // Returns error if email already exists.
    Create(ctx context.Context, user *User) error
}
```

### 3. 行内注释最佳实践

**✅ 好：解释 WHY**
```go
// Use exponential backoff to prevent thundering herd problem
// when multiple clients retry simultaneously.
time.Sleep(backoff)
```

**❌ 差：重复 WHAT**
```go
// Increment counter by 1
counter++
```

### 4. 步骤注释

```go
func (s *Service) ProcessOrder(ctx context.Context, orderID string) error {
    // 1. Load order aggregate
    order, err := s.orderRepo.FindByID(ctx, orderID)
    if err != nil {
        return err
    }

    // 2. Validate order state
    if !order.CanProcess() {
        return ErrOrderInvalidState
    }

    // 3. Execute business logic
    if err := order.Process(); err != nil {
        return err
    }

    // 4. Persist changes
    return s.orderRepo.Save(ctx, order)
}
```

### 5. TODO 和 FIXME 规范

```go
// TODO(shenfay): Add rate limiting for login endpoint
// TODO(shenfay): Implement email verification flow

// FIXME(shenfay): Memory leak when connection pool exhausted
// FIXME(shenfay): Race condition in token refresh logic
```

**格式**：`// TODO(username): Description`

---

## 注释检查清单

### Handler 层
- [ ] 包含 Swagger 注解（@Summary, @Param, @Success, @Failure）
- [ ] 函数注释说明业务功能
- [ ] 错误码说明
- [ ] 行内注释解释复杂逻辑

### Service 层
- [ ] 函数注释包含参数说明
- [ ] 函数注释包含返回值说明
- [ ] 函数注释包含错误场景
- [ ] 复杂业务逻辑使用步骤注释

### Domain 层
- [ ] 聚合根有注释
- [ ] 值对象有注释
- [ ] 领域事件有注释
- [ ] 公开方法有简洁注释

### Infrastructure 层
- [ ] 实现类有注释
- [ ] 技术选型有说明
- [ ] 特殊配置有注释

---

## 工具支持

### 生成文档
```bash
# 生成 Go doc
go doc ./internal/app/authentication

# 生成 Swagger 文档
swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal

# 查看函数文档
go doc authentication.Service.Register
```

### 检查注释覆盖率
```bash
# 检查未注释的公开标识符
golint ./...
```

---

## 示例文件

完整示例参考：
- Handler: `internal/transport/http/handlers/auth.go`
- Service: `internal/app/authentication/service.go`
- Domain: `internal/domain/user/entity.go`
- Infrastructure: `internal/infra/repository/user.go`
