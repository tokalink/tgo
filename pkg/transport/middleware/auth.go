package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type userContextKey string

const (
	UserClaimsContextKey userContextKey = "tgo_user_claims"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// UserClaims holds the authenticated user payload.
type UserClaims struct {
	Subject     string   `json:"sub"`
	Tenant      string   `json:"tenant,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Email       string   `json:"email,omitempty"`
	Name        string   `json:"name,omitempty"`
}

// WithUserClaims injects UserClaims into the context.
func WithUserClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, UserClaimsContextKey, claims)
}

// UserClaimsFromContext extracts UserClaims from the context.
func UserClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*UserClaims)
	return claims, ok && claims != nil
}

// AuthOptions configures the Auth middleware.
type AuthOptions struct {
	Validator func(token string) (*UserClaims, error)
}

// RequireAuth creates a middleware that enforces presence and validity of an authentication token.
func RequireAuth(opts ...AuthOptions) func(http.Handler) http.Handler {
	var validator func(token string) (*UserClaims, error)
	if len(opts) > 0 && opts[0].Validator != nil {
		validator = opts[0].Validator
	} else {
		// Default parser decodes JWT claims payload
		validator = defaultJWTValidator
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"invalid authorization format, expected Bearer <token>"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimSpace(parts[1])
			claims, err := validator(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := WithUserClaims(r.Context(), claims)
			// If tenant is in claims and not already in context, set it
			if claims.Tenant != "" && GetTenant(ctx) == "" {
				ctx = WithTenant(ctx, claims.Tenant)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func defaultJWTValidator(token string) (*UserClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, ErrUnauthorized
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		payloadBytes, err = base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, ErrUnauthorized
		}
	}

	var claims UserClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrUnauthorized
	}

	if claims.Subject == "" {
		return nil, ErrUnauthorized
	}

	return &claims, nil
}
