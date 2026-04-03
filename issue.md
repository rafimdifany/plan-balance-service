# Project Setup Plan

## Objective
Membangun fondasi awal project backend menggunakan Golang dengan stack teknologi yang telah ditentukan. Dokumen ini berisi instruksi high-level untuk melakukan setup project.

## Tech Stack & Dependencies
- **Bahasa**: Golang
- **HTTP Router**: `github.com/gin-gonic/gin`
- **Database Driver**: `github.com/jackc/pgx/v5` (gunakan `pgxpool` untuk connection pooling)
- **Query/ORM**: Tanpa ORM, cukup manual SQL query menggunakan `pgx`
- **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Password Hashing**: `golang.org/x/crypto/bcrypt`
- **Google Auth Verify**: `google.golang.org/api/idtoken` (untuk memverifikasi token login dari Google)
- **DB Migration**: `github.com/golang-migrate/migrate/v4`
- **Logging**: `go.uber.org/zap`
- **Validator**: `github.com/go-playground/validator/v10` (integrasi dengan Gin)
- **CORS Middleware**: `github.com/gin-contrib/cors`

## Phase 1: Project Initialization
1. Inisialisasi Go module: Jalankan `go mod init <nama-module>` di root folder.
2. Download Dependencies: Lakukan `go get` untuk semua package yang telah disebutkan di bagian Tech Stack.
3. Struktur Folder: Buat struktur project Go standar (contoh: `cmd/api`, `internal`, `pkg`, `migrations`, `configs`).

## Phase 2: Configuration & Database Setup
1. **Environment Variables**: Buat sistem untuk membaca konfigurasi dari `.env` file atau environment OS. Konfigurasi yang dibutuhkan minimal: PORT, Database URL, JWT Secret, dan Google Client ID.
2. **Database Connection**: Buat koneksi ke PostgreSQL menggunakan `pgxpool`. Pastikan koneksi dikelola dengan baik dan mendukung connection pooling.
3. **DB Migration Utility**: Buat script (bisa menggunakan Makefile atau root CLI) yang memanfaatkan command line `golang-migrate` untuk menjalankan file SQL yang ada di dalam folder `migrations`.
4. **Logger Config**: Inisialisasi konfigurasi `zap` logger (json format untuk production, console format untuk environment local/testing).

## Phase 3: Core Application Components
1. **Security & Utilities**:
   - Buat fungsi helper untuk mengenkripsi password dan memvalidasi hash menggunakan `bcrypt`.
   - Buat fungsi helper/Service untuk meng-generate dan memverifikasi JWT token.
   - Buat fungsi helper/Service yang diintegrasikan dengan `idtoken` untuk memvalidasi token login yang didapat dari sisi client/Google SSO.
2. **Router & Middleware**:
   - Inisialisasi router `gin`.
   - Pasang `cors` middleware sesuai kebutuhan.
   - Buat custom middleware logging menggunakan `zap` untuk mencatat setiap request HTTP yang masuk ke server.
3. **Validator**: Daftarkan engine `validator/v10` ke dalam Gin (umumnya Gin sudah memiliki binding bawaan untuk ini, cukup kustomisasi terjemahan error bila perlu).

## Phase 4: Server Entry Point (`main.go`)
1. Buat file `main.go` di `cmd/api/main.go`.
2. Lakukan wiring/dependency injection sederhana (baca config -> init logger -> connect DB -> init router).
3. Buat rute `GET /ping` sederhana untuk memastikan web server berjalan normal.
4. Jalankan server HTTP, dan jangan lupa implementasikan **Graceful Shutdown** agar server meng-close koneksi DB dan network dengan aman ketika menerima signal stop / interupsi dari OS.
