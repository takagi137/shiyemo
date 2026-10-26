package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockCommandRunner 测试用命令执行器
// Backup 成功时需创建备份文件，Restore 成功时直接返回 nil

type mockCommandRunner struct {
	combinedOutputFn func(name string, args []string, env []string) ([]byte, error)
}

func (m *mockCommandRunner) CombinedOutput(name string, args []string, env []string) ([]byte, error) {
	if m.combinedOutputFn != nil {
		return m.combinedOutputFn(name, args, env)
	}
	return []byte(""), nil
}

func TestNewBackupService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tempDir := filepath.Join(os.TempDir(), "backup-test-new")
		defer os.RemoveAll(tempDir)

		service, err := NewBackupService(tempDir, "localhost", "5432", "user", "pass", "db")
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("invalid_dir", func(t *testing.T) {
		_, err := NewBackupService("Z:\\impossible\\path\\that\\does\\not\\exist\\deeply", "localhost", "5432", "user", "pass", "db")
		assert.Error(t, err)
	})
}

func TestBackupService_DeleteBackup_NotFound(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-del-nf")
	defer os.RemoveAll(tempDir)

	service, _ := NewBackupService(tempDir, "localhost", "5432", "test", "test", "test")
	err := service.DeleteBackup("/nonexistent/file.sql")
	assert.Error(t, err)
}

func TestBackupService_ListBackups_Empty(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-empty")
	defer os.RemoveAll(tempDir)

	service, _ := NewBackupService(tempDir, "localhost", "5432", "test", "test", "test")
	backups, err := service.ListBackups()
	assert.NoError(t, err)
	assert.Len(t, backups, 0)
}

func TestBackupService_Restore_NotFound(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-restore-nf")
	defer os.RemoveAll(tempDir)

	service, _ := NewBackupService(tempDir, "localhost", "5432", "test", "test", "test")
	err := service.Restore("/nonexistent/backup.sql")
	assert.Error(t, err)
}

func TestBackupService_ListBackups(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "backup-test-list")
	defer os.RemoveAll(tempDir)

	// 创建备份服务
	service, err := NewBackupService(tempDir, "localhost", "5432", "test", "test", "test")
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// 创建一些测试文件
	testFiles := []string{
		"backup_20260101_120000.sql",
		"backup_20260102_120000.sql",
		"backup_20260103_120000.sql",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		assert.NoError(t, err)
	}

	// 创建一些非备份文件
	nonBackupFiles := []string{
		"readme.txt",
		"backup_20260104_120000.txt",
	}

	for _, file := range nonBackupFiles {
		filePath := filepath.Join(tempDir, file)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		assert.NoError(t, err)
	}

	// 测试列出备份
	backups, err := service.ListBackups()
	assert.NoError(t, err)
	assert.Len(t, backups, 3)

	// 验证备份信息
	for _, backup := range backups {
		assert.NotEmpty(t, backup.Name)
		assert.NotEmpty(t, backup.Path)
		assert.True(t, backup.Size > 0)
		assert.False(t, backup.CreatedAt.IsZero())
	}
}

func TestBackupService_DeleteBackup(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "backup-test-delete")
	defer os.RemoveAll(tempDir)

	// 创建备份服务
	service, err := NewBackupService(tempDir, "localhost", "5432", "test", "test", "test")
	assert.NoError(t, err)

	// 创建测试文件
	testFile := filepath.Join(tempDir, "backup_20260101_120000.sql")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	assert.NoError(t, err)

	// 验证文件存在
	_, err = os.Stat(testFile)
	assert.NoError(t, err)

	// 删除备份
	err = service.DeleteBackup(testFile)
	assert.NoError(t, err)

	// 验证文件已删除
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestBackupService_CleanupOldBackups(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "backup-test-cleanup")
	defer os.RemoveAll(tempDir)

	// 创建备份服务
	service, err := NewBackupService(tempDir, "localhost", "5432", "test", "test", "test")
	assert.NoError(t, err)

	// 创建一些测试文件，设置不同的修改时间
	testFiles := []struct {
		name    string
		age     time.Duration
		content string
	}{
		{"backup_20260101_120000.sql", 10 * 24 * time.Hour, "old backup 1"},
		{"backup_20260102_120000.sql", 8 * 24 * time.Hour, "old backup 2"},
		{"backup_20260103_120000.sql", 2 * 24 * time.Hour, "recent backup"},
		{"backup_20260104_120000.sql", 1 * 24 * time.Hour, "newest backup"},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		err := os.WriteFile(filePath, []byte(tf.content), 0644)
		assert.NoError(t, err)

		// 设置修改时间
		modTime := time.Now().Add(-tf.age)
		err = os.Chtimes(filePath, modTime, modTime)
		assert.NoError(t, err)
	}

	// 清理 7 天前的备份
	deleted, err := service.CleanupOldBackups(7 * 24 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, 2, deleted)

	// 验证剩余文件
	backups, err := service.ListBackups()
	assert.NoError(t, err)
	assert.Len(t, backups, 2)

	// 验证剩余的是最近的文件
	remainingNames := make(map[string]bool)
	for _, b := range backups {
		remainingNames[b.Name] = true
	}
	assert.True(t, remainingNames["backup_20260103_120000.sql"])
	assert.True(t, remainingNames["backup_20260104_120000.sql"])
}

// === 新增测试：CommandRunner 抽象后的 Backup/Restore 覆盖 ===

func TestNewBackupServiceWithRunner(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-runner")
	defer os.RemoveAll(tempDir)

	mock := &mockCommandRunner{}
	service, err := NewBackupServiceWithRunner(tempDir, "localhost", "5432", "user", "pass", "db", mock)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestBackup_Success(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-success")
	defer os.RemoveAll(tempDir)

	var capturedArgs []string
	mock := &mockCommandRunner{
		combinedOutputFn: func(name string, args []string, env []string) ([]byte, error) {
			assert.Equal(t, "pg_dump", name)
			capturedArgs = args
			// 模拟 pg_dump 创建备份文件：从参数中提取 -f 路径
			for i, arg := range args {
				if arg == "-f" && i+1 < len(args) {
					err := os.WriteFile(args[i+1], []byte("-- SQL dump content"), 0644)
					assert.NoError(t, err)
				}
			}
			return []byte(""), nil
		},
	}

	service, err := NewBackupServiceWithRunner(tempDir, "127.0.0.1", "5432", "admin", "secret", "testdb", mock)
	assert.NoError(t, err)

	backupFile, err := service.Backup()
	assert.NoError(t, err)
	assert.NotEmpty(t, backupFile)

	// 验证传递给 pg_dump 的参数
	assert.Contains(t, capturedArgs, "-h")
	assert.Contains(t, capturedArgs, "127.0.0.1")
	assert.Contains(t, capturedArgs, "-d")
	assert.Contains(t, capturedArgs, "testdb")
	assert.Contains(t, capturedArgs, "--no-owner")

	// 验证备份文件确实存在
	_, err = os.Stat(backupFile)
	assert.NoError(t, err)
}

func TestBackup_RunError(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-runerr")
	defer os.RemoveAll(tempDir)

	mock := &mockCommandRunner{
		combinedOutputFn: func(name string, args []string, env []string) ([]byte, error) {
			return []byte("pg_dump: connection refused"), fmt.Errorf("exit status 1")
		},
	}

	service, err := NewBackupServiceWithRunner(tempDir, "localhost", "5432", "user", "pass", "db", mock)
	assert.NoError(t, err)

	backupFile, err := service.Backup()
	assert.Error(t, err)
	assert.Empty(t, backupFile)
	assert.Contains(t, err.Error(), "pg_dump failed")
}

func TestRestore_Success(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-restore-ok")
	defer os.RemoveAll(tempDir)

	var capturedCmd string
	mock := &mockCommandRunner{
		combinedOutputFn: func(name string, args []string, env []string) ([]byte, error) {
			capturedCmd = name
			return []byte(""), nil
		},
	}

	service, err := NewBackupServiceWithRunner(tempDir, "127.0.0.1", "5432", "admin", "secret", "testdb", mock)
	assert.NoError(t, err)

	// 先创建 service（会 MkdirAll），再写备份文件
	backupFile := filepath.Join(tempDir, "backup_restore_test.sql")
	err = os.WriteFile(backupFile, []byte("CREATE TABLE test(id int);"), 0644)
	assert.NoError(t, err)

	err = service.Restore(backupFile)
	assert.NoError(t, err)
	assert.Equal(t, "psql", capturedCmd)
}

func TestRestore_RunError(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup-test-restore-err")
	defer os.RemoveAll(tempDir)

	mock := &mockCommandRunner{
		combinedOutputFn: func(name string, args []string, env []string) ([]byte, error) {
			return []byte("syntax error near line 1"), fmt.Errorf("exit status 1")
		},
	}

	service, err := NewBackupServiceWithRunner(tempDir, "localhost", "5432", "user", "pass", "db", mock)
	assert.NoError(t, err)

	backupFile := filepath.Join(tempDir, "backup_restore_err.sql")
	err = os.WriteFile(backupFile, []byte("bad sql"), 0644)
	assert.NoError(t, err)

	err = service.Restore(backupFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "psql restore failed")
}

// === 边界路径覆盖 ===

func TestBackup_ListBackups_InfoError(t *testing.T) {
	// 测试非 .sql 文件和目录的过滤逻辑
	tempDir := filepath.Join(os.TempDir(), "backup-test-list-filter")
	defer os.RemoveAll(tempDir)

	service, err := NewBackupServiceWithRunner(tempDir, "localhost", "5432", "test", "test", "test", &mockCommandRunner{})
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "valid.sql"), []byte("data"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("data"), 0644)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(tempDir, "subdir"), 0755)
	assert.NoError(t, err)

	backups, err := service.ListBackups()
	assert.NoError(t, err)
	assert.Len(t, backups, 1)
	assert.Equal(t, "valid.sql", backups[0].Name)
}

func TestCleanupOldBackups_DeletePartialFailure(t *testing.T) {
	// 测试 ListBackups 失败时 CleanupOldBackups 的错误路径
	tempDir := filepath.Join(os.TempDir(), "backup-test-cleanup-delete-dir")
	defer os.RemoveAll(filepath.Dir(tempDir))

	service, err := NewBackupServiceWithRunner(tempDir, "localhost", "5432", "test", "test", "test", &mockCommandRunner{})
	assert.NoError(t, err)

	// 删除目录，使后续 ListBackups 读取失败
	err = os.Remove(tempDir)
	assert.NoError(t, err)

	deleted, err := service.CleanupOldBackups(time.Hour)
	assert.Error(t, err)
	assert.Equal(t, 0, deleted)
}