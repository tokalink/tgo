package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Logger returns an HTTP middleware that logs structured details of incoming HTTP and RPC requests.
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			tenant := GetTenant(r.Context())
			if tenant == "" {
				tenant = "public"
			}

			log.Printf("[HTTP] %s %s %d %s tenant=%s size=%dB",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
				tenant,
				wrapped.written,
			)
		})
	}
}
