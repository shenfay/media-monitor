package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/user"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// ResetTokenPO 密码重置令牌持久化对象
type ResetTokenPO struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	UserID    string    `gorm:"type:varchar(50);not null;index"`
	Token     string    `gorm:"uniqueIndex;type:varchar(255);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 表名
func (ResetTokenPO) TableName() string {
	return "password_reset_tokens"
}

// ToDomain 转换为领域实体
func (po *ResetTokenPO) ToDomain() *user.ResetToken {
	return &user.ResetToken{
		ID:        po.ID,
		UserID:    po.UserID,
		Token:     po.Token,
		ExpiresAt: po.ExpiresAt,
		Used:      po.Used,
		CreatedAt: po.CreatedAt,
	}
}

// resetTokenRepository GORM 实现
type resetTokenRepository struct {
	db *gorm.DB
}

// NewResetTokenRepository 创建重置令牌仓储
func NewResetTokenRepository(db *gorm.DB) user.ResetTokenRepository {
	return &resetTokenRepository{db: db}
}

func (r *resetTokenRepository) Create(ctx context.Context, token *user.ResetToken) error {
	now := utils.Now()
	po := &ResetTokenPO{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Used:      token.Used,
		CreatedAt: now,
	}
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, token string) (*user.ResetToken, error) {
	var po ResetTokenPO
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&po).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *resetTokenRepository) MarkAsUsed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&ResetTokenPO{}).
		Where("id = ?", id).
		Update("used", true).Error
}
