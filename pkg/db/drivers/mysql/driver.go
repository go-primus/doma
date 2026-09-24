package mysql

import (
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Driver struct{}

func (m *Driver) Name() string {
	return "mysql"
}

func (m *Driver) Initialize(dsn string) error {
	slog.Debug("mysql driver initialize", "dsn", dsn)
	return initialize(dsn)
}

func (m *Driver) Connect(dsn string) (*gorm.DB, error) {
	slog.Debug("mysql driver connect", "dsn", dsn)
	gormLogger := slogGorm.New(
		slogGorm.WithTraceAll(),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
}
