package app

import (
	"fmt"
	"net/http"

	db "github.com/jhamayank02/AQI-Route-Optimizer/config/db"
	env "github.com/jhamayank02/AQI-Route-Optimizer/config/env"
)

type Config struct {
	Addr string
}

func NewConfig() Config {
	port := env.GetString("PORT", ":8080")
	return Config{
		Addr: port,
	}
}

type App struct {
	config Config
}

func NewApp(cfg Config) App {
	return App{
		config: cfg,
	}
}

func (app *App) Run() error {
	db := db.DB

	fmt.Println("app.go", db)

	server := &http.Server{
		Addr: app.config.Addr,
	}
	return server.ListenAndServe()
}
