package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type MiddlewareConfig struct {
	logger *slog.Logger
}

func NewMiddlewareConfig(logger *slog.Logger) *MiddlewareConfig {
	return &MiddlewareConfig{
		logger: logger,
	}
}

func (m *MiddlewareConfig) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				m.logger.Error(
					"panic recovered",
					"error", rec,
					"stack", string(debug.Stack()),
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}()

		c.Next()
	}
}
