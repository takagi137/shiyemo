package email

import (
	"testing"

	"github.com/shiyemo/server/internal/common/config"
	"github.com/stretchr/testify/assert"
)

func TestMockEmailSender_Send(t *testing.T) {
	sender := NewMockEmailSender()

	err := sender.Send("test@example.com", "测试邮件", "测试内容")
	assert.NoError(t, err)
	assert.Len(t, sender.SentEmails, 1)
	assert.Equal(t, "test@example.com", sender.SentEmails[0].To)
	assert.Equal(t, "测试邮件", sender.SentEmails[0].Subject)
	assert.Equal(t, "测试内容", sender.SentEmails[0].Body)
}

func TestMockEmailSender_SendVerificationEmail(t *testing.T) {
	sender := NewMockEmailSender()

	err := sender.SendVerificationEmail("test@example.com", "test-token-123", "http://localhost:3000")
	assert.NoError(t, err)
	assert.Len(t, sender.SentEmails, 1)
	assert.Equal(t, "test@example.com", sender.SentEmails[0].To)
	assert.Contains(t, sender.SentEmails[0].Body, "test-token-123")
	assert.Contains(t, sender.SentEmails[0].Body, "http://localhost:3000")
}

func TestSMTPEmailSender_Send_Disabled(t *testing.T) {
	cfg := &config.SMTPConfig{
		Enabled: false,
	}
	sender := NewSMTPEmailSenderWithConfig(cfg)

	// SMTP 未启用时应该跳过发送，不返回错误
	err := sender.Send("test@example.com", "测试邮件", "测试内容")
	assert.NoError(t, err)
}

func TestSMTPEmailSender_SendVerificationEmail_Disabled(t *testing.T) {
	cfg := &config.SMTPConfig{
		Enabled: false,
	}
	sender := NewSMTPEmailSenderWithConfig(cfg)

	// SMTP 未启用时应该跳过发送，不返回错误
	err := sender.SendVerificationEmail("test@example.com", "test-token-123", "http://localhost:3000")
	assert.NoError(t, err)
}
