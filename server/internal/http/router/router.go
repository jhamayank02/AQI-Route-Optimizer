package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/http/handlers"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/middlewares"
)

func Register(r *gin.Engine, h *handlers.Handler, m *middlewares.MiddlewareConfig) {
	api := r.Group("/api")

	api.GET("/health", h.HealthCheck)
	api.GET("/locations/search", m.RateLimitMiddleware(10, 1*time.Minute), h.SearchLocation)
	api.GET("/routes/recommendation", m.RateLimitMiddleware(10, 1*time.Minute), h.RecommendRoute)
}
