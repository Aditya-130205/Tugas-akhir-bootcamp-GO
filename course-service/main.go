package main

import (
	"database/sql"
	"log"

	"course-service/internal/handler"
	"course-service/internal/repository"
	"course-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Konfigurasi Koneksi Database PostgreSQL
	dsn := "host=localhost user=user_course password=password_course dbname=db_course port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka database: %v", err)
	}
	defer db.Close()

	// Tes koneksi database
	if err := db.Ping(); err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}
	log.Println("Berhasil terhubung ke database PostgreSQL!")

	// 2. Dependency Injection (Merakit Layer)
	courseRepo := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)

	// 3. Setup Gin Router
	router := gin.Default()

	// 4. Mendaftarkan Route dari Handler
	handler.NewCourseHandler(router, courseService)

	// 5. Menjalankan Server di Port 8081
	log.Println("🚀 Course Service berjalan di http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
