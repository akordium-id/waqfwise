package cache

import (
	"context"
	"fmt"
	"time"
)

// RateLimiter implements rate limiting using Redis
type RateLimiter struct {
	redis  *Redis
	prefix string
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redis *Redis) *RateLimiter {
	return &RateLimiter{
		redis:  redis,
		prefix: "ratelimit:",
	}
}

// Allow checks if a request is allowed based on rate limiting
// Returns true if allowed, false if rate limit exceeded
func (rl *RateLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	redisKey := rl.getKey(key)

	// Increment counter
	count, err := rl.redis.Increment(ctx, redisKey)
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// Set expiration on first request
	if count == 1 {
		if err := rl.redis.Expire(ctx, redisKey, window); err != nil {
			return false, fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	return count <= limit, nil
}

// AllowN checks if n requests are allowed
func (rl *RateLimiter) AllowN(ctx context.Context, key string, n int64, limit int64, window time.Duration) (bool, error) {
	redisKey := rl.getKey(key)

	// Increment counter by n
	count, err := rl.redis.IncrementBy(ctx, redisKey, n)
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// Set expiration on first request
	if count == n {
		if err := rl.redis.Expire(ctx, redisKey, window); err != nil {
			return false, fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	return count <= limit, nil
}

// Reset resets the rate limit for a key
func (rl *RateLimiter) Reset(ctx context.Context, key string) error {
	redisKey := rl.getKey(key)
	return rl.redis.Delete(ctx, redisKey)
}

// GetCount returns the current count for a key
func (rl *RateLimiter) GetCount(ctx context.Context, key string) (int64, error) {
	redisKey := rl.getKey(key)
	result, err := rl.redis.Get(ctx, redisKey)
	if err != nil {
		return 0, nil // Key doesn't exist yet
	}

	var count int64
	_, err = fmt.Sscanf(result, "%d", &count)
	if err != nil {
		return 0, fmt.Errorf("failed to parse count: %w", err)
	}

	return count, nil
}

func (rl *RateLimiter) getKey(key string) string {
	return rl.prefix + key
}
