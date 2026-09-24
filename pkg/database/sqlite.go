package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tgo-framework/tgo/pkg/transport/middleware"
)

type SQLiteAdapter struct {
	db       *sql.DB
	isolator TenantIsolator
}

func NewSQLite(filepath string) (*SQLiteAdapter, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	return &SQLiteAdapter{db: db}, nil
}

// WithTenantIsolator attaches a tenant isolator to the adapter.
func (a *SQLiteAdapter) WithTenantIsolator(isolator TenantIsolator) *SQLiteAdapter {
	a.isolator = isolator
	return a
}

func (a *SQLiteAdapter) Conn() *sql.DB {
	return a.db
}

func (a *SQLiteAdapter) Close() error {
	if a.isolator != nil {
		_ = a.isolator.Close()
	}
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *SQLiteAdapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

func (a *SQLiteAdapter) FromContext(ctx context.Context) (DBEngine, error) {
	slug, ok := middleware.TenantFromContext(ctx)
	if !ok || slug == "" {
		return a, nil
	}

	if a.isolator == nil {
		return a, nil
	}

	tenantDB, err := a.isolator.DBForTenant(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant db for %q: %w", slug, err)
	}

	return &SQLiteAdapter{
		db:       tenantDB,
		isolator: a.isolator,
	}, nil
}
