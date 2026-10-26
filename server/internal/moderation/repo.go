package moderation

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// ModerationRepo 审核仓库
type ModerationRepo struct {
	db *gorm.DB
}

// NewModerationRepo 创建审核仓库
func NewModerationRepo() *ModerationRepo {
	return &ModerationRepo{db: database.GetDB()}
}

// NewModerationRepoWithDB 使用注入的 DB 创建审核仓库（测试用）
func NewModerationRepoWithDB(db *gorm.DB) *ModerationRepo {
	return &ModerationRepo{db: db}
}

// CreateWord 创建敏感词
func (r *ModerationRepo) CreateWord(word *SensitiveWord) error {
	return r.db.Create(word).Error
}

// GetWordByID 根据ID获取敏感词
func (r *ModerationRepo) GetWordByID(id int64) (*SensitiveWord, error) {
	var word SensitiveWord
	err := r.db.First(&word, id).Error
	return &word, err
}

// GetWordByName 根据名称获取敏感词
func (r *ModerationRepo) GetWordByName(name string) (*SensitiveWord, error) {
	var word SensitiveWord
	err := r.db.Where("word = ?", name).First(&word).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &word, err
}

// ListWords 获取敏感词列表
func (r *ModerationRepo) ListWords(page, size int) ([]SensitiveWord, int64, error) {
	var words []SensitiveWord
	var total int64

	r.db.Model(&SensitiveWord{}).Count(&total)

	offset := (page - 1) * size
	err := r.db.Offset(offset).Limit(size).Order("id DESC").Find(&words).Error

	return words, total, err
}

// DeleteWord 删除敏感词
func (r *ModerationRepo) DeleteWord(id int64) error {
	return r.db.Delete(&SensitiveWord{}, id).Error
}

// GetAllWords 获取所有敏感词（用于加载到内存）
func (r *ModerationRepo) GetAllWords() ([]string, error) {
	var words []string
	err := r.db.Model(&SensitiveWord{}).Pluck("word", &words).Error
	return words, err
}
