// Package handlers provides HTTP request handlers for the application's API endpoints
package handlers

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"wg-edu-server/models"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// HandleLogin processes login requests
//
// Parameters:
//   - c: Gin context containing the request and response
//
// The function expects a JSON body with username and password.
// It returns a 200 OK with JWT token on success, or appropriate error status code
func (h *Handler) HandleLogin(c *gin.Context) {
	// Parse the request body
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the user
	user, err := h.DB.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check the password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Create a JWT token
	token, err := createToken(user, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token and user info
	resp := LoginResponse{
		Token: token,
		User:  *user,
	}

	c.JSON(http.StatusOK, resp)
}

// HandleProtected processes protected route requests
//
// Parameters:
//   - c: Gin context containing the request and response
//
// The function validates the JWT token in the Authorization header.
// It returns a 200 OK with protected data on success, or 401 Unauthorized on failure
func (h *Handler) HandleProtected(c *gin.Context) {
	// Extract and validate JWT token from Authorization header
	tokenString := extractToken(c)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse and validate the token
	claims, err := validateToken(tokenString, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Return protected data
	c.JSON(http.StatusOK, gin.H{
		"message": "Protected data accessed successfully",
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}

// Claims represents the JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Helper function to create a JWT token
func createToken(user *models.User, secret string) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Helper function to extract token from request
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// Helper function to validate JWT token
func validateToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// validateAdminToken validates that the request contains a valid admin token
//
// Parameters:
//   - c: Gin context containing the request and response
//
// Returns:
//   - *Claims: The JWT claims if token is valid and has admin role
//   - error: Error if token is invalid or does not have admin role
func (h *Handler) validateAdminToken(c *gin.Context) (*Claims, error) {
	tokenString := extractToken(c)
	if tokenString == "" {
		return nil, fmt.Errorf("no token provided")
	}

	claims, err := validateToken(tokenString, h.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Verify admin role
	if claims.Role != "admin" {
		return nil, fmt.Errorf("admin role required")
	}

	return claims, nil
} 