package setting

import (
	"errors"

	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingRepoInterface 系统设置仓储接口
type SettingRepoInterface interface {
	GetAll() (map[string]string, error)
	Get(key string) (string, error)
	Set(key, value string) error
	SetMany(kv map[string]string) error
}

// SettingRepo 系统设置仓储实现
type SettingRepo struct {
	db *gorm.DB
}

// NewSettingRepo 创建系统设置仓储
func NewSettingRepo() *SettingRepo {
	return &SettingRepo{db: database.GetDB()}
}

// NewSettingRepoWithDB 使用注入的 DB 创建仓储（测试用）
func NewSettingRepoWithDB(db *gorm.DB) *SettingRepo {
	return &SettingRepo{db: db}
}

// GetAll 取回全部设置，返回 key -> value
func (r *SettingRepo) GetAll() (map[string]string, error) {
	var rows []SystemSetting
	if err := r.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[row.Key] = row.Value
	}
	return out, nil
}

// Get 取单个设置，key 不存在时返回空字符串且不报错
func (r *SettingRepo) Get(key string) (string, error) {
	var row SystemSetting
	err := r.db.Where("key = ?", key).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

// Set 写入单个设置（存在则更新）
func (r *SettingRepo) Set(key, value string) error {
	return r.upsert(r.db, key, value)
}

// SetMany 在单个事务内批量写入，避免部分成功导致配置处于中间状态
func (r *SettingRepo) SetMany(kv map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for k, v := range kv {
			if err := r.upsert(tx, k, v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *SettingRepo) upsert(tx *gorm.DB, key, value string) error {
	row := SystemSetting{Key: key, Value: value}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"value": value, "updated_at": gorm.Expr("now()")}),
	}).Create(&row).Error
}
