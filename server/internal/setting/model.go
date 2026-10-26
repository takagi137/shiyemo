package setting

import "time"

// SystemSetting 系统设置（键值对）
type SystemSetting struct {
	Key       string    `json:"key" gorm:"primaryKey;size:100"`
	Value     string    `json:"value" gorm:"type:text;not null;default:''"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (SystemSetting) TableName() string {
	return "system_settings"
}

// 设置键常量
const (
	KeySMTPEnabled  = "smtp.enabled"
	KeySMTPHost     = "smtp.host"
	KeySMTPPort     = "smtp.port"
	KeySMTPUsername = "smtp.username"
	KeySMTPPassword = "smtp.password"
	KeySMTPFrom     = "smtp.from"

	// KeyEmailVerifyRequired 是否强制要求邮箱验证后才能发帖/评论。
	// 与 smtp.enabled 独立：SMTP 关闭时若仍强制验证，
	// 新用户将永远收不到验证邮件，也就永远无法发帖（死锁）。
	// 因此 service 层会在 SMTP 关闭时强制返回 false。
	KeyEmailVerifyRequired = "email.verify_required"
)

// SMTPSettings SMTP 配置（对外结构）
type SMTPSettings struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	// Password 读取时始终为空字符串，仅写入时使用。
	// 避免管理后台把已保存的密码明文回显到浏览器。
	Password string `json:"password"`
	// PasswordSet 告诉前端「已存过密码」，用于区分「未配置」和「不回显」
	PasswordSet bool   `json:"password_set"`
	From        string `json:"from"`

	EmailVerifyRequired bool `json:"email_verify_required"`
}

// UpdateSMTPRequest 更新 SMTP 配置请求
type UpdateSMTPRequest struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host" binding:"omitempty,max=255"`
	Port     int    `json:"port" binding:"omitempty,min=1,max=65535"`
	Username string `json:"username" binding:"omitempty,max=255"`
	// Password 留空表示「不修改已保存的密码」
	Password string `json:"password" binding:"omitempty,max=255"`
	From     string `json:"from" binding:"omitempty,email"`

	EmailVerifyRequired bool `json:"email_verify_required"`
}

// TestSendRequest 测试发信请求
type TestSendRequest struct {
	To string `json:"to" binding:"required,email"`
}
