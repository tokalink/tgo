# ADR-001: ConnectRPC sebagai Transport Utama

## Status
Accepted

## Konteks
TGo memerlukan satu layer transport terpadu yang dapat melayani:
1. Native gRPC (HTTP/2 binary Protobuf) untuk microservices dan high-throughput calls.
2. RESTful HTTP/1.1 (JSON) langsung dari web browser dan mobile tanpa perantara `grpc-web proxy` atau gateway terpisah.
3. In-Memory direct execution untuk aplikasi Desktop & Mobile (Wails v3) tanpa latency network socket.

## Keputusan
Menggunakan `connectrpc.com/connect` sebagai fondasi transport. ConnectRPC secara native mendukung protokol gRPC, gRPC-Web, dan Connect JSON di satu HTTP handler.

## Konsekuensi
- Type-safety yang kuat dan dokumentasi skema otomatis melalui kontrak Protobuf.
- Eliminasi arsitektur ganda antara REST controller dan gRPC handler.
- Mengharuskan definisi kontrak `.proto` untuk RPC services.

## Alternatif Ditolak
- **grpc-gateway**: Memerlukan gateway runtime terpisah dan kompilasi plugin ganda.
- **Pure HTTP REST (Gin/Fiber)**: Tidak memiliki type-safety antar-layanan dan serialisasi binary efisien.
