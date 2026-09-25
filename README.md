# ⚡ TGo Framework

> **Next-Gen Go Application Framework** for Cloud SaaS, Unified Wails Desktop/Mobile Clients, and High-Performance Microservices.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![ConnectRPC](https://img.shields.io/badge/Protocol-ConnectRPC-6366f1.svg)](https://connectrpc.com)
[![Multi-Tenancy](https://img.shields.io/badge/Multi--Tenancy-Kernel--Level-10b981.svg)](#-native-multi-tenancy--dual-architecture)

---

## 🌟 Fitur Utama

- **🚀 Triple Mode Transport (ConnectRPC)**: Satu handler Go melayani 3 moda sekaligus:
  1. **gRPC (HTTP/2 Binary Protobuf)** — untuk komunikasi microservices berkecepatan tinggi.
  2. **REST HTTP/1.1 & HTTP/2 (JSON)** — kompatibel langsung dengan browser `fetch()`, cURL, dan Swagger.
  3. **In-Memory IPC Bridge (0.006 ms)** — nol network hop / zero socket overhead untuk desktop app (Wails v3) dan mobile.
- **🏢 Dual Architecture (Multi-Tenant SaaS & Standard Monolith)**:
  - **Multi-Tenant Mode**: Resolusi tenant dinamis (Subdomain, Header `X-Tenant-Slug`, JWT Claim), isolasi database tingkat kernel (**PostgreSQL Schema-per-tenant** via `search_path` & **SQLite File-per-tenant**).
  - **Non-Tenancy (Standard Mode)**: Aplikasi monolit / single-tenant standar tanpa overhead routing tenant.
- **💾 Universal Database Engine**:
  - Dukungan agnostik untuk **PostgreSQL**, **SQLite**, dan **MySQL**.
  - Context-aware database router (`db.FromContext(ctx)`).
- **🛠️ Developer Ergonomics ("Craft" CLI)**:
  - Scaffolding proyek baru: `craft new my-app`.
  - Generator code: `craft make:model`, `craft make:service`, `craft make:handler`, `craft make:migration`, `craft make:frontend`.
  - Runner & Migrator: `craft serve`, `craft migrate`, `craft tenant:create`.

---

## 📦 Instalasi Global Craft CLI

Install binary `craft` secara global ke sistem Anda:

```bash
# Pastikan Go bin ada di PATH Anda ($env:USERPROFILE\go\bin di Windows)
go install github.com/tokalink/tgo/cmd/craft@latest
```

---

## 🚀 Membuat Proyek Web Baru (< 1 Menit)

Gunakan perintah `craft new` untuk membuat proyek baru dengan **Git repository yang 100% baru dan bersih**:

```bash
# Buat aplikasi standard (non-tenancy)
craft new my-web-app

# Atau buat aplikasi multi-tenant SaaS
craft new my-saas-app --mode=multi
```

Masuk ke folder proyek dan jalankan:
```bash
cd my-web-app
go run main.go
```
Buka browser di **`http://localhost:8080`**.

---

## 🛠️ Panduan Perintah Craft CLI

| Perintah | Deskripsi |
|---|---|
| `craft new <project-name>` | Scaffold proyek baru lengkap dengan `go.mod`, config, UI, dan `git init` fresh |
| `craft serve [-p 8080]` | Menjalankan server aplikasi TGo ConnectRPC & HTTP |
| `craft make:model <Name> [-m]` | Membuat struct model entity & repository (`-m` untuk sekaligus buat migrasi SQL) |
| `craft make:migration <Name>` | Membuat file migrasi SQL pair (`up.sql` dan `down.sql`) |
| `craft make:service <Name>` | Membuat action business logic domain service |
| `craft make:handler <Name>` | Membuat controller handler HTTP & RPC |
| `craft make:frontend` | Men-generate starter frontend modern (HTML/CSS/JS) |
| `craft migrate` | Menjalankan migrasi database ke schema aktif |
| `craft migrate:tenant --tenant=all` | Menjalankan migrasi ke seluruh schema tenant |
| `craft tenant:create <slug>` | Mem-provisioning schema database tenant baru secara otomatis |

---

## 💻 Contoh Penggunaan Transport Layer

### 1. In-Memory IPC (Untuk Wails Desktop / Mobile)
```go
// 0.006 ms per invocation tanpa network socket
invoker := app.Server().ServeInMemory(ctx)
req := &userv1.GetProfileRequest{UserId: "usr_101"}
var resp userv1.GetProfileResponse
err := invoker.Invoke(ctx, "/user.v1.UserService/GetProfile", req, &resp)
```

### 2. Standard REST HTTP (cURL / Browser)
```bash
curl -X POST http://localhost:8080/user.v1.UserService/GetProfile \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Slug: acme_corp" \
  -d '{"user_id": "usr_101"}'
```

---

## 🎨 Menjalankan Demo Interaktif (Login $\rightarrow$ Dashboard)

TGo menyediakan branch khusus **`demo`** yang berisi antarmuka lengkap (Login $\rightarrow$ Dashboard dengan Sidebar, Tenant Switcher, dan Live ConnectRPC Playground):

```bash
git checkout demo
go run starter/main.go
```
Buka [http://localhost:8080](http://localhost:8080) untuk mencoba demo interaktif.

---

## 📂 Struktur Repositori

```text
tgo/
├── cmd/
│   └── craft/                  # CLI Tool ("Artisan" for Go)
│       └── commands/           # new, serve, make, migrate, tenant
├── pkg/
│   ├── app/                    # Application Container & Lifecycle
│   ├── config/                 # Viper + .env Structured Config Loader
│   ├── database/               # Universal DB & Multi-Tenant Isolators
│   └── transport/
│       ├── connect/            # ConnectRPC Server & Interceptors
│       ├── inmemory/           # In-Memory IPC Bridge (Wails v3)
│       ├── middleware/         # Tenant, Auth, Logger, CORS, Recover
│       └── welcome/            # Embedded Static Frontend Server
├── starter/                    # Starter Core Template Application
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

Jalankan latency benchmark ConnectRPC & In-Memory IPC:
```bash
go test -bench="." -benchmem ./pkg/transport/connect
```

---

## 📜 Lisensi
MIT License © 2026 [Tokalink](https://github.com/tokalink).
