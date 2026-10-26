package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	App     AppConfig
	DB      DBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Storage StorageConfig
	Log     LogConfig
	Backup  BackupConfig
	SMTP    SMTPConfig
}

// AppConfig 应用配置
type AppConfig struct {
	Port      string
	Mode      string // debug/release
	PublicURL string // browser-facing origin used in emails
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string
	ExpireHour int
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type      string // local / minio
	UploadDir string
	BaseURL   string
	MinIO     MinIOConfig
}

// LogConfig 日志配置
type LogConfig struct {
	Dir   string
	Level string
}

// BackupConfig 备份配置
type BackupConfig struct {
	Dir string
}

// MinIOConfig MinIO 配置
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	BaseURL   string
}

// SMTPConfig SMTP 配置
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Enabled  bool
}

// Load 从环境变量加载配置
func Load() *Config {
	return &Config{
		App: AppConfig{
			Port:      getEnv("APP_PORT", "8080"),
			Mode:      getEnv("APP_MODE", "debug"),
			PublicURL: getEnv("APP_PUBLIC_URL", "http://localhost:3000"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "shiyemo"),
			Password: getEnv("DB_PASSWORD", "shiyemo123"),
			DBName:   getEnv("DB_NAME", "shiyemo"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "shiyemo-jwt-secret-change-in-production"),
			ExpireHour: 168, // 7 天
		},
		Storage: StorageConfig{
			Type:      getEnv("STORAGE_TYPE", "local"),
			UploadDir: getEnv("STORAGE_UPLOAD_DIR", "uploads"),
			BaseURL:   getEnv("STORAGE_BASE_URL", "/uploads"),
			MinIO: MinIOConfig{
				Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
				AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
				SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
				Bucket:    getEnv("MINIO_BUCKET", "shiyemo"),
				UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
				BaseURL:   getEnv("MINIO_BASE_URL", "http://localhost:9000"),
			},
		},
		Log: LogConfig{
			Dir:   getEnv("LOG_DIR", "logs"),
			Level: getEnv("LOG_LEVEL", "info"),
		},
		Backup: BackupConfig{
			Dir: getEnv("BACKUP_DIR", "backups"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.example.com"),
			Port:     587,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@shiyemo.com"),
			Enabled:  getEnv("SMTP_ENABLED", "false") == "true",
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
