package repository

import (
	"context"
	"database/sql"
	"errors"

	"enrollment-service/internal/domain"
)

type EnrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create menyimpan data pendaftaran baru
func (r *EnrollmentRepository) Create(ctx context.Context, e *domain.Enrollment) error {
	query := `
		INSERT INTO enrollments (student_id, course_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	if e.Status == "" {
		e.Status = domain.StatusActive
	}

	return r.db.QueryRowContext(ctx, query, e.StudentID, e.CourseID, e.Status).Scan(&e.ID, &e.CreatedAt)
}

// GetByID mengambil pendaftaran berdasarkan ID
func (r *EnrollmentRepository) GetByID(ctx context.Context, id int) (*domain.Enrollment, error) {
	query := `
		SELECT id, student_id, course_id, status, created_at
		FROM enrollments
		WHERE id = $1
	`
	var e domain.Enrollment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.StudentID,
		&e.CourseID,
		&e.Status,
		&e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// GetByStudentIDAndCourseID mengecek apakah siswa sudah daftar di kelas tertentu
func (r *EnrollmentRepository) GetByStudentIDAndCourseID(ctx context.Context, studentID string, courseID int) (*domain.Enrollment, error) {
	query := `
		SELECT id, student_id, course_id, status, created_at
		FROM enrollments
		WHERE student_id = $1 AND course_id = $2 AND status = 'ACTIVE'
	`
	var e domain.Enrollment
	err := r.db.QueryRowContext(ctx, query, studentID, courseID).Scan(
		&e.ID,
		&e.StudentID,
		&e.CourseID,
		&e.Status,
		&e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}
