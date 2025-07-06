package handlers

import (
	"log"
	"net/http"
	"strconv"

	"wg-edu-server/models"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error message response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success message response
type SuccessResponse struct {
	Message string `json:"message"`
}

// GetAllSubjects handles GET request to retrieve all subjects
// @Summary Get all subjects
// @Description Retrieves a list of all subjects
// @Tags subjects
// @Produce json
// @Success 200 {array} models.Subject
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects [get]
func (h *Handler) GetAllSubjects(c *gin.Context) {
	subjects, err := h.DB.GetAllSubjects()
	if err != nil {
		log.Printf("Error getting subjects: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve subjects"})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// GetSubjectsByGrade handles GET request to retrieve subjects by grade
// @Summary Get subjects by grade
// @Description Retrieves a list of subjects for a specific grade
// @Tags subjects
// @Produce json
// @Param grade path string true "Grade (PIB, IB1, IB2)"
// @Success 200 {array} models.Subject
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/{grade} [get]
func (h *Handler) GetSubjectsByGrade(c *gin.Context) {
	grade := c.Param("grade")
	if grade != "PIB" && grade != "IB1" && grade != "IB2" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid grade. Must be PIB, IB1, or IB2"})
		return
	}

	subjects, err := h.DB.GetSubjectsByGrade(grade)
	if err != nil {
		log.Printf("Error getting subjects for grade %s: %v", grade, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve subjects"})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// GetSubjectByID handles GET request to retrieve a subject by ID
// @Summary Get subject by ID
// @Description Retrieves a subject by its ID
// @Tags subjects
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} models.Subject
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/id/{id} [get]
func (h *Handler) GetSubjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid subject ID"})
		return
	}

	subject, err := h.DB.GetSubjectByID(id)
	if err != nil {
		log.Printf("Error getting subject with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Subject not found"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// GetAllTeachers handles GET request to retrieve all teachers with their subjects
// @Summary Get all teachers
// @Description Retrieves a list of all teachers with their assigned subjects
// @Tags teachers
// @Produce json
// @Success 200 {array} models.Teacher
// @Failure 500 {object} ErrorResponse
// @Router /api/teachers [get]
func (h *Handler) GetAllTeachers(c *gin.Context) {
	teachers, err := h.DB.GetAllTeachers()
	if err != nil {
		log.Printf("Error getting teachers: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve teachers"})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// GetTeacherByID handles GET request to retrieve a teacher by ID with their subjects
// @Summary Get teacher by ID
// @Description Retrieves a teacher by their ID with assigned subjects
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} models.Teacher
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/teachers/{id} [get]
func (h *Handler) GetTeacherByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid teacher ID"})
		return
	}

	teacher, err := h.DB.GetTeacherByID(id)
	if err != nil {
		log.Printf("Error getting teacher with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// AssignSubjectRequest represents the request body for assigning a subject
type AssignSubjectRequest struct {
	SubjectID int `json:"subject_id"`
}

// AssignSubjectToTeacher handles POST request to assign a subject to a teacher
// @Summary Assign subject to teacher
// @Description Assigns a subject to a teacher
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param request body AssignSubjectRequest true "Subject assignment request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/teachers/{id}/subjects [post]
func (h *Handler) AssignSubjectToTeacher(c *gin.Context) {
	// Verify admin role
	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized: Admin access required"})
		return
	}

	idStr := c.Param("id")
	teacherID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid teacher ID"})
		return
	}

	var req AssignSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	err = h.DB.AssignSubjectToTeacher(teacherID, req.SubjectID)
	if err != nil {
		log.Printf("Error assigning subject %d to teacher %d: %v", req.SubjectID, teacherID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to assign subject to teacher"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Subject assigned to teacher successfully"})
}

// RemoveSubjectFromTeacher handles DELETE request to remove a subject from a teacher
// @Summary Remove subject from teacher
// @Description Removes a subject assignment from a teacher
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Param subjectId path int true "Subject ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/teachers/{id}/subjects/{subjectId} [delete]
func (h *Handler) RemoveSubjectFromTeacher(c *gin.Context) {
	// Verify admin role
	if !isAdmin(c) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized: Admin access required"})
		return
	}

	teacherIDStr := c.Param("id")
	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid teacher ID"})
		return
	}

	subjectIDStr := c.Param("subjectId")
	subjectID, err := strconv.Atoi(subjectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid subject ID"})
		return
	}

	err = h.DB.RemoveSubjectFromTeacher(teacherID, subjectID)
	if err != nil {
		log.Printf("Error removing subject %d from teacher %d: %v", subjectID, teacherID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to remove subject from teacher"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Subject removed from teacher successfully"})
}

// GetAllSubjectsGrouped handles GET request to retrieve all subjects grouped by grade
// @Summary Get all subjects grouped by grade
// @Description Retrieves a list of all subjects organized by grade level
// @Tags subjects
// @Produce json
// @Success 200 {object} map[string][]models.Subject
// @Failure 500 {object} ErrorResponse
// @Router /api/subjects/grouped [get]
func (h *Handler) GetAllSubjectsGrouped(c *gin.Context) {
	subjects, err := h.DB.GetAllSubjects()
	if err != nil {
		log.Printf("Error getting subjects: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve subjects"})
		return
	}

	// Group subjects by grade
	grouped := make(map[string][]*models.Subject)
	for _, subject := range subjects {
		grade := subject.Grade
		grouped[grade] = append(grouped[grade], subject)
	}

	c.JSON(http.StatusOK, grouped)
}

// Helper function to check if user has admin role
func isAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}
	return role == "admin"
}
