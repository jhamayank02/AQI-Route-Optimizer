package config

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	env "github.com/jhamayank02/AQI-Route-Optimizer/config/env"
)

var DB *sql.DB

func init() {
	fmt.Println("Initializing db configuration...")
	DB, err := setupDB()
	if err != nil {
		fmt.Println("Error initializing db:", err)
	}
	fmt.Println("Db initialized successfully", DB)
}

func setupDB() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = env.GetString("DB_USER", "root")
	cfg.Passwd = env.GetString("DB_PASSWORD", "root")
	cfg.Net = env.GetString("DB_NETWORK", "tcp")
	cfg.Addr = env.GetString("DB_ADDRESS", "127.0.0.1:3306")
	cfg.DBName = env.GetString("DB_NAME", "api_route_optimizer")

	fmt.Println("Connecting to db", cfg.DBName, cfg.FormatDSN())

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil || db == nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}
