package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tokalink/tgo/pkg/transport/connect"
	"github.com/tokalink/tgo/pkg/transport/middleware"
)

// App is the main application container for TGo
type App struct {
	server connect.Server
	addr   string
}

// New creates a new instance of TGo Application
func New() *App {
	srv := connect.NewServer()
	// Attach default global middleware pipeline
	srv.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORS(),
		middleware.TenantMiddleware(middleware.TenantMiddlewareOptions{}),
	)

	return &App{
		server: srv,
		addr:   ":8080",
	}
}

// SetAddr sets the listening address for the server.
func (a *App) SetAddr(addr string) *App {
	a.addr = addr
	return a
}

// Server returns the unified transport server.
func (a *App) Server() connect.Server {
	return a.server
}

// Boot initializes the application
func (a *App) Boot() error {
	log.Println("[TGo Kernel] Booting application container...")
	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("[TGo Kernel] Shutting down application gracefully...")
	return nil
}

// Run executes the application and blocks until an interrupt signal is received
func (a *App) Run() error {
	if err := a.Boot(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		if err := a.server.Serve(ctx, a.addr); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or server failure
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("\n[TGo Kernel] Interrupt signal received...")
		cancel()
		return a.Shutdown(ctx)
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}
