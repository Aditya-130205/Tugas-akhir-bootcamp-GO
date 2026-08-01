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
