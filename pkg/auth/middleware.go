package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader is the header name for authorization
	AuthorizationHeader = "Authorization"
	// BearerSchema is the authorization schema
	BearerSchema = "Bearer "
	// UserIDKey is the context key for user ID
	UserIDKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey = "user_email"
	// UserRoleKey is the context key for user role
	UserRoleKey = "user_role"
)

// JWTMiddleware returns a middleware that validates JWT tokens
func JWTMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, BearerSchema)

		// Validate token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole returns a middleware that requires a specific role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(UserRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role not found",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid user role type",
			})
			c.Abort()
			return
		}

		// Check if user has required role
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
		c.Abort()
	}
}

// GetUserID retrieves the user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail retrieves the user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}

	e, ok := email.(string)
	return e, ok
}

// GetUserRole retrieves the user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}

	r, ok := role.(string)
	return r, ok
}
