# Project Restructuring Plan

## Objective
Melakukan restrukturisasi folder dan file pada project Golang agar sesuai dengan arsitektur standar yang lebih rapi (berbasis Domain-Driven Design / Clean Architecture).

## Panduan Restrukturisasi

Tolong ubah struktur dan pindahkan file-file yang ada pada project ini agar sesuai dengan struktur direktori di bawah ini:

```text
.
├── cmd/
│   └── api/
│       └── main.go          # entry point
│
├── internal/
│   ├── config/              # load env/config
│   ├── handler/             # HTTP handler (controller)
│   ├── middleware/          # auth, logging, rate limit
│   ├── service/             # business logic
│   ├── repository/          # DB access (pgx/sql)
│   ├── model/               # struct entity (User, etc)
│   └── dto/                 # request/response struct
│
├── pkg/                     # optional (helper reusable)
│   ├── logger/              # Pindahkan logger dari internal/logger ke sini
│   └── utils/               # Pindahkan helper security (jwt, bcrypt, google auth) ke sini
│
├── db/
│   ├── migrations/          # Pindahkan folder migrations yang ada di root ke dalam db/migrations
│   └── seed/                # optional dummy data
│
├── scripts/                 # bash scripts (deploy, migrate)
│
├── .env                     # (sudah ada)
├── .env.example             # (sudah ada)
├── go.mod                   # (sudah ada)
├── Makefile                 # (sudah ada)
└── README.md                # Buatkan file README sederhana jika belum ada
```

## Detail Task / Langkah Kerja (High-Level)

1. **Pembuatan Folder Baru**:
   Buat folder-folder yang belum ada di dalam root project:
   - `internal/handler/`
   - `internal/service/`
   - `internal/repository/`
   - `internal/model/`
   - `internal/dto/`
   - `pkg/logger/`
   - `pkg/utils/`
   - `db/migrations/`
   - `db/seed/`
   - `scripts/`

2. **Pemindahan File (Migration & Packaging)**:
   - Pindahkan folder `migrations/` yang ada di luar (root) ke dalam `db/`. Jangan lupa perbarui path migrasi pada file seperti `Makefile` jika diperlukan.
   - Pindahkan package `logger` dari `internal/logger` ke `pkg/logger`. Lakukan update `import path` di `main.go` dan `middleware`.
   - Pindahkan helper package yang bersifat utility (seperti fungsi JWT, password hash bcrypt, dan Verifikasi Google di `internal/security`) ke dalam `pkg/utils` atau folder relevan yang sesuai (misal: `pkg/security`).
   - Pindahkan `internal/validator` ke folder yang sesuai (misal `pkg/validator` atau bagian dari middleware/utils).
   - Pastikan dependencies yang memanggil file hasil pindahan di-update *import* path-nya (yaitu di `cmd/api/main.go`).

3. **Verifikasi**:
   - Jalankan `go mod tidy` untuk memastikan tidak ada import yang rusak.
   - Jalankan `go build -o tmp/api cmd/api/main.go` untuk memastikan kode dapat dikompilasi dengan lancar setelah restrukturisasi.

Pastikan proses refactoring ini tidak mengubah business logic apapun dan hanya mengubah layout direktori serta package imports.
