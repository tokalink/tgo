# ADR-003: `context.Context` sebagai Carrier State Tenant

## Status
Accepted

## Konteks
Dibutuhkan mekanisme idiomatis Go untuk membawa informasi tenant aktif dari transport middleware ke service logic dan database adapter secara aman lintas goroutine.

## Keputusan
Menyimpan tenant slug dan data identitas di dalam `context.Context` menggunakan unexported typed context key (`type contextKey string`).

## Konsekuensi
- **Zero Global State**: Menjamin request safety saat menangani ribuan goroutine bersamaan.
- **Testability**: Sangat mudah melakukan unit test dengan menyuntikkan `context.WithValue`.
- Fungsi database dapat meresolusi koneksi yang sesuai melalui `db.FromContext(ctx)`.
