package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrTenantSlugInvalid = errors.New("tenant slug must contain only alphanumeric characters, underscores, or hyphens")
	ErrTenantNotFound    = errors.New("tenant not found")
	ErrTenantAlreadyExists = errors.New("tenant already exists")
)

var validSlugRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateTenantSlug ensures the tenant slug contains only safe characters.
func ValidateTenantSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("%w: slug cannot be empty", ErrTenantSlugInvalid)
	}
	if !validSlugRegex.MatchString(slug) {
		return fmt.Errorf("%w: %q", ErrTenantSlugInvalid, slug)
	}
	return nil
}

// TenantIsolator defines the contract for managing multi-tenant database isolation and lifecycle.
type TenantIsolator interface {
	// DBForTenant returns the database connection handle for the given tenant slug.
	DBForTenant(ctx context.Context, slug string) (*sql.DB, error)

	// ProvisionTenant provisions any required resources (e.g. schema or db file) for a new tenant.
	ProvisionTenant(ctx context.Context, slug string) error

	// DeprovisionTenant tears down resources for a tenant.
	DeprovisionTenant(ctx context.Context, slug string) error

	// Close cleans up connections held by the isolator.
	Close() error
}
