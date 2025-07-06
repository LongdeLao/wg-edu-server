// Package handlers provides HTTP request handlers for the application's API endpoints
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// HandleHealth processes health check requests
//
// Parameters:
//   - c: Gin context containing the request and response
//
// Returns:
//   - 200 OK with health status information
func (h *Handler) HandleHealth(c *gin.Context) {
	resp := HealthResponse{
		Status:    "OK",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, resp)
}
