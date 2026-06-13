package aqi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type openMeteoResponse struct {
	Current struct {
		USAQI float64 `json:"us_aqi"`
	} `json:"current"`
}

type Client struct {
	logger     *slog.Logger
	httpClient *http.Client
	baseURL    string
}

func NewClient(logger *slog.Logger, baseURL string) *Client {
	return &Client{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) GetAQI(ctx context.Context, lat float64, lng float64) (float64, error) {
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=us_aqi", c.baseURL, lat, lng)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to fetch AQI", "error", err, "lat", lat, "lng", lng)
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("aqi provider returned status %d", res.StatusCode)
	}

	var payload openMeteoResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return 0, err
	}

	return payload.Current.USAQI, nil
}
