# Product Requirement Document (PRD)
## Project: Next-Gen Go Application Framework — **TGo** (Codename: "AuraGo" / "OmniGo")

> **Dokumen Versi**: 2.0.0
> **Status**: Ready for Planning & Independent Repository Bootstrap
> **Terakhir Diperbarui**: 2026-09-23
> **Penulis**: Tim Core Platform
> **Target Penggunaan**: Framework fondasi untuk membangun aplikasi SaaS Cloud, Unified Frontend Desktop/Mobile (Wails v3), Mobile Backend, dan Microservices baru.

---

## Daftar Isi

1. [Executive Summary & Visi Produk](#1-executive-summary--visi-produk)
2. [Tujuan & Non-Tujuan](#2-tujuan--non-tujuan-goals--non-goals)
3. [Arsitektur Inti](#3-arsitektur-inti-core-pillars)
4. [Struktur Folder Standar](#4-struktur-folder-standar-framework-skeleton)
5. [Spesifikasi Komponen Teknis](#5-spesifikasi-komponen-teknis)
6. [Contoh Alur Eksekusi Kode](#6-contoh-alur-eksekusi-kode)
7. [Architecture Decision Records (ADR)](#7-architecture-decision-records-adr)
8. [Model Keamanan (Security Model)](#8-model-keamanan-security-model)
9. [Observabilitas & Monitoring](#9-observabilitas--monitoring)
10. [Strategi Pengujian (Testing Strategy)](#10-strategi-pengujian-testing-strategy)
11. [Rencana Tahapan Implementasi (Roadmap)](#11-rencana-tahapan-implementasi-roadmap)
12. [Success Metrics & KPI](#12-success-metrics--kpi)
13. [Risiko & Mitigasi](#13-risiko--mitigasi)
14. [Open Questions](#14-open-questions)
15. [Glosarium](#15-glosarium)
16. [Panduan Memulai](#16-panduan-memulai-di-repositori-baru)

---

## 1. Executive Summary & Visi Produk

### 1.1 Latar Belakang
Pengembangan aplikasi modern dengan Go sering kali terbentur dua ekstrem:
1. **Go Standard / Minimalist (Gin, Chi, Fiber)**: Terlalu *low-level*. Setiap pembuatan aplikasi baru menuntut developer merakit manual koneksi database, routing, multi-tenancy, konfigurasi, migrasi, dan autentikasi dari nol. Ini memakan waktu lama dan menghasilkan inkonsistensi antar proyek.
2. **Goravel / Laravel Ports**: Berusaha memindahkan pola PHP monolitik ke Go tanpa memanfaatkan kekuatan sejati Go (Goroutines, High-Throughput gRPC/ConnectRPC, dan In-Memory binary embedding).

**Gap yang ada**: Tidak ada framework Go yang menggabungkan ergonomis developer-friendly, performa RPC modern, dan multi-tenancy tingkat kernel dalam satu solusi yang siap produksi.

### 1.2 Visi Produk
Membangun framework Go generasi baru yang menggabungkan:
- **Developer Ergonomics ala Laravel**: Struktur folder bersih, CLI generator ("Craft"), middleware pipeline, migrasi skema, dan konfigurasi deklaratif.
- **Modern High-Performance RPC (ConnectRPC)**: Satu service otomatis melayani **gRPC (Protobuf)**, **REST HTTP/1.1 (JSON)**, dan **In-Memory IPC (Direct Go Call)**.
- **Native Multi-Tenancy**: Dukungan multi-tenant tingkat kernel (Schema-per-tenant di Postgres, DB-per-tenant di SQLite/MySQL) *out-of-the-box*.
- **Hybrid Deployment**: Dapat dijalankan sebagai **Cloud SaaS Daemon** ataupun di-*embed* langsung ke dalam **Desktop & Mobile Klien (Wails v3 Beta)** menggunakan UI Web yang seragam (HTML/CSS/JS) tanpa overhead jaringan.

### 1.3 Target Pengguna

| Persona | Kebutuhan Utama | Bagaimana TGo Menjawab |
|---|---|---|
| **Backend Developer Go** | Setup cepat, struktur konsisten | CLI scaffolding + starter template |
| **SaaS Developer** | Multi-tenancy bawaan, isolasi data | Tenant middleware + schema isolator |
| **Desktop & Mobile Developer** | 1 Codebase UI Web untuk multi-platform tanpa overhead jaringan | Wails In-Memory IPC Bridge (0ms latency) |
| **DevOps / Platform Eng.** | Observabilitas, tracing, health check | OpenTelemetry integration bawaan |

---

## 2. Tujuan & Non-Tujuan (Goals & Non-Goals)

### 2.1 Tujuan Utama (Goals)
1. **Unified Transport**: 1 implementasi logika bisnis dapat diakses via:
   - Native gRPC (HTTP/2 binary Protobuf) untuk microservices dan sinkronisasi cepat.
   - RESTful HTTP/JSON langsung dari web browser/mobile tanpa *grpc-web proxy*.
   - In-memory function call untuk desktop client (0ms network latency).
2. **SaaS Multi-Tenant Bawaan**:
   - Deteksi tenant otomatis melalui subdomain, HTTP header (`X-Tenant-ID`), atau JWT claim.
   - Isolasi data dinamis tanpa perlu menulis klausul `WHERE tenant_id` manual di setiap query.
   - Provisioning tenant baru otomatis (schema creation + seed data).
3. **Database Agnostik**:
   - Mendukung SQLite (offline/local), MySQL, dan PostgreSQL secara transparan.
   - Migrasi SQL kompatibel lintas dialek.
4. **Developer Experience (DX) Superior**:
   - Tool CLI (`craft`) untuk scaffolding kode: model, service, handler, migration, dan seeder.
   - Konfigurasi berbasis environment (`.env`) dan YAML/JSON terpadu.
   - Error message yang jelas, actionable, dan disertai tautan dokumentasi.
5. **Production-Ready dari Hari Pertama**:
   - OpenTelemetry tracing, structured logging (JSON), dan health check endpoint bawaan.
   - Graceful shutdown dengan drain request yang sedang berjalan.

### 2.2 Bukan Tujuan (Non-Goals)
- Tidak meniru PHP secara 1-ke-1 (tidak menggunakan reflection berat yang memperlambat startup Go).
- Tidak mengikat developer pada satu library UI/Frontend tertentu.
- Tidak menyediakan ORM penuh bawaan — menggunakan adapter yang dapat diganti (`database/sql` raw atau GORM).
- Tidak bertanggung jawab atas infrastruktur deployment (Kubernetes, Docker Compose) — menyediakan binary tunggal saja.
- Tidak mendukung non-Go runtime (Node.js, Python, dll.).

---

## 3. Arsitektur Inti (Core Pillars)

### 3.1 Diagram Arsitektur Berlapis

```text
+------------------------------------------------------------------------+
|                        APLIKASI PENGGUNA                               |
|            (app/handlers, app/actions, app/models)                     |
+------------------------------------------------------------------------+
|                           CORE KERNEL                                  |
|  +------------------+  +-----------------+  +----------------------+  |
|  |  App Container   |  |  Config Loader  |  | DI / Context Registry|  |
|  |   (Lifecycle)    |  |  (Viper+.env)   |  | (Tenant-aware ctx)   |  |
|  +------------------+  +-----------------+  +----------------------+  |
+------------------------------------------------------------------------+
|                    TRANSPORT LAYER  (Triple Mode)                      |
|  +------------------+  +-----------------+  +----------------------+  |
|  |   ConnectRPC     |  | REST HTTP Gate  |  | In-Memory IPC Bridge |  |
|  |   (gRPC+JSON)    |  | (HTTP/1.1+2)    |  | (Wails Desktop/Mobile)|  |
|  +------------------+  +-----------------+  +----------------------+  |
+------------------------------------------------------------------------+
|                   PIPELINE & MIDDLEWARE ENGINE                         |
|  +------------------+  +-----------------+  +----------------------+  |
|  |  TenantResolver  |  |  Auth & RBAC    |  | RateLimiter / Logger |  |
|  | (sub/header/jwt) |  |  (JWT + Roles)  |  | CORS/Recover/OTEL    |  |
|  +------------------+  +-----------------+  +----------------------+  |
+------------------------------------------------------------------------+
|                     DATA & STORAGE LAYER                               |
|  +------------------+  +-----------------+  +----------------------+  |
|  |  Universal DB    |  | Schema Migrator |  |  Tenant Isolator     |  |
|  | (sql + GORM)     |  | (Up/Down/Redo)  |  | (PG SearchPath /     |  |
|  | SQLite/MySQL/PG  |  | per-tenant      |  |  SQLite per-file)    |  |
|  +------------------+  +-----------------+  +----------------------+  |
+------------------------------------------------------------------------+
|                   CROSS-CUTTING CONCERNS                               |
|   OpenTelemetry Traces  *  Structured Logs (slog/zap)  *  Health Probes|
|   Graceful Shutdown  *  Panic Recovery  *  Circuit Breaker             |
+------------------------------------------------------------------------+
```

### 3.2 Prinsip Desain Kunci

| Prinsip | Implementasi |
|---|---|
| **Convention over Configuration** | Struktur folder standar, nama file konvensional |
| **Zero Global State** | Semua state lewat `context.Context` atau DI container |
| **Interface First** | Setiap komponen core didefinisikan sebagai Go interface |
| **Composable Middleware** | Middleware dapat di-stack secara deklaratif di routes |
| **Fail Fast & Loud** | Konfigurasi invalid menyebabkan panic saat startup, bukan runtime |

---

## 4. Struktur Folder Standar (Framework Skeleton)

Struktur repositori yang rapi, modular, dan memisahkan framework core dengan aplikasi:

```text
tgo/                                   # Root repository framework
+-- cmd/
|   +-- craft/                         # CLI Tool ("Artisan" for Go)
|       +-- main.go
|       +-- commands/                  # Sub-commands (serve, make:*, migrate, seed)
|
+-- pkg/                               # Framework Core Packages (Reusable Library)
|   +-- app/                           # Application Container & Lifecycle
|   |   +-- app.go                     # App struct, Boot(), Shutdown()
|   |   +-- provider.go                # ServiceProvider interface
|   |
|   +-- config/                        # Config & .env management
|   |   +-- config.go                  # Loader: Viper + godotenv
|   |   +-- schema.go                  # Typed config structs per domain
|   |
|   +-- transport/
|   |   +-- connect/                   # ConnectRPC engine
|   |   |   +-- server.go              # Mux setup (gRPC + JSON)
|   |   |   +-- interceptor.go         # ConnectRPC interceptor chain
|   |   +-- inmemory/                  # In-memory RPC bridge
|   |   |   +-- invoker.go             # Direct Go function call adapter
|   |   +-- middleware/                # Shared HTTP/RPC middleware
|   |       +-- tenant.go              # Tenant detection & context injection
|   |       +-- auth.go                # JWT validation & claims extraction
|   |       +-- rbac.go                # Role-based permission gate
|   |       +-- ratelimit.go           # Token bucket rate limiter
|   |       +-- logger.go              # Structured request/response logger
|   |       +-- recover.go             # Panic recovery middleware
|   |       +-- cors.go                # CORS policy middleware
|   |
|   +-- database/
|   |   +-- engine.go                  # DBEngine interface (abstract)
|   |   +-- connection.go              # Connection pool manager
|   |   +-- dialect.go                 # SQLite / MySQL / PostgreSQL adapters
|   |   +-- migrator.go                # Migration engine (up/down/redo/status)
|   |   +-- tenant.go                  # Schema & database isolation manager
|   |
|   +-- auth/
|   |   +-- jwt.go                     # Token signer & validator (RS256/HS256)
|   |   +-- rbac.go                    # Role & Permission evaluator
|   |
|   +-- event/                         # Event dispatcher & Async Job Queue
|   |   +-- dispatcher.go              # Sync/Async event bus
|   |   +-- queue.go                   # Background job runner (goroutine pool)
|   |
|   +-- telemetry/                     # Observability
|       +-- otel.go                    # OpenTelemetry tracer & meter setup
|       +-- logger.go                  # slog/zap structured logger factory
|       +-- health.go                  # /healthz & /readyz HTTP handlers
|
+-- starter/                           # Template Project Baru (Boilerplate)
|   +-- app/
|   |   +-- actions/                   # Domain Business Logic (Services)
|   |   +-- handlers/                  # RPC & HTTP Handlers (Controllers)
|   |   +-- models/                    # Structs & Data Entities
|   |   +-- middleware/                # Custom App Middlewares
|   +-- config/
|   |   +-- app.yaml
|   |   +-- database.yaml
|   +-- database/
|   |   +-- migrations/                # SQL / Go migration files
|   |   +-- seeders/                   # Seed scripts
|   +-- proto/                         # Protocol Buffer contracts
|   |   +-- v1/
|   |       +-- auth.proto
|   |       +-- example.proto
|   +-- routes/
|   |   +-- api.go                     # Route bindings & handler registration
|   +-- .env.example
|   +-- go.mod
|   +-- main.go                        # Application entrypoint
|
+-- docs/                              # Framework documentation
|   +-- adr/                           # Architecture Decision Records
|   +-- guides/                        # Getting started, tutorials
|
+-- README.md
```

---

## 5. Spesifikasi Komponen Teknis

### 5.1 Transport Layer (ConnectRPC + HTTP/JSON + In-Memory)

**Teknologi & Versi**:
- `connectrpc.com/connect` v1.x
- `google.golang.org/protobuf` v1.x
- `golang.org/x/net/http2` (via ConnectRPC)

**Fitur**:
- Server menerima request `Content-Type: application/proto` (gRPC binary, HTTP/2) dan `application/json` (REST, HTTP/1.1) **pada port yang sama** tanpa konfigurasi tambahan.
- `InProcessTransport`: adapter yang memungkinkan aplikasi Desktop & Mobile (Wails v3) memanggil handler RPC secara *in-memory* **tanpa membuka port network**, memberikan latency 0ms.
- Semua transport melewati middleware pipeline yang sama (autentikasi, logging, tracing).

**Interface Kunci**:
```go
// pkg/transport/connect/server.go
type Server interface {
    Register(path string, handler http.Handler)
    Serve(ctx context.Context, addr string) error
    ServeInMemory(ctx context.Context) InMemoryInvoker
}

// pkg/transport/inmemory/invoker.go
type InMemoryInvoker interface {
    Invoke(ctx context.Context, procedure string, req, resp proto.Message) error
}
```

---

### 5.2 Multi-Tenant Engine

**Strategi Resolusi Tenant** (dieksekusi berurutan, short-circuit pada match pertama):

| Prioritas | Sumber | Contoh |
|:---------:|--------|--------|
| 1 | Subdomain | `client1.app.com` -> tenant `client1` |
| 2 | HTTP Header | `X-Tenant-Slug: client1` |
| 3 | JWT Claim | Payload: `"tenant": "client1"` |
| 4 | Default | Fallback ke tenant `public` (jika dikonfigurasi) |

**Strategi Isolasi Database**:

| Database | Strategi | Detail |
|---|---|---|
| **PostgreSQL** | Schema-per-tenant | `CREATE SCHEMA IF NOT EXISTS "tenant_{slug}"` + `SET search_path = tenant_{slug}, public` |
| **SQLite** | File-per-tenant | `{data_dir}/{slug}.db` |
| **MySQL** | Database-per-tenant | `CREATE DATABASE IF NOT EXISTS tenant_{slug}` atau table prefix |

**Tenant Lifecycle**:
```
Provisioning -> [create schema] -> [run migrations] -> [seed data] -> Active
Suspension   -> [revoke access] -> Suspended (data retained)
Deletion     -> [archive] -> [drop schema] -> Deleted (soft by default)
```

**Interface Kunci**:
```go
// pkg/transport/middleware/tenant.go
type TenantResolver interface {
    Resolve(r *http.Request) (tenantSlug string, err error)
}

// pkg/database/tenant.go
type TenantIsolator interface {
    DBForTenant(ctx context.Context, slug string) (*sql.DB, error)
    ProvisionTenant(ctx context.Context, slug string) error
    DeprovisionTenant(ctx context.Context, slug string) error
}
```

---

### 5.3 Database & Migrator

**Abstraksi**: Layer `DBEngine` (interface) mendukung `database/sql` raw (performa tinggi) dan GORM (produktivitas) dengan adapter dialek yang dapat ditukar. Injeksi ke handler/service melalui `context.Context` — tidak ada global DB variable.

**Kompatibilitas SQL Lintas Dialek**:

| Fitur | SQLite | MySQL | PostgreSQL |
|---|---|---|---|
| Parameter placeholder | `?` | `?` | `$1, $2, ...` |
| Auto-increment | `INTEGER PRIMARY KEY` | `AUTO_INCREMENT` | `SERIAL` / `BIGSERIAL` |
| Boolean | `0/1` (integer) | `TINYINT(1)` | `BOOLEAN` |
| Timestamp now | `CURRENT_TIMESTAMP` | `NOW()` | `NOW()` |
| JSON column | `TEXT` (json validated) | `JSON` | `JSONB` |

**Migrator Engine**:
- File migrasi: `{timestamp}_{name}.up.sql` dan `{timestamp}_{name}.down.sql`.
- State tracking via tabel `schema_migrations` pada setiap schema/database.
- Mendukung migrasi per-tenant: `craft migrate --tenant=all` atau `--tenant=slug`.
- Transaksi per-migrasi: rollback otomatis jika migrasi gagal.

---

### 5.4 Auth & RBAC

**JWT**:
- Mendukung algoritma **RS256** (production, asymmetric key) dan **HS256** (development).
- Claims standar: `sub`, `exp`, `iat`, `tenant`, `roles`, `permissions`.
- Token refresh dengan sliding window (configurable TTL).

**RBAC Model**:
```
User -> N Roles -> N Permissions
Permission = "resource:action" (contoh: "orders:read", "invoices:write")
```

---

### 5.5 CLI Scaffolding (`craft`)

**Runtime**: `github.com/spf13/cobra` v1.x

| Perintah | Deskripsi | Output |
|---|---|---|
| `craft serve` | Jalankan server dengan hot-reload | Process berjalan di foreground |
| `craft make:service <Name>` | Buat template business logic | `app/actions/{name}_action.go` |
| `craft make:handler <Name>` | Buat RPC/HTTP controller | `app/handlers/{name}_handler.go` |
| `craft make:model <Name>` | Buat struct model + repo | `app/models/{name}.go` |
| `craft make:frontend` | Generate Wails Unified UI starter (Vue/Svelte + Tailwind) | `frontend/` (untuk Desktop & Mobile) |
| `craft make:migration <name>` | Buat file migrasi dengan timestamp | `database/migrations/{ts}_{name}.up.sql` |
| `craft migrate [--tenant=all\|slug]` | Jalankan migrasi pending | Status per migration file |
| `craft migrate:rollback [--step=N]` | Rollback N migrasi terakhir | Status rollback |
| `craft migrate:status` | Tampilkan status semua migrasi | Table: applied / pending |
| `craft seed [--class=Name]` | Jalankan seeder | Log output seeding |
| `craft tenant:create <slug>` | Provision tenant baru | Schema + migration + seed |
| `craft tenant:list` | List semua tenant aktif | Table format |

---

### 5.6 Event & Job Queue

- **Event Bus**: Sinkron dan asinkron (goroutine pool dengan bounded concurrency).
- **Job Queue**: Menggunakan channel-based worker pool. Mendukung retry dengan exponential backoff.
- Extensible: dapat diganti dengan adapter Redis Streams / NATS untuk produksi skala besar.

```go
// Dispatch event asinkron
app.Events().Dispatch(ctx, &events.UserRegistered{UserID: "123"})

// Register listener
app.Events().Listen(&events.UserRegistered{}, handlers.SendWelcomeEmail)
```

---

## 6. Contoh Alur Eksekusi Kode

### 6.1 Definisi Proto Contract (`proto/v1/user.proto`)

```protobuf
syntax = "proto3";
package user.v1;
option go_package = "app/proto/user/v1;userv1";

service UserService {
  rpc GetProfile (GetProfileRequest) returns (GetProfileResponse);
  rpc UpdateProfile (UpdateProfileRequest) returns (UpdateProfileResponse);
}

message GetProfileRequest {
  string user_id = 1;
}

message GetProfileResponse {
  string user_id = 1;
  string name = 2;
  string email = 3;
  string role = 4;
  string tenant_slug = 5;
}
```

### 6.2 Business Logic Service (`app/actions/user_action.go`)

```go
package actions

import (
	"context"
	"fmt"
	"my-app/app/models"
	"tgo/pkg/database"
)

type UserAction struct {
	db database.DBEngine
}

func NewUserAction(db database.DBEngine) *UserAction {
	return &UserAction{db: db}
}

// GetProfile otomatis menggunakan database tenant yang sesuai context.
// Tidak perlu menulis WHERE tenant_id = ? secara manual.
func (a *UserAction) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	tenantDB, err := a.db.FromContext(ctx) // resolves tenant-specific DB/schema
	if err != nil {
		return nil, fmt.Errorf("resolving tenant db: %w", err)
	}

	var user models.User
	if err := tenantDB.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
		return nil, fmt.Errorf("fetching user %s: %w", userID, err)
	}
	return &user, nil
}
```

### 6.3 RPC Handler (`app/handlers/user_handler.go`)

```go
package handlers

import (
	"context"
	userv1 "my-app/proto/user/v1"
	"connectrpc.com/connect"
)

type UserHandler struct {
	action *actions.UserAction
}

func (h *UserHandler) GetProfile(
	ctx context.Context,
	req *connect.Request[userv1.GetProfileRequest],
) (*connect.Response[userv1.GetProfileResponse], error) {
	user, err := h.action.GetProfile(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(&userv1.GetProfileResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}), nil
}
```

### 6.4 Registrasi Routes (`routes/api.go`)

```go
package routes

import (
	"tgo/pkg/transport/middleware"
	userv1connect "my-app/proto/user/v1/userv1connect"
)

func Register(app *tgo.App) {
	mux := app.Mux()

	// Global middleware stack
	mux.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORS(app.Config().CORSPolicy()),
		middleware.TenantResolver(),
		middleware.RateLimiter(100, "per_minute"),
	)

	// Authenticated routes
	auth := mux.Group("/", middleware.RequireAuth())
	{
		path, handler := userv1connect.NewUserServiceHandler(&handlers.UserHandler{})
		auth.Handle(path, handler)
	}
}
```

---

## 7. Architecture Decision Records (ADR)

### ADR-001: ConnectRPC sebagai Transport Utama

- **Status**: Accepted
- **Konteks**: Perlu satu transport yang mendukung gRPC binary dan REST JSON tanpa proxy tambahan.
- **Keputusan**: Menggunakan `connectrpc.com/connect`. ConnectRPC secara native mendukung protocol gRPC, gRPC-Web, dan Connect JSON di satu handler.
- **Konsekuensi**: Developer harus mendefinisikan kontrak lewat `.proto` file. Trade-off: lebih verbose di awal, tapi memberi type-safety dan schema documentation otomatis.
- **Alternatif Ditolak**: grpc-gateway (memerlukan gateway service terpisah), pure HTTP REST (tidak ada type-safety).

---

### ADR-002: Schema-per-Tenant untuk PostgreSQL

- **Status**: Accepted
- **Konteks**: Perlu isolasi data antar tenant yang kuat tanpa overhead operasional per-database instance.
- **Keputusan**: Menggunakan PostgreSQL schema isolation (`SET search_path`). Satu database instance, schema terpisah per tenant.
- **Konsekuensi**: Query otomatis ter-scope ke schema tenant. Cross-tenant query memerlukan explicit schema prefix.
- **Limitasi**: Tidak ideal jika jumlah tenant > 10.000. Untuk kasus tersebut, gunakan strategi database-per-tenant dengan connection pool.
- **Alternatif Ditolak**: Row-level isolation (risiko data leakage jika `WHERE` clause terlewat), Database-per-tenant (overhead operasional tinggi untuk < 1.000 tenant).

---

### ADR-003: `context.Context` sebagai Carrier State Tenant

- **Status**: Accepted
- **Konteks**: Perlu cara yang aman dan idiomatis Go untuk membawa informasi tenant sepanjang request lifecycle.
- **Keputusan**: Tenant slug, resolved DB connection, dan claims JWT disimpan di `context.Context` menggunakan typed key (bukan string literal) untuk menghindari collision.
- **Konsekuensi**: Zero global state. Setiap goroutine memiliki context yang jelas. Unit testing lebih mudah karena context dapat di-mock.

---

### ADR-004: slog sebagai Structured Logger Standar

- **Status**: Accepted
- **Konteks**: `log/slog` (Go 1.21+) adalah standard library structured logger.
- **Keputusan**: Default menggunakan `slog` dengan JSON handler untuk production. Dapat diganti dengan `zap` atau `zerolog` via adapter interface.
- **Alternatif Ditolak**: logrus (sudah maintenance mode), zap (dependency eksternal sebagai default).

---

## 8. Model Keamanan (Security Model)

### 8.1 Autentikasi

| Layer | Mekanisme | Detail |
|---|---|---|
| **API External** | JWT Bearer Token | RS256 (asymmetric) untuk production |
| **Service-to-Service** | JWT dengan scope terbatas | Issuer berbeda, TTL pendek (5 menit) |
| **In-Memory (Desktop)** | Context injection langsung | Tidak melalui network, auth dilakukan saat login desktop |

### 8.2 Otorisasi (RBAC)

- Permission granular: `resource:action` (contoh: `invoices:read`, `users:delete`).
- Roles dapat di-compose (user bisa punya multiple roles).
- Permission check dilakukan di middleware layer, **bukan** di business logic.

### 8.3 Isolasi Tenant

- **Tenant context injection wajib** untuk semua request yang memerlukan data.
- Request tanpa tenant context yang valid akan di-reject di middleware (HTTP 403).
- Cross-tenant access hanya dimungkinkan melalui **Super Admin role** dengan explicit audit log.

### 8.4 Security Headers & Hardening

```go
// Middleware yang diaktifkan secara default
middleware.SecurityHeaders()          // X-Frame-Options, X-Content-Type-Options, HSTS
middleware.RequestSizeLimit(8 << 20)  // Max 8MB request body
middleware.InputSanitizer()           // Reject null bytes & invalid UTF-8
```

### 8.5 Secret Management

- Secrets (DB password, JWT private key) **tidak boleh** di-hardcode.
- Urutan resolusi konfigurasi: Environment Variable -> `.env` file -> Vault (via plugin).
- Private key JWT disimpan di file `.pem` di luar repository (dicantumkan di `.gitignore`).

---

## 9. Observabilitas & Monitoring

### 9.1 Structured Logging

- Format: JSON ke stdout (di-pickup oleh log aggregator).
- Field wajib setiap log entry: `timestamp`, `level`, `tenant_slug`, `request_id`, `service`, `message`.
- Trace ID otomatis disisipkan dari OpenTelemetry span.

```json
{
  "timestamp": "2026-09-23T12:00:00Z",
  "level": "INFO",
  "service": "user-service",
  "tenant_slug": "acme-corp",
  "request_id": "req_xyz123",
  "trace_id": "abc...def",
  "message": "GetProfile called",
  "user_id": "usr_456"
}
```

### 9.2 Distributed Tracing

- **Provider**: OpenTelemetry SDK (`go.opentelemetry.io/otel`)
- **Exporter**: OTLP ke Jaeger / Grafana Tempo (configurable via env).
- Span otomatis dibuat untuk: incoming HTTP/gRPC request, DB query, event dispatch.

### 9.3 Metrics

**Provider**: OpenTelemetry Metrics -> Prometheus exporter. Metric bawaan:

| Metric | Type | Label |
|---|---|---|
| `tgo_requests_total` | Counter | `method`, `path`, `status`, `tenant` |
| `tgo_request_duration_seconds` | Histogram | `method`, `path`, `tenant` |
| `tgo_db_query_duration_seconds` | Histogram | `operation`, `table`, `tenant` |
| `tgo_active_tenants` | Gauge | - |

### 9.4 Health Probes

- `GET /healthz` — Liveness probe: selalu 200 jika server up.
- `GET /readyz` — Readiness probe: cek koneksi DB, dependency kritis.

---

## 10. Strategi Pengujian (Testing Strategy)

### 10.1 Piramida Testing

```
        +-------+
        |  E2E  |   <- Sedikit, test alur bisnis utama
      +-----------+
      | Integration|  <- Test transport + DB real
    +---------------+
    |     Unit       |  <- Mayoritas, test logic murni
    +---------------+
```

### 10.2 Unit Testing

- Setiap komponen `pkg/` harus memiliki file `_test.go` dengan coverage minimal **80%**.
- Menggunakan Go standard `testing` package + `testify/assert`.
- Tenant context di-mock via helper `testutil.ContextWithTenant(ctx, "test-tenant")`.

### 10.3 Integration Testing

- **Database**: Menggunakan `testcontainers-go` untuk spin up PostgreSQL/MySQL real instance per test suite.
- **Transport**: Test handler menggunakan `httptest.NewRecorder()` dan ConnectRPC test client.
- Setiap test integration membuat schema tenant terpisah dan di-cleanup setelah test selesai.

### 10.4 E2E Testing

- Menggunakan `starter/` template sebagai aplikasi referensi.
- Test script `make test-e2e` menjalankan server lokal, melakukan serangkaian request via gRPC dan HTTP, dan memvalidasi response.

### 10.5 Benchmark Target

| Komponen | Target P99 |
|---|---|
| ConnectRPC JSON handler (no DB) | < 1ms |
| In-Memory IPC | < 0.1ms |
| DB Query single row (PostgreSQL, local) | < 5ms |
| App startup time | < 500ms |

---

## 11. Rencana Tahapan Implementasi (Roadmap)

### Fase 1: Repository Bootstrap & Core Engine
**Target**: 2 minggu

| Task | Acceptance Criteria |
|---|---|
| Init repo + go.mod | `go build ./...` sukses tanpa error |
| `pkg/app` Container | `app.New()` dapat boot dan shutdown gracefully |
| `pkg/config` Loader | `.env` dan YAML ter-load, panic jika key wajib tidak ada |
| Starter template | `go run ./starter/main.go` berhasil start server |

---

### Fase 2: Universal Database & Migrator
**Target**: 2 minggu

| Task | Acceptance Criteria |
|---|---|
| `pkg/database/engine.go` interface | Interface dapat di-implement oleh adapter manapun |
| SQLite adapter | CRUD operasi berjalan, test pass |
| MySQL adapter | CRUD operasi berjalan, test pass |
| PostgreSQL adapter | CRUD operasi berjalan, test pass |
| Migrator `up/down/status` | Migrasi berjalan, state tersimpan di `schema_migrations` |
| `craft migrate` CLI | Perintah dapat dijalankan dari terminal |

---

### Fase 3: Multi-Tenancy Engine
**Target**: 2 minggu

| Task | Acceptance Criteria |
|---|---|
| `TenantResolver` (subdomain) | Request ke `x.app.com` -> context berisi `tenant_slug=x` |
| `TenantResolver` (header/JWT) | Header `X-Tenant-Slug` ter-resolve dengan benar |
| PostgreSQL schema isolator | Query otomatis menggunakan schema tenant yang tepat |
| SQLite file isolator | Setiap tenant memiliki file `.db` terpisah |
| `craft tenant:create` | Provisioning schema + migrasi + seed berjalan |
| Cross-tenant isolation test | Request tenant A tidak dapat melihat data tenant B |

---

### Fase 4: ConnectRPC & Triple Transport
**Target**: 3 minggu

| Task | Acceptance Criteria |
|---|---|
| ConnectRPC server setup | Handler menerima `application/proto` dan `application/json` |
| Middleware chain integration | Tenant + Auth middleware berjalan di semua transport |
| In-Memory IPC bridge | Desktop client dapat invoke handler tanpa network |
| Starter proto + handler | `GetProfile` via gRPC dan REST mengembalikan data yang sama |
| Benchmark transport layer | Target performa P99 terpenuhi |

---

### Fase 5: CLI Tool ("Craft") & Dokumentasi
**Target**: 2 minggu

| Task | Acceptance Criteria |
|---|---|
| `craft make:*` generators | File-file yang di-generate memiliki struktur yang benar |
| `craft serve` dengan hot-reload | Perubahan file me-restart server secara otomatis |
| README & Getting Started Guide | Developer baru dapat setup dalam < 15 menit |
| ADR docs tersimpan di `docs/adr/` | Semua keputusan arsitektur terdokumentasi |

---

## 12. Success Metrics & KPI

### Developer Experience

| Metric | Target |
|---|---|
| Time-to-first-running-server (dari `go mod init`) | < 15 menit |
| Lines of boilerplate untuk membuat 1 endpoint baru | < 30 baris |
| Test coverage `pkg/` core | >= 80% |

### Performa Teknis

| Metric | Target |
|---|---|
| ConnectRPC throughput (simple handler, no DB) | >= 10.000 RPS per core |
| DB query P99 (single row, PostgreSQL) | < 5ms |
| Application startup time | < 500ms |
| Binary size (starter app) | < 30MB |

### Kualitas & Maintainability

| Metric | Target |
|---|---|
| Go vet & staticcheck | Zero issues (wajib di CI) |
| Cyclomatic complexity per function | <= 10 |
| Direct dependency jumlah | < 20 package |

---

## 13. Risiko & Mitigasi

| ID | Risiko | Kemungkinan | Dampak | Mitigasi |
|---|---|---|---|---|
| R-01 | ConnectRPC API breaking change | Rendah | Tinggi | Pin versi minor, monitor changelog |
| R-02 | PostgreSQL schema > 10K tenant degradation | Sedang | Sedang | Benchmark di Fase 3, siapkan DB-per-tenant fallback |
| R-03 | Kompleksitas tinggi -> adoption rendah | Sedang | Tinggi | Prioritaskan DX, starter template, dan docs di Fase 5 |
| R-04 | Go version compatibility (min Go 1.21 karena slog) | Rendah | Sedang | Dokumentasikan min Go version, CI test Go 1.21 & 1.22 |
| R-05 | Context misuse menyebabkan data leakage antar tenant | Rendah | Kritis | Integration test cross-tenant isolation wajib di Fase 3 |
| R-06 | In-Memory IPC tidak thread-safe untuk desktop | Sedang | Tinggi | Race detector aktif di test, benchmark concurrent calls |

---

## 14. Open Questions

> Item berikut perlu keputusan tim sebelum implementasi dimulai.

| # | Pertanyaan | Dampak | Status |
|---|---|---|---|
| OQ-01 | Apakah menggunakan GORM sebagai default ORM, atau pure `database/sql`? | DX vs Performa | Belum Diputuskan |
| OQ-02 | Strategi hot-reload di `craft serve`: `air` sebagai dependency atau implementasi sendiri? | Kompleksitas CLI | Belum Diputuskan |
| OQ-03 | Apakah framework mendukung multi-region (tenant di region berbeda)? | Scope Fase 1 | Di-defer ke post-v1 |
| OQ-04 | Format config: YAML saja, atau juga mendukung TOML dan JSON? | DX | Butuh Input Tim |
| OQ-05 | Apakah `craft` perlu plugin system (3rd party command)? | Extensibility | Di-defer ke post-v1 |
| OQ-06 | Apakah ada kebutuhan GraphQL transport selain gRPC + REST? | Scope | Belum Diputuskan |

---

## 15. Glosarium

| Istilah | Definisi |
|---|---|
| **ConnectRPC** | Protocol RPC modern yang kompatibel dengan gRPC dan mendukung JSON over HTTP/1.1 secara native. |
| **DBEngine** | Interface abstraksi database TGo yang dapat di-implementasi oleh berbagai adapter. |
| **DI Container** | Dependency Injection Container — mekanisme untuk meregister dan me-resolve service dependencies. |
| **In-Memory IPC** | Inter-Process Communication dalam satu binary, tanpa network socket. Digunakan untuk embed backend di aplikasi desktop. |
| **Middleware Pipeline** | Rantai fungsi yang dieksekusi berurutan sebelum request sampai ke handler utama. |
| **Multi-Tenancy** | Arsitektur di mana satu instance aplikasi melayani banyak pelanggan (tenant) dengan isolasi data. |
| **Proto / Protobuf** | Protocol Buffers — format serialisasi data biner dari Google, digunakan sebagai kontrak API di ConnectRPC/gRPC. |
| **RBAC** | Role-Based Access Control — sistem otorisasi berdasarkan peran (role) yang dimiliki pengguna. |
| **Schema Isolation** | Teknik isolasi data PostgreSQL di mana setiap tenant memiliki schema terpisah dalam satu database. |
| **Tenant** | Satu pelanggan/klien dalam sistem multi-tenant. Diidentifikasi dengan `slug` unik. |
| **Transport Layer** | Layer yang bertanggung jawab atas protokol komunikasi (gRPC, HTTP, In-Memory). |

---

## 16. Panduan Memulai di Repositori Baru

### Prerequisites

```bash
# Pastikan Go 1.21+ terinstall
go version  # go version go1.21.x ...

# Install buf (Protocol Buffer toolchain)
# https://buf.build/docs/installation

# Install tools
go install github.com/air-verse/air@latest       # Hot reload
```

### Setup Langkah-demi-Langkah

```bash
# 1. Buat folder proyek baru
mkdir d:\Projects\tgo && cd d:\Projects\tgo

# 2. Salin starter template dari framework
# Ikuti instruksi di README framework

# 3. Inisialisasi Go module
go mod init my-app-name

# 4. Sesuaikan konfigurasi
copy .env.example .env
# Edit .env sesuai environment lokal Anda

# 5. Jalankan migrasi awal
go run ./cmd/craft migrate

# 6. Jalankan server development
go run ./cmd/craft serve
# Server berjalan di http://localhost:8080
# gRPC: grpc://localhost:8080
# Health: http://localhost:8080/healthz
```

### Membuat Endpoint Pertama

```bash
# 1. Tambahkan definisi ke proto file
# Edit proto/v1/product.proto

# 2. Generate Go code dari proto
buf generate

# 3. Buat handler
go run ./cmd/craft make:handler Product

# 4. Buat service/action
go run ./cmd/craft make:service Product

# 5. Buat migrasi database
go run ./cmd/craft make:migration create_products_table

# 6. Register route di routes/api.go
```

---

*Dokumen ini merupakan living document. Update versi pada metadata header setiap kali ada perubahan signifikan.*

