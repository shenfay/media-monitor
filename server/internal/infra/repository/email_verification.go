package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/shenfay/go-react-admin/internal/domain/user"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// VerificationTokenPO 邮箱验证令牌持久化对象
type VerificationTokenPO struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	UserID    string    `gorm:"type:varchar(50);not null;index"`
	Token     string    `gorm:"uniqueIndex;type:varchar(255);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 表名
func (VerificationTokenPO) TableName() string {
	return "email_verification_tokens"
}

// ToDomain 转换为领域实体
func (po *VerificationTokenPO) ToDomain() *user.VerificationToken {
	return &user.VerificationToken{
		ID:        po.ID,
		UserID:    po.UserID,
		Token:     po.Token,
		ExpiresAt: po.ExpiresAt,
		Used:      po.Used,
		CreatedAt: po.CreatedAt,
	}
}

// verificationTokenRepository GORM 实现
type verificationTokenRepository struct {
	db *gorm.DB
}

// NewVerificationTokenRepository 创建验证令牌仓储
func NewVerificationTokenRepository(db *gorm.DB) user.VerificationTokenRepository {
	return &verificationTokenRepository{db: db}
}

func (r *verificationTokenRepository) Create(ctx context.Context, token *user.VerificationToken) error {
	now := utils.Now()
	po := &VerificationTokenPO{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Used:      token.Used,
		CreatedAt: now,
	}
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *verificationTokenRepository) FindByToken(ctx context.Context, token string) (*user.VerificationToken, error) {
	var po VerificationTokenPO
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&po).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *verificationTokenRepository) MarkAsUsed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&VerificationTokenPO{}).
		Where("id = ?", id).
		Update("used", true).Error
}
