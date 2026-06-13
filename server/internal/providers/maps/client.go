package maps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/jhamayank02/AQI-Route-Optimizer/internal/domain"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/redis"
)

type orsGeoJSONResponse struct {
	Features []orsFeature `json:"features"`
}

type orsFeature struct {
	Properties orsProperties `json:"properties"`
	Geometry   orsGeometry   `json:"geometry"`
}

type orsProperties struct {
	Summary orsSummary `json:"summary"`
}

type orsSummary struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

type orsGeometry struct {
	Coordinates [][]float64 `json:"coordinates"`
}

type ProviderError struct {
	StatusCode int
	Message    string
	Retryable  bool
}

func (e *ProviderError) Error() string {
	return e.Message
}

type geocodeResponse struct {
	Features []geocodeFeature `json:"features"`
}

type geocodeFeature struct {
	Geometry   geocodeGeometry   `json:"geometry"`
	Properties geocodeProperties `json:"properties"`
}

type geocodeGeometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type geocodeProperties struct {
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	Confidence float64 `json:"confidence"`
}

type Client struct {
	logger            *slog.Logger
	httpClient        *http.Client
	apiKey            string
	findRouteURL      string
	searchLocationURL string
	redisProvider     *redis.RedisConfig
}

const routeCacheTTL = 15 * time.Minute

func NewClient(logger *slog.Logger, apiKey string, findRouteURL string, searchLocationURL string, redisProvider *redis.RedisConfig) *Client {
	return &Client{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey:            apiKey,
		findRouteURL:      findRouteURL,
		searchLocationURL: searchLocationURL,
		redisProvider:     redisProvider,
	}
}

func (c *Client) FindRoutes(ctx context.Context, start domain.Coordinates, dest domain.Coordinates) ([]domain.Route, error) {
	var cacheKey string
	if c.redisProvider != nil {
		cacheKey = c.redisProvider.RouteCacheKey(start.Lat, start.Lng, dest.Lat, dest.Lng)
		routes, err := c.redisProvider.GetRoutes(ctx, cacheKey)
		if err == nil {
			c.logger.Debug("using cached routes", "key", cacheKey)
			return routes, nil
		}

		c.logger.Debug("route cache unavailable or missed", "key", cacheKey, "error", err)
	}

	payload := map[string]any{
		"coordinates": [][]float64{
			{start.Lng, start.Lat},
			{dest.Lng, dest.Lat},
		},
		"alternative_routes": map[string]any{
			"target_count":  3,
			"share_factor":  0.6,
			"weight_factor": 1.4,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.findRouteURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to send route request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := handleORSStatus(res.StatusCode); err != nil {
		return nil, err
	}

	var orsResp orsGeoJSONResponse
	if err := json.Unmarshal(resBody, &orsResp); err != nil {
		return nil, err
	}

	if len(orsResp.Features) == 0 {
		return nil, fmt.Errorf("no route features returned")
	}

	routes := make([]domain.Route, 0, len(orsResp.Features))
	for _, feature := range orsResp.Features {
		coordinates := make([]domain.Coordinates, 0, len(feature.Geometry.Coordinates))
		for _, item := range feature.Geometry.Coordinates {
			if len(item) < 2 {
				continue
			}

			coordinates = append(coordinates, domain.Coordinates{
				Lng: item[0],
				Lat: item[1],
			})
		}

		routes = append(routes, domain.Route{
			DistanceKM:      feature.Properties.Summary.Distance / 1000,
			DurationMinutes: feature.Properties.Summary.Duration / 60,
			Coordinates:     coordinates,
		})
	}

	if c.redisProvider != nil && cacheKey != "" {
		if err := c.redisProvider.SetRoutes(ctx, cacheKey, routes, routeCacheTTL); err != nil {
			c.logger.Warn("failed to cache routes", "key", cacheKey, "error", err)
		}
	}

	return routes, nil
}

func (c *Client) SearchLocation(query string) ([]domain.LocationSuggestion, error) {
	targetURL := fmt.Sprintf("%s?api_key=%s&text=%s", c.searchLocationURL, c.apiKey, url.QueryEscape(query))

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to send geocode request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	var geoResp geocodeResponse
	if err := json.NewDecoder(res.Body).Decode(&geoResp); err != nil {
		return nil, err
	}

	suggestions := make([]domain.LocationSuggestion, 0, len(geoResp.Features))
	for _, feature := range geoResp.Features {
		if len(feature.Geometry.Coordinates) < 2 {
			continue
		}

		suggestions = append(suggestions, domain.LocationSuggestion{
			Label:      feature.Properties.Label,
			Name:       feature.Properties.Name,
			Lat:        feature.Geometry.Coordinates[1],
			Lng:        feature.Geometry.Coordinates[0],
			Country:    feature.Properties.Country,
			Region:     feature.Properties.Region,
			Confidence: feature.Properties.Confidence,
		})
	}

	return suggestions, nil
}

func handleORSStatus(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return &ProviderError{StatusCode: statusCode, Message: "invalid route request. please check coordinates or request body"}
	case http.StatusNotFound:
		return &ProviderError{StatusCode: statusCode, Message: "route not found for given coordinates"}
	case http.StatusMethodNotAllowed:
		return &ProviderError{StatusCode: statusCode, Message: "invalid HTTP method used for OpenRouteService"}
	case http.StatusRequestEntityTooLarge:
		return &ProviderError{StatusCode: statusCode, Message: "route request is too large"}
	case http.StatusInternalServerError:
		return &ProviderError{StatusCode: statusCode, Message: "OpenRouteService internal server error", Retryable: true}
	case http.StatusNotImplemented:
		return &ProviderError{StatusCode: statusCode, Message: "requested OpenRouteService feature is not supported"}
	case http.StatusServiceUnavailable:
		return &ProviderError{StatusCode: statusCode, Message: "OpenRouteService is temporarily unavailable", Retryable: true}
	default:
		return &ProviderError{StatusCode: statusCode, Message: "unexpected OpenRouteService error", Retryable: statusCode >= 500}
	}
}
