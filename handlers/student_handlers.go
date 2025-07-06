// Package handlers provides HTTP request handlers for the application's API endpoints
package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"wg-edu-server/models"
)

// HandleGetAllStudents retrieves all students
//
// Parameters:
//   - c: Gin context containing the request and response
//
// Returns:
//   - 200 OK with array of students on success
//   - 401 Unauthorized if not authenticated as admin
//   - 500 Internal Server Error on database failure
func (h *Handler) HandleGetAllStudents(c *gin.Context) {
	// Validate admin role
	claims, err := h.validateAdminToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Admin access required"})
		return
	}

	students, err := h.DB.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"admin_id": claims.UserID,
	})
}

// HandleGetStudent retrieves a specific student by ID
//
// Parameters:
//   - c: Gin context containing the request and response
//   - id: Student ID parameter from the URL
//
// Returns:
//   - 200 OK with student object on success
//   - 400 Bad Request if student ID is invalid
//   - 401 Unauthorized if not authenticated as admin
//   - 404 Not Found if student doesn't exist
//   - 500 Internal Server Error on database failure
func (h *Handler) HandleGetStudent(c *gin.Context) {
	// Validate admin role
	_, err := h.validateAdminToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Admin access required"})
		return
	}

	// Parse student ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	student, err := h.DB.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// HandleCreateStudent creates a new student
//
// Parameters:
//   - c: Gin context containing the request and response
//
// Expected Request Body:
//   - first_name: Student's first name
//   - last_name: Student's last name
//   - email: Student's email address
//   - grade: Student's grade/class
//   - username: Login username for the student
//   - password: Login password for the student
//
// Returns:
//   - 201 Created with the new student object on success
//   - 400 Bad Request if request data is invalid
//   - 401 Unauthorized if not authenticated as admin
//   - 500 Internal Server Error on database failure
func (h *Handler) HandleCreateStudent(c *gin.Context) {
	// Validate admin role
	_, err := h.validateAdminToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Admin access required"})
		return
	}

	var req models.StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	student, err := h.DB.CreateStudent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// HandleUpdateStudent updates an existing student
//
// Parameters:
//   - c: Gin context containing the request and response
//   - id: Student ID parameter from the URL
//
// Expected Request Body:
//   - first_name: Student's updated first name
//   - last_name: Student's updated last name
//   - email: Student's updated email address
//   - grade: Student's updated grade/class
//   - password: Updated password (optional)
//
// Returns:
//   - 200 OK with the updated student object on success
//   - 400 Bad Request if request data or ID is invalid
//   - 401 Unauthorized if not authenticated as admin
//   - 500 Internal Server Error on database failure
func (h *Handler) HandleUpdateStudent(c *gin.Context) {
	// Validate admin role
	_, err := h.validateAdminToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Admin access required"})
		return
	}

	// Parse student ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var req models.StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate required fields (except password which is optional for updates)
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	student, err := h.DB.UpdateStudent(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// HandleDeleteStudent deletes a student
//
// Parameters:
//   - c: Gin context containing the request and response
//   - id: Student ID parameter from the URL
//
// Returns:
//   - 200 OK with success message on successful deletion
//   - 400 Bad Request if student ID is invalid
//   - 401 Unauthorized if not authenticated as admin
//   - 500 Internal Server Error on database failure
func (h *Handler) HandleDeleteStudent(c *gin.Context) {
	// Validate admin role
	_, err := h.validateAdminToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Admin access required"})
		return
	}

	// Parse student ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	err = h.DB.DeleteStudent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
} 