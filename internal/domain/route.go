package domain

type RouteRequest struct {
	SrcLat float64 `form:"src_lat" binding:"required"`
	SrcLng float64 `form:"src_lng" binding:"required"`
	DstLat float64 `form:"dst_lat" binding:"required"`
	DstLng float64 `form:"dst_lng" binding:"required"`
}

func (r RouteRequest) Source() Coordinates {
	return Coordinates{Lat: r.SrcLat, Lng: r.SrcLng}
}

func (r RouteRequest) Destination() Coordinates {
	return Coordinates{Lat: r.DstLat, Lng: r.DstLng}
}

type Route struct {
	DistanceKM      float64       `json:"distance_km"`
	DurationMinutes float64       `json:"duration_minutes"`
	Coordinates     []Coordinates `json:"coordinates"`
}

type AQISample struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	AQI float64 `json:"aqi"`
}

type EvaluatedRoute struct {
	Route            Route       `json:"route"`
	AQISamples       []AQISample `json:"aqi_samples"`
	AverageAQI       float64     `json:"average_aqi"`
	MaxAQI           float64     `json:"max_aqi"`
	RouteScore       float64     `json:"route_score"`
	Recommendation   string      `json:"recommendation"`
	SamplingStrategy string      `json:"sampling_strategy"`
	IsSelected       bool        `json:"is_selected"`
	SelectionReason  string      `json:"selection_reason,omitempty"`
}

type RouteRecommendation struct {
	SelectedRouteIndex int              `json:"selected_route_index"`
	SelectedRoute      EvaluatedRoute   `json:"selected_route"`
	AllRoutes          []EvaluatedRoute `json:"all_routes"`
}
