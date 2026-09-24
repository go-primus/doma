package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func initialize(dsn string) error {
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		dsn = "postgres://" + dsn
	}
	u, err := url.Parse(dsn)
	if err != nil {
		slog.Error("postgres parse url failed", "dsn", dsn, "err", err)
		return err
	}
	dbName := strings.TrimPrefix(u.Path, "/")

	dsnWithoutDB := strings.Replace(dsn, u.Path, "/template1", 1)

	db, err := sql.Open("pgx", dsnWithoutDB)
	if err != nil {
		slog.Error("postgres open connection failed", "dsn", dsnWithoutDB, "err", err)
		return err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", dbName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		slog.Warn("postgres check database existence failed", "db", dbName, "err", err)
	}

	if !exists {
		stmt := fmt.Sprintf("CREATE DATABASE %s", dbName)
		slog.Debug("postgres create database", "stmt", stmt)
		if _, err = db.Exec(stmt); err != nil {
			slog.Error("postgres create database failed", "db", dbName, "err", err)
			return err
		}
	}
	slog.Info("postgres database initialized", "db", dbName)
	return nil
}
