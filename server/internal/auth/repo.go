package auth

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// AuthRepository 认证仓储接口
type AuthRepository interface {
	CreateUser(user *User) error
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByPhone(phone string) (*User, error)
	EmailExists(email string) bool
	PhoneExists(phone string) bool
	UpdateEmailVerified(userID int64, verified bool) error
	UpdatePassword(userID int64, hashedPassword string) error
}

// AuthRepo 认证仓储（实现 AuthRepository）
type AuthRepo struct {
	db *gorm.DB
}

// NewAuthRepo 创建认证仓储
func NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		db: database.GetDB(),
	}
}

// NewAuthRepoWithDB 使用注入的 DB 创建认证仓储（测试用）
func NewAuthRepoWithDB(db *gorm.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

// CreateUser 创建用户
func (r *AuthRepo) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// GetByID 根据 ID 查询用户
func (r *AuthRepo) GetByID(id int64) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱查询用户
func (r *AuthRepo) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号查询用户
func (r *AuthRepo) GetByPhone(phone string) (*User, error) {
	var user User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// EmailExists 邮箱是否已注册
func (r *AuthRepo) EmailExists(email string) bool {
	var count int64
	r.db.Model(&User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// PhoneExists 手机号是否已注册
func (r *AuthRepo) PhoneExists(phone string) bool {
	var count int64
	r.db.Model(&User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// NicknameExists 昵称是否已存在
func (r *AuthRepo) NicknameExists(nickname string) bool {
	var count int64
	r.db.Model(&User{}).Where("nickname = ?", nickname).Count(&count)
	return count > 0
}

// GetAll 获取所有用户
func (r *AuthRepo) GetAll() ([]User, error) {
	var users []User
	err := r.db.Order("id ASC").Find(&users).Error
	return users, err
}

// Count 统计用户数量
func (r *AuthRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&User{}).Count(&count).Error
	return count, err
}

// UpdateEmailVerified 更新用户邮箱验证状态
func (r *AuthRepo) UpdateEmailVerified(userID int64, verified bool) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("email_verified", verified).Error
}

// UpdatePassword 更新用户密码
func (r *AuthRepo) UpdatePassword(userID int64, hashedPassword string) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
