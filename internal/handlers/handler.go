package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/services/maps"
)

type Handler struct {
	logger    *slog.Logger
	mapConfig *maps.MapConfig
}

func NewHandler(logger *slog.Logger, mapConfig *maps.MapConfig) *Handler {
	return &Handler{
		logger:    logger,
		mapConfig: mapConfig,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func (h *Handler) GiveRoutes(c *gin.Context) {

}

func (h *Handler) FindRoutes(c *gin.Context) {
	// Get lat and lng from query params
	src_lat := c.Query("src_lat")
	src_lng := c.Query("src_lng")

	dst_lat := c.Query("dst_lat")
	dst_lng := c.Query("dst_lng")

	if strings.TrimSpace(src_lat) == "" || strings.TrimSpace(src_lng) == "" || strings.TrimSpace(dst_lat) == "" || strings.TrimSpace(dst_lng) == "" {
		c.JSON(400, gin.H{
			"message": "src_lat, src_lng, dst_lat and dst_lng are required",
		})
		return
	}

	src_lat_float, err := strconv.ParseFloat(src_lat, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid src_lat",
		})
		return
	}

	src_lng_float, err := strconv.ParseFloat(src_lng, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid src_lng",
		})
		return
	}

	dst_lat_float, err := strconv.ParseFloat(dst_lat, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid dst_lat",
		})
		return
	}

	dst_lng_float, err := strconv.ParseFloat(dst_lng, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid dst_lng",
		})
		return
	}

	srcCoordinates := maps.Coordinates{
		Lat: src_lat_float,
		Lng: src_lng_float,
	}
	dstCoordinates := maps.Coordinates{
		Lat: dst_lat_float,
		Lng: dst_lng_float,
	}

	response, err := h.mapConfig.FindRoutes(srcCoordinates, dstCoordinates)

	if err != nil {
		var providerErr *maps.ProviderError

		if errors.As(err, &providerErr) {
			status := http.StatusBadGateway

			if providerErr.StatusCode == 400 || providerErr.StatusCode == 404 {
				status = http.StatusBadRequest
			}

			c.JSON(status, gin.H{
				"message":   providerErr.Message,
				"retryable": providerErr.Retryable,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get route",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "OK",
		"data":    response,
	})
}
