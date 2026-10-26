package verify

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// VerifyRepoInterface 验证仓储接口
type VerifyRepoInterface interface {
	Create(verification *EmailVerification) error
	GetByToken(token string) (*EmailVerification, error)
	GetByUserID(userID int64) (*EmailVerification, error)
	Update(verification *EmailVerification) error
	Delete(id int64) error
}

// VerifyRepo 验证仓储实现
type VerifyRepo struct {
	db *gorm.DB
}

// NewVerifyRepo 创建验证仓储
func NewVerifyRepo() *VerifyRepo {
	return &VerifyRepo{
		db: database.GetDB(),
	}
}

// NewVerifyRepoWithDB 使用注入的 DB 创建验证仓储（测试用）
func NewVerifyRepoWithDB(db *gorm.DB) *VerifyRepo {
	return &VerifyRepo{
		db: db,
	}
}

// Create 创建验证记录
func (r *VerifyRepo) Create(verification *EmailVerification) error {
	return r.db.Create(verification).Error
}

// GetByToken 根据 token 查询验证记录
func (r *VerifyRepo) GetByToken(token string) (*EmailVerification, error) {
	var verification EmailVerification
	err := r.db.Where("token = ?", token).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByUserID 根据用户 ID 查询验证记录
func (r *VerifyRepo) GetByUserID(userID int64) (*EmailVerification, error) {
	var verification EmailVerification
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// Update 更新验证记录
func (r *VerifyRepo) Update(verification *EmailVerification) error {
	return r.db.Save(verification).Error
}

// Delete 删除验证记录
func (r *VerifyRepo) Delete(id int64) error {
	return r.db.Delete(&EmailVerification{}, id).Error
}
