// Package handlers provides HTTP request handlers for the application's API endpoints
package handlers

import (
	"wg-edu-server/models"
)

// Handler holds dependencies for the handlers
type Handler struct {
	DB        *models.DB
	JWTSecret string
}
