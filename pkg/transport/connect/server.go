package connect

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tgo-framework/tgo/pkg/transport/inmemory"
	"github.com/tgo-framework/tgo/pkg/transport/welcome"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Server defines the unified transport server contract.
type Server interface {
	// Register registers a handler at the specified path prefix.
	Register(path string, handler http.Handler)

	// Use appends middlewares to the global request pipeline.
	Use(middlewares ...func(http.Handler) http.Handler)

	// Handler returns the fully assembled HTTP handler (with all middlewares applied).
	Handler() http.Handler

	// Serve starts the HTTP/1.1 and HTTP/2 (h2c) server on the given address.
	Serve(ctx context.Context, addr string) error

	// ServeInMemory returns an In-Memory Invoker for desktop/mobile in-process calls.
	ServeInMemory(ctx context.Context) inmemory.InMemoryInvoker
}

type connectServer struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

// NewServer creates a new unified ConnectRPC & HTTP transport server.
func NewServer() Server {
	mux := http.NewServeMux()
	// Mount interactive welcome dashboard on root path
	mux.Handle("/", welcome.Handler())

	return &connectServer{
		mux:         mux,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Register mounts a path and handler to the server's router.
func (s *connectServer) Register(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

// Use appends middlewares to the global middleware pipeline.
func (s *connectServer) Use(middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Handler returns the root HTTP handler wrapped in all registered middlewares.
func (s *connectServer) Handler() http.Handler {
	var handler http.Handler = s.mux
	// Apply middlewares in reverse so first added is executed first
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}
	return handler
}

// ServeInMemory provides an in-memory invoker targeting this server's handlers and middleware stack.
func (s *connectServer) ServeInMemory(ctx context.Context) inmemory.InMemoryInvoker {
	return inmemory.NewInvoker(s.Handler())
}

// Serve starts the server supporting both HTTP/1.1 and HTTP/2 cleartext (h2c).
func (s *connectServer) Serve(ctx context.Context, addr string) error {
	rootHandler := s.Handler()

	// Support HTTP/2 cleartext (h2c) for high-performance gRPC clients
	h2s := &http2.Server{}
	h2cHandler := h2c.NewHandler(rootHandler, h2s)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           h2cHandler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("[ConnectServer] Listening on http://127.0.0.1%s (HTTP/1.1, HTTP/2 h2c, ConnectRPC)\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server listen error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("[ConnectServer] Shutting down gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}
