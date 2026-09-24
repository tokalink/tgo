package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// PostgresTenantIsolator manages schema-per-tenant PostgreSQL database isolation.
type PostgresTenantIsolator struct {
	db     *sql.DB
	prefix string
}

// NewPostgresTenantIsolator creates a new PostgreSQL tenant isolator.
func NewPostgresTenantIsolator(db *sql.DB, schemaPrefix ...string) *PostgresTenantIsolator {
	prefix := "tenant_"
	if len(schemaPrefix) > 0 && schemaPrefix[0] != "" {
		prefix = schemaPrefix[0]
	}
	return &PostgresTenantIsolator{
		db:     db,
		prefix: prefix,
	}
}

// SchemaName returns the full schema name for a tenant.
func (i *PostgresTenantIsolator) SchemaName(slug string) string {
	cleanSlug := strings.ToLower(strings.TrimSpace(slug))
	return fmt.Sprintf("%s%s", i.prefix, cleanSlug)
}

// DBForTenant returns the shared DB connection pool for PostgreSQL.
func (i *PostgresTenantIsolator) DBForTenant(ctx context.Context, slug string) (*sql.DB, error) {
	if err := ValidateTenantSlug(slug); err != nil {
		return nil, err
	}
	return i.db, nil
}

// ProvisionTenant creates the schema for the tenant if it doesn't already exist and initializes migration table.
func (i *PostgresTenantIsolator) ProvisionTenant(ctx context.Context, slug string) error {
	if err := ValidateTenantSlug(slug); err != nil {
		return err
	}

	schemaName := i.SchemaName(slug)
	query := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, schemaName)
	if _, err := i.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create schema %q: %w", schemaName, err)
	}

	// Initialize migration table inside the tenant schema
	initQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s".schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`, schemaName)
	if _, err := i.db.ExecContext(ctx, initQuery); err != nil {
		return fmt.Errorf("failed to init schema_migrations in %q: %w", schemaName, err)
	}

	return nil
}

// DeprovisionTenant drops the tenant schema and all objects within it.
func (i *PostgresTenantIsolator) DeprovisionTenant(ctx context.Context, slug string) error {
	if err := ValidateTenantSlug(slug); err != nil {
		return err
	}

	schemaName := i.SchemaName(slug)
	query := fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE`, schemaName)
	if _, err := i.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to drop schema %q: %w", schemaName, err)
	}

	return nil
}

// Close closes the underlying DB connection.
func (i *PostgresTenantIsolator) Close() error {
	if i.db != nil {
		return i.db.Close()
	}
	return nil
}
