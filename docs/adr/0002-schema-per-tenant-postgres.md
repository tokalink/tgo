# ADR-002: Schema-per-Tenant untuk PostgreSQL

## Status
Accepted

## Konteks
SaaS modern membutuhkan pemisahan data antar-tenant yang kuat untuk mencegah kebocoran data (cross-tenant leakage) tanpa menimbulkan beban operasional berlebih mengelola ribuan database server terpisah.

## Keputusan
Mengadopsi strategi **Schema-per-tenant** pada PostgreSQL menggunakan skema bernama `tenant_{slug}` dan mengatur `search_path = tenant_{slug}, public`.

## Konsekuensi
- Isolasi logis yang kuat di tingkat database engine.
- Query SQL tidak perlu menambahkan klausul `WHERE tenant_id = ?` manual pada setiap query.
- Backup dan migrasi dapat dieksekusi per-tenant.

## Alternatif Ditolak
- **Row-Level Security / Shared Table**: Risiko tinggi data leak jika ada query developer yang tidak menyertakan filter tenant.
- **Database-per-Tenant di RDS/Postgres**: Overhead resource dan batas jumlah connection pool yang tinggi untuk jumlah tenant menengah.
