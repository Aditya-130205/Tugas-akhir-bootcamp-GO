package handler

import (
	"errors"
	"net/http"
	"strconv"

	"enrollment-service/internal/domain"
	"enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	service *service.EnrollmentService
}

type CreateEnrollmentInput struct {
	StudentID string `json:"student_id" binding:"required"`
	CourseID  int    `json:"course_id" binding:"required"`
}

func NewEnrollmentHandler(r *gin.Engine, s *service.EnrollmentService) {
	h := &EnrollmentHandler{service: s}

	routes := r.Group("/enrollments")
	{
		routes.POST("", h.CreateEnrollment)
		routes.GET("/:id", h.GetEnrollmentByID)
	}
}

func (h *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
	var input CreateEnrollmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid: " + err.Error()})
		return
	}

	enrollment, err := h.service.EnrollStudent(c.Request.Context(), input.StudentID, input.CourseID)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyEnrolled) || errors.Is(err, domain.ErrCourseFull) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

func (h *EnrollmentHandler) GetEnrollmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pendaftaran tidak valid"})
		return
	}

	enrollment, err := h.service.GetEnrollmentByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrEnrollmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, enrollment)
}
