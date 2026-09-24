package connect

import (
	"context"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/tgo-framework/tgo/pkg/transport/middleware"
)

// NewLoggingInterceptor returns a ConnectRPC unary interceptor for logging RPC calls.
func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			tenant := middleware.GetTenant(ctx)
			if tenant == "" {
				tenant = "public"
			}

			resp, err := next(ctx, req)
			duration := time.Since(start)

			status := "OK"
			if err != nil {
				status = connect.CodeOf(err).String()
			}

			log.Printf("[ConnectRPC] %s status=%s duration=%s tenant=%s",
				req.Spec().Procedure,
				status,
				duration,
				tenant,
			)

			return resp, err
		}
	}
}

// NewTenantInterceptor ensures tenant metadata in Connect headers is pushed into context if not already set.
func NewTenantInterceptor(headerNames ...string) connect.UnaryInterceptorFunc {
	if len(headerNames) == 0 {
		headerNames = []string{"x-tenant-slug", "x-tenant-id"}
	}

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if middleware.GetTenant(ctx) == "" {
				for _, h := range headerNames {
					if val := req.Header().Get(h); val != "" {
						ctx = middleware.WithTenant(ctx, val)
						break
					}
				}
			}
			return next(ctx, req)
		}
	}
}
