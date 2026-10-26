package backup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CommandRunner 命令执行器抽象，方便测试 Mock
//go:generate mockery --name CommandRunner --output . --outpkg backup --filename mock_runner_test.go

type CommandRunner interface {
	CombinedOutput(name string, args []string, env []string) ([]byte, error)
}

// OsCommandRunner 真实命令执行器，调用系统二进制

type OsCommandRunner struct{}

func (r *OsCommandRunner) CombinedOutput(name string, args []string, env []string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	return cmd.CombinedOutput()
}

// BackupService 备份服务
type BackupService struct {
	backupDir  string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	runner     CommandRunner
}

// NewBackupService 创建备份服务（使用真实命令执行器）
func NewBackupService(backupDir, dbHost, dbPort, dbUser, dbPassword, dbName string) (*BackupService, error) {
	return NewBackupServiceWithRunner(backupDir, dbHost, dbPort, dbUser, dbPassword, dbName, &OsCommandRunner{})
}

// NewBackupServiceWithRunner 创建备份服务（注入指定命令执行器，供测试使用）
func NewBackupServiceWithRunner(backupDir, dbHost, dbPort, dbUser, dbPassword, dbName string, runner CommandRunner) (*BackupService, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("create backup directory: %w", err)
	}

	return &BackupService{
		backupDir:  backupDir,
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		dbName:     dbName,
		runner:     runner,
	}, nil
}

// Backup 执行备份
func (s *BackupService) Backup() (string, error) {
	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(s.backupDir, fmt.Sprintf("backup_%s.sql", timestamp))

	// 构建 pg_dump 参数
	args := []string{
		"-h", s.dbHost,
		"-p", s.dbPort,
		"-U", s.dbUser,
		"-d", s.dbName,
		"-f", backupFile,
		"--no-owner",
		"--no-privileges",
	}

	// 构建环境变量（密码）
	env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.dbPassword))

	// 执行备份
	output, err := s.runner.CombinedOutput("pg_dump", args, env)
	if err != nil {
		return "", fmt.Errorf("pg_dump failed: %w, output: %s", err, string(output))
	}

	// 验证备份文件
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return "", fmt.Errorf("backup file not created")
	}

	// 获取文件大小
	stat, err := os.Stat(backupFile)
	if err != nil {
		return "", fmt.Errorf("stat backup file: %w", err)
	}

	log.Printf("备份成功: %s (大小: %d bytes)", backupFile, stat.Size())
	return backupFile, nil
}

// Restore 恢复备份
func (s *BackupService) Restore(backupFile string) error {
	// 检查备份文件是否存在
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFile)
	}

	// 构建 psql 参数
	args := []string{
		"-h", s.dbHost,
		"-p", s.dbPort,
		"-U", s.dbUser,
		"-d", s.dbName,
		"-f", backupFile,
		"--no-psqlrc",
	}

	// 构建环境变量（密码）
	env := append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.dbPassword))

	// 执行恢复
	output, err := s.runner.CombinedOutput("psql", args, env)
	if err != nil {
		return fmt.Errorf("psql restore failed: %w, output: %s", err, string(output))
	}

	log.Printf("恢复成功: %s", backupFile)
	return nil
}

// ListBackups 列出所有备份文件
func (s *BackupService) ListBackups() ([]BackupInfo, error) {
	var backups []BackupInfo

	// 读取备份目录
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return nil, fmt.Errorf("read backup directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 检查是否是备份文件
		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		// 获取文件信息
		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			Name:      entry.Name(),
			Path:      filepath.Join(s.backupDir, entry.Name()),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return backups, nil
}

// DeleteBackup 删除备份文件
func (s *BackupService) DeleteBackup(backupFile string) error {
	// 检查文件是否存在
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFile)
	}

	// 删除文件
	if err := os.Remove(backupFile); err != nil {
		return fmt.Errorf("delete backup file: %w", err)
	}

	log.Printf("删除备份: %s", backupFile)
	return nil
}

// CleanupOldBackups 清理旧备份
func (s *BackupService) CleanupOldBackups(maxAge time.Duration) (int, error) {
	backups, err := s.ListBackups()
	if err != nil {
		return 0, err
	}

	deleted := 0
	for _, backup := range backups {
		if time.Since(backup.CreatedAt) > maxAge {
			if err := s.DeleteBackup(backup.Path); err != nil {
				log.Printf("删除旧备份失败: %s, 错误: %v", backup.Name, err)
				continue
			}
			deleted++
		}
	}

	return deleted, nil
}

// BackupInfo 备份信息
type BackupInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}