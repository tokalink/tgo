package middleware

import (
	"net/http"
	"strings"
)

// CORSOptions configures the CORS middleware.
type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSOptions provides secure sensible defaults for development and production.
func DefaultCORSOptions() CORSOptions {
	return CORSOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Tenant-Slug", "X-Tenant-ID", "Connect-Protocol-Version"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORS returns an HTTP middleware handling Cross-Origin Resource Sharing according to opts.
func CORS(opts ...CORSOptions) func(http.Handler) http.Handler {
	cfg := DefaultCORSOptions()
	if len(opts) > 0 {
		cfg = opts[0]
	}

	originsHeader := strings.Join(cfg.AllowedOrigins, ", ")
	methodsHeader := strings.Join(cfg.AllowedMethods, ", ")
	headersHeader := strings.Join(cfg.AllowedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				// If wildcard or matching origin
				if originsHeader == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					for _, o := range cfg.AllowedOrigins {
						if o == origin || o == "*" {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}

				if cfg.AllowCredentials && originsHeader != "*" {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				w.Header().Set("Access-Control-Allow-Methods", methodsHeader)
				w.Header().Set("Access-Control-Allow-Headers", headersHeader)
			}

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
