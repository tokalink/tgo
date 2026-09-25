package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tgo-framework/tgo/pkg/transport/middleware"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.GetTenant(r.Context())
	isSingleTenant := tenant == "public" || tenant == "single" || tenant == "none" || tenant == ""

	schemaName := fmt.Sprintf("tenant_%s", tenant)
	dbIsolation := "PostgreSQL Schema & SQLite File-Per-Tenant"
	modeName := "Multi-Tenant (SaaS Isolation)"
	if isSingleTenant {
		schemaName = "public (Default Standard Schema)"
		dbIsolation = "Standard Single Database (Monolith / Desktop Mode)"
		modeName = "Single-Tenant (Standard App)"
		tenant = "public"
	}

	stats := map[string]interface{}{
		"mode":            modeName,
		"is_single_tenant": isSingleTenant,
		"tenant":          tenant,
		"schema_name":     schemaName,
		"db_isolation":    dbIsolation,
		"total_revenue":   "$128,450.00",
		"active_users":    2840,
		"rpc_throughput":  "172,500 req/s",
		"uptime":          "99.99%",
		"server_latency":  "0.4 ms (In-Memory IPC: 0.006 ms)",
		"status":          "Healthy",
		"available_tenants": []string{
			"acme_corp",
			"enterprise_client",
			"tokalink_saas",
			"fintech_ltd",
			"public",
		},
		"recent_activity": []map[string]string{
			{"id": "EVT-101", "event": "ConnectRPC Call: UserService/GetProfile", "tenant": tenant, "time": "Just now", "status": "200 OK"},
			{"id": "EVT-102", "event": "DB Query executed on " + schemaName, "tenant": tenant, "time": "2 mins ago", "status": "Success"},
			{"id": "EVT-103", "event": "JWT Session Verified", "tenant": tenant, "time": "5 mins ago", "status": "Success"},
			{"id": "EVT-104", "event": "HTTP/2 Health Check", "tenant": tenant, "time": "12 mins ago", "status": "Healthy"},
		},
		"team_members": []map[string]string{
			{"name": "Budi Santoso", "email": "budi@example.com", "role": "Super Admin", "status": "Active"},
			{"name": "Siti Rahma", "email": "siti@example.com", "role": "Lead Architect", "status": "Active"},
			{"name": "Alex Chen", "email": "alex@example.com", "role": "Software Engineer", "status": "Active"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
