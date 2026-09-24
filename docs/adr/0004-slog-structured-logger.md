# ADR-004: Standard Structured Logging dengan Tenant & Request Trace

## Status
Accepted

## Konteks
Observabilitas dan kemudahan debugging request multi-tenant membutuhkan structured logging (JSON) yang menyertakan informasi tenant, durasi eksekusi, status code, dan ukuran respons.

## Keputusan
Menyediakan structured request/response logging middleware yang secara konsisten mencatat konteks tenant (`tenant={slug}`), method, path, latency, dan HTTP status code.

## Konsekuensi
- Log dapat langsung diproses oleh log aggregators (Datadog, Loki, CloudWatch).
- Developer dapat melacak anomali per-tenant secara instan.
