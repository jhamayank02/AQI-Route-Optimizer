package main

import (
	env "github.com/jhamayank02/AQI-Route-Optimizer/config/env"

	"github.com/jhamayank02/AQI-Route-Optimizer/app"
)

func main() {
	// Load env variables
	env.Load()

	// Initialize app config
	cfg := app.NewConfig()
	// Initialize app
	app := app.NewApp(cfg)

	// Run http server
	app.Run()
}
