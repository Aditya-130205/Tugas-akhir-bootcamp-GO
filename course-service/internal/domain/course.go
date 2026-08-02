package domain

import (
	"errors"
	"time"
)

// Course merepresentasikan tabel courses di database
type Course struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Lecturer    string    `json:"lecturer"`
	Capacity    int       `json:"capacity"`
	Taken       int       `json:"taken"`
	CreatedAt   time.Time `json:"created_at"`
}

// Sentinel Errors untuk layer Service sesuai kontrak (BR-04, BR-01)
var (
	ErrCourseNotFound = errors.New("course tidak ditemukan")
	ErrNoSeat         = errors.New("kuota course sudah penuh")
	ErrInvalidInput   = errors.New("input tidak valid")
)

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
