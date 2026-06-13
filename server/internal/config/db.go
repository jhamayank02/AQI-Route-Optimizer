package config

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	DB *sql.DB
}

func NewDBConfig(logger *slog.Logger) (*DBConfig, error) {
	logger.Info("initializing database")

	db, err := setupDB(logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		return nil, err
	}

	logger.Info("database initialized successfully")

	return &DBConfig{
		DB: db,
	}, nil
}

func setupDB(logger *slog.Logger) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = GetString("DB_USER", "root", logger)
	cfg.Passwd = GetString("DB_PASSWORD", "root", logger)
	cfg.Net = GetString("DB_NETWORK", "tcp", logger)
	cfg.Addr = GetString("DB_ADDRESS", "127.0.0.1:3306", logger)
	cfg.DBName = GetString("DB_NAME", "api_route_optimizer", logger)

	logger.Debug("connecting to database", "database", cfg.DBName, "address", cfg.Addr)

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Error("open mysql connection", "error", err)
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("ping mysql database", "error", err)
		return nil, fmt.Errorf("ping mysql database: %w", err)
	}

	return db, nil
}
