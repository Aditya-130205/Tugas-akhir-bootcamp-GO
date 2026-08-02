package main

import (
	"database/sql"
	"log"

	"enrollment-service/internal/handler"
	"enrollment-service/internal/repository"
	"enrollment-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Koneksi ke db_enrollment di port 5433 (Port Docker db-enrollment)
	dsn := "host=localhost user=user_enrollment password=password_enrollment dbname=db_enrollment port=5433 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka database enrollment: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database enrollment tidak merespon: %v", err)
	}
	log.Println("Berhasil terhubung ke database Enrollment PostgreSQL (Port 5433)!")

	// URL Course Service yang berjalan di port 8081
	courseServiceURL := "http://localhost:8081"

	// Dependency Injection
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseServiceURL)

	router := gin.Default()
	handler.NewEnrollmentHandler(router, enrollmentService)

	// Enrollment Service berjalan di port 8082 agar tidak bentrok dengan Course Service (8081)
	log.Println("🚀 Enrollment Service berjalan di http://localhost:8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Gagal menjalankan server enrollment: %v", err)
	}
}
