package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteTenantIsolator manages file-per-tenant SQLite database isolation.
type SQLiteTenantIsolator struct {
	dataDir string
	dbs     map[string]*sql.DB
	mu      sync.RWMutex
}

// NewSQLiteTenantIsolator initializes a tenant isolator with the specified storage directory.
func NewSQLiteTenantIsolator(dataDir string) (*SQLiteTenantIsolator, error) {
	if dataDir != ":memory:" {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create tenant data dir: %w", err)
		}
	}
	return &SQLiteTenantIsolator{
		dataDir: dataDir,
		dbs:     make(map[string]*sql.DB),
	}, nil
}

// tenantFilePath returns the file path for a tenant database.
func (i *SQLiteTenantIsolator) tenantFilePath(slug string) string {
	if i.dataDir == ":memory:" {
		return fmt.Sprintf("file:%s?mode=memory&cache=shared", slug)
	}
	return filepath.Join(i.dataDir, fmt.Sprintf("%s.db", slug))
}

// DBForTenant returns a pooled connection for the specified tenant, opening it if not already cached.
func (i *SQLiteTenantIsolator) DBForTenant(ctx context.Context, slug string) (*sql.DB, error) {
	if err := ValidateTenantSlug(slug); err != nil {
		return nil, err
	}

	i.mu.RLock()
	db, ok := i.dbs[slug]
	i.mu.RUnlock()
	if ok {
		return db, nil
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check after lock
	if db, ok := i.dbs[slug]; ok {
		return db, nil
	}

	dbPath := i.tenantFilePath(slug)
	newDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tenant db %q: %w", slug, err)
	}

	if err := newDB.PingContext(ctx); err != nil {
		newDB.Close()
		return nil, fmt.Errorf("failed to ping tenant db %q: %w", slug, err)
	}

	i.dbs[slug] = newDB
	return newDB, nil
}

// ProvisionTenant creates the tenant database and initializes default table structure.
func (i *SQLiteTenantIsolator) ProvisionTenant(ctx context.Context, slug string) error {
	if err := ValidateTenantSlug(slug); err != nil {
		return err
	}

	db, err := i.DBForTenant(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to provision tenant %q: %w", slug, err)
	}

	// Initialize migration table or basic structure
	migrator := NewMigrator(&SQLiteAdapter{db: db})
	if err := migrator.InitSchema(ctx); err != nil {
		return fmt.Errorf("failed to initialize schema for tenant %q: %w", slug, err)
	}

	return nil
}

// DeprovisionTenant closes active connections and removes the tenant database file.
func (i *SQLiteTenantIsolator) DeprovisionTenant(ctx context.Context, slug string) error {
	if err := ValidateTenantSlug(slug); err != nil {
		return err
	}

	i.mu.Lock()
	if db, ok := i.dbs[slug]; ok {
		_ = db.Close()
		delete(i.dbs, slug)
	}
	i.mu.Unlock()

	if i.dataDir != ":memory:" {
		dbPath := i.tenantFilePath(slug)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete tenant db file %q: %w", dbPath, err)
		}
	}

	return nil
}

// Close closes all open tenant connections.
func (i *SQLiteTenantIsolator) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	var firstErr error
	for slug, db := range i.dbs {
		if err := db.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("error closing tenant db %q: %w", slug, err)
		}
	}
	i.dbs = make(map[string]*sql.DB)
	return firstErr
}
