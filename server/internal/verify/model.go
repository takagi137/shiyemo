package verify

import "time"

// EmailVerification 邮箱验证记录
type EmailVerification struct {
	ID         int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64      `json:"user_id" gorm:"not null;index"`
	Token      string     `json:"token" gorm:"not null;uniqueIndex;size:100"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"not null"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TableName 表名
func (EmailVerification) TableName() string {
	return "email_verifications"
}

// IsExpired 是否已过期
func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.ExpiresAt)
}

// IsVerified 是否已验证
func (ev *EmailVerification) IsVerified() bool {
	return ev.VerifiedAt != nil
}

// Verify 标记为已验证
func (ev *EmailVerification) Verify() {
	now := time.Now()
	ev.VerifiedAt = &now
}
