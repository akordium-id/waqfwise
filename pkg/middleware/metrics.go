package middleware

import (
	"strconv"
	"time"

	"github.com/akordium-id/waqfwise/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// Metrics returns a gin middleware for collecting metrics
func Metrics(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		m.RecordHTTPRequest(
			c.Request.Method,
			path,
			status,
			duration,
			c.Request.ContentLength,
			int64(c.Writer.Size()),
		)
	}
}
