package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubdomainResolver(t *testing.T) {
	tests := []struct {
		name       string
		rootDomain string
		host       string
		wantTenant string
		wantErr    bool
	}{
		{
			name:       "valid subdomain with configured root domain",
			rootDomain: "example.com",
			host:       "acme.example.com",
			wantTenant: "acme",
			wantErr:    false,
		},
		{
			name:       "valid subdomain with port",
			rootDomain: "example.com",
			host:       "tenant-123.example.com:8080",
			wantTenant: "tenant-123",
			wantErr:    false,
		},
		{
			name:       "root domain itself should fail",
			rootDomain: "example.com",
			host:       "example.com",
			wantTenant: "",
			wantErr:    true,
		},
		{
			name:       "www subdomain should fail",
			rootDomain: "example.com",
			host:       "www.example.com",
			wantTenant: "",
			wantErr:    true,
		},
		{
			name:       "unrelated domain should fail",
			rootDomain: "example.com",
			host:       "acme.other.org",
			wantTenant: "",
			wantErr:    true,
		},
		{
			name:       "auto fallback without root domain",
			rootDomain: "",
			host:       "myorg.app.local",
			wantTenant: "myorg",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewSubdomainResolver(tt.rootDomain)
			req := httptest.NewRequest("GET", "http://"+tt.host+"/api", nil)
			req.Host = tt.host

			tenant, err := resolver.Resolve(req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if tenant != tt.wantTenant {
				t.Fatalf("expected tenant: %s, got: %s", tt.wantTenant, tenant)
			}
		})
	}
}

func TestHeaderResolver(t *testing.T) {
	resolver := NewHeaderResolver("X-Tenant-Slug", "X-Tenant-ID")

	t.Run("resolves primary header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Set("X-Tenant-Slug", "acme-corp")

		tenant, err := resolver.Resolve(req)
		if err != nil || tenant != "acme-corp" {
			t.Fatalf("expected acme-corp, got tenant: %s, err: %v", tenant, err)
		}
	})

	t.Run("resolves secondary header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Set("X-Tenant-ID", "tenant-999")

		tenant, err := resolver.Resolve(req)
		if err != nil || tenant != "tenant-999" {
			t.Fatalf("expected tenant-999, got tenant: %s, err: %v", tenant, err)
		}
	})

	t.Run("fails when headers missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		_, err := resolver.Resolve(req)
		if err == nil {
			t.Fatalf("expected error when headers missing")
		}
	})
}

func TestJWTClaimResolver(t *testing.T) {
	resolver := NewJWTClaimResolver()

	t.Run("resolves claim from valid JWT", func(t *testing.T) {
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"123","tenant":"tenant-jwt-org","exp":9999999999}`))
		sig := "fakesignature"
		rawToken := fmt.Sprintf("%s.%s.%s", header, payload, sig)

		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Set("Authorization", "Bearer "+rawToken)

		tenant, err := resolver.Resolve(req)
		if err != nil || tenant != "tenant-jwt-org" {
			t.Fatalf("expected tenant-jwt-org, got: %s, err: %v", tenant, err)
		}
	})

	t.Run("fails on missing auth header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		_, err := resolver.Resolve(req)
		if err == nil {
			t.Fatalf("expected error on missing auth header")
		}
	})
}

func TestChainResolverAndMiddleware(t *testing.T) {
	chain := NewChainResolver(
		NewHeaderResolver("X-Tenant-Slug"),
		NewDefaultTenantResolver("public"),
	)

	var capturedTenant string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedTenant = GetTenant(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	mw := TenantMiddleware(TenantMiddlewareOptions{
		Resolver: chain,
	})

	// Test 1: Header takes precedence
	req1 := httptest.NewRequest("GET", "/api", nil)
	req1.Header.Set("X-Tenant-Slug", "custom-org")
	w1 := httptest.NewRecorder()
	mw(handler).ServeHTTP(w1, req1)

	if capturedTenant != "custom-org" {
		t.Fatalf("expected custom-org, got: %s", capturedTenant)
	}

	// Test 2: Fallback to default tenant
	req2 := httptest.NewRequest("GET", "/api", nil)
	w2 := httptest.NewRecorder()
	mw(handler).ServeHTTP(w2, req2)

	if capturedTenant != "public" {
		t.Fatalf("expected public fallback, got: %s", capturedTenant)
	}
}
