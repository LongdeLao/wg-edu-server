// Package routes provides API route definitions for the WG Education platform.
//
// This package configures the Gin router with all application endpoints,
// middleware, and route groups.
package routes

import (
	"wg-edu-server/handlers"
	"wg-edu-server/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
//
// Parameters:
//   - router: Gin router instance
//   - handler: Handler containing dependencies and endpoint handlers
//
// This function organizes routes into logical groups and applies middleware
func SetupRoutes(router *gin.Engine, handler *handlers.Handler) {
	// Add CORS middleware
	router.Use(CORSMiddleware())

	// API routes group
	api := router.Group("/api")
	{
		// Health check endpoint (public)
		api.GET("/health", handler.HandleHealth)

		// Login endpoint (public)
		api.POST("/login", handler.HandleLogin)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(handler.JWTSecret))
		{
			// General protected endpoint
			protected.GET("/protected", handler.HandleProtected)

			// Subject routes (available to all authenticated users)
			subjects := protected.Group("/subjects")
			{
				subjects.GET("", handler.GetAllSubjects)                // Get all subjects
				subjects.GET("/grouped", handler.GetAllSubjectsGrouped) // Get subjects grouped by grade
				subjects.GET("/:grade", handler.GetSubjectsByGrade)     // Get subjects by grade
				subjects.GET("/id/:id", handler.GetSubjectByID)         // Get subject by ID
			}

			// Teacher routes (available to teachers and admins)
			teachers := protected.Group("/teachers")
			teachers.Use(middleware.TeacherOrAdmin())
			{
				teachers.GET("", handler.GetAllTeachers)     // Get all teachers
				teachers.GET("/:id", handler.GetTeacherByID) // Get teacher by ID

				// Subject assignment (admin only)
				teacherAdmin := teachers.Group("")
				teacherAdmin.Use(middleware.AdminOnly())
				{
					teacherAdmin.POST("/:id/subjects", handler.AssignSubjectToTeacher)                // Assign subject
					teacherAdmin.DELETE("/:id/subjects/:subjectId", handler.RemoveSubjectFromTeacher) // Remove subject
				}
			}

			// Admin routes group
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				// Student management
				students := admin.Group("/students")
				{
					students.GET("", handler.HandleGetAllStudents)       // Get all students
					students.GET("/:id", handler.HandleGetStudent)       // Get specific student
					students.POST("", handler.HandleCreateStudent)       // Create new student
					students.PUT("/:id", handler.HandleUpdateStudent)    // Update student
					students.DELETE("/:id", handler.HandleDeleteStudent) // Delete student
				}
			}
		}
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
//
// Returns:
//   - gin.HandlerFunc: Middleware function for the Gin router
//
// This middleware adds appropriate CORS headers to allow cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
