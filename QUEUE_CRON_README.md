# Redis Queue & Cron Job System

Sistem Queue dan Cron Job yang terintegrasi dengan Redis untuk Golang Starter Kit 2025.

## Fitur Utama

### Queue System
- **Job Processing**: Proses job secara asynchronous menggunakan Redis
- **Worker Pool**: Multiple workers untuk memproses job secara concurrent
- **Job Retry**: Automatic retry untuk job yang gagal
- **Delayed Jobs**: Schedule job untuk dijalankan di masa depan
- **Job Registry**: Registrasi job types secara modular

### Cron Job System
- **Scheduled Jobs**: Job yang berjalan berdasarkan schedule tertentu
- **Cron Expressions**: Support untuk basic cron expressions
- **Interval-based**: Support untuk interval-based scheduling
- **Job Management**: Start/stop job secara individual

## Setup & Konfigurasi

### 1. Environment Variables

Tambahkan konfigurasi Redis ke file `.env`:

```bash
# Redis Configuration
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Queue Configuration
QUEUE_WORKERS=3
QUEUE_DEFAULT=default

# Cron Configuration
CRON_ENABLED=true
```

### 2. Docker Setup

Jalankan Redis menggunakan Docker Compose:

```bash
docker-compose up redis -d
```

### 3. Install Dependencies

```bash
go mod tidy
```

## Penggunaan

### CLI Commands

#### Queue Workers

```bash
# Jalankan queue workers
go run main.go queue:work

# Jalankan dengan konfigurasi custom
go run main.go queue:work --workers=5 --queue=emails

# Lihat statistik queue
go run main.go queue:stats

# Kirim test jobs
go run main.go queue:test --count=10
```

#### Cron Jobs

```bash
# Jalankan cron scheduler
go run main.go cron:run

# Lihat status cron jobs
go run main.go cron:status
```

### API Endpoints

#### Queue Management

```http
# Get queue statistics
GET /api/queue/stats

# Dispatch email job
POST /api/queue/email
{
  "to": "user@example.com",
  "subject": "Test Email",
  "body": "This is a test email"
}

# Dispatch file cleanup job
POST /api/queue/file-cleanup
{
  "directory": "./tmp",
  "older_than_hours": 24
}

# Send test jobs
POST /api/queue/test?count=5
```

#### Cron Management

```http
# Get cron job status
GET /api/cron/status
```

#### Health Checks

```http
# Check queue health
GET /api/health/queue

# Check cron health
GET /api/health/cron
```

## Membuat Job Baru

### 1. Queue Job

Buat file baru di `app/jobs/`:

```go
package jobs

import (
    "golang_starter_kit_2025/app/queue"
    "log"
)

type MyCustomJob struct {
    queue.BaseJob
}

func NewMyCustomJob() queue.JobInterface {
    return &MyCustomJob{
        BaseJob: queue.BaseJob{
            Name: "my_custom_job",
        },
    }
}

func (j *MyCustomJob) Handle() error {
    data := j.GetData()
    
    // Process your job here
    log.Printf("Processing custom job with data: %+v", data)
    
    return nil
}
```

Daftarkan job di `app/services/queue_cron_service.go`:

```go
func (qs *QueueService) registerJobs() {
    // ... existing jobs
    queue.RegisterJob("my_custom_job", jobs.NewMyCustomJob)
}
```

### 2. Cron Job

Buat cron job baru:

```go
package jobs

import (
    "golang_starter_kit_2025/app/cron"
    "log"
)

type MyCustomCronJob struct {
    cron.BaseCronJob
}

func NewMyCustomCronJob() cron.CronJobInterface {
    return &MyCustomCronJob{
        BaseCronJob: cron.BaseCronJob{
            Name:     "my_custom_cron",
            Schedule: "0 */6 * * *", // Every 6 hours
        },
    }
}

func (j *MyCustomCronJob) Execute() error {
    log.Println("Executing custom cron job")
    
    // Your cron job logic here
    
    return nil
}
```

Daftarkan cron job:

```go
func (cs *CronService) registerCronJobs() {
    // ... existing jobs
    customJob := jobs.NewMyCustomCronJob()
    cs.Manager.AddJob(customJob, cron.Every6Hours)
}
```

## Penggunaan dalam Code

### Dispatch Queue Job

```go
import "golang_starter_kit_2025/app/services"

// Buat queue service
queueService := services.NewQueueService(3)

// Dispatch email job
err := queueService.DispatchEmailJob(
    "user@example.com",
    "Welcome!",
    "Welcome to our platform",
)

// Dispatch custom job
job := jobs.NewMyCustomJob()
job.SetData(map[string]interface{}{
    "user_id": 123,
    "action":  "process_data",
})
err := queueService.Queue.Push(job)
```

### Delayed Jobs

```go
import "time"

// Schedule job untuk 1 jam ke depan
delay := 1 * time.Hour
err := queueService.Queue.Push(job, delay)
```

## Monitoring & Debugging

### Log Output

Queue dan Cron jobs menghasilkan log yang informatif:

```
🚀 Starting worker pool with 3 workers for queue: default
👷 Worker 1 started
🔄 Worker 1 processing job: email_job (ID: abc123)
✅ Worker 1 job email_job completed successfully (took 150ms)
```

### Error Handling

- Job yang gagal akan di-retry sesuai konfigurasi `MaxTries`
- Job yang melebihi retry limit akan dipindahkan ke `failed` queue
- Error logs akan menampilkan detail kesalahan

### Statistics

Gunakan API endpoint atau CLI command untuk melihat statistik:

```bash
go run main.go queue:stats
# Output:
# 📊 Queue Statistics for 'default':
#    Size: 5 jobs
#    Workers: 3
```

## Best Practices

1. **Job Idempotency**: Pastikan job dapat dijalankan berulang kali tanpa efek samping
2. **Error Handling**: Selalu handle error dengan baik dan berikan log yang informatif
3. **Data Validation**: Validasi data input sebelum memproses job
4. **Resource Management**: Tutup koneksi database/file setelah selesai
5. **Monitoring**: Monitor queue size dan worker performance secara regular

## Troubleshooting

### Redis Connection Issues

```bash
# Test Redis connection
redis-cli ping
# Should return: PONG
```

### Queue Not Processing

1. Pastikan Redis berjalan
2. Periksa konfigurasi REDIS_HOST
3. Pastikan workers sedang berjalan
4. Periksa log untuk error messages

### Cron Jobs Not Running

1. Pastikan cron service dijalankan dengan `cron:run`
2. Periksa cron expressions menggunakan online cron calculator
3. Periksa log untuk error messages

## Contoh Production Setup

```bash
# Terminal 1: Jalankan queue workers
go run main.go queue:work --workers=5

# Terminal 2: Jalankan cron scheduler
go run main.go cron:run

# Terminal 3: Jalankan web server
go run main.go
```

Untuk production, gunakan process manager seperti systemd atau supervisor untuk menjalankan services ini secara otomatis.
