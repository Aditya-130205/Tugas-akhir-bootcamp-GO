package repository

import (
	"context"
	"database/sql"
	"errors"

	"course-service/internal/domain"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create menyimpan data course baru ke database
func (r *CourseRepository) Create(ctx context.Context, c *domain.Course) error {
	query := `
		INSERT INTO courses (code, name, description, lecturer, capacity, taken) 
		VALUES ($1, $2, $3, $4, $5, 0) 
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query, c.Code, c.Name, c.Description, c.Lecturer, c.Capacity).
		Scan(&c.ID, &c.CreatedAt)
}

// GetByID mengambil satu course berdasarkan ID
func (r *CourseRepository) GetByID(ctx context.Context, id int) (*domain.Course, error) {
	query := `SELECT id, code, name, description, lecturer, capacity, taken, created_at FROM courses WHERE id = $1`

	var c domain.Course
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Code, &c.Name, &c.Description, &c.Lecturer, &c.Capacity, &c.Taken, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCourseNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ReserveSeat mengamankan 1 kursi (Atomic Update untuk BR-01)
func (r *CourseRepository) ReserveSeat(ctx context.Context, courseID int) error {
	query := `
		UPDATE courses 
		SET taken = taken + 1 
		WHERE id = $1 AND capacity - taken > 0 
		RETURNING id;
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNoSeat // Kuota habis atau ID tidak ada
		}
		return err
	}
	return nil
}

// ReleaseSeat mengembalikan 1 kursi (batal enroll)
func (r *CourseRepository) ReleaseSeat(ctx context.Context, courseID int) error {
	query := `
		UPDATE courses 
		SET taken = taken - 1 
		WHERE id = $1 AND taken > 0 
		RETURNING id;
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrCourseNotFound
		}
		return err
	}
	return nil
}

// GetAvailableCourses mengambil daftar course yang kuotanya masih ada
func (r *CourseRepository) GetAvailableCourses(ctx context.Context) ([]domain.Course, error) {
	query := `
		SELECT id, code, name, description, lecturer, capacity, taken, created_at 
		FROM courses 
		WHERE capacity - taken > 0
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var c domain.Course
		err := rows.Scan(
			&c.ID, &c.Code, &c.Name, &c.Description, &c.Lecturer, &c.Capacity, &c.Taken, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
