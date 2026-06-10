package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/handlers"
)

func Register(r *gin.Engine, h *handlers.Handler) {
	rg := r.Group("/api")

	rg.GET("/health", h.HealthCheck)
	// rg.GET("/routes", handlers.GiveRoutes)
	rg.GET("/find-routes", h.FindRoutes)
	rg.GET("/search", h.SearchLocation)
}
