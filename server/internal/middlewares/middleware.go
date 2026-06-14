package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/redis"
)

type MiddlewareConfig struct {
	logger *slog.Logger
	redis  *redis.RedisConfig
}

func NewMiddlewareConfig(logger *slog.Logger, redis *redis.RedisConfig) *MiddlewareConfig {
	return &MiddlewareConfig{
		logger: logger,
		redis:  redis,
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

func (m *MiddlewareConfig) RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.redis == nil || m.redis.Client == nil {
			m.logger.Warn(
				"rate limiter skipped because Redis is unavailable",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
			)
			c.Next()
			return
		}

		ctx := c.Request.Context()
		ip := c.ClientIP()

		key := fmt.Sprintf("rate_limit:%s:%s", c.Request.URL.Path, ip)

		count, err := m.redis.Client.Incr(ctx, key).Result()
		if err != nil {
			m.logger.Error("failed to increment rate limit count", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			c.Abort()
			return
		}

		if count == 1 {
			if err := m.redis.Client.Expire(ctx, key, window).Err(); err != nil {
				m.logger.Warn("failed to set rate limit expiry", "error", err, "key", key)
			}
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
