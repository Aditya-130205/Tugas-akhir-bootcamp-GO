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

// ==========================================
// Struct Response Amplop (Envelope)
// ==========================================

// SuccessEnvelope untuk format JSON saat HTTP 2xx
type SuccessEnvelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ErrorDetail penampung objek error
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorEnvelope untuk format JSON saat HTTP 4xx/5xx
type ErrorEnvelope struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}
