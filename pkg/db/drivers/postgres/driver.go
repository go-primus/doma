package postgres

import (
	"log/slog"
	"strings"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Driver struct{}

func (p *Driver) Name() string {
	return "postgres"
}

func (p *Driver) Initialize(dsn string) error {
	slog.Debug("postgres driver initialize", "dsn", dsn)
	return initialize(dsn)
}

func (p *Driver) Connect(dsn string) (*gorm.DB, error) {
	slog.Debug("postgres driver connect", "dsn", dsn)
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		dsn = "postgres://" + dsn
	}
	gormLogger := slogGorm.New(
		slogGorm.WithTraceAll(),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
}
