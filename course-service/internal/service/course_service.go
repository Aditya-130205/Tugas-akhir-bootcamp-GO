package service

import (
	"context"
	"fmt"

	"course-service/internal/domain"
)

// CourseRepository adalah interface yang diminta oleh spesifikasi tugas akhir.
// Interface diletakkan di layer Service agar layer ini tidak bergantung pada detail database (SQL).
type CourseRepository interface {
	ReserveSeat(ctx context.Context, courseID int) error
	ReleaseSeat(ctx context.Context, courseID int) error
	Create(ctx context.Context, c *domain.Course) error
	GetByID(ctx context.Context, id int) (*domain.Course, error)
	GetAvailableCourses(ctx context.Context) ([]domain.Course, error)
}

type CourseService struct {
	repo CourseRepository // Menyuntikkan interface, bukan struct (Dependency Injection)
}

func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

// CreateCourse memvalidasi input sebelum diteruskan ke repository
func (s *CourseService) CreateCourse(ctx context.Context, c *domain.Course) error {
	if c.Capacity <= 0 {
		return fmt.Errorf("%w: kapasitas harus lebih dari 0", domain.ErrInvalidInput)
	}
	return s.repo.Create(ctx, c)
}

func (s *CourseService) GetAvailableCourses(ctx context.Context) ([]domain.Course, error) {
	return s.repo.GetAvailableCourses(ctx)
}

// Reserve mengamankan kursi dengan error handling domain
func (s *CourseService) Reserve(ctx context.Context, courseID int) error {
	return s.repo.ReserveSeat(ctx, courseID)
}

// Release melepaskan kursi dengan error handling domain
func (s *CourseService) Release(ctx context.Context, courseID int) error {
	return s.repo.ReleaseSeat(ctx, courseID)
}
