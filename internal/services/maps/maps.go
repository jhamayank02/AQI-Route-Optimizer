package maps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ORSGeoJSONResponse struct {
	Features []ORSFeature `json:"features"`
}

type ORSFeature struct {
	Properties ORSProperties `json:"properties"`
	Geometry   ORSGeometry   `json:"geometry"`
}

type ORSProperties struct {
	Summary ORSSummary `json:"summary"`
}

type ORSSummary struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

type ORSGeometry struct {
	Coordinates [][]float64 `json:"coordinates"`
}

type RouteResponse struct {
	DistanceKM      float64     `json:"distance_km"`
	DurationMinutes float64     `json:"duration_minutes"`
	Coordinates     [][]float64 `json:"coordinates"`
}

type ProviderError struct {
	StatusCode int
	Message    string
	Retryable  bool
}

func (e *ProviderError) Error() string {
	return e.Message
}

type Coordinates struct {
	Lat float64
	Lng float64
}

type MapConfig struct {
	logger       *slog.Logger
	apiKey       string
	findRouteUrl string
}

func NewMapConfig(logger *slog.Logger, apiKey string, findRouteUrl string) *MapConfig {
	return &MapConfig{
		logger:       logger,
		apiKey:       apiKey,
		findRouteUrl: findRouteUrl,
	}
}

func (mc *MapConfig) FindRoutes(start Coordinates, dest Coordinates) (*RouteResponse, error) {
	// targetUrl := fmt.Sprintf("%s?api_key=%s&start=%f,%f&end=%f,%f", mc.findRouteUrl, mc.apiKey, start.Lat, start.Lng, dest.Lat, dest.Lng)
	targetUrl := mc.findRouteUrl

	mc.logger.Info("targetUrl", "url", targetUrl)

	payload := map[string]interface{}{
		"coordinates": [][]float64{{start.Lat, start.Lng}, {dest.Lat, dest.Lng}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		mc.logger.Error("failed to marshal payload", "error", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, targetUrl, bytes.NewBuffer(body))
	if err != nil {
		mc.logger.Error("failed to create request", "error", err)
		return nil, err
	}

	req.Header.Set("Authorization", mc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		mc.logger.Error("failed to send request", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := handleORSStatus(res.StatusCode, resBody); err != nil {
		return nil, err
	}

	var orsResp ORSGeoJSONResponse

	if err := json.Unmarshal(resBody, &orsResp); err != nil {
		return nil, err
	}

	if len(orsResp.Features) == 0 {
		return nil, fmt.Errorf("no route features returned")
	}

	feature := orsResp.Features[0]

	result := RouteResponse{
		DistanceKM:      feature.Properties.Summary.Distance / 1000,
		DurationMinutes: feature.Properties.Summary.Duration / 60,
		Coordinates:     feature.Geometry.Coordinates,
	}

	return &result, nil
}

func handleORSStatus(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusOK:
		return nil

	case http.StatusBadRequest:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "invalid route request. please check coordinates or request body",
			Retryable:  false,
		}

	case http.StatusNotFound:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "route not found for given coordinates",
			Retryable:  false,
		}

	case http.StatusMethodNotAllowed:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "invalid HTTP method used for OpenRouteService",
			Retryable:  false,
		}

	case http.StatusRequestEntityTooLarge:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "route request is too large",
			Retryable:  false,
		}

	case http.StatusInternalServerError:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "OpenRouteService internal server error",
			Retryable:  true,
		}

	case http.StatusNotImplemented:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "requested OpenRouteService feature is not supported",
			Retryable:  false,
		}

	case http.StatusServiceUnavailable:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "OpenRouteService is temporarily unavailable",
			Retryable:  true,
		}

	default:
		return &ProviderError{
			StatusCode: statusCode,
			Message:    "unexpected OpenRouteService error",
			Retryable:  statusCode >= 500,
		}
	}
}
