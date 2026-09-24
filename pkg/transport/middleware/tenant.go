package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type contextKey string

const (
	// TenantContextKey is the context key for storing the active tenant slug.
	TenantContextKey contextKey = "tgo_tenant_slug"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrInvalidTenant  = errors.New("invalid tenant identifier")
)

// WithTenant returns a new context with the tenant slug attached.
func WithTenant(ctx context.Context, slug string) context.Context {
	return context.WithValue(ctx, TenantContextKey, slug)
}

// GetTenant retrieves the tenant slug from context, returning empty string if not found.
func GetTenant(ctx context.Context) string {
	if val, ok := ctx.Value(TenantContextKey).(string); ok {
		return val
	}
	return ""
}

// TenantFromContext returns the tenant slug and a boolean indicating if it was present.
func TenantFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(TenantContextKey).(string)
	return val, ok && val != ""
}

// TenantResolver defines the contract for resolving a tenant from an incoming HTTP request.
type TenantResolver interface {
	Resolve(r *http.Request) (string, error)
}

// SubdomainResolver extracts the tenant slug from the Host header based on a root domain.
// Example: If RootDomain is "app.com" and Host is "client1.app.com:8080", resolved tenant is "client1".
type SubdomainResolver struct {
	RootDomain string
}

func NewSubdomainResolver(rootDomain string) *SubdomainResolver {
	return &SubdomainResolver{
		RootDomain: strings.ToLower(strings.TrimSpace(rootDomain)),
	}
}

func (s *SubdomainResolver) Resolve(r *http.Request) (string, error) {
	host := r.Host
	if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}
	host = strings.ToLower(strings.TrimSpace(host))

	if host == "localhost" || net.ParseIP(host) != nil {
		return "", ErrTenantNotFound
	}

	if s.RootDomain == "" {
		// If no root domain configured, try standard 3-part hostname (subdomain.domain.tld)
		parts := strings.Split(host, ".")
		if len(parts) >= 3 && parts[0] != "www" {
			return parts[0], nil
		}
		return "", ErrTenantNotFound
	}

	root := strings.TrimPrefix(s.RootDomain, ".")
	if !strings.HasSuffix(host, root) || host == root {
		return "", ErrTenantNotFound
	}

	sub := strings.TrimSuffix(host, "."+root)
	if sub == "" || sub == "www" || sub == host {
		return "", ErrTenantNotFound
	}

	// In case of multiple nested subdomains, take the leftmost part
	parts := strings.Split(sub, ".")
	return parts[0], nil
}

// HeaderResolver extracts the tenant slug from custom HTTP headers (e.g. X-Tenant-Slug or X-Tenant-ID).
type HeaderResolver struct {
	HeaderNames []string
}

func NewHeaderResolver(headerNames ...string) *HeaderResolver {
	if len(headerNames) == 0 {
		headerNames = []string{"X-Tenant-Slug", "X-Tenant-ID"}
	}
	return &HeaderResolver{HeaderNames: headerNames}
}

func (h *HeaderResolver) Resolve(r *http.Request) (string, error) {
	for _, name := range h.HeaderNames {
		if val := r.Header.Get(name); strings.TrimSpace(val) != "" {
			return strings.TrimSpace(val), nil
		}
	}
	return "", ErrTenantNotFound
}

// JWTClaimResolver extracts tenant from a JWT token claim in the Authorization header.
type JWTClaimResolver struct {
	ClaimKeys []string
}

func NewJWTClaimResolver(claimKeys ...string) *JWTClaimResolver {
	if len(claimKeys) == 0 {
		claimKeys = []string{"tenant", "tenant_slug", "tenant_id"}
	}
	return &JWTClaimResolver{ClaimKeys: claimKeys}
}

func (j *JWTClaimResolver) Resolve(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrTenantNotFound
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrTenantNotFound
	}

	tokenStr := strings.TrimSpace(parts[1])
	jwtParts := strings.Split(tokenStr, ".")
	if len(jwtParts) < 2 {
		return "", ErrTenantNotFound
	}

	// Payload is the second segment (base64url encoded)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		// Fallback to standard URL or StdEncoding if padded
		payloadBytes, err = base64.URLEncoding.DecodeString(jwtParts[1])
		if err != nil {
			return "", fmt.Errorf("%w: invalid jwt payload base64", ErrTenantNotFound)
		}
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", fmt.Errorf("%w: invalid jwt claims json", ErrTenantNotFound)
	}

	for _, key := range j.ClaimKeys {
		if val, exists := claims[key]; exists {
			if strVal, ok := val.(string); ok && strings.TrimSpace(strVal) != "" {
				return strings.TrimSpace(strVal), nil
			}
		}
	}

	return "", ErrTenantNotFound
}

// DefaultTenantResolver provides a static fallback tenant if configured.
type DefaultTenantResolver struct {
	DefaultTenant string
}

func NewDefaultTenantResolver(defaultTenant string) *DefaultTenantResolver {
	return &DefaultTenantResolver{DefaultTenant: defaultTenant}
}

func (d *DefaultTenantResolver) Resolve(r *http.Request) (string, error) {
	if d.DefaultTenant != "" {
		return d.DefaultTenant, nil
	}
	return "", ErrTenantNotFound
}

// ChainResolver evaluates multiple resolvers in priority order until one succeeds.
type ChainResolver struct {
	Resolvers []TenantResolver
}

func NewChainResolver(resolvers ...TenantResolver) *ChainResolver {
	return &ChainResolver{Resolvers: resolvers}
}

func (c *ChainResolver) Resolve(r *http.Request) (string, error) {
	for _, resolver := range c.Resolvers {
		slug, err := resolver.Resolve(r)
		if err == nil && slug != "" {
			return slug, nil
		}
	}
	return "", ErrTenantNotFound
}

// TenantMiddlewareOptions defines configuration options for TenantMiddleware.
type TenantMiddlewareOptions struct {
	Resolver        TenantResolver
	RequireTenant   bool
	OnMissingTenant func(w http.ResponseWriter, r *http.Request)
}

// TenantMiddleware creates a standard HTTP middleware that resolves the tenant and injects it into the request context.
func TenantMiddleware(opts TenantMiddlewareOptions) func(http.Handler) http.Handler {
	if opts.Resolver == nil {
		opts.Resolver = NewChainResolver(
			NewSubdomainResolver(""),
			NewHeaderResolver(),
			NewJWTClaimResolver(),
		)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slug, err := opts.Resolver.Resolve(r)
			if err != nil || slug == "" {
				if opts.RequireTenant {
					if opts.OnMissingTenant != nil {
						opts.OnMissingTenant(w, r)
					} else {
						http.Error(w, `{"error": "tenant required or not found"}`, http.StatusNotFound)
					}
					return
				}
				// If not strictly required, continue with original context
				next.ServeHTTP(w, r)
				return
			}

			ctx := WithTenant(r.Context(), slug)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
