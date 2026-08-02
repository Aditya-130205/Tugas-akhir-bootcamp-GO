# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency modules
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Compile binary Go (sesuaikan cmd/main.go dengan letak file main.go kamu)
RUN CGO_ENABLED=0 GOOS=linux go build -o course-service main.go

# Final Stage (Gunakan alpine ringan)
FROM alpine:3.18

WORKDIR /app

# Copy binary dari builder stage
COPY --from=builder /app/course-service .

# Expose port Course Service (8081)
EXPOSE 8081

# Jalankan aplikasi
CMD ["./course-service"]