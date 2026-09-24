package sqlite

import (
	"log/slog"
	"os"
	"path/filepath"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Driver struct{}

func (s *Driver) Name() string {
	return "sqlite"
}

func (s *Driver) Initialize(dsn string) error {
	slog.Debug("sqlite driver initialize", "path", dsn)
	indexpath := dsn
	if err := os.MkdirAll(filepath.Dir(indexpath), os.ModePerm); err != nil {
		slog.Error("sqlite create directory failed", "path", filepath.Dir(indexpath), "err", err)
		return err
	}
	slog.Info("sqlite database initialized", "path", indexpath)
	return nil
}

func (s *Driver) Connect(dsn string) (*gorm.DB, error) {
	slog.Debug("sqlite driver connect", "path", dsn)
	indexpath := dsn
	gormLogger := slogGorm.New(
		slogGorm.WithTraceAll(),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)
	return gorm.Open(sqlite.Open(indexpath), &gorm.Config{
		Logger: gormLogger,
	})
}
