package service

import (
	"context"
	"fmt"

	"course-service/internal/domain"
)

// CourseRepository adalah interface yang diletakkan di layer Service
type CourseRepository interface {
	ReserveSeat(ctx context.Context, courseID int) error
	ReleaseSeat(ctx context.Context, courseID int) error
	Create(ctx context.Context, c *domain.Course) error
	GetByID(ctx context.Context, id int) (*domain.Course, error)
	GetAvailableCourses(ctx context.Context) ([]domain.Course, error)
}

type CourseService struct {
	repo CourseRepository
}

// NewCourseService adalah constructor (Dependency Injection)
func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

// CreateCourse memvalidasi input lalu menyimpannya
func (s *CourseService) CreateCourse(ctx context.Context, c *domain.Course) error {
	if c.Capacity <= 0 {
		return fmt.Errorf("%w: kapasitas harus lebih dari 0", domain.ErrInvalidInput)
	}
	return s.repo.Create(ctx, c)
}

// GetCourseByID mengambil course dan menangani error jika tidak ditemukan
func (s *CourseService) GetCourseByID(ctx context.Context, id int) (*domain.Course, error) {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, domain.ErrCourseNotFound
	}
	return course, nil
}

// GetAvailableCourses mengambil daftar course yang kuotanya masih ada
func (s *CourseService) GetAvailableCourses(ctx context.Context) ([]domain.Course, error) {
	return s.repo.GetAvailableCourses(ctx)
}

// Reserve mengamankan kursi (untuk dipakai service lain nanti)
func (s *CourseService) Reserve(ctx context.Context, courseID int) error {
	return s.repo.ReserveSeat(ctx, courseID)
}

// Release melepaskan kursi (untuk dipakai service lain nanti)
func (s *CourseService) Release(ctx context.Context, courseID int) error {
	return s.repo.ReleaseSeat(ctx, courseID)
}
