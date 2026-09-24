package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

// Recover returns an HTTP middleware that recovers from panics, logs the stack trace, and writes a 500 error response.
func Recover() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()
					log.Printf("[PANIC RECOVERED] %v\nStack trace:\n%s", rec, string(stack))

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"internal server error","detail":%q}`, fmt.Sprintf("%v", rec))))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
