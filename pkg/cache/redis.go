package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/akordium-id/waqfwise/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Redis wraps the Redis client
type Redis struct {
	client *redis.Client
	logger *zap.Logger
}

// New creates a new Redis client
func New(cfg *config.RedisConfig, log *zap.Logger) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	if log != nil {
		log.Info("redis connection established",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.Int("db", cfg.DB),
		)
	}

	return &Redis{
		client: client,
		logger: log,
	}, nil
}

// Client returns the underlying Redis client
func (r *Redis) Client() *redis.Client {
	return r.client
}

// Close closes the Redis connection
func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	if r.logger != nil {
		r.logger.Info("redis connection closed")
	}

	return nil
}

// HealthCheck checks if Redis connection is healthy
func (r *Redis) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// Set sets a key-value pair with expiration
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete deletes a key
func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *Redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire sets an expiration on a key
func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Increment increments a counter
func (r *Redis) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrementBy increments a counter by a specific amount
func (r *Redis) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// GetStats returns Redis statistics
func (r *Redis) GetStats() map[string]interface{} {
	stats := r.client.PoolStats()
	return map[string]interface{}{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}
}
