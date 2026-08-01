CREATE TABLE enrollments (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(30) NOT NULL,
    course_id INT NOT NULL, -- tanpa foreign key: tabel course ada di database lain
    status VARCHAR(15) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Penegakan BR-02 di level database
CREATE UNIQUE INDEX satu_pendaftaran_aktif 
ON enrollments (student_id, course_id) 
WHERE status = 'ACTIVE';