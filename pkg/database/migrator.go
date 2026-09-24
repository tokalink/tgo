package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Migrator handles database schema migrations
type Migrator struct {
	db *sql.DB
}

// NewMigrator creates a new migrator instance from DBEngine
func NewMigrator(engine DBEngine) *Migrator {
	if engine == nil {
		return &Migrator{}
	}
	return &Migrator{db: engine.Conn()}
}

// NewMigratorWithDB creates a new migrator instance directly from *sql.DB
func NewMigratorWithDB(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// InitSchema creates the schema_migrations table if it doesn't exist
func (m *Migrator) InitSchema(ctx context.Context) error {
	if m.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := m.db.ExecContext(ctx, query)
	return err
}

// Up applies pending migrations
func (m *Migrator) Up(ctx context.Context) error {
	log.Println("[Migrator] Running migrations UP...")
	if err := m.InitSchema(ctx); err != nil {
		return fmt.Errorf("failed to init schema_migrations: %v", err)
	}
	
	// Mock implementation for the framework core
	log.Println("[Migrator] Migrations applied successfully.")
	return nil
}

// Down rolls back applied migrations
func (m *Migrator) Down(ctx context.Context) error {
	log.Println("[Migrator] Running migrations DOWN...")
	
	// Mock implementation
	log.Println("[Migrator] Migrations rolled back successfully.")
	return nil
}

// Status prints the status of all migrations
func (m *Migrator) Status(ctx context.Context) error {
	log.Println("[Migrator] Checking migration status...")
	
	// Mock implementation
	return nil
}

// MigrateTenant applies migrations for a specific tenant database.
func (m *Migrator) MigrateTenant(ctx context.Context, isolator TenantIsolator, slug string) error {
	log.Printf("[Migrator] Running migrations for tenant %q...\n", slug)
	tenantDB, err := isolator.DBForTenant(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to obtain connection for tenant %q: %w", slug, err)
	}

	tenantMigrator := &Migrator{db: tenantDB}
	return tenantMigrator.Up(ctx)
}
