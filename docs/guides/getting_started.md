# Getting Started with TGo Framework

Panduan ini membantu developer baru memulai proyek menggunakan framework **TGo** dari nol hingga siap production.

---

## 1. Persyaratan Sistem
- **Go 1.25** atau yang lebih baru
- **Git**
- Database: **SQLite** (bawaan development) atau **PostgreSQL** / **MySQL**

---

## 2. Instalasi Craft CLI

Craft adalah command-line tool resmi untuk TGo (seperti *Artisan* pada Laravel atau *Cargo* pada Rust).

```bash
go install github.com/tokalink/tgo/cmd/craft@latest
```

> **Tips Windows**: Pastikan folder `$env:USERPROFILE\go\bin` sudah terdaftar di environment variable `PATH` Anda.

---

## 3. Membuat Proyek Baru (`craft new`)

Gunakan perintah `craft new` untuk membuat scaffolding proyek lengkap dengan **Git repository yang 100% baru dan bersih**:

### A. Mode Standard (Single-Tenant / Monolith)
Cocok untuk aplikasi web internal, standalone tools, desktop apps, atau microservices:
```bash
craft new my-app
```

### B. Mode Multi-Tenant SaaS
Cocok untuk platform SaaS B2B dengan isolasi database schema per tenant:
```bash
craft new my-saas --mode=multi
```

---

## 4. Struktur Proyek yang Dihasilkan

Setelah scaffolding selesai, struktur folder proyek Anda adalah sebagai berikut:

```text
my-app/
├── app/
│   ├── handlers/       # HTTP & ConnectRPC Controller Handlers
│   ├── models/         # Entity Structs & Repository Database Queries
│   └── services/       # Domain Business Logic Actions
├── config/
│   └── app.yaml        # Konfigurasi aplikasi & database
├── database/
│   └── migrations/     # File migrasi SQL (up.sql & down.sql)
├── frontend/           # Starter UI (HTML, CSS, JS)
├── proto/              # Protobuf Service Definitions
├── .air.toml           # Konfigurasi Live Reload
├── .env                # Local Environment Variables
├── go.mod              # Module Go mandiri
└── main.go             # Application Entrypoint
```

---

## 5. Menjalankan Aplikasi

Masuk ke folder proyek dan jalankan server:

```bash
cd my-app
go run main.go
# atau menggunakan Craft:
craft serve -p 8080
```

Buka browser di **`http://localhost:8080`**.

---

## 6. Scaffolding dengan Craft CLI

TGo menyediakan generator praktis untuk mempercepat pembuatan kode:

### A. Membuat Model & Migrasi Database
```bash
# Buat Model User dan otomatis buat file migrasi SQL-nya
craft make:model User -m
```

### B. Membuat Business Logic Action
```bash
# Buat service RegisterUserAction
craft make:service RegisterUser
```

### C. Membuat Controller Handler
```bash
# Buat UserHandler
craft make:handler UserHandler
```

### D. Menjalankan Migrasi Database
```bash
# Mode Standard:
craft migrate

# Mode Multi-Tenant (migrasi seluruh schema tenant):
craft migrate:tenant --tenant=all
```

---

## 7. Mengakses Transport Layer

Setiap handler TGo secara otomatis mendukung 3 moda transport:

### A. Melalui REST HTTP (JSON)
```bash
curl -X POST http://localhost:8080/user.v1.UserService/GetProfile \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Slug: acme_corp" \
  -d '{"user_id": "usr_99"}'
```

### B. Melalui ConnectRPC / gRPC
Klien frontend (TypeScript/React/Vue) atau backend microservices dapat memanggil RPC via `@connectrpc/connect-web` atau client library Go native.

### C. Melalui In-Memory IPC (Wails Desktop / Mobile)
Untuk aplikasi desktop lokal tanpa network socket overhead (0.006 ms latency):
```go
invoker := application.Server().ServeInMemory(ctx)
req := &userv1.GetProfileRequest{UserId: "101"}
var resp userv1.GetProfileResponse
err := invoker.Invoke(ctx, "/user.v1.UserService/GetProfile", req, &resp)
```

---

## 8. Mencoba Demo Interaktif

Jika Anda ingin melihat implementasi lengkap aplikasi TGo dengan antarmuka Login $\rightarrow$ Dashboard (dengan Sidebar, Tenant Switcher, dan Live ConnectRPC Tester):

```bash
git checkout demo
go run starter/main.go
```
Buka browser di [http://localhost:8080](http://localhost:8080).
