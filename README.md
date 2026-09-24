# ⚡ TGo Framework

> **Next-Gen Go Application Framework** for Cloud SaaS, Unified Wails Desktop/Mobile Clients, and High-Performance Microservices.

---

## 🌟 Fitur Utama

- **🚀 Triple Mode Transport (ConnectRPC)**: Satu service secara otomatis melayani **gRPC (HTTP/2 binary Protobuf)**, **REST HTTP/1.1 (JSON)**, dan **In-Memory IPC Bridge** (0 network hop, ideal untuk Wails Desktop/Mobile).
- **🏢 Native Multi-Tenancy**:
  - Resolusi tenant otomatis melalui **Subdomain**, **HTTP Header** (`X-Tenant-Slug`), atau **JWT Claim**.
  - Isolasi database tingkat kernel (**PostgreSQL Schema-per-tenant** via `search_path` & **SQLite File-per-tenant**).
  - Provisioning tenant otomatis (`craft tenant:create <slug>`).
- **💾 Universal Database & Migrator**:
  - Dukungan agnostik untuk **SQLite**, **PostgreSQL**, dan **MySQL**.
  - Context-aware database connection (`db.FromContext(ctx)`).
- **🛠️ Developer Ergonomics ("Craft" CLI)**:
  - Scaffolding cepat: `craft make:model`, `craft make:service`, `craft make:handler`, `craft make:frontend`.
  - Migrasi skema: `craft migrate`.
  - Server runner: `craft serve`.

---

## 🚀 Quick Start (< 5 Menit)

### 1. Jalankan Aplikasi Starter
```bash
go run ./starter
```
Server akan aktif di port `:8080` dan melayani ConnectRPC + REST JSON.

### 2. Test Endpoint via cURL
```bash
curl -X POST http://localhost:8080/user.v1.UserService/GetProfile \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Slug: acme" \
  -d '{"user_id": "usr_123"}'
```

Output:
```json
{
  "userId": "usr_123",
  "name": "User usr_123",
  "email": "user_usr_123@acme.example.com",
  "role": "member",
  "tenantSlug": "acme"
}
```

### 3. CLI Craft Scaffolding
```bash
# Provision Tenant Baru
go run ./cmd/craft tenant:create my_company

# Generate Model
go run ./cmd/craft make:model Order

# Generate Service Action
go run ./cmd/craft make:service Order

# Generate Handler
go run ./cmd/craft make:handler Order

# Generate Frontend
go run ./cmd/craft make:frontend
```

---

## 📂 Struktur Repositori

```text
tgo/
├── cmd/
│   └── craft/                  # CLI Tool ("Artisan" for Go)
├── pkg/
│   ├── app/                    # Application Container & Lifecycle
│   ├── config/                 # Viper + .env Config Loader
│   ├── database/               # Universal DB & Multi-Tenant Isolators
│   └── transport/
│       ├── connect/            # ConnectRPC Server & Interceptors
│       ├── inmemory/           # In-Memory IPC Bridge (Wails v3)
│       └── middleware/         # Tenant, Auth, Logger, CORS, Recover
├── starter/                    # Starter Template Application
│   ├── app/handlers/
│   ├── proto/v1/
│   └── main.go
└── docs/
    ├── adr/                    # Architecture Decision Records
    └── guides/                 # Getting Started Guides
```

---

## 🧪 Testing & Benchmarks

Jalankan seluruh test suite:
```bash
go test -v ./...
```

Jalankan transport latency & throughput benchmark:
```bash
go test -bench="." -benchmem ./pkg/transport/connect
```

---

## 📜 Lisensi
MIT License © 2026 TGo Platform Team.
