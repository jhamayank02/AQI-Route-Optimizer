# AQI Route Optimizer

AQI Route Optimizer is a small full-stack app for finding routes that are not only fast, but also cleaner to travel through.

You enter a source and destination, the app searches for possible driving routes, samples air quality along each route, scores the options, and highlights the best one based on both travel time and AQI exposure. It is built as a Vue client backed by a Go API, with OpenRouteService handling maps/geocoding and Open-Meteo providing air-quality data.

## Why This Exists

Most route planners optimize for time, distance, tolls, or traffic. That is useful, but it leaves out something people in polluted cities care about every day: the air they are moving through.

This project tries to answer a more human question:

> If two routes get me there, which one is easier on my lungs?

It is especially useful as a portfolio project because it touches several real-world concerns at once: external APIs, route geometry, AQI sampling, ranking logic, caching, typed frontend state, and Dockerized deployment.

## What It Does

- Searches places from free-text input.
- Finds multiple driving route alternatives between two points.
- Samples air quality at evenly spaced points along each route.
- Scores every route using AQI and travel duration.
- Selects the route with the best combined score.
- Shows the selected route and AQI details in the client.
- Caches route responses through Redis when Redis is available.

## Screenshots

Add screenshots here once the UI is ready to show.

### Home / Search

<!-- Screenshot placeholder: add an image of the search screen here. -->

```md
![Home search screen](docs/screenshots/home-search.png)
```

### Route Results

<!-- Screenshot placeholder: add an image of the map and recommended route here. -->

```md
![Route results screen](docs/screenshots/route-results.png)
```

### How It Works

<!-- Screenshot placeholder: add an image of the explanation page here. -->

```md
![How it works screen](docs/screenshots/how-it-works.png)
```

## Tech Stack

### Client

- Vue 3
- TypeScript
- Vite
- Vue Router
- Pinia
- Leaflet and `@vue-leaflet/vue-leaflet`
- Tailwind CSS
- Vue Toastification
- Nginx for the production container

### Server

- Go
- Gin
- `slog` structured logging
- OpenRouteService API
- Open-Meteo Air Quality API
- Redis for optional route caching
- Docker multi-stage build

### Infrastructure

- Docker
- Docker Compose
- Redis container
- Nginx reverse proxy from the client container to the API container

## Folder Structure

```text
.
|-- client/
|   |-- public/
|   |   |-- favicon files and web manifest
|   |-- src/
|   |   |-- components/
|   |   |   |-- AboutPage.vue
|   |   |   |-- Home.vue
|   |   |   |-- HowItWorksPage.vue
|   |   |   |-- LoadingPage.vue
|   |   |   |-- Map.vue
|   |   |   `-- SearchBox.vue
|   |   |-- constants/
|   |   |   `-- api.ts
|   |   |-- router/
|   |   |   `-- router.ts
|   |   |-- services/
|   |   |   |-- api.service.ts
|   |   |   |-- location.service.ts
|   |   |   `-- route.service.ts
|   |   |-- store/
|   |   |   `-- store.ts
|   |   |-- utils/
|   |   |   `-- utils.ts
|   |   |-- App.vue
|   |   |-- index.css
|   |   `-- main.ts
|   |-- Dockerfile
|   |-- nginx.conf
|   |-- package.json
|   |-- tailwind.config.ts
|   |-- tsconfig*.json
|   `-- vite.config.ts
|
|-- server/
|   |-- cmd/
|   |   `-- api/
|   |       `-- main.go
|   |-- deployment/
|   |   `-- Dockerfile
|   |-- internal/
|   |   |-- bootstrap/
|   |   |   `-- app.go
|   |   |-- config/
|   |   |   |-- db.go
|   |   |   `-- env.go
|   |   |-- domain/
|   |   |   |-- location.go
|   |   |   `-- route.go
|   |   |-- http/
|   |   |   |-- handlers/
|   |   |   |   `-- handler.go
|   |   |   `-- router/
|   |   |       `-- router.go
|   |   |-- providers/
|   |   |   |-- aqi/
|   |   |   |   `-- client.go
|   |   |   |-- maps/
|   |   |   |   `-- client.go
|   |   |   `-- redis/
|   |   |       `-- redis.go
|   |   `-- services/
|   |       `-- routeplanner/
|   |           `-- service.go
|   |-- .env.example
|   |-- go.mod
|   `-- go.sum
|
|-- docker-compose.yml
`-- README.md
```

## How The App Works

1. The user searches for a source and destination in the Vue client.
2. The client calls the Go API with both coordinates.
3. The server asks OpenRouteService for route alternatives.
4. Each route is sampled at a fixed number of points.
5. The AQI provider fetches air-quality data for those sampled coordinates.
6. The route planner calculates average AQI, max AQI, travel duration, and a combined score.
7. The best-scoring route is returned to the client and shown on the map.

The scoring is intentionally simple and readable. It penalizes both higher AQI and longer travel time, then picks the route with the highest score. If two routes score the same, the shorter one wins.

## Best Practices Used

- Clear client/server split, so the UI and API can evolve independently.
- Standard Go service layout with `cmd/` for the entrypoint and `internal/` for application code.
- Business logic kept in `internal/services/routeplanner` instead of being buried in HTTP handlers.
- Provider code isolated under `internal/providers`, making external APIs easier to replace or mock later.
- Domain models kept in `internal/domain` so request, response, and route types stay consistent.
- Context-aware HTTP requests on the server.
- HTTP client timeouts for external API calls.
- Structured logging with Go's `slog`.
- Redis caching is optional; the app still runs if Redis is unavailable.
- Concurrent AQI lookups for sampled route points.
- Typed frontend services and state with TypeScript and Pinia.
- Environment-based configuration instead of hard-coded secrets.
- Docker multi-stage builds for smaller production images.
- Nginx proxying in the client container, so the frontend can call `/api` in production.

## Environment Variables

Create a server environment file:

```bash
cp server/.env.example server/.env
```

Then update the values:

```env
PORT=":8000"
OPENROUTE_SERVICE_API_KEY="your_openrouteservice_api_key"
FIND_ROUTE_URL="https://api.openrouteservice.org/v2/directions/driving-car/geojson"
SEARCH_LOCATION_URL="https://api.openrouteservice.org/geocode/search"
GET_AQI_URL="https://air-quality-api.open-meteo.com/v1/air-quality"

# Optional, used for route caching
REDIS_ADDR="localhost:6379"
REDIS_PASS=""
REDIS_DB="0"
REDIS_PROTOCOL="2"
```

The `.env.example` also includes database values. They are present for future persistence work, but the current route recommendation flow does not require a database.

For local frontend development, you can create a client env file if your API is not available through `/api`:

```env
VITE_API_BASE_URL="http://localhost:8000/api"
```

## Run Without Docker

Run the server:

```bash
cd server
go mod download
go run ./cmd/api
```

Run the client in another terminal:

```bash
cd client
npm install
npm run dev
```

By default:

- API runs on `http://localhost:8000`
- Vite client runs on `http://localhost:5173`
- Client API calls use `VITE_API_BASE_URL` when provided, otherwise `/api`

Redis is optional in local development. If `REDIS_ADDR` is empty or Redis is not running, the server logs a warning and continues without route caching.

## Run With Docker

Create `server/.env` first, then run:

```bash
docker compose up --build
```

The Compose setup starts:

- `client` on `http://localhost`
- `server` on `http://localhost:8000`
- `redis` on `localhost:6379`

For caching inside Docker, set this in `server/.env`:

```env
REDIS_ADDR="redis:6379"
REDIS_PASS=""
REDIS_DB="0"
REDIS_PROTOCOL="2"
```

Stop the stack:

```bash
docker compose down
```

If you also want to remove the Redis volume:

```bash
docker compose down -v
```

## API Endpoints

### Health Check

```http
GET /api/health
```

### Search Locations

```http
GET /api/locations/search?query=delhi
```

### Get Route Recommendation

```http
GET /api/routes/recommendation?src_lat=28.6139&src_lng=77.2090&dst_lat=28.4595&dst_lng=77.0266
```

Example response shape:

```json
{
  "message": "OK",
  "data": {
    "selected_route_index": 0,
    "selected_route": {
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
      "sampling_strategy": "evenly distributed points across the route polyline",
      "is_selected": true,
      "selection_reason": "selected because it has the best combined AQI and travel-time score"
    },
    "all_routes": []
  }
}
```
