package mysql

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

func initialize(dsn string) error {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		slog.Error("mysql parse dsn failed", "dsn", dsn, "err", err)
		return err
	}
	dbName := config.DBName

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("mysql open connection failed", "dsn", dsn, "err", err)
		return err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT 1 FROM information_schema.SCHEMATA WHERE schema_name = ?", dbName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		slog.Warn("mysql check database existence failed", "db", dbName, "err", err)
	}

	if !exists {
		stmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)
		slog.Debug("mysql create database", "stmt", stmt)
		if _, err = db.Exec(stmt); err != nil {
			if mysqlError, ok := err.(*mysql.MySQLError); !ok || mysqlError.Number != 1049 {
				slog.Error("mysql create database failed", "db", dbName, "err", err)
				return err
			}
			config.DBName = ""
			db, err = sql.Open("mysql", config.FormatDSN())
			if err != nil {
				slog.Error("mysql reconnect failed", "dsn", config.FormatDSN(), "err", err)
				return err
			}
			defer db.Close()
			if _, err = db.Exec(stmt); err != nil {
				slog.Error("mysql create database after reconnect failed", "db", dbName, "err", err)
				return err
			}
		}
	}
	slog.Info("mysql database initialized", "db", dbName)
	return nil
}
