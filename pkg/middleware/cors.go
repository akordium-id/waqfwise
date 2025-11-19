package middleware

import (
	"github.com/akordium-id/waqfwise/pkg/config"
	"github.com/gin-gonic/gin"
)

// CORS returns a gin middleware for CORS
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else if len(cfg.AllowedOrigins) > 0 {
				c.Header("Access-Control-Allow-Origin", cfg.AllowedOrigins[0])
			}

			c.Header("Access-Control-Allow-Credentials", "true")

			// Set allowed methods
			methods := "GET,POST,PUT,DELETE,OPTIONS"
			if len(cfg.AllowedMethods) > 0 {
				methods = ""
				for i, method := range cfg.AllowedMethods {
					if i > 0 {
						methods += ","
					}
					methods += method
				}
			}
			c.Header("Access-Control-Allow-Methods", methods)

			// Set allowed headers
			headers := "Content-Type,Authorization"
			if len(cfg.AllowedHeaders) > 0 {
				headers = ""
				for i, header := range cfg.AllowedHeaders {
					if i > 0 {
						headers += ","
					}
					headers += header
				}
			}
			c.Header("Access-Control-Allow-Headers", headers)
			c.Header("Access-Control-Max-Age", "86400")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
