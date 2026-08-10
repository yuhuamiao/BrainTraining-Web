package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSQLiteUsesSafeConnectionSettings(t *testing.T) {
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_PATH", filepath.Join(t.TempDir(), "brain-training.db"))

	db, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB() error: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() error: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	if got := sqlDB.Stats().MaxOpenConnections; got != 1 {
		t.Fatalf("MaxOpenConnections = %d, want 1", got)
	}

	var journalMode string
	if err := db.Raw("PRAGMA journal_mode").Scan(&journalMode).Error; err != nil {
		t.Fatalf("read journal_mode: %v", err)
	}
	if !strings.EqualFold(journalMode, "wal") {
		t.Fatalf("journal_mode = %q, want wal", journalMode)
	}

	var busyTimeout int
	if err := db.Raw("PRAGMA busy_timeout").Scan(&busyTimeout).Error; err != nil {
		t.Fatalf("read busy_timeout: %v", err)
	}
	if busyTimeout != 5000 {
		t.Fatalf("busy_timeout = %d, want 5000", busyTimeout)
	}
}

func TestInitDBRejectsUnknownDriver(t *testing.T) {
	t.Setenv("DB_DRIVER", "myql")
	if _, err := InitDB(); err == nil {
		t.Fatal("InitDB() accepted an unknown DB_DRIVER")
	}
}
