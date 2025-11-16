# Dokumentasi Lengkap: CV AI Evaluator Backend

## 📋 Tentang Project

**CV AI Evaluator** adalah sistem backend berbasis Go (Golang) yang mengotomasi proses screening kandidat dengan AI. Sistem ini menerima CV dan laporan proyek kandidat, kemudian mengevaluasinya menggunakan Large Language Model (LLM) lokal dengan pendekatan RAG (Retrieval-Augmented Generation).

### Fitur Utama
- ✅ Upload dokumen CV dan Project Report (PDF/Markdown)
- ✅ Evaluasi otomatis menggunakan AI lokal (Ollama)
- ✅ RAG pipeline dengan ChromaDB untuk vector search
- ✅ Asynchronous processing menggunakan Goroutines (tanpa Redis)
- ✅ RESTful API dengan 3 endpoint utama
- ✅ Penyimpanan hasil evaluasi di MySQL

### Tech Stack
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP Router)
- **Database**: MySQL 8.0
- **ORM**: GORM
- **Vector DB**: ChromaDB (chromem-go embedded)
- **LLM**: Ollama (model gemma3:4b)
- **PDF Processing**: UniPDF

---

## 📁 Struktur Project

```
cv-ai-evaluator/
│
├── cmd/
│   └── api/
│       └── main.go                          # Entry point aplikasi
│
├── internal/
│   ├── models/                              # Database models (GORM)
│   │   ├── uploaded_document.go             # Model dokumen upload
│   │   ├── evaluation_job.go                # Model evaluation job
│   │   └── ground_truth_document.go         # Model ground truth
│   │
│   ├── database/                            # Database connection
│   │   └── mysql.go                         # MySQL config & connection
│   │
│   ├── handlers/                            # HTTP handlers
│   │   ├── upload_handler.go                # Handler POST /upload
│   │   ├── evaluate_handler.go              # Handler POST /evaluate
│   │   └── result_handler.go                # Handler GET /result/{id}
│   │
│   ├── services/                            # Business logic layer
│   │   ├── document_service.go              # Service untuk dokumen
│   │   └── evaluation_service.go            # Service untuk evaluasi
│   │
│   └── worker/                              # Background worker
│       └── evaluation_worker.go             # Worker pool & AI pipeline
│
├── pkg/                                     # Shared packages
│   ├── utils/                               # Utilities
│   │   ├── pdf_extractor.go                 # Extract text dari PDF
│   │   └── text_reader.go                   # Read text/markdown files
│   │
│   ├── llm/                                 # LLM integration
│   │   └── ollama_client.go                 # Ollama client
│   │
│   └── vectordb/                            # Vector database
│       └── chroma_client.go                 # ChromaDB client
│
├── storage/
│   ├── uploads/                             # Uploaded files (CV & Report)
│   └── groundtruth/                         # Ground truth documents
│       ├── job_description_backend.md
│       ├── case_study_brief.md
│       ├── cv_scoring_rubric.md
│       └── project_scoring_rubric.md
│
├── scripts/
│   └── ingest_groundtruth.go                # Script ingest ground truth
│
├── config/
│   └── config.go                            # Configuration management
│
├── chroma_data/                             # ChromaDB persistent storage
│
├── .env                                     # Environment variables
├── go.mod                                   # Go module dependencies
├── go.sum                                   # Dependency checksums
└── README.md                                # Project documentation
```

### Penjelasan Layer

**1. CMD Layer** (`cmd/`)
- Entry point aplikasi
- Inisialisasi dependencies
- Setup router dan server

**2. Internal Layer** (`internal/`)
- **Models**: Struktur data database (GORM models)
- **Database**: Koneksi dan konfigurasi database
- **Handlers**: HTTP request handlers (controller)
- **Services**: Business logic (validasi, operasi data)
- **Worker**: Background processing (AI evaluation pipeline)

**3. PKG Layer** (`pkg/`)
- Shared utilities yang reusable
- External integrations (LLM, VectorDB)
- Pure functions tanpa side effects

**4. Storage Layer** (`storage/`)
- File uploads dari user
- Ground truth documents untuk RAG

**5. Config Layer** (`config/`)
- Environment configuration
- Database connection strings
- Service URLs

***

## 🚀 Cara Run Project

### Persiapan Environment

#### 1. Prerequisites
Pastikan sudah terinstall:
```bash
# Check Go version
go version
# Output: go version go1.21.x windows/amd64

# Check MySQL
mysql --version
# Output: mysql  Ver 8.0.44 for Win64

# Check Ollama
ollama --version
# Output: ollama version is 0.12.11

# Check model Ollama
ollama list
# Output: gemma3:4b    a2af6cc3eb7f    3.3 GB    ...
```

#### 2. Clone/Setup Project
```bash
# clone direktori project
git clone https://github.com/Adrian463588/go-cv-ai-evaluator.git
cd go-cv-ai-evaluator/cv-ai-evaluator

# Init Go module
go mod init cv-ai-evaluator


```

#### 3. Install Dependencies
```bash
# Framework & Router
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors

# Database
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# Utilities
go get github.com/google/uuid
go get github.com/joho/godotenv

# PDF Processing
go get github.com/unidoc/unipdf/v3

# ChromaDB (Embedded)
go get github.com/philippgille/chromem-go

# Tidy dependencies
go mod tidy
```

#### 4. Setup Database
```bash
# Masuk ke MySQL
mysql -u root -p

# Jalankan SQL
CREATE DATABASE cv_ai_evaluator CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE cv_ai_evaluator;

# Copy paste SQL dari STEP 3 di dokumentasi sebelumnya
# (Table creation scripts)
```

#### 5. Konfigurasi Environment
Buat file `.env` di root project:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=cv_ai_evaluator
SERVER_PORT=8080
OLLAMA_URL=http://localhost:11434
CHROMA_URL=http://localhost:8000
UPLOAD_DIR=./storage/uploads
```

#### 6. Persiapkan Ground Truth Documents
```bash
# Buat 4 file markdown di storage/groundtruth/:
# 1. job_description_backend.md
# 2. case_study_brief.md
# 3. cv_scoring_rubric.md
# 4. project_scoring_rubric.md

# Copy template dari dokumentasi sebelumnya
```

### Menjalankan Aplikasi

#### Step 1: Start Ollama
```bash
# Di terminal terpisah, pastikan Ollama running
ollama serve

# Test model
ollama run gemma3:4b "Hello, test"
```

#### Step 2: Ingest Ground Truth Documents
```bash
# Jalankan script ingestion
go run scripts/ingest_groundtruth.go
```

Expected output:
```
Database connected successfully!
Ingesting: job_description_backend.md
✓ Successfully ingested: job_description_backend.md
Ingesting: case_study_brief.md
✓ Successfully ingested: case_study_brief.md
Ingesting: cv_scoring_rubric.md
✓ Successfully ingested: cv_scoring_rubric.md
Ingesting: project_scoring_rubric.md
✓ Successfully ingested: project_scoring_rubric.md
Ground truth ingestion completed!
```

#### Step 3: Run Main Application
```bash
# Run dari root project
go run cmd/api/main.go
```

Expected output:
```
Database connected successfully!
Started 3 workers
Worker 1 started
Worker 2 started
Worker 3 started
[GIN-debug] POST   /upload                   --> cv-ai-evaluator/internal/handlers.(*UploadHandler).Upload-fm
[GIN-debug] POST   /evaluate                 --> cv-ai-evaluator/internal/handlers.(*EvaluateHandler).Evaluate-fm
[GIN-debug] GET    /result/:id               --> cv-ai-evaluator/internal/handlers.(*ResultHandler).GetResult-fm
[GIN-debug] GET    /health                   --> main.main.func1
Server starting on port 8080
[GIN-debug] Listening and serving HTTP on :8080
```

#### Alternative: Build Binary
```bash
# Build executable
go build -o cv-evaluator.exe cmd/api/main.go

# Run
./cv-evaluator.exe
```

***

## 🧪 Cara Testing dengan Postman

### Setup Postman Collection

#### 1. Create New Collection
- Nama: `CV AI Evaluator`
- Base URL Variable: `{{base_url}}` = `http://localhost:8080`

### Test Endpoint 1: Health Check

**Request:**
```
GET http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok"
}
```

### Test Endpoint 2: Upload Documents

**Request:**
```
POST http://localhost:8080/upload
Content-Type: multipart/form-data

Body (form-data):
- Key: cv, Type: File, Value: [pilih file CV.pdf]
- Key: report, Type: File, Value: [pilih file Report.pdf]
```

**Expected Response (200 OK):**
```json
{
  "cv_document_id": "550e8400-e29b-41d4-a716-446655440000",
  "report_document_id": "660e8400-e29b-41d4-a716-446655440001",
  "message": "Files uploaded successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "CV file is required"
}
```

### Test Endpoint 3: Start Evaluation

**Request:**
```
POST http://localhost:8080/evaluate
Content-Type: application/json

Body (raw JSON):
{
  "job_title": "Backend Engineer",
  "cv_id": "550e8400-e29b-41d4-a716-446655440000",
  "report_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

**Expected Response (200 OK):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "status": "queued"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "CV document not found with id: 550e8400-e29b-41d4-a716-446655440000"
}
```

### Test Endpoint 4: Get Evaluation Result

**Request (Status: Queued/Processing):**
```
GET http://localhost:8080/result/770e8400-e29b-41d4-a716-446655440002
```

**Response (Processing):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "status": "processing"
}
```

**Request (Status: Completed):**
```
GET http://localhost:8080/result/770e8400-e29b-41d4-a716-446655440002
```

**Response (Completed - after 2-5 minutes):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "status": "completed",
  "result": {
    "cv_match_rate": 0.82,
    "cv_feedback": "Strong technical skills in Go and MySQL with 4+ years backend experience. Demonstrated expertise in REST APIs and microservices. Limited exposure to AI/ML integration. Good communication skills evident from documentation.",
    "project_score": 4.5,
    "project_feedback": "Excellent implementation of RAG pipeline with proper error handling. Clean code structure following best practices. Comprehensive documentation. Minor improvement needed in retry logic for LLM failures.",
    "overall_summary": "Strong hire recommendation. Candidate demonstrates solid backend engineering capabilities with relevant experience in required tech stack. Project shows good understanding of AI workflows and production-level code quality. Minor gaps in advanced error handling can be addressed through mentoring. Overall well-qualified for the Backend Engineer position."
  }
}
```

**Response (Failed):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "status": "failed",
  "error": "LLM call failed: connection timeout after 3 retries"
}
```

### Testing Flow Lengkap

**1. Test Upload → Evaluate → Result**
```bash
# Step 1: Upload files
POST /upload
-> Save cv_document_id and report_document_id

# Step 2: Start evaluation
POST /evaluate with saved IDs
-> Save job_id

# Step 3: Poll result (tunggu 2-5 menit)
GET /result/{job_id}
-> Check status until "completed"
```

**2. Test Validation**
```bash
# Invalid CV ID
POST /evaluate
{
  "job_title": "Backend Engineer",
  "cv_id": "invalid-id",
  "report_id": "660e8400-e29b-41d4-a716-446655440001"
}
-> Expected: 404 Not Found

# Missing required field
POST /evaluate
{
  "cv_id": "550e8400-e29b-41d4-a716-446655440000",
  "report_id": "660e8400-e29b-41d4-a716-446655440001"
}
-> Expected: 400 Bad Request (job_title required)
```

***

## 🐛 Cara Memperbaiki Issues

### Issue 1: Connection Refused - Port Already in Use

**Error:**
```
listen tcp :8080: bind: address already in use
```

**Diagnosis:**
```bash
# Check port yang digunakan
netstat -ano | findstr :8080
```

**Solution:**
```bash
# Option 1: Kill process yang menggunakan port
taskkill /PID <PID_NUMBER> /F

# Option 2: Ganti port di .env
SERVER_PORT=8081
```

### Issue 2: Ollama Not Responding

**Error:**
```
Worker 1 failed to process job: LLM call failed: Post "http://localhost:11434/api/generate": dial tcp [::1]:11434: connect: connection refused
```

**Diagnosis:**
```bash
# Check Ollama service
ollama list

# Test Ollama
curl http://localhost:11434/api/generate -d '{"model":"gemma3:4b","prompt":"test"}'
```

**Solution:**
```bash
# Restart Ollama service
# Di terminal baru
ollama serve

# Atau restart Windows service
net stop ollama
net start ollama
```

### Issue 3: PDF Extraction Failed

**Error:**
```
failed to extract CV text: failed to open PDF file: open storage/uploads/file.pdf: The system cannot find the file specified
```

**Diagnosis:**
```bash
# Check file exists
dir storage\uploads\

# Check file permissions
icacls storage\uploads\
```

**Solution:**
```bash
# Ensure directory exists
mkdir storage\uploads

# Check file path di database
mysql -u root -p cv_ai_evaluator
SELECT id, file_path FROM uploaded_documents WHERE id = 'xxx';

# Verify file di filesystem matches database path
```

### Issue 4: Worker Timeout

**Error:**
```
Worker 1 failed to process job: CV evaluation failed: LLM call failed: context deadline exceeded
```

**Diagnosis:**
- LLM processing terlalu lama (model besar, hardware lambat)

**Solution:**
```go
// Di pkg/llm/ollama_client.go
// Increase timeout
Client: &http.Client{
    Timeout: 600 * time.Second, // Dari 300s ke 600s (10 menit)
}
```

### Issue 5: ChromaDB Query Returns No Results

**Error:**
```
Warning: failed to get job description context: no relevant documents found for type: job_description
```

**Diagnosis:**
```bash
# Check if ground truth documents exist
mysql -u root -p cv_ai_evaluator
SELECT * FROM ground_truth_documents;
```

**Solution:**
```bash
# Re-ingest ground truth documents
go run scripts/ingest_groundtruth.go

# Verify ingestion
# Check chroma_data/ folder created
dir chroma_data

# Check database records
mysql -u root -p cv_ai_evaluator
SELECT COUNT(*) FROM ground_truth_documents;
```

### Issue 6: JSON Parse Error

**Error:**
```
failed to parse CV response: could not parse CV response: (raw LLM output)
```

**Diagnosis:**
- LLM tidak return valid JSON
- Model kecil (gemma3:4b) kadang tidak konsisten

**Solution:**
```go
// Worker sudah punya fallback regex parsing
// Tapi bisa improve prompt:

// Di evaluation_worker.go, tambahkan di prompt:
CRITICAL: You MUST respond ONLY with valid JSON. 
No explanation text before or after the JSON.
Example format:
{"match_rate": 0.85, "feedback": "candidate shows..."}
```

### Issue 7: Memory Leak di Worker

**Symptoms:**
- RAM usage terus naik
- Worker tidak release memory

**Diagnosis:**
```bash
# Monitor dengan Task Manager atau:
# Install pprof
go get github.com/pkg/profile
```

**Solution:**
```go
// Di evaluation_worker.go
// Ensure proper cleanup di processJob:

defer func() {
    // Cleanup variables
    cvText = ""
    reportText = ""
    runtime.GC() // Force garbage collection
}()
```

***

## 🗄️ Cara Memperbaiki Database Issues

### Issue 1: Can't Connect to MySQL

**Error:**
```
Failed to connect to database: dial tcp 127.0.0.1:3306: connectex: No connection could be made because the target machine actively refused it.
```

**Diagnosis:**
```bash
# Check MySQL service
net start | findstr -i mysql

# Try connect manually
mysql -u root -p
```

**Solution:**
```bash
# Start MySQL service
net start MySQL80

# Or restart
net stop MySQL80
net start MySQL80

# Check port
netstat -ano | findstr :3306
```

### Issue 2: Access Denied for User

**Error:**
```
Error 1045: Access denied for user 'root'@'localhost' (using password: YES)
```

**Diagnosis:**
```bash
# Test credentials
mysql -u root -p
Enter password: [your_password]
```

**Solution:**
```sql
-- Reset password jika lupa
-- Di MySQL command line (sebagai admin):
ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;

-- Update .env dengan password baru
DB_PASSWORD=new_password
```

### Issue 3: Database Not Found

**Error:**
```
Error 1049: Unknown database 'cv_ai_evaluator'
```

**Solution:**
```sql
-- Create database
CREATE DATABASE cv_ai_evaluator CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Verify
SHOW DATABASES;
```

### Issue 4: Table Doesn't Exist

**Error:**
```
Error 1146: Table 'cv_ai_evaluator.uploaded_documents' doesn't exist
```

**Solution:**
```sql
-- Connect to database
USE cv_ai_evaluator;

-- Check tables
SHOW TABLES;

-- If empty, run all CREATE TABLE statements dari STEP 3

-- Verify
DESC uploaded_documents;
DESC evaluation_jobs;
DESC ground_truth_documents;
```

### Issue 5: Foreign Key Constraint Failed

**Error:**
```
Error 1452: Cannot add or update a child row: a foreign key constraint fails
```

**Diagnosis:**
```sql
-- Check foreign key constraints
SELECT 
    CONSTRAINT_NAME, 
    TABLE_NAME, 
    REFERENCED_TABLE_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = 'cv_ai_evaluator'
AND REFERENCED_TABLE_NAME IS NOT NULL;

-- Check if referenced record exists
SELECT * FROM uploaded_documents WHERE id = 'xxx';
```

**Solution:**
```sql
-- Option 1: Ensure parent record exists
INSERT INTO uploaded_documents VALUES (...);

-- Option 2: Temporarily disable FK checks (dev only!)
SET FOREIGN_KEY_CHECKS=0;
-- Your operation
SET FOREIGN_KEY_CHECKS=1;

-- Option 3: Delete orphaned records
DELETE FROM evaluation_jobs 
WHERE cv_document_id NOT IN (SELECT id FROM uploaded_documents);
```

### Issue 6: Duplicate Entry

**Error:**
```
Error 1062: Duplicate entry 'xxx' for key 'PRIMARY'
```

**Solution:**
```sql
-- Check existing record
SELECT * FROM uploaded_documents WHERE id = 'xxx';

-- Option 1: Use different ID (UUID harus unique)
-- Pastikan UUID generator working correctly

-- Option 2: Update instead of insert
INSERT INTO uploaded_documents (...) 
ON DUPLICATE KEY UPDATE 
file_path = VALUES(file_path),
original_filename = VALUES(original_filename);
```

### Issue 7: Character Encoding Issues

**Error:**
```
Incorrect string value: '\xF0\x9F\x98\x80...' for column 'cv_feedback'
```

**Solution:**
```sql
-- Check database charset
SHOW CREATE DATABASE cv_ai_evaluator;

-- If not utf8mb4, recreate:
DROP DATABASE cv_ai_evaluator;
CREATE DATABASE cv_ai_evaluator 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- Ensure tables use utf8mb4
ALTER TABLE uploaded_documents 
CONVERT TO CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;
```

### Issue 8: Connection Timeout

**Error:**
```
Error 2013: Lost connection to MySQL server during query
```

**Solution:**
```sql
-- Increase timeout di MySQL config
-- Edit my.ini (Windows) or my.cnf (Linux)
[mysqld]
max_allowed_packet=256M
wait_timeout=28800
interactive_timeout=28800

-- Restart MySQL
net stop MySQL80
net start MySQL80
```

### Issue 9: Too Many Connections

**Error:**
```
Error 1040: Too many connections
```

**Diagnosis:**
```sql
-- Check current connections
SHOW PROCESSLIST;

-- Check max connections
SHOW VARIABLES LIKE 'max_connections';
```

**Solution:**
```sql
-- Increase max connections
SET GLOBAL max_connections = 200;

-- Or in my.ini
[mysqld]
max_connections=200

-- Also check connection pooling di internal/database/mysql.go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
```

### Issue 10: Slow Queries

**Symptoms:**
- API response lambat
- Worker processing lama

**Diagnosis:**
```sql
-- Enable slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;

-- Check slow queries
SHOW GLOBAL STATUS LIKE 'Slow_queries';

-- Analyze query
EXPLAIN SELECT * FROM evaluation_jobs WHERE status = 'queued';
```

**Solution:**
```sql
-- Add indexes untuk frequently queried columns
CREATE INDEX idx_status ON evaluation_jobs(status);
CREATE INDEX idx_document_type ON uploaded_documents(document_type);
CREATE INDEX idx_created_at ON evaluation_jobs(created_at);

-- Check indexes
SHOW INDEX FROM evaluation_jobs;

-- Analyze table
ANALYZE TABLE evaluation_jobs;
```

***

## 🔍 Database Maintenance Commands

### Backup Database
```bash
# Backup semua data
mysqldump -u root -p cv_ai_evaluator > backup_2025_11_17.sql

# Backup struktur saja
mysqldump -u root -p --no-data cv_ai_evaluator > structure.sql

# Backup data saja
mysqldump -u root -p --no-create-info cv_ai_evaluator > data.sql
```

### Restore Database
```bash
# Restore dari backup
mysql -u root -p cv_ai_evaluator < backup_2025_11_17.sql
```

### Clean Up Old Data
```sql
-- Delete completed jobs older than 30 days
DELETE FROM evaluation_jobs 
WHERE status = 'completed' 
AND completed_at < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- Delete orphaned documents
DELETE FROM uploaded_documents 
WHERE id NOT IN (
    SELECT cv_document_id FROM evaluation_jobs
    UNION
    SELECT report_document_id FROM evaluation_jobs
);
```


