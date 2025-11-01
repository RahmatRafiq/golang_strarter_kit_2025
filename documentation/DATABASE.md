# Database Management

## Perintah CLI untuk Migrasi

### 1. Membuat File Migrasi Baru
```bash
go run main.go make:migration <prefix_nama_migrasi>
```
Membuat satu file kosong dengan format:
- `YYYYMMDDHHMMSS_<prefix_nama_migrasi>.sql`

📌 **Rekomendasi**:
- Gunakan prefix seperti `create_` atau `alter_` untuk mempermudah identifikasi jenis migrasi.
- Contoh:
    - `create_users_table`
    - `alter_products_table`

### 2. Menjalankan Satu File Migrasi
```bash
go run main.go migrate --file <nama_file_migration> [--connection=mysql]
```
Menjalankan bagian `UP` dari `<nama_file_migration>.sql`.

### 3. Menjalankan Semua Migrasi yang Tertunda
```bash
go run main.go migrate:all [--connection=mysql]
```
- Membuat batch baru.
- Menjalankan bagian `UP` dari semua file `.sql` yang belum tercatat di tabel `migrations`.
- Mencatat setiap file ke batch tersebut.

### 4. Melihat Status Migrasi ✨ NEW
```bash
go run main.go migrate:status [--connection=mysql]
```
Menampilkan tabel status semua migrasi:
- ✅ Migrasi yang sudah dijalankan (dengan nomor batch)
- ⏳ Migrasi yang belum dijalankan (pending)
- Ringkasan total migrasi

**Contoh Output:**
```
================================================================================
Migration                                          Batch      Status
--------------------------------------------------------------------------------
20250426184415_create_roles_table                  1          ✅ Ran
20250426184424_create_permissions_table            1          ✅ Ran
20250508221607_create_test                         -          ⏳ Pending
================================================================================
Total: 12 migrations (11 ran, 1 pending)
```

### 5. Rollback Satu File Migrasi
```bash
go run main.go rollback --file=<nama_file_migration> [--connection=mysql]
```
Menjalankan bagian `DOWN` dari `<nama_file_migration>.sql` (tanpa mengubah batch).

### 6. Rollback Batch Terakhir (Default)
```bash
go run main.go rollback:batch [--connection=mysql]
```
Jika flag `--batch` tidak diset, akan otomatis meng-rollback batch terakhir.

### 7. Rollback N Batch Terakhir ✨ NEW
```bash
go run main.go rollback:batch --step=3 [--connection=mysql]
```
Rollback N batch terakhir (seperti Laravel):
- `--step=1` - Rollback 1 batch terakhir
- `--step=3` - Rollback 3 batch terakhir
- Otomatis menghitung batch target

### 8. Rollback Batch Tertentu
```bash
go run main.go rollback:batch --batch=<nomor_batch> [--connection=mysql]
```
Meng-rollback hanya migrasi di batch `<nomor_batch>`, lalu menghapus catatannya.

### 9. Rollback Semua Batch
```bash
go run main.go rollback:all [--connection=mysql]
# Atau gunakan alias
go run main.go migrate:reset [--connection=mysql]  # ✨ NEW
```
- Loop dari batch tertinggi → 1.
- Menjalankan bagian `DOWN` dari semua file `.sql` per batch.
- Menghapus seluruh catatan di tabel `migrations`.

### 10. Fresh Migration (Reset & Re-run)
```bash
go run main.go migrate:fresh [--connection=mysql]
```
Rollback semua migrasi kemudian jalankan ulang semua migrasi.

### 11. Fresh Migration dengan Seeding ✨ NEW
```bash
go run main.go migrate:fresh --seed [--connection=mysql]
```
Rollback semua, migrate semua, lalu jalankan seeder:
- Rollback semua migrasi
- Jalankan ulang semua migrasi
- Otomatis menjalankan semua seeder

### 12. Drop Semua Tabel ✨ NEW
```bash
go run main.go db:wipe [--connection=mysql] [--force]
```
Menghapus SEMUA tabel dari database:
- ⚠️ **WARNING**: Perintah berbahaya! Membutuhkan konfirmasi "yes"
- `--force` - Skip konfirmasi (untuk CI/CD)
- Otomatis handle foreign key constraints
- Support MySQL dan PostgreSQL

---

## Perintah CLI untuk Seeder

### 1. Membuat File Seeder Baru
```bash
go run main.go make:seeder --name=<nama_seeder>
```
Membuat file seeder baru di direktori `app/database/seeds/` dengan format nama file:
- `YYYYMMDDHHMMSS_<nama_seeder>.go`

### 2. Menjalankan Semua Seeder
```bash
go run main.go db:seed [--connection=mysql]
```
Menjalankan semua seeder yang ada di direktori `app/database/seeds/`.

### 3. Menjalankan Seeder Spesifik ✨ NEW
```bash
go run main.go db:seed --class=UserSeeder [--connection=mysql]
```
Menjalankan satu seeder tertentu saja (seperti Laravel):
- Cek apakah seeder sudah pernah dijalankan
- Error jika seeder tidak ditemukan
- Support multi-database

### 4. Rollback Batch Seeder Terakhir (Default)
```bash
go run main.go rollback:seeder [--connection=mysql]
```
Menghapus data yang dimasukkan oleh batch seeder terakhir.

### 5. Rollback Batch Seeder Tertentu
```bash
go run main.go rollback:seeder --batch=<nomor_batch> [--connection=mysql]
```
Menghapus data yang dimasukkan oleh batch seeder dengan nomor `<nomor_batch>`.

---

## Multi-Database Support

Semua perintah migrasi dan seeder sekarang mendukung flag `--connection` untuk menentukan database connection yang ingin digunakan:

### Koneksi yang Tersedia:
- `mysql` (default)
- `postgres`
- `mysql_secondary`

### Contoh Penggunaan:
```bash
# Migrate di PostgreSQL
go run main.go migrate:all --connection=postgres

# Seed di database sekunder
go run main.go db:seed --connection=mysql_secondary

# Status migrasi di PostgreSQL
go run main.go migrate:status --connection=postgres

# Wipe PostgreSQL database
go run main.go db:wipe --connection=postgres --force
```

---

## Format File Migrasi

📌 **Catatan**:
- Tabel `migrations` akan otomatis dibuat saat pertama kali menjalankan `migrate:all` atau `rollback:batch`.
- Pastikan setiap file `.sql` memiliki bagian `UP` dan `DOWN` yang jelas sebelum menjalankan migrate/rollback.
- Contoh format file migrasi:

```sql
-- +++ UP Migration
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

### PostgreSQL Format:
```sql
-- +++ UP Migration
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

---

## Best Practices

### 1. Gunakan migrate:status Secara Berkala
Cek status migrasi sebelum deploy:
```bash
go run main.go migrate:status
```

### 2. Testing dengan Fresh Migration
Untuk testing lengkap dengan data seeder:
```bash
go run main.go migrate:fresh --seed
```

### 3. Rollback Bertahap
Gunakan `--step` untuk rollback bertahap:
```bash
# Rollback 1 batch terakhir saja
go run main.go rollback:batch --step=1

# Rollback 3 batch terakhir
go run main.go rollback:batch --step=3
```

### 4. CI/CD Pipeline
Untuk automation, gunakan `--force`:
```bash
go run main.go db:wipe --force --connection=test_db
go run main.go migrate:fresh --seed --connection=test_db
```

### 5. Development Workflow
```bash
# 1. Buat migration baru
go run main.go make:migration create_posts_table

# 2. Edit file migration yang dibuat

# 3. Jalankan migration
go run main.go migrate:all

# 4. Cek status
go run main.go migrate:status

# 5. Jika ada error, rollback
go run main.go rollback:batch

# 6. Fix migration file, lalu migrate lagi
go run main.go migrate:all
```

---

## Troubleshooting

### Migration Error
Jika migration gagal di tengah jalan:
```bash
# Cek status
go run main.go migrate:status

# Rollback yang error
go run main.go rollback:batch

# Perbaiki file SQL
# Jalankan ulang
go run main.go migrate:all
```

### Seeder Conflict
Jika seeder sudah pernah dijalankan:
```bash
# Cek apakah sudah ada di database
# Rollback jika perlu
go run main.go rollback:seeder

# Jalankan ulang
go run main.go db:seed
```

### Database Wipe Not Working
Jika ada foreign key issues:
```bash
# Gunakan db:wipe dengan force
go run main.go db:wipe --force

# Atau manual disable foreign keys (MySQL)
# SET FOREIGN_KEY_CHECKS = 0;
# DROP TABLE ...;
# SET FOREIGN_KEY_CHECKS = 1;
```
