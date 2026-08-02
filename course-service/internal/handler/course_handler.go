package handler

import (
	"errors"
	"net/http"
	"strconv"

	"course-service/internal/domain"
	"course-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service *service.CourseService
}

// Helper internal untuk mengirim response sukses (Amplop Sukses)
func respondSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, domain.SuccessEnvelope{
		Success: true,
		Data:    data,
	})
}

// Helper internal untuk mengirim response error (Amplop Error)
func respondError(c *gin.Context, statusCode int, errorCode string, message string) {
	c.JSON(statusCode, domain.ErrorEnvelope{
		Success: false,
		Error: domain.ErrorDetail{
			Code:    errorCode,
			Message: message,
		},
	})
}

// NewCourseHandler berfungsi untuk mendaftarkan rute (routes) API
func NewCourseHandler(router *gin.Engine, service *service.CourseService) {
	handler := &CourseHandler{service: service}

	// Membuat grup endpoint /courses
	courseRoutes := router.Group("/courses")
	{
		courseRoutes.POST("", handler.CreateCourse)
		courseRoutes.GET("/available", handler.GetAvailableCourses)
		courseRoutes.GET("/:id", handler.GetCourseByID)
	}
}

// CreateCourse menangani pembuatan course baru
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var course domain.Course
	// Membaca JSON dari request body
	if err := c.ShouldBindJSON(&course); err != nil {
		respondError(c, http.StatusBadRequest, "ERR_VALIDATION", "Format input JSON tidak valid")
		return
	}

	// Memanggil layer Service
	if err := h.service.CreateCourse(c.Request.Context(), &course); err != nil {
		respondError(c, http.StatusInternalServerError, "ERR_INTERNAL", err.Error())
		return
	}

	// Response Sukses (201 Created)
	respondSuccess(c, http.StatusCreated, course)
}

// GetAvailableCourses mengambil course yang kuotanya masih ada
func (h *CourseHandler) GetAvailableCourses(c *gin.Context) {
	courses, err := h.service.GetAvailableCourses(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Gagal mengambil data course")
		return
	}

	// Response Sukses (200 OK)
	respondSuccess(c, http.StatusOK, courses)
}

// GetCourseByID mengambil detail satu course
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	// Mengambil ID dari URL (misal: /courses/1)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondError(c, http.StatusBadRequest, "ERR_VALIDATION", "Format ID harus berupa angka")
		return
	}

	course, err := h.service.GetCourseByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCourseNotFound) {
			respondError(c, http.StatusNotFound, "ERR_COURSE_NOT_FOUND", "course tidak ditemukan")
			return
		}
		respondError(c, http.StatusInternalServerError, "ERR_INTERNAL", "Terjadi kesalahan pada server")
		return
	}

	// Response Sukses (200 OK)
	respondSuccess(c, http.StatusOK, course)
}
