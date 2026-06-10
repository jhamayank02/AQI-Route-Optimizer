# AQI Route Optimizer

AQI Route Optimizer is a Go backend service that helps users find a cleaner route between a source and destination. It combines route geometry from OpenRouteService with air-quality data from Open-Meteo, then scores the route based on AQI exposure and travel duration.

The project is now organized using a more standard Go service layout:

- `cmd/` contains application entrypoints
- `internal/bootstrap` wires dependencies and starts the server
- `internal/http` contains transport concerns like handlers and route registration
- `internal/providers` contains external API clients
- `internal/services` contains business logic
- `internal/domain` contains shared request/response and domain models

## Features

- Search locations from free-text input
- Fetch route geometry between source and destination
- Sample AQI along the route path
- Compute a simple route score using AQI and duration
- Return route metadata and AQI observations as JSON

## Tech Stack

- Go
- Gin
- OpenRouteService API
- Open-Meteo Air Quality API

## Project Structure

```text
.
|-- cmd/
|   `-- api/
|       `-- main.go
|-- internal/
|   |-- bootstrap/
|   |   `-- app.go
|   |-- config/
|   |   |-- db.go
|   |   `-- env.go
|   |-- domain/
|   |   |-- location.go
|   |   `-- route.go
|   |-- http/
|   |   |-- handlers/
|   |   |   `-- handler.go
|   |   `-- router/
|   |       `-- router.go
|   |-- providers/
|   |   |-- aqi/
|   |   |   `-- client.go
|   |   `-- maps/
|   |       `-- client.go
|   `-- services/
|       `-- routeplanner/
|           `-- service.go
|-- .air.toml
|-- .env.example
|-- go.mod
|-- go.sum
`-- README.md
```

## Request Flow

1. Client sends source and destination coordinates.
2. The route planner asks OpenRouteService for the route polyline.
3. The planner samples points across the route.
4. The AQI provider fetches air-quality data for those points.
5. The service returns route details, AQI samples, and a score.

## Environment Variables

Copy `.env.example` to `.env` and update the values:

```env
PORT=":8000"
DB_USER="root"
DB_PASSWORD="root"
DB_NETWORK="tcp"
DB_ADDRESS="127.0.0.1:3306"
DB_NAME="aqi_route_optimizer"
OPENROUTE_SERVICE_API_KEY="your_openrouteservice_api_key"
FIND_ROUTE_URL="https://api.openrouteservice.org/v2/directions/driving-car/geojson"
SEARCH_LOCATION_URL="https://api.openrouteservice.org/geocode/search"
GET_AQI_URL="https://air-quality-api.open-meteo.com/v1/air-quality"
```

Note: the database config exists in the repository, but the current route recommendation flow does not require persistence yet.

## Run Locally

```bash
go mod tidy
go run ./cmd/api
```

Server starts on the configured `PORT`.

## API Endpoints

### Health Check

```http
GET /api/health
```

### Search Locations

```http
GET /api/locations/search?query=delhi
```

### Get AQI-Based Route Recommendation

```http
GET /api/routes/recommendation?src_lat=28.6139&src_lng=77.2090&dst_lat=28.4595&dst_lng=77.0266
```

### Example Response

```json
{
  "message": "OK",
  "data": {
    "route": {
      "distance_km": 32.4,
      "duration_minutes": 58.7,
      "coordinates": [
        {
          "lat": 28.6139,
          "lng": 77.209
        }
      ]
    },
    "aqi_samples": [
      {
        "lat": 28.6139,
        "lng": 77.209,
        "aqi": 84
      }
    ],
    "average_aqi": 84,
    "max_aqi": 102,
    "route_score": 62.22,
    "recommendation": "moderate air quality; acceptable for most users",
    "sampling_strategy": "evenly distributed points across the route polyline"
  }
}
```

## Why This Structure Is Better

- Business logic is separated from HTTP and provider details
- External APIs are isolated behind provider clients
- Shared models live in one predictable place
- Entry-point code is minimal and easy to test or extend
- The route recommendation logic has a dedicated service instead of being embedded in handlers

## Suggested Next Improvements

- Add unit tests for `routeplanner.Service`
- Support multiple route alternatives and rank them
- Persist route searches and AQI history in a database
- Add request validation and structured API error responses
- Introduce interfaces where mocking is needed for tests
