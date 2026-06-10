package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/domain"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/maps"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/services/routeplanner"
)

type Handler struct {
	logger    *slog.Logger
	planner   *routeplanner.Service
	mapClient *maps.Client
}

func NewHandler(logger *slog.Logger, planner *routeplanner.Service, mapClient *maps.Client) *Handler {
	return &Handler{
		logger:    logger,
		planner:   planner,
		mapClient: mapClient,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (h *Handler) SearchLocation(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query is required",
		})
		return
	}

	suggestions, err := h.mapClient.SearchLocation(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to search location",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    suggestions,
	})
}

func (h *Handler) RecommendRoute(c *gin.Context) {
	var req domain.RouteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "src_lat, src_lng, dst_lat and dst_lng are required as valid numbers",
		})
		return
	}

	result, err := h.planner.Recommend(c.Request.Context(), req)
	if err != nil {
		var providerErr *maps.ProviderError
		if errors.As(err, &providerErr) {
			status := http.StatusBadGateway
			if providerErr.StatusCode == http.StatusBadRequest || providerErr.StatusCode == http.StatusNotFound {
				status = http.StatusBadRequest
			}

			c.JSON(status, gin.H{
				"message":   providerErr.Message,
				"retryable": providerErr.Retryable,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to recommend route",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    result,
	})
}
