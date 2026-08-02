# QuickCourse - Sistem Microservices Pendaftaran Mata Kuliah 🎓

Sistem backend *microservices* untuk pendaftaran mata kuliah online (*enrollment*) berbasis Go, PostgreSQL, dan Docker Compose.

---

## 🚀 Cara Menjalankan

Cukup jalankan satu perintah dari direktori *root* repository untuk menyalakan seluruh service dan database:

```bash
# Jalankan dari keadaan bersih (Build & Up)
docker compose down -v
docker compose up --build