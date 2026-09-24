package database

import (
	"context"
	"database/sql"
)

// DBEngine defines the core contract for any database adapter
type DBEngine interface {
	// Conn returns the underlying *sql.DB connection pool
	Conn() *sql.DB

	// Close terminates the database connection
	Close() error

	// Ping verifies if the connection to the database is still alive
	Ping(ctx context.Context) error

	// FromContext retrieves a tenant-scoped database connection from the context.
	// For single-tenant setups, it just returns itself.
	FromContext(ctx context.Context) (DBEngine, error)
}

// Config holds the configuration needed to connect to a database
type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
