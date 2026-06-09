package main

import (
	"log/slog"
	"os"

	env "github.com/jhamayank02/AQI-Route-Optimizer/config/env"

	"github.com/jhamayank02/AQI-Route-Optimizer/app"
)

func main() {
	// Initialize logger
	slog := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	// Load env variables
	env.Load(slog)

	// Initialize app config
	cfg := app.NewConfig(slog)
	// Initialize app
	app := app.NewApp(cfg, slog)

	// Run http server
	app.Run()
}
