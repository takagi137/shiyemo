package setting

import (
	"strconv"

	"github.com/shiyemo/server/internal/common/config"
	"github.com/shiyemo/server/internal/common/email"
)

// SettingService 系统设置服务
type SettingService struct {
	repo SettingRepoInterface
}

// SettingServiceInterface 系统设置服务接口（供 mock 测试）
type SettingServiceInterface interface {
	GetSMTP() (*SMTPSettings, error)
	GetSMTPWithPassword() (*SMTPSettings, error)
	UpdateSMTP(req *UpdateSMTPRequest) error
	SMTPConfigProvider() email.ConfigProvider
	TestSend(to string) error
	IsEmailVerifyRequired() bool
}

// 编译期断言：*SettingService 必须实现 SettingServiceInterface
var _ SettingServiceInterface = (*SettingService)(nil)

// NewSettingService 创建系统设置服务
func NewSettingService() *SettingService {
	return &SettingService{repo: NewSettingRepo()}
}

// NewSettingServiceWithDeps 创建系统设置服务（注入依赖，用于测试）
func NewSettingServiceWithDeps(repo SettingRepoInterface) *SettingService {
	return &SettingService{repo: repo}
}

// GetSMTP 读取 SMTP 配置。
// 优先级：数据库 > 环境变量 > 默认值。
// 数据库中没有对应键时回落到环境变量，保证首次部署（表为空）仍可用。
func (s *SettingService) GetSMTP() (*SMTPSettings, error) {
	kv, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	envCfg := config.Load().SMTP

	out := &SMTPSettings{
		Enabled:  boolOr(kv, KeySMTPEnabled, envCfg.Enabled),
		Host:     strOr(kv, KeySMTPHost, envCfg.Host),
		Port:     intOr(kv, KeySMTPPort, envCfg.Port),
		Username: strOr(kv, KeySMTPUsername, envCfg.Username),
		From:     strOr(kv, KeySMTPFrom, envCfg.From),

		EmailVerifyRequired: boolOr(kv, KeyEmailVerifyRequired, false),
	}

	// 密码永不回显，只告知是否已设置
	out.Password = ""
	out.PasswordSet = strOr(kv, KeySMTPPassword, envCfg.Password) != ""

	// 死锁保护：SMTP 关闭时验证邮件根本发不出去，
	// 若仍强制邮箱验证，新用户将永远无法发帖。
	if !out.Enabled {
		out.EmailVerifyRequired = false
	}

	return out, nil
}

// GetSMTPWithPassword 供邮件发送器使用，返回含明文密码的配置。
// 与 GetSMTP 的区别仅在于不抹除 Password，切勿用于 HTTP 响应。
func (s *SettingService) GetSMTPWithPassword() (*SMTPSettings, error) {
	out, err := s.GetSMTP()
	if err != nil {
		return nil, err
	}
	kv, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	out.Password = strOr(kv, KeySMTPPassword, config.Load().SMTP.Password)
	return out, nil
}

// UpdateSMTP 更新 SMTP 配置。
// req.Password 为空表示保留原密码，避免管理员每次保存都要重输。
func (s *SettingService) UpdateSMTP(req *UpdateSMTPRequest) error {
	kv := map[string]string{
		KeySMTPEnabled:         strconv.FormatBool(req.Enabled),
		KeySMTPHost:            req.Host,
		KeySMTPUsername:        req.Username,
		KeySMTPFrom:            req.From,
		KeyEmailVerifyRequired: strconv.FormatBool(req.EmailVerifyRequired && req.Enabled),
	}
	if req.Port > 0 {
		kv[KeySMTPPort] = strconv.Itoa(req.Port)
	}
	if req.Password != "" {
		kv[KeySMTPPassword] = req.Password
	}
	return s.repo.SetMany(kv)
}

// SMTPConfigProvider 返回一个配置提供者，供 email 发送器在每次发信时解析当前配置。
// 这样管理后台保存设置后立即生效，无需重启服务。
func (s *SettingService) SMTPConfigProvider() email.ConfigProvider {
	return func() *config.SMTPConfig {
		cfg, err := s.GetSMTPWithPassword()
		if err != nil {
			return nil
		}
		return &config.SMTPConfig{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			From:     cfg.From,
			Enabled:  cfg.Enabled,
		}
	}
}

// TestSend 用当前配置发一封测试邮件，供管理员验证 SMTP 是否配置正确。
// 与普通发信不同：即使 Enabled 为 false 也会真实投递，
// 否则管理员无法在正式开启前先验证配置。
func (s *SettingService) TestSend(to string) error {
	cfg, err := s.GetSMTPWithPassword()
	if err != nil {
		return err
	}
	sender := email.NewSMTPEmailSenderWithConfig(&config.SMTPConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
		Enabled:  true, // 强制启用，仅用于本次测试投递
	})
	return sender.Send(to,
		"【失业么】SMTP 配置测试",
		"<p>这是一封测试邮件。</p><p>如果你收到了它，说明管理后台的 SMTP 配置正确可用。</p>")
}

// IsEmailVerifyRequired 是否强制邮箱验证（供中间件调用）
func (s *SettingService) IsEmailVerifyRequired() bool {
	cfg, err := s.GetSMTP()
	if err != nil {
		// 读取失败时放行，避免设置表异常导致全站无法发帖
		return false
	}
	return cfg.EmailVerifyRequired
}

func strOr(kv map[string]string, key, fallback string) string {
	if v, ok := kv[key]; ok && v != "" {
		return v
	}
	return fallback
}

func boolOr(kv map[string]string, key string, fallback bool) bool {
	if v, ok := kv[key]; ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func intOr(kv map[string]string, key string, fallback int) int {
	if v, ok := kv[key]; ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
