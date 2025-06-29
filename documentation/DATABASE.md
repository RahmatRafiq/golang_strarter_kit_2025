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
go run main.go migrate --file <nama_file_migration>
```
Menjalankan bagian `UP` dari `<nama_file_migration>.sql`.

### 3. Menjalankan Semua Migrasi yang Tertunda
```bash
go run main.go migrate:all
```
- Membuat batch baru.
- Menjalankan bagian `UP` dari semua file `.sql` yang belum tercatat di tabel `migrations`.
- Mencatat setiap file ke batch tersebut.

### 4. Rollback Satu File Migrasi
```bash
go run main.go rollback --file=<nama_file_migration>
```
Menjalankan bagian `DOWN` dari `<nama_file_migration>.sql` (tanpa mengubah batch).

### 5. Rollback Semua Batch
```bash
go run main.go rollback:all
```
- Loop dari batch tertinggi → 1.
- Menjalankan bagian `DOWN` dari semua file `.sql` per batch.
- Menghapus seluruh catatan di tabel `migrations`.

### 6. Rollback Batch Tertentu
```bash
go run main.go rollback:batch --batch=<nomor_batch>
```
Meng-rollback hanya migrasi di batch `<nomor_batch>`, lalu menghapus catatannya.

### 7. Rollback Batch Terakhir (Default)
```bash
go run main.go rollback:batch
```
Jika flag `--batch` tidak diset, akan otomatis meng-rollback batch terakhir.

---

📌 **Catatan**:
- Tabel `migrations` akan otomatis dibuat saat pertama kali menjalankan `migrate:all` atau `rollback:batch`.
- Pastikan setiap file `.sql` memiliki bagian `UP` dan `DOWN` yang jelas sebelum menjalankan migrate/rollback.
- Contoh format file migrasi:
    ```sql
    -- UP
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    );

    -- DOWN
    DROP TABLE users;
    ```

## Perintah CLI untuk Seeder

### 1. Membuat File Seeder Baru
```bash
go run main.go make:seeder --name=<nama_seeder>
```
Membuat file seeder baru di direktori `app/database/seeds/` dengan format nama file:
- `YYYYMMDDHHMMSS_<nama_seeder>.go`

### 2. Menjalankan Semua Seeder
```bash
go run main.go db:seed
```
Menjalankan semua seeder yang ada di direktori `app/database/seeds/`.

### 3. Rollback Batch Seeder Terakhir (Default)
```bash
go run main.go rollback:seeder
```
Menghapus data yang dimasukkan oleh batch seeder terakhir.

### 4. Rollback Batch Seeder Tertentu
```bash
go run main.go rollback:seeder --batch=<nomor_batch>
```
Menghapus data yang dimasukkan oleh batch seeder dengan nomor `<nomor_batch>`.

📌 **Catatan**:
- Seeder file yang dibuat akan memiliki template dasar untuk mempermudah implementasi.
- Pastikan untuk menyesuaikan isi file seeder dengan kebutuhan data aplikasi Anda.
- Gunakan perintah rollback untuk menghapus data yang tidak diperlukan atau untuk pengujian ulang.
