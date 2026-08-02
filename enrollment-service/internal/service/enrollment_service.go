package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"enrollment-service/internal/domain"
	"enrollment-service/internal/repository"
)

// 1. Tambahkan struct pembungkus "data" agar sesuai dengan JSON Course Service
type CourseDetail struct {
	ID       int `json:"id"`
	Capacity int `json:"capacity"`
	Taken    int `json:"taken"`
}

type CourseResponse struct {
	Data CourseDetail `json:"data"` // <-- Ini kuncinya!
}

type EnrollmentService struct {
	repo             *repository.EnrollmentRepository
	courseServiceURL string
	httpClient       *http.Client
}

func NewEnrollmentService(repo *repository.EnrollmentRepository, courseServiceURL string) *EnrollmentService {
	return &EnrollmentService{
		repo:             repo,
		courseServiceURL: courseServiceURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// EnrollStudent memproses pendaftaran kelas oleh siswa
func (s *EnrollmentService) EnrollStudent(ctx context.Context, studentID string, courseID int) (*domain.Enrollment, error) {
	// 1. Cek apakah siswa sudah terdaftar di kelas ini sebelumnya
	existing, err := s.repo.GetByStudentIDAndCourseID(ctx, studentID, courseID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrAlreadyEnrolled
	}

	// 2. Komunikasi HTTP ke Course Service untuk verifikasi data & kuota
	url := fmt.Sprintf("%s/courses/%d", s.courseServiceURL, courseID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request ke course service: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi course service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("course dengan ID %d tidak ditemukan", courseID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("course service merespon dengan status: %d", resp.StatusCode)
	}

	var course CourseResponse
	if err := json.NewDecoder(resp.Body).Decode(&course); err != nil {
		return nil, fmt.Errorf("gagal membaca data course: %w", err)
	}

	// 3. Validasi kuota kelas (Diubah jadi course.Data.Taken dan course.Data.Capacity)
	if course.Data.Taken >= course.Data.Capacity {
		return nil, domain.ErrCourseFull
	}

	// 4. Simpan pendaftaran baru ke database enrollment
	enrollment := &domain.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    domain.StatusActive,
	}

	if err := s.repo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

// GetEnrollmentByID mengambil data pendaftaran berdasarkan ID
func (s *EnrollmentService) GetEnrollmentByID(ctx context.Context, id int) (*domain.Enrollment, error) {
	enrollment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, domain.ErrEnrollmentNotFound
	}
	return enrollment, nil
}

// CancelEnrollment membatalkan pendaftaran siswa & mengembalikan kuota
func (s *EnrollmentService) CancelEnrollment(ctx context.Context, id int) error {
	// 1. Cek apakah enrollment ada
	enrollment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if enrollment == nil {
		return domain.ErrEnrollmentNotFound
	}

	// 2. Jika sudah dibatalkan sebelumnya, lewati saja (Idempotent)
	if enrollment.Status == domain.StatusCancelled {
		return nil
	}

	// 3. Update status pendaftaran menjadi CANCELLED di database
	if err := s.repo.UpdateStatus(ctx, id, domain.StatusCancelled); err != nil {
		return fmt.Errorf("gagal mengubah status pendaftaran: %w", err)
	}

	// 4. (Opsional) Mengembalikan kuota ke Course Service jika endpoint release-nya ada
	// url := fmt.Sprintf("%s/courses/%d/release", s.courseServiceURL, enrollment.CourseID)
	// req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	// _, _ = s.httpClient.Do(req)

	return nil
}
