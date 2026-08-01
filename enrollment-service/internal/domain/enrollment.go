package domain

import (
	"errors"
	"time"
)

// Pesan error standar untuk domain Enrollment
var (
	ErrEnrollmentNotFound = errors.New("data pendaftaran tidak ditemukan")
	ErrCourseFull         = errors.New("kuota kelas sudah penuh")
	ErrAlreadyEnrolled    = errors.New("siswa sudah terdaftar di kelas ini")
)

// Status Enrollment
const (
	StatusActive    = "ACTIVE"
	StatusCancelled = "CANCELLED"
)

// Enrollment merepresentasikan tabel pendaftaran di database
type Enrollment struct {
	ID        int       `json:"id"`
	StudentID string    `json:"student_id"` // Menggunakan string (VARCHAR)
	CourseID  int       `json:"course_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
