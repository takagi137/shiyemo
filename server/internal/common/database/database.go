package database

import (
	"fmt"
	"log"

	"github.com/shiyemo/server/internal/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(cfg *config.DBConfig) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// SetDB 设置数据库实例（用于测试）
func SetDB(db *gorm.DB) {
	DB = db
}
