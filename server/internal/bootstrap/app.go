package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/config"
	httpHandlers "github.com/jhamayank02/AQI-Route-Optimizer/internal/http/handlers"
	httpRouter "github.com/jhamayank02/AQI-Route-Optimizer/internal/http/router"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/middlewares"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/aqi"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/maps"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/redis"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/services/routeplanner"
)

type Config struct {
	Addr string
}

func NewConfig(logger *slog.Logger) Config {
	return Config{
		Addr: config.GetString("PORT", ":8080", logger),
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCfg := redis.NewRedisConfig(
		config.GetString("REDIS_ADDR", "", app.logger),
		config.GetString("REDIS_PASS", "", app.logger),
		config.GetInt("REDIS_DB", 0, app.logger),
		config.GetInt("REDIS_PROTOCOL", 2, app.logger),
		app.logger,
	)
	_, err := redisCfg.SetupRedisClient(ctx)
	if err != nil {
		app.logger.Warn("Redis unavailable; continuing without route cache", "error", err)
	} else {
		defer func() {
			if err := redisCfg.Client.Close(); err != nil {
				app.logger.Warn("failed to close Redis client", "error", err)
			}
		}()
	}

	mapClient := maps.NewClient(
		app.logger,
		config.GetString("OPENROUTE_SERVICE_API_KEY", "", app.logger),
		config.GetString("FIND_ROUTE_URL", "", app.logger),
		config.GetString("SEARCH_LOCATION_URL", "", app.logger),
		redisCfg,
	)

	aqiClient := aqi.NewClient(
		app.logger,
		config.GetString("GET_AQI_URL", "", app.logger),
	)

	planner := routeplanner.NewService(app.logger, mapClient, aqiClient)
	handler := httpHandlers.NewHandler(app.logger, planner, mapClient)

	middleware := middlewares.NewMiddlewareConfig(app.logger, redisCfg)

	engine := gin.New()
	engine.Use(gin.Logger(), middleware.RecoveryMiddleware())
	engine.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://localhost:5173"},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		},
	))
	httpRouter.Register(engine, handler, middleware)

	server := &http.Server{
		Addr:              app.config.Addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}
