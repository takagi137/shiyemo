package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/shiyemo/server/internal/common/config"
)

// EmailSender 邮件发送器接口
type EmailSender interface {
	Send(to, subject, body string) error
	SendVerificationEmail(to, token, baseURL string) error
	SendResetPasswordEmail(to, token, baseURL string) error
}

// ConfigProvider 在每次发信时解析当前 SMTP 配置。
// 用提供者而非固定结构体，是为了让管理后台修改配置后立即生效，无需重启容器。
type ConfigProvider func() *config.SMTPConfig

// SMTPEmailSender SMTP 邮件发送器
type SMTPEmailSender struct {
	provider ConfigProvider
}

// NewSMTPEmailSender 创建 SMTP 邮件发送器（读环境变量，配置在构造时固化）
func NewSMTPEmailSender() *SMTPEmailSender {
	cfg := config.Load()
	return NewSMTPEmailSenderWithConfig(&cfg.SMTP)
}

// NewSMTPEmailSenderWithConfig 使用指定配置创建 SMTP 邮件发送器（配置固定不变）
func NewSMTPEmailSenderWithConfig(cfg *config.SMTPConfig) *SMTPEmailSender {
	return &SMTPEmailSender{
		provider: func() *config.SMTPConfig { return cfg },
	}
}

// NewSMTPEmailSenderWithProvider 使用配置提供者创建发送器。
// 每次发信都会调用 provider 重新取配置，因此管理后台的改动即时生效。
func NewSMTPEmailSenderWithProvider(provider ConfigProvider) *SMTPEmailSender {
	return &SMTPEmailSender{provider: provider}
}

// Send 发送邮件
func (s *SMTPEmailSender) Send(to, subject, body string) error {
	cfg := s.provider()
	if cfg == nil {
		log.Printf("[邮件] SMTP 配置不可用，跳过发送: to=%s", to)
		return nil
	}
	if !cfg.Enabled {
		log.Printf("[邮件] SMTP 未启用，跳过发送: to=%s, subject=%s", to, subject)
		return nil
	}

	// 构建邮件消息
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", cfg.From, to, subject, body)

	// SMTP 认证
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	err := smtp.SendMail(addr, auth, cfg.From, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("[邮件] 发送失败: to=%s, error=%v", to, err)
		return err
	}

	log.Printf("[邮件] 发送成功: to=%s, subject=%s", to, subject)
	return nil
}

// SendVerificationEmail 发送验证邮件
func (s *SMTPEmailSender) SendVerificationEmail(to, token, baseURL string) error {
	subject := "【失业么】邮箱验证"
	verifyURL := fmt.Sprintf("%s/verify-email/confirm?token=%s", baseURL, token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 1px solid #eee; }
        .content { padding: 30px 0; }
        .button { display: inline-block; padding: 12px 30px; background-color: #4F46E5; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px 0; border-top: 1px solid #eee; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>失业么</h1>
            <p>AI 时代个体的"能力重建与收入重建"社区</p>
        </div>
        <div class="content">
            <h2>验证你的邮箱</h2>
            <p>你好！感谢注册失业么社区。</p>
            <p>请点击下方按钮验证你的邮箱地址：</p>
            <p style="text-align: center;">
                <a href="%s" class="button">验证邮箱</a>
            </p>
            <p>如果按钮无法点击，请复制以下链接到浏览器打开：</p>
            <p style="word-break: break-all; background-color: #f5f5f5; padding: 10px; border-radius: 4px;">
                %s
            </p>
            <p>此链接将在 24 小时后失效。</p>
            <p>如果这不是你的操作，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>© 2024 失业么社区. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, verifyURL, verifyURL)

	return s.Send(to, subject, body)
}

// SendResetPasswordEmail 发送密码重置邮件
func (s *SMTPEmailSender) SendResetPasswordEmail(to, token, baseURL string) error {
	subject := "【失业么】密码重置"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; padding: 20px 0; border-bottom: 1px solid #eee; }
        .content { padding: 30px 0; }
        .button { display: inline-block; padding: 12px 30px; background-color: #4F46E5; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px 0; border-top: 1px solid #eee; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>失业么</h1>
        </div>
        <div class="content">
            <h2>重置你的密码</h2>
            <p>你好！我们收到了你的密码重置请求。</p>
            <p>请点击下方按钮重置密码：</p>
            <p style="text-align: center;">
                <a href="%s" class="button">重置密码</a>
            </p>
            <p>如果按钮无法点击，请复制以下链接到浏览器打开：</p>
            <p style="word-break: break-all; background-color: #f5f5f5; padding: 10px; border-radius: 4px;">
                %s
            </p>
            <p>此链接将在 15 分钟后失效。</p>
            <p>如果这不是你的操作，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>© 2024 失业么社区. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, resetURL, resetURL)

	return s.Send(to, subject, body)
}

// MockEmailSender Mock 邮件发送器（用于测试）
type MockEmailSender struct {
	SentEmails []MockEmail
}

// MockEmail Mock 邮件
type MockEmail struct {
	To      string
	Subject string
	Body    string
}

// NewMockEmailSender 创建 Mock 邮件发送器
func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{
		SentEmails: make([]MockEmail, 0),
	}
}

// Send 发送邮件（Mock）
func (m *MockEmailSender) Send(to, subject, body string) error {
	m.SentEmails = append(m.SentEmails, MockEmail{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}

// SendVerificationEmail 发送验证邮件（Mock）
func (m *MockEmailSender) SendVerificationEmail(to, token, baseURL string) error {
	return m.Send(to, "验证邮箱", "验证链接: "+baseURL+"/verify-email/confirm?token="+token)
}

// SendResetPasswordEmail 发送密码重置邮件（Mock）
func (m *MockEmailSender) SendResetPasswordEmail(to, token, baseURL string) error {
	return m.Send(to, "密码重置", "重置链接: "+baseURL+"/reset-password?token="+token)
}
