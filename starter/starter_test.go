package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tokalink/tgo/pkg/app"
	"github.com/tokalink/tgo/pkg/transport/middleware"
	"github.com/tokalink/tgo/starter/app/handlers"
	userv1 "github.com/tokalink/tgo/starter/proto/v1/userv1"
	"github.com/tokalink/tgo/starter/proto/v1/userv1/userconnect"
	"google.golang.org/protobuf/proto"
)

func TestStarter_EndToEnd(t *testing.T) {
	// 1. Initialize App Container
	application := app.New()

	// 2. Register ConnectRPC Services
	userHandler := handlers.NewUserHandler()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	application.Server().Register(path, handler)

	ts := httptest.NewServer(application.Server().Handler())
	defer ts.Close()

	// 3. Test REST JSON Call with Tenant Header
	t.Run("REST JSON Request with Tenant Header", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"user_id": "usr_starter_01",
		})

		req, err := http.NewRequest("POST", ts.URL+"/user.v1.UserService/GetProfile", bytes.NewReader(reqBody))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Tenant-Slug", "production_tenant")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("expected 200 OK, got %d: %s", resp.StatusCode, string(body))
		}

		var jsonResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
			t.Fatalf("failed to decode json: %v", err)
		}

		if jsonResp["userId"] != "usr_starter_01" && jsonResp["user_id"] != "usr_starter_01" {
			t.Fatalf("expected userId usr_starter_01, got: %v", jsonResp)
		}
		if jsonResp["tenantSlug"] != "production_tenant" && jsonResp["tenant_slug"] != "production_tenant" {
			t.Fatalf("expected tenant production_tenant, got: %v", jsonResp)
		}
	})

	// 4. Test In-Memory IPC Invoker (Wails Desktop Mode)
	t.Run("In-Memory IPC Bridge Invocation", func(t *testing.T) {
		invoker := application.Server().ServeInMemory(context.Background())

		ctx := middleware.WithTenant(context.Background(), "desktop_offline_tenant")
		req := &userv1.GetProfileRequest{UserId: "usr_wails_99"}
		var resp userv1.GetProfileResponse

		err := invoker.Invoke(ctx, "/user.v1.UserService/GetProfile", req, &resp)
		if err != nil {
			t.Fatalf("in-memory invoke failed: %v", err)
		}

		if resp.GetUserId() != "usr_wails_99" {
			t.Fatalf("expected usr_wails_99, got: %s", resp.GetUserId())
		}
		if resp.GetTenantSlug() != "desktop_offline_tenant" {
			t.Fatalf("expected desktop_offline_tenant, got: %s", resp.GetTenantSlug())
		}
	})

	// 5. Test Binary Protobuf Request
	t.Run("Binary Protobuf Request", func(t *testing.T) {
		protoReq := &userv1.GetProfileRequest{UserId: "usr_proto_88"}
		protoBytes, _ := proto.Marshal(protoReq)

		req, err := http.NewRequest("POST", ts.URL+"/user.v1.UserService/GetProfile", bytes.NewReader(protoBytes))
		if err != nil {
			t.Fatalf("failed to create proto request: %v", err)
		}
		req.Header.Set("Content-Type", "application/proto")
		req.Header.Set("X-Tenant-Slug", "grpc_tenant")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("proto request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
		}

		respBytes, _ := io.ReadAll(resp.Body)
		var protoResp userv1.GetProfileResponse
		if err := proto.Unmarshal(respBytes, &protoResp); err != nil {
			t.Fatalf("failed to unmarshal proto response: %v", err)
		}

		if protoResp.GetUserId() != "usr_proto_88" {
			t.Fatalf("expected usr_proto_88, got: %s", protoResp.GetUserId())
		}
		if protoResp.GetTenantSlug() != "grpc_tenant" {
			t.Fatalf("expected grpc_tenant, got: %s", protoResp.GetTenantSlug())
		}
	})
}
