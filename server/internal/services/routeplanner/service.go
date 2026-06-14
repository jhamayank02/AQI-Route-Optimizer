package routeplanner

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"sync"

	"github.com/jhamayank02/AQI-Route-Optimizer/internal/domain"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/aqi"
	"github.com/jhamayank02/AQI-Route-Optimizer/internal/providers/maps"
)

type Service struct {
	logger    *slog.Logger
	mapClient *maps.Client
	aqiClient *aqi.Client
}

func NewService(logger *slog.Logger, mapClient *maps.Client, aqiClient *aqi.Client) *Service {
	return &Service{
		logger:    logger,
		mapClient: mapClient,
		aqiClient: aqiClient,
	}
}

func (s *Service) Recommend(ctx context.Context, req domain.RouteRequest) (*domain.RouteRecommendation, error) {
	routes, err := s.mapClient.FindRoutes(ctx, req.Source(), req.Destination())
	if err != nil {
		return nil, err
	}

	evaluatedRoutes := make([]domain.EvaluatedRoute, 0, len(routes))
	for _, route := range routes {
		evaluatedRoute, err := s.evaluateRoute(ctx, route)
		if err != nil {
			return nil, err
		}
		evaluatedRoutes = append(evaluatedRoutes, evaluatedRoute)
	}

	sort.SliceStable(evaluatedRoutes, func(i int, j int) bool {
		if evaluatedRoutes[i].RouteScore == evaluatedRoutes[j].RouteScore {
			return evaluatedRoutes[i].Route.DurationMinutes < evaluatedRoutes[j].Route.DurationMinutes
		}
		return evaluatedRoutes[i].RouteScore > evaluatedRoutes[j].RouteScore
	})

	if len(evaluatedRoutes) == 0 {
		return nil, fmt.Errorf("no routes available to evaluate")
	}

	evaluatedRoutes[0].IsSelected = true
	evaluatedRoutes[0].SelectionReason = "selected because it has the best combined AQI and travel-time score"

	return &domain.RouteRecommendation{
		SelectedRouteIndex: 0,
		SelectedRoute:      evaluatedRoutes[0],
		AllRoutes:          evaluatedRoutes,
	}, nil
}

func (s *Service) evaluateRoute(ctx context.Context, route domain.Route) (domain.EvaluatedRoute, error) {
	samples := sampleCoordinates(route.Coordinates, 8)
	aqiSamples := make([]domain.AQISample, 0, len(samples))

	var total float64
	var maxAQI float64

	var wg *sync.WaitGroup
	wg = &sync.WaitGroup{}
	resultChan := make(chan domain.AQIResult, len(samples))

	for _, coordinate := range samples {
		wg.Add(1)
		go func(cord domain.Coordinates) {
			defer wg.Done()
			resultSent := false
			defer func() {
				if rec := recover(); rec != nil {
					s.logger.Error(
						"panic recovered while fetching AQI sample",
						"error", rec,
						"lat", cord.Lat,
						"lng", cord.Lng,
					)
					if !resultSent {
						resultChan <- domain.AQIResult{
							Error: fmt.Errorf("internal server error"),
						}
					}
				}
			}()

			aqiValue, err := s.aqiClient.GetAQI(ctx, cord.Lat, cord.Lng)

			resultChan <- domain.AQIResult{
				Sample: domain.AQISample{
					Lat: cord.Lat,
					Lng: cord.Lng,
					AQI: aqiValue,
				},
				Error: err,
			}
			resultSent = true
		}(coordinate)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {

		if result.Error != nil {
			return domain.EvaluatedRoute{}, result.Error
		}

		aqiSamples = append(
			aqiSamples,
			result.Sample,
		)

		total += result.Sample.AQI

		maxAQI = math.Max(
			maxAQI,
			result.Sample.AQI,
		)
	}

	averageAQI := 0.0
	if len(aqiSamples) > 0 {
		averageAQI = total / float64(len(aqiSamples))
	}

	return domain.EvaluatedRoute{
		Route:            route,
		AQISamples:       aqiSamples,
		AverageAQI:       averageAQI,
		MaxAQI:           maxAQI,
		RouteScore:       calculateScore(averageAQI, route.DurationMinutes),
		Recommendation:   recommendationText(averageAQI),
		SamplingStrategy: "evenly distributed points across the route polyline",
	}, nil
}

func calculateScore(averageAQI float64, durationMinutes float64) float64 {
	aqiPenalty := math.Min(averageAQI, 300) / 3
	durationPenalty := math.Min(durationMinutes, 180) / 6
	score := 100 - aqiPenalty - durationPenalty
	if score < 0 {
		return 0
	}
	return math.Round(score*100) / 100
}

func recommendationText(averageAQI float64) string {
	switch {
	case averageAQI <= 50:
		return "good air quality along most of the route"
	case averageAQI <= 100:
		return "moderate air quality; acceptable for most users"
	case averageAQI <= 150:
		return "unhealthy for sensitive groups; use caution"
	default:
		return "poor air quality; consider delaying travel or using protection"
	}
}

func sampleCoordinates(points []domain.Coordinates, limit int) []domain.Coordinates {
	if len(points) <= limit || limit <= 0 {
		return points
	}

	samples := make([]domain.Coordinates, 0, limit)
	lastIndex := len(points) - 1

	for i := 0; i < limit; i++ {
		index := int(math.Round(float64(i*lastIndex) / float64(limit-1)))
		samples = append(samples, points[index])
	}

	return samples
}
