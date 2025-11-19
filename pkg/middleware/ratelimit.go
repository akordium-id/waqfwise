package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/akordium-id/waqfwise/pkg/cache"
	"github.com/gin-gonic/gin"
)

// RateLimit returns a gin middleware for rate limiting
func RateLimit(rateLimiter *cache.RateLimiter, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as the key
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())

		// Check if request is allowed
		allowed, err := rateLimiter.Allow(c.Request.Context(), key, int64(requestsPerMinute), time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limit check failed",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
