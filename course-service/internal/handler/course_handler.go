package handler

import (
	"net/http"
	"strconv"

	"course-service/internal/domain"
	"course-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service *service.CourseService
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input JSON tidak valid"})
		return
	}

	// Memanggil layer Service
	if err := h.service.CreateCourse(c.Request.Context(), &course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course berhasil dibuat", "data": course})
}

// GetAvailableCourses mengambil course yang kuotanya masih ada
func (h *CourseHandler) GetAvailableCourses(c *gin.Context) {
	courses, err := h.service.GetAvailableCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": courses})
}

// GetCourseByID mengambil detail satu course
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	// Mengambil ID dari URL (misal: /courses/1)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID harus berupa angka"})
		return
	}

	course, err := h.service.GetCourseByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrCourseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}
