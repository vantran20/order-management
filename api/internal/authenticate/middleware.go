package authenticate

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiddlewareAuth struct {
	authService AuthService
}

func NewAuthMiddleware(authService AuthService) *MiddlewareAuth {
	return &MiddlewareAuth{
		authService: authService,
	}
}

func (m *MiddlewareAuth) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the token is in the correct format (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			switch {
			case errors.Is(err, ErrTokenExpired):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// RateLimiter middleware limits request frequency
func RateLimiter() gin.HandlerFunc {
	// In a real application, you would use a proper rate limiting library
	// This is a simplified example
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Check if rate limit exceeded (implementation omitted)
		if isRateLimitExceeded(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}

		c.Next()
	}
}

// Helper function for rate limiting
func isRateLimitExceeded(clientIP string) bool {
	// In a real app, you would implement proper rate limiting
	// using Redis or another store to track requests
	return false
}
