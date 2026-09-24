# Getting Started with TGo Framework

Panduan ini membantu developer baru memulai proyek menggunakan framework **TGo** dalam waktu kurang dari 15 menit.

---

## 1. Persyaratan Sistem
- Go 1.22 atau yang lebih baru
- SQLite3 / PostgreSQL (opsional untuk production)

---

## 2. Menjalankan Server Starter

Jalankan template starter yang sudah disediakan di folder `starter/`:

```bash
go run ./starter/main.go
```

Server akan aktif pada `http://localhost:8080`.

---

## 3. Mengakses Transport Layer

### A. Melalui REST HTTP (JSON)
```bash
curl -X POST http://localhost:8080/user.v1.UserService/GetProfile \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Slug: acme" \
  -d '{"user_id": "usr_99"}'
```

### B. Melalui ConnectRPC / gRPC
Klien dapat langsung memanggil endpoint ini menggunakan standard Connect client library (Go, TypeScript, Swift, Kotlin, Dart).

### C. Melalui In-Memory IPC (Wails v3)
Untuk aplikasi desktop atau mobile lokal tanpa network overhead:
```go
invoker := app.Server().ServeInMemory(ctx)
req := &userv1.GetProfileRequest{UserId: "101"}
var resp userv1.GetProfileResponse
err := invoker.Invoke(ctx, "/user.v1.UserService/GetProfile", req, &resp)
```

---

## 4. Mengelola Multi-Tenancy

### A. Membuat Tenant Baru
```bash
go run ./cmd/craft tenant:create acme_corp
```
Perintah ini akan membuat database schema baru dan menjalankan migrasi awal.

### B. Resolusi Tenant di Endpoint
TGo otomatis mendeteksi tenant dari request melalui urutan berikut:
1. **Subdomain**: `acme.yourdomain.com`
2. **HTTP Header**: `X-Tenant-Slug: acme`
3. **JWT Token Claims**: payload `"tenant": "acme"`
4. **Fallback Default**: `public`

---

## 5. Scaffolding dengan Craft CLI

Gunakan perintah `craft` untuk mempercepat development:

```bash
# Buat Model & Repository
go run ./cmd/craft make:model Customer

# Buat Business Logic Service Action
go run ./cmd/craft make:service Customer

# Buat RPC & HTTP Handler
go run ./cmd/craft make:handler Customer

# Buat Starter UI Frontend
go run ./cmd/craft make:frontend
```
