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
	if tenant == "" {
		tenant = "public"
	}

	stats := map[string]interface{}{
		"tenant":          tenant,
		"schema_name":     fmt.Sprintf("tenant_%s", tenant),
		"db_isolation":    "PostgreSQL Schema & SQLite File-Per-Tenant",
		"total_revenue":   "$128,450.00",
		"active_users":    2840,
		"rpc_throughput":  "172,500 req/s",
		"uptime":          "99.99%",
		"server_latency":  "0.6 ms (In-Memory IPC: 0.006 ms)",
		"status":          "Healthy",
		"available_tenants": []string{
			"acme_corp",
			"enterprise_client",
			"tokalink_saas",
			"fintech_ltd",
			"public",
		},
		"recent_activity": []map[string]string{
			{"id": "EVT-101", "event": "RPC Call: UserService/GetProfile", "tenant": tenant, "time": "Just now", "status": "200 OK"},
			{"id": "EVT-102", "event": "Schema Migrated: users_table", "tenant": tenant, "time": "2 mins ago", "status": "Applied"},
			{"id": "EVT-103", "event": "JWT Token Issued", "tenant": tenant, "time": "5 mins ago", "status": "Success"},
			{"id": "EVT-104", "event": "Tenant Isolation Ping", "tenant": tenant, "time": "12 mins ago", "status": "Verified"},
		},
		"team_members": []map[string]string{
			{"name": "Budi Santoso", "email": "budi@" + tenant + ".example.com", "role": "Owner / Admin", "status": "Active"},
			{"name": "Siti Rahma", "email": "siti@" + tenant + ".example.com", "role": "Backend Engineer", "status": "Active"},
			{"name": "Alex Chen", "email": "alex@" + tenant + ".example.com", "role": "Product Manager", "status": "Active"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
