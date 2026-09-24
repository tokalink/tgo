package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tgo-framework/tgo/pkg/transport/middleware"
	"github.com/tgo-framework/tgo/starter/app/handlers"
	userv1 "github.com/tgo-framework/tgo/starter/proto/v1/userv1"
	"github.com/tgo-framework/tgo/starter/proto/v1/userv1/userconnect"
	"google.golang.org/protobuf/proto"
)

func setupTestServer() (Server, http.Handler) {
	srv := NewServer()
	srv.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORS(),
		middleware.TenantMiddleware(middleware.TenantMiddlewareOptions{}),
	)

	userHandler := handlers.NewUserHandler()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	srv.Register(path, handler)

	return srv, srv.Handler()
}

func TestServer_JSON_Request(t *testing.T) {
	_, handler := setupTestServer()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"user_id": "usr_101",
	})

	req, err := http.NewRequest("POST", ts.URL+"/user.v1.UserService/GetProfile", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Slug", "alpha_corp")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(body))
	}

	var jsonResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	if jsonResp["userId"] != "usr_101" && jsonResp["user_id"] != "usr_101" {
		t.Fatalf("unexpected user_id: %v", jsonResp)
	}
	if jsonResp["tenantSlug"] != "alpha_corp" && jsonResp["tenant_slug"] != "alpha_corp" {
		t.Fatalf("unexpected tenant_slug: %v", jsonResp)
	}
}

func TestServer_Proto_Request(t *testing.T) {
	_, handler := setupTestServer()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	protoReq := &userv1.GetProfileRequest{UserId: "usr_202"}
	protoBytes, err := proto.Marshal(protoReq)
	if err != nil {
		t.Fatalf("failed to marshal proto: %v", err)
	}

	req, err := http.NewRequest("POST", ts.URL+"/user.v1.UserService/GetProfile", bytes.NewReader(protoBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/proto")
	req.Header.Set("X-Tenant-Slug", "beta_org")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(body))
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	var protoResp userv1.GetProfileResponse
	if err := proto.Unmarshal(respBytes, &protoResp); err != nil {
		t.Fatalf("failed to unmarshal proto response: %v", err)
	}

	if protoResp.GetUserId() != "usr_202" {
		t.Fatalf("expected usr_202, got %s", protoResp.GetUserId())
	}
	if protoResp.GetTenantSlug() != "beta_org" {
		t.Fatalf("expected beta_org, got %s", protoResp.GetTenantSlug())
	}
}

func TestServer_InMemoryInvoker(t *testing.T) {
	srv, _ := setupTestServer()
	invoker := srv.ServeInMemory(context.Background())

	ctx := middleware.WithTenant(context.Background(), "inmemory_tenant")
	req := &userv1.GetProfileRequest{UserId: "usr_303"}
	var resp userv1.GetProfileResponse

	err := invoker.Invoke(ctx, "/user.v1.UserService/GetProfile", req, &resp)
	if err != nil {
		t.Fatalf("in-memory invoke failed: %v", err)
	}

	if resp.GetUserId() != "usr_303" {
		t.Fatalf("expected usr_303, got %s", resp.GetUserId())
	}
	if resp.GetTenantSlug() != "inmemory_tenant" {
		t.Fatalf("expected inmemory_tenant, got %s", resp.GetTenantSlug())
	}
}

func TestServer_InMemoryInvoker_JSON(t *testing.T) {
	srv, _ := setupTestServer()
	invoker := srv.ServeInMemory(context.Background())

	ctx := middleware.WithTenant(context.Background(), "json_tenant")
	rawJSON := []byte(`{"user_id": "usr_404"}`)

	respBytes, err := invoker.InvokeJSON(ctx, "/user.v1.UserService/GetProfile", rawJSON)
	if err != nil {
		t.Fatalf("in-memory json invoke failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	if parsed["userId"] != "usr_404" && parsed["user_id"] != "usr_404" {
		t.Fatalf("expected usr_404, got %v", parsed)
	}
}

func TestServer_WelcomePage(t *testing.T) {
	_, handler := setupTestServer()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("failed to fetch welcome page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	bodyStr := string(body)
	if !bytes.Contains(body, []byte("TGo Framework")) {
		t.Fatalf("expected welcome page to contain 'TGo Framework', got: %s", bodyStr)
	}
}
