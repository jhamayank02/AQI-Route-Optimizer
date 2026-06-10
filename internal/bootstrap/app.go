package bootstrap

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/config"
	httpHandlers "github.com/jhamayank02/AQI-Route-Optimizer/internal/http/handlers"
	httpRouter "github.com/jhamayank02/AQI-Route-Optimizer/internal/http/router"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/aqi"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/maps"
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
	mapClient := maps.NewClient(
		app.logger,
		config.GetString("OPENROUTE_SERVICE_API_KEY", "", app.logger),
		config.GetString("FIND_ROUTE_URL", "", app.logger),
		config.GetString("SEARCH_LOCATION_URL", "", app.logger),
	)

	aqiClient := aqi.NewClient(
		app.logger,
		config.GetString("GET_AQI_URL", "", app.logger),
	)

	planner := routeplanner.NewService(app.logger, mapClient, aqiClient)
	handler := httpHandlers.NewHandler(app.logger, planner, mapClient)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	httpRouter.Register(engine, handler)

	server := &http.Server{
		Addr:              app.config.Addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}
