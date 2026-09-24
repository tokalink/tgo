package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tgo-framework/tgo/pkg/transport/middleware"
)

func TestTenantSlugValidation(t *testing.T) {
	validSlugs := []string{"acme", "tenant-123", "company_org", "clientA"}
	for _, slug := range validSlugs {
		if err := ValidateTenantSlug(slug); err != nil {
			t.Errorf("expected slug %q to be valid, got: %v", slug, err)
		}
	}

	invalidSlugs := []string{"", "acme corp", "tenant;drop", "test/slug", "abc$123"}
	for _, slug := range invalidSlugs {
		if err := ValidateTenantSlug(slug); err == nil {
			t.Errorf("expected slug %q to be invalid, but validation passed", slug)
		}
	}
}

func TestSQLiteTenantIsolator_CrossTenantIsolation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tgo-tenant-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	isolator, err := NewSQLiteTenantIsolator(tmpDir)
	if err != nil {
		t.Fatalf("failed to create sqlite tenant isolator: %v", err)
	}
	defer isolator.Close()

	ctx := context.Background()
	tenantA := "tenant_a"
	tenantB := "tenant_b"

	// 1. Provision both tenants
	if err := isolator.ProvisionTenant(ctx, tenantA); err != nil {
		t.Fatalf("failed to provision tenant A: %v", err)
	}
	if err := isolator.ProvisionTenant(ctx, tenantB); err != nil {
		t.Fatalf("failed to provision tenant B: %v", err)
	}

	// 2. Setup user tables in both tenant databases
	dbA, err := isolator.DBForTenant(ctx, tenantA)
	if err != nil {
		t.Fatalf("failed to get DB for tenant A: %v", err)
	}
	_, err = dbA.ExecContext(ctx, "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);")
	if err != nil {
		t.Fatalf("failed to create table in tenant A: %v", err)
	}

	dbB, err := isolator.DBForTenant(ctx, tenantB)
	if err != nil {
		t.Fatalf("failed to get DB for tenant B: %v", err)
	}
	_, err = dbB.ExecContext(ctx, "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);")
	if err != nil {
		t.Fatalf("failed to create table in tenant B: %v", err)
	}

	// 3. Write record into Tenant A
	_, err = dbA.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (1, 'Alice from Tenant A');")
	if err != nil {
		t.Fatalf("failed to insert user in tenant A: %v", err)
	}

	// 4. Write different record with SAME ID into Tenant B
	_, err = dbB.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (1, 'Bob from Tenant B');")
	if err != nil {
		t.Fatalf("failed to insert user in tenant B: %v", err)
	}

	// 5. Verify Tenant A data
	var nameA string
	err = dbA.QueryRowContext(ctx, "SELECT name FROM users WHERE id = 1").Scan(&nameA)
	if err != nil || nameA != "Alice from Tenant A" {
		t.Fatalf("tenant A data corrupted: expected 'Alice from Tenant A', got %q, err: %v", nameA, err)
	}

	// 6. Verify Tenant B data
	var nameB string
	err = dbB.QueryRowContext(ctx, "SELECT name FROM users WHERE id = 1").Scan(&nameB)
	if err != nil || nameB != "Bob from Tenant B" {
		t.Fatalf("tenant B data corrupted: expected 'Bob from Tenant B', got %q, err: %v", nameB, err)
	}

	// 7. Verify file presence on disk
	fileA := filepath.Join(tmpDir, tenantA+".db")
	if _, err := os.Stat(fileA); os.IsNotExist(err) {
		t.Fatalf("expected tenant A db file to exist at %q", fileA)
	}

	// 8. Test Deprovisioning
	if err := isolator.DeprovisionTenant(ctx, tenantA); err != nil {
		t.Fatalf("failed to deprovision tenant A: %v", err)
	}
	if _, err := os.Stat(fileA); !os.IsNotExist(err) {
		t.Fatalf("expected tenant A db file to be removed after deprovisioning")
	}
}

func TestDBEngine_FromContext_Resolution(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tgo-dbengine-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	isolator, err := NewSQLiteTenantIsolator(tmpDir)
	if err != nil {
		t.Fatalf("failed to create sqlite tenant isolator: %v", err)
	}
	defer isolator.Close()

	primaryAdapter, err := NewSQLite(filepath.Join(tmpDir, "primary.db"))
	if err != nil {
		t.Fatalf("failed to create primary sqlite adapter: %v", err)
	}
	defer primaryAdapter.Close()

	engine := primaryAdapter.WithTenantIsolator(isolator)

	ctx := context.Background()

	// Provision tenant "org_1"
	if err := isolator.ProvisionTenant(ctx, "org_1"); err != nil {
		t.Fatalf("failed to provision org_1: %v", err)
	}

	// Context without tenant should return the primary engine
	engNoTenant, err := engine.FromContext(ctx)
	if err != nil {
		t.Fatalf("failed FromContext with no tenant: %v", err)
	}
	if engNoTenant.Conn() != primaryAdapter.Conn() {
		t.Fatalf("expected primary connection when no tenant in context")
	}

	// Context with tenant should return tenant-scoped engine
	tenantCtx := middleware.WithTenant(ctx, "org_1")
	engTenant, err := engine.FromContext(tenantCtx)
	if err != nil {
		t.Fatalf("failed FromContext with tenant: %v", err)
	}
	if engTenant.Conn() == primaryAdapter.Conn() {
		t.Fatalf("expected separate tenant-scoped connection")
	}
}

func TestPostgresTenantIsolator_SchemaNaming(t *testing.T) {
	isolator := NewPostgresTenantIsolator(nil, "tenant_")
	if name := isolator.SchemaName("AcmeCorp"); name != "tenant_acmecorp" {
		t.Fatalf("expected schema tenant_acmecorp, got %s", name)
	}
}
