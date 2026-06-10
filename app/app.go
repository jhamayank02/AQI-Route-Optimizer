package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	env "github.com/jhamayank02/AQI-Route-Optimizer/config/env"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/handlers"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/services/maps"
	"github.com/jhamayank02/AQI-Route-Optimizer/router"
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
	// db, err := db.NewDBConfig(app.logger)
	// if err != nil {
	// 	app.logger.Error("failed to initialize db", "error", err)
	// 	return err
	// }

	// Initialize maps
	mapApiKey := env.GetString("OPENROUTE_SERVICE_API_KEY", "", app.logger)
	findRouteUrl := env.GetString("FIND_ROUTE_URL", "", app.logger)
	_ = maps.NewMapConfig(app.logger, mapApiKey, findRouteUrl)

	// Initialize handlers
	mapConfig := maps.NewMapConfig(app.logger, mapApiKey, findRouteUrl)
	handlers := handlers.NewHandler(app.logger, mapConfig)

	// Initialize gin router
	r := gin.Default()

	// Register routes
	router.Register(r, handlers)

	server := &http.Server{
		Addr:    app.config.Addr,
		Handler: r,
	}
	return server.ListenAndServe()
}
