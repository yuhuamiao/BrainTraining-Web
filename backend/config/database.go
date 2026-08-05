package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var dialector gorm.Dialector
	if strings.EqualFold(envOrDefault("DB_DRIVER", "sqlite"), "mysql") {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			envOrDefault("DB_HOST", "127.0.0.1"),
			envOrDefault("DB_PORT", "3306"),
			envOrDefault("DB_NAME", "brain_training"),
		)
		dialector = mysql.Open(dsn)
	} else {
		dbPath := envOrDefault("DB_PATH", "./data/brain-training.db")
		if dbPath != ":memory:" {
			if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
				return nil, fmt.Errorf("create database directory: %w", err)
			}
		}
		dialector = sqlite.Open(dbPath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
