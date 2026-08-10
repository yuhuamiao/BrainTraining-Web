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
	driver := strings.ToLower(envOrDefault("DB_DRIVER", "sqlite"))
	isSQLite := driver == "sqlite"
	switch driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			envOrDefault("DB_HOST", "127.0.0.1"),
			envOrDefault("DB_PORT", "3306"),
			envOrDefault("DB_NAME", "brain_training"),
		)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dbPath := envOrDefault("DB_PATH", "./data/brain-training.db")
		if dbPath != ":memory:" {
			if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
				return nil, fmt.Errorf("create database directory: %w", err)
			}
		}
		dsn := dbPath
		if dbPath != ":memory:" {
			separator := "?"
			if strings.Contains(dbPath, "?") {
				separator = "&"
			}
			dsn += separator + "_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"
		}
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q: expected sqlite or mysql", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if isSQLite {
		// SQLite has a single writer. Serializing access avoids intermittent lock errors under load.
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(0)
	} else {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	return db, nil
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
