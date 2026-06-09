package app

import (
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/jhamayank02/AQI-Route-Optimizer/config/db"
	env "github.com/jhamayank02/AQI-Route-Optimizer/config/env"
)

type Config struct {
	Addr string
}

func NewConfig(logger *slog.Logger) Config {
	port := env.GetString("PORT", ":8080", logger)
	return Config{
		Addr: port,
	}
}

type App struct {
	config Config
	logger *slog.Logger
}

func NewApp(cfg Config, logger *slog.Logger) App {
	return App{
		config: cfg,
		logger: logger,
	}
}

func (app *App) Run() error {
	db, err := db.NewDBConfig(app.logger)
	if err != nil {
		app.logger.Error("failed to initialize db", "error", err)
		return err
	}

	fmt.Println("app.go", db)

	server := &http.Server{
		Addr: app.config.Addr,
	}
	return server.ListenAndServe()
}
