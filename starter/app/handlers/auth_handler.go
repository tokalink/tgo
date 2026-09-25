package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TenantSlug string `json:"tenant_slug"`
}

type UserInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	TenantSlug string `json:"tenant_slug"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json request"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email and password are required"}`, http.StatusBadRequest)
		return
	}

	tenant := strings.TrimSpace(req.TenantSlug)
	if tenant == "" {
		tenant = "acme_corp"
	}

	// Generate standard JWT mock payload for demo
	name := strings.Title(strings.Split(req.Email, "@")[0])
	if name == "" {
		name = "Demo User"
	}

	user := UserInfo{
		ID:         fmt.Sprintf("usr_%d", time.Now().Unix()%10000),
		Name:       name,
		Email:      req.Email,
		Role:       "Administrator",
		TenantSlug: tenant,
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claimsJSON, _ := json.Marshal(map[string]interface{}{
		"sub":         user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"role":        user.Role,
		"tenant":      user.TenantSlug,
		"tenant_slug": user.TenantSlug,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	})
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)
	sig := "tgo_demo_signature_valid"
	token := fmt.Sprintf("%s.%s.%s", header, payload, sig)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  user,
	})
}
