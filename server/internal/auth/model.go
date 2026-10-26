package auth

import "time"

// User 用户模型
type User struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Phone          *string    `json:"phone,omitempty" gorm:"uniqueIndex;size:20"`
	Email          *string    `json:"email,omitempty" gorm:"uniqueIndex;size:255"`
	EmailVerified  bool       `json:"email_verified" gorm:"default:false"`
	Password       string     `json:"-" gorm:"not null;size:255"`
	Nickname       string     `json:"nickname" gorm:"not null;size:50;index"`
	Bio            string     `json:"bio,omitempty" gorm:"size:500;default:''"`
	Company        string     `json:"company,omitempty" gorm:"size:100;default:''"`
	Occupation     string     `json:"occupation,omitempty" gorm:"size:100;default:''"`
	Avatar         string     `json:"avatar,omitempty" gorm:"size:500;default:''"`
	City           string     `json:"city,omitempty" gorm:"size:50;default:'';index"`
	Industry       string     `json:"industry,omitempty" gorm:"size:100;default:'';index"`
	LeftJobAt      *time.Time `json:"left_job_at,omitempty"`
	LastWorkDate   *time.Time `json:"last_work_date,omitempty"`
	FollowingCount int        `json:"following_count" gorm:"default:0"`
	FollowerCount  int        `json:"follower_count" gorm:"default:0"`
	IsAdmin        bool       `json:"is_admin" gorm:"default:false"`
	Status         int8       `json:"status" gorm:"default:1;index"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// IsActive 是否正常状态
func (u *User) IsActive() bool {
	return u.Status == 1
}

// HasPhone 是否绑定了手机号
func (u *User) HasPhone() bool {
	return u.Phone != nil && *u.Phone != ""
}

// HasEmail 是否绑定了邮箱
func (u *User) HasEmail() bool {
	return u.Email != nil && *u.Email != ""
}

// IsEmailVerified 邮箱是否已验证
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}

// FreedomDays 计算自由天数（从最后工作日到今天的天数）
func (u *User) FreedomDays() *int {
	if u.LastWorkDate == nil {
		return nil
	}
	days := int(time.Since(*u.LastWorkDate).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return &days
}
