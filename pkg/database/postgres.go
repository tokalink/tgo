package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/tgo-framework/tgo/pkg/transport/middleware"
)

type PostgresAdapter struct {
	db       *sql.DB
	isolator TenantIsolator
}

func NewPostgres(cfg Config) (*PostgresAdapter, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresAdapter{db: db}, nil
}

// WithTenantIsolator attaches a tenant isolator to the Postgres adapter.
func (a *PostgresAdapter) WithTenantIsolator(isolator TenantIsolator) *PostgresAdapter {
	a.isolator = isolator
	return a
}

func (a *PostgresAdapter) Conn() *sql.DB {
	return a.db
}

func (a *PostgresAdapter) Close() error {
	if a.isolator != nil {
		_ = a.isolator.Close()
	}
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *PostgresAdapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

func (a *PostgresAdapter) FromContext(ctx context.Context) (DBEngine, error) {
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

	return &PostgresAdapter{
		db:       tenantDB,
		isolator: a.isolator,
	}, nil
}
