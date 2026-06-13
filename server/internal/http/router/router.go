package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/http/handlers"
)

func Register(r *gin.Engine, h *handlers.Handler) {
	api := r.Group("/api")

	api.GET("/health", h.HealthCheck)
	api.GET("/locations/search", h.SearchLocation)
	api.GET("/routes/recommendation", h.RecommendRoute)
}
