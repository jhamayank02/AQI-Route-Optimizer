package domain

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type LocationSuggestion struct {
	Label      string  `json:"label"`
	Name       string  `json:"name"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	Confidence float64 `json:"confidence"`
}
