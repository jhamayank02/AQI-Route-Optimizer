package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jhamayank02/AQI-Route-Optimizer/internal/domain"
	redisclient "github.com/redis/go-redis/v9"
)

var ErrClientNotInitialized = errors.New("redis client is not initialized")

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Protocol int
	logger   *slog.Logger
	Client   *redisclient.Client
}

func NewRedisConfig(addr string, password string, db int, protocol int, logger *slog.Logger) *RedisConfig {
	return &RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: protocol,
		logger:   logger,
	}
}

func (r *RedisConfig) SetupRedisClient(ctx context.Context) (*redisclient.Client, error) {
	redisClient := redisclient.NewClient(&redisclient.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
		Protocol: r.Protocol,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		_ = redisClient.Close()
		r.logger.Error("failed to connect to Redis", "error", err)
		return nil, err
	}

	r.Client = redisClient
	r.logger.Info("connected to Redis", "address", r.Addr)
	return redisClient, nil
}

func (r *RedisConfig) RouteCacheKey(srcLat, srcLng, dstLat, dstLng float64) string {
	key := fmt.Sprintf("route:%.5f:%.5f:%.5f:%.5f", srcLat, srcLng, dstLat, dstLng)
	r.logger.Debug("generated Redis route cache key", "key", key)
	return key
}

func (r *RedisConfig) SetRoutes(ctx context.Context, key string, routes []domain.Route, ttl time.Duration) error {
	if r.Client == nil {
		return ErrClientNotInitialized
	}

	data, err := json.Marshal(routes)
	if err != nil {
		r.logger.Error("failed to marshal routes", "error", err)
		return err
	}

	return r.Client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisConfig) GetRoutes(ctx context.Context, key string) ([]domain.Route, error) {
	if r.Client == nil {
		return nil, ErrClientNotInitialized
	}

	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redisclient.Nil) {
			r.logger.Debug("route cache miss", "key", key)
			return nil, err
		}

		r.logger.Error("failed to get routes", "error", err)
		return nil, err
	}

	var routes []domain.Route
	if err := json.Unmarshal(data, &routes); err != nil {
		r.logger.Error("failed to unmarshal routes", "error", err)
		return nil, err
	}

	r.logger.Debug("fetched routes from Redis", "key", key)
	return routes, nil
}
