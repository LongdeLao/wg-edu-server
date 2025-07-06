// Package middleware provides HTTP middleware functions for the application.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth middleware validates JWT tokens and sets user information in the context
//
// Parameters:
//   - jwtSecret: Secret key for JWT validation
//
// Returns:
//   - gin.HandlerFunc: Middleware function for Gin router
//
// The middleware extracts the Bearer token from the Authorization header,
// validates it, and sets the user_id and role in the Gin context
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := extractToken(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// Set claims in context for use in handlers
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// AdminOnly middleware ensures the user has admin role
//
// Returns:
//   - gin.HandlerFunc: Middleware function for Gin router
//
// This middleware should be used after JWTAuth to check if the authenticated
// user has admin privileges
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if roleStr, ok := role.(string); !ok || roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TeacherOrAdmin middleware ensures the user has teacher or admin role
//
// Returns:
//   - gin.HandlerFunc: Middleware function for Gin router
//
// This middleware should be used after JWTAuth to check if the authenticated
// user has teacher or admin privileges
func TeacherOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if roleStr, ok := role.(string); !ok || (roleStr != "teacher" && roleStr != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Teacher or admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Helper function to extract token from Authorization header
func extractToken(authHeader string) string {
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}
