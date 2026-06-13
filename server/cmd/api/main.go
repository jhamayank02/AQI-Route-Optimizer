package main

import (
	"log/slog"
	"os"

	"github.com/jhamayank02/AQI-Route-Optimizer/internal/bootstrap"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/config"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	config.LoadEnv(logger)

	cfg := bootstrap.NewConfig(logger)
	app := bootstrap.NewApp(cfg, logger)

	if err := app.Run(); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
