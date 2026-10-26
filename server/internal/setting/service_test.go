package setting

import (
	"errors"
	"testing"
)

// mockRepo 内存实现，避免测试依赖 PostgreSQL
type mockRepo struct {
	data   map[string]string
	getErr error
	setErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{data: make(map[string]string)}
}

func (m *mockRepo) GetAll() (map[string]string, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	out := make(map[string]string, len(m.data))
	for k, v := range m.data {
		out[k] = v
	}
	return out, nil
}

func (m *mockRepo) Get(key string) (string, error) {
	if m.getErr != nil {
		return "", m.getErr
	}
	return m.data[key], nil
}

func (m *mockRepo) Set(key, value string) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.data[key] = value
	return nil
}

func (m *mockRepo) SetMany(kv map[string]string) error {
	if m.setErr != nil {
		return m.setErr
	}
	for k, v := range kv {
		m.data[k] = v
	}
	return nil
}

func TestGetSMTP_数据库值覆盖环境变量(t *testing.T) {
	repo := newMockRepo()
	repo.data[KeySMTPEnabled] = "true"
	repo.data[KeySMTPHost] = "smtp.qq.com"
	repo.data[KeySMTPPort] = "465"
	repo.data[KeySMTPFrom] = "hi@example.com"

	cfg, err := NewSettingServiceWithDeps(repo).GetSMTP()
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if !cfg.Enabled {
		t.Error("Enabled 应为 true")
	}
	if cfg.Host != "smtp.qq.com" {
		t.Errorf("Host 应为 smtp.qq.com，实际 %s", cfg.Host)
	}
	if cfg.Port != 465 {
		t.Errorf("Port 应为 465，实际 %d", cfg.Port)
	}
}

func TestGetSMTP_密码永不回显(t *testing.T) {
	repo := newMockRepo()
	repo.data[KeySMTPEnabled] = "true"
	repo.data[KeySMTPPassword] = "super-secret"

	cfg, err := NewSettingServiceWithDeps(repo).GetSMTP()
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if cfg.Password != "" {
		t.Errorf("Password 必须为空，实际泄露了: %s", cfg.Password)
	}
	if !cfg.PasswordSet {
		t.Error("PasswordSet 应为 true，用于告知前端已配置过密码")
	}
}

func TestGetSMTPWithPassword_返回明文供发信使用(t *testing.T) {
	repo := newMockRepo()
	repo.data[KeySMTPEnabled] = "true"
	repo.data[KeySMTPPassword] = "super-secret"

	cfg, err := NewSettingServiceWithDeps(repo).GetSMTPWithPassword()
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if cfg.Password != "super-secret" {
		t.Errorf("应返回明文密码，实际 %q", cfg.Password)
	}
}

// 死锁保护：SMTP 关闭时不得强制邮箱验证，
// 否则新用户收不到验证邮件也就永远无法发帖。
func TestGetSMTP_SMTP关闭时强制取消邮箱验证(t *testing.T) {
	repo := newMockRepo()
	repo.data[KeySMTPEnabled] = "false"
	repo.data[KeyEmailVerifyRequired] = "true" // 即使库里存了 true

	cfg, err := NewSettingServiceWithDeps(repo).GetSMTP()
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if cfg.EmailVerifyRequired {
		t.Error("SMTP 关闭时 EmailVerifyRequired 必须为 false，否则形成死锁")
	}
}

func TestUpdateSMTP_空密码保留原值(t *testing.T) {
	repo := newMockRepo()
	repo.data[KeySMTPPassword] = "old-password"

	err := NewSettingServiceWithDeps(repo).UpdateSMTP(&UpdateSMTPRequest{
		Enabled:  true,
		Host:     "smtp.163.com",
		Port:     994,
		Password: "", // 留空 = 不修改
	})
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if repo.data[KeySMTPPassword] != "old-password" {
		t.Errorf("空密码应保留原值，实际变成 %q", repo.data[KeySMTPPassword])
	}
	if repo.data[KeySMTPHost] != "smtp.163.com" {
		t.Errorf("Host 应更新，实际 %q", repo.data[KeySMTPHost])
	}
}

func TestUpdateSMTP_非空密码覆盖(t *testing.T) {
	repo := newMockRepo()
	repo.data[KeySMTPPassword] = "old-password"

	err := NewSettingServiceWithDeps(repo).UpdateSMTP(&UpdateSMTPRequest{
		Enabled:  true,
		Password: "new-password",
	})
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if repo.data[KeySMTPPassword] != "new-password" {
		t.Errorf("密码应被覆盖，实际 %q", repo.data[KeySMTPPassword])
	}
}

// SMTP 关闭时不允许把「强制邮箱验证」存成 true
func TestUpdateSMTP_关闭SMTP时不得写入强制验证(t *testing.T) {
	repo := newMockRepo()

	err := NewSettingServiceWithDeps(repo).UpdateSMTP(&UpdateSMTPRequest{
		Enabled:             false,
		EmailVerifyRequired: true,
	})
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if repo.data[KeyEmailVerifyRequired] != "false" {
		t.Errorf("SMTP 关闭时应写入 false，实际 %q", repo.data[KeyEmailVerifyRequired])
	}
}

func TestIsEmailVerifyRequired_读取失败时放行(t *testing.T) {
	repo := newMockRepo()
	repo.getErr = errors.New("数据库炸了")

	if NewSettingServiceWithDeps(repo).IsEmailVerifyRequired() {
		t.Error("设置表读取失败时应放行，避免全站无法发帖")
	}
}

func TestGetSMTP_库中无值时回落环境变量默认(t *testing.T) {
	// 空仓储 → 应回落到 config.Load() 的默认值，且不报错
	cfg, err := NewSettingServiceWithDeps(newMockRepo()).GetSMTP()
	if err != nil {
		t.Fatalf("预期无错误，实际: %v", err)
	}
	if cfg.Port <= 0 {
		t.Errorf("Port 应回落到环境变量默认值，实际 %d", cfg.Port)
	}
}
