package connect

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tgo-framework/tgo/starter/app/handlers"
	userv1 "github.com/tgo-framework/tgo/starter/proto/v1/userv1"
	"github.com/tgo-framework/tgo/starter/proto/v1/userv1/userconnect"
)

func BenchmarkInMemoryInvoker_Proto(b *testing.B) {
	srv := NewServer()
	userHandler := handlers.NewUserHandler()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	srv.Register(path, handler)

	invoker := srv.ServeInMemory(context.Background())
	ctx := context.Background()
	req := &userv1.GetProfileRequest{UserId: "usr_bench"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var resp userv1.GetProfileResponse
		if err := invoker.Invoke(ctx, "/user.v1.UserService/GetProfile", req, &resp); err != nil {
			b.Fatalf("invoke failed: %v", err)
		}
	}
}

func BenchmarkInMemoryInvoker_JSON(b *testing.B) {
	srv := NewServer()
	userHandler := handlers.NewUserHandler()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	srv.Register(path, handler)

	invoker := srv.ServeInMemory(context.Background())
	ctx := context.Background()
	reqJSON := []byte(`{"user_id":"usr_bench"}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := invoker.InvokeJSON(ctx, "/user.v1.UserService/GetProfile", reqJSON); err != nil {
			b.Fatalf("invoke json failed: %v", err)
		}
	}
}

func BenchmarkHTTPServer_JSON(b *testing.B) {
	srv := NewServer()
	userHandler := handlers.NewUserHandler()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	srv.Register(path, handler)

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	client := ts.Client()
	url := ts.URL + "/user.v1.UserService/GetProfile"
	reqJSON := []byte(`{"user_id":"usr_bench"}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", url, bytes.NewReader(reqJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("http request failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
