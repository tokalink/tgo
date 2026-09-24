# TGo Framework - Implementation Todo List

Ini adalah daftar tugas implementasi framework TGo berdasarkan Roadmap di dalam PRD. 
Beri tanda centang (`[x]`) pada tugas yang sudah diselesaikan.

## Fase 1: Repository Bootstrap & Core Engine
*Target: 2 Minggu*

- [x] Inisialisasi struktur repositori (`go mod init`).
- [x] Buat package `pkg/app` (Application Container & Lifecycle: `Boot()`, `Shutdown()`).
- [x] Buat package `pkg/config` (Config Loader menggunakan Viper + `.env`).
- [x] Buat *Starter template* sederhana (`starter/main.go`).
- [x] Pastikan `go build ./...` sukses tanpa error dan `starter/main.go` bisa di-*run*.

## Fase 2: Universal Database & Migrator
*Target: 2 Minggu*

- [x] Definisikan interface abstrak `pkg/database/engine.go` (`DBEngine`).
- [x] Implementasikan Adapter SQLite (CRUD berjalan & unit test pass).
- [x] Implementasikan Adapter MySQL (CRUD berjalan & unit test pass).
- [x] Implementasikan Adapter PostgreSQL (CRUD berjalan & unit test pass).
- [x] Bangun engine Migrator (`up`, `down`, `status` dengan tabel `schema_migrations`).
- [x] Buat command `craft migrate` agar bisa dijalankan dari terminal.

## Fase 3: Multi-Tenancy Engine
*Target: 2 Minggu*

- [x] Implementasikan middleware `TenantResolver` berbasis **Subdomain**.
- [x] Implementasikan middleware `TenantResolver` berbasis **HTTP Header** / **JWT Claim**.
- [x] Buat *Tenant Isolator* untuk PostgreSQL (menggunakan fitur `search_path`).
- [x] Buat *Tenant Isolator* untuk SQLite (file `.db` terpisah per tenant).
- [x] Buat command `craft tenant:create` (menjalankan provision schema, migrasi, dan seed).
- [x] Tulis *Integration Test* untuk memastikan isolasi antar-tenant aman (Cross-tenant isolation test).

## Fase 4: ConnectRPC & Triple Transport
*Target: 3 Minggu*

- [x] Konfigurasi ConnectRPC Server (Menerima `application/proto` & `application/json` di satu port).
- [x] Integrasikan *Middleware Chain* (Tenant, Auth, Logger) pada semua transport.
- [x] Bangun **In-Memory IPC Bridge** (Agar Wails bisa invoke handler tanpa network TCP).
- [x] Buat implementasi `GetProfile` pada file `.proto` dan *handler* di starter template.
- [x] Lakukan benchmark layer transport untuk memastikan target latency & RPS terpenuhi.

## Fase 5: CLI Tool ("Craft"), Frontend & Dokumentasi
*Target: 2 Minggu*

- [x] Implementasikan CLI `craft make:model` (Struct + Repository).
- [x] Implementasikan CLI `craft make:service` & `craft make:handler`.
- [x] Implementasikan CLI `craft make:frontend` (Generate starter Web/Wails UI dengan Vue/Svelte + Tailwind).
- [x] Implementasikan CLI `craft serve` dengan fitur *hot-reload*.
- [x] Lengkapi *README* dan *Getting Started Guide* (Target: Developer baru bisa setup < 15 menit).
- [x] Pindahkan hasil keputusan arsitektur (ADR) ke folder `docs/adr/`.
