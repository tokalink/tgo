package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	makeModelWithMigration bool
)

var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Scaffold models, migrations, services, handlers, and frontend boilerplate",
}

// craft make:model <Name> [-m / --migration]
var MakeModelCmd = &cobra.Command{
	Use:   "model <Name>",
	Short: "Generate a new model struct and repository (use -m to also generate migration)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		pascalName := toPascalCase(name)
		snakeName := toSnakeCase(name)
		targetDir := "app/models"
		_ = os.MkdirAll(targetDir, 0755)

		filePath := filepath.Join(targetDir, fmt.Sprintf("%s.go", snakeName))
		content := fmt.Sprintf(`package models

import (
	"context"
	"time"

	"github.com/tokalink/tgo/pkg/database"
)

// %s represents the entity model
type %s struct {
	ID        string    `+"`json:\"id\" db:\"id\"`"+`
	Name      string    `+"`json:\"name\" db:\"name\"`"+`
	CreatedAt time.Time `+"`json:\"created_at\" db:\"created_at\"`"+`
	UpdatedAt time.Time `+"`json:\"updated_at\" db:\"updated_at\"`"+`
}

// %sRepository handles database queries for %s
type %sRepository struct {
	db database.DBEngine
}

func New%sRepository(db database.DBEngine) *%sRepository {
	return &%sRepository{db: db}
}

func (r *%sRepository) FindByID(ctx context.Context, id string) (*%s, error) {
	scopedDB, err := r.db.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	row := scopedDB.Conn().QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM %s WHERE id = ?", id)
	var m %s
	if err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
`, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, snakeName+"s", pascalName)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Fatalf("[Craft] Failed to write model file: %v", err)
		}
		log.Printf("[Craft] Model created: %s\n", filePath)

		// If -m / --migration is provided, generate corresponding migration files
		if makeModelWithMigration {
			tableName := snakeName + "s"
			migrationName := fmt.Sprintf("create_%s_table", tableName)
			createMigration(migrationName, tableName)
		}
	},
}

// craft make:migration <Name>
var MakeMigrationCmd = &cobra.Command{
	Use:   "migration <name>",
	Short: "Generate a new SQL migration up and down pair",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := toSnakeCase(args[0])
		tableName := strings.TrimSuffix(strings.TrimPrefix(name, "create_"), "_table")
		createMigration(name, tableName)
	},
}

func createMigration(name, tableName string) {
	targetDir := filepath.Join("database", "migrations")
	_ = os.MkdirAll(targetDir, 0755)

	ts := time.Now().Format("20060102150405")
	upFileName := fmt.Sprintf("%s_%s.up.sql", ts, name)
	downFileName := fmt.Sprintf("%s_%s.down.sql", ts, name)

	upPath := filepath.Join(targetDir, upFileName)
	downPath := filepath.Join(targetDir, downFileName)

	upContent := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`, tableName)

	downContent := fmt.Sprintf(`DROP TABLE IF EXISTS %s;
`, tableName)

	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		log.Fatalf("[Craft] Failed to write migration UP file: %v", err)
	}
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		log.Fatalf("[Craft] Failed to write migration DOWN file: %v", err)
	}

	log.Printf("[Craft] Migration created: %s\n", upPath)
	log.Printf("[Craft] Migration created: %s\n", downPath)
}

// craft make:service <Name>
var MakeServiceCmd = &cobra.Command{
	Use:   "service <Name>",
	Short: "Generate a new business logic service action",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		pascalName := toPascalCase(name)
		snakeName := toSnakeCase(name)
		targetDir := "app/actions"
		_ = os.MkdirAll(targetDir, 0755)

		filePath := filepath.Join(targetDir, fmt.Sprintf("%s_action.go", snakeName))
		content := fmt.Sprintf(`package actions

import (
	"context"
	"fmt"

	"github.com/tokalink/tgo/pkg/database"
)

// %sAction encapsulates domain business logic
type %sAction struct {
	db database.DBEngine
}

func New%sAction(db database.DBEngine) *%sAction {
	return &%sAction{db: db}
}

func (a *%sAction) Execute(ctx context.Context, id string) error {
	scopedDB, err := a.db.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("resolving tenant db: %%w", err)
	}
	_ = scopedDB
	return nil
}
`, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Fatalf("[Craft] Failed to write service action file: %v", err)
		}
		log.Printf("[Craft] Service created: %s\n", filePath)
	},
}

// craft make:handler <Name>
var MakeHandlerCmd = &cobra.Command{
	Use:   "handler <Name>",
	Short: "Generate a new RPC & HTTP controller handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		pascalName := toPascalCase(name)
		snakeName := toSnakeCase(name)
		targetDir := "app/handlers"
		_ = os.MkdirAll(targetDir, 0755)

		filePath := filepath.Join(targetDir, fmt.Sprintf("%s_handler.go", snakeName))
		content := fmt.Sprintf(`package handlers

import (
	"net/http"
)

// %sHandler handles incoming transport requests
type %sHandler struct{}

func New%sHandler() *%sHandler {
	return &%sHandler{}
}

func (h *%sHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`+"`{\"message\": \"Hello from %sHandler\"}`"+`))
}
`, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName, pascalName)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Fatalf("[Craft] Failed to write handler file: %v", err)
		}
		log.Printf("[Craft] Handler created: %s\n", filePath)
	},
}

// craft make:frontend
var MakeFrontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Generate unified frontend starter (HTML/CSS/JS + Wails ready)",
	Run: func(cmd *cobra.Command, args []string) {
		frontendDir := "frontend"
		_ = os.MkdirAll(filepath.Join(frontendDir, "src"), 0755)

		indexHTML := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>TGo App</title>
  <link rel="stylesheet" href="style.css" />
</head>
<body>
  <div id="app">
    <header class="navbar">
      <div class="logo">⚡ TGo Framework</div>
      <div class="badge" id="tenant-badge">Tenant: default</div>
    </header>
    <main class="hero">
      <h1>Modern Go Cloud & Desktop Architecture</h1>
      <p class="subtitle">Unified ConnectRPC, Native Multi-Tenancy & Wails In-Memory IPC</p>
      <div class="card">
        <button id="btn-fetch" class="primary-btn">Fetch Profile via RPC</button>
        <pre id="output">Click button to call backend...</pre>
      </div>
    </main>
  </div>
  <script src="main.js"></script>
</body>
</html>`

		styleCSS := `* { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
body { background: #0f172a; color: #f8fafc; min-height: 100vh; display: flex; flex-direction: column; }
.navbar { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem 2rem; border-bottom: 1px solid #334155; background: #1e293b; }
.logo { font-size: 1.25rem; font-weight: 700; color: #38bdf8; }
.badge { background: #0369a1; padding: 0.35rem 0.75rem; border-radius: 9999px; font-size: 0.85rem; font-weight: 600; }
.hero { max-width: 800px; margin: 4rem auto; text-align: center; padding: 0 1.5rem; }
h1 { font-size: 2.5rem; margin-bottom: 1rem; color: #ffffff; }
.subtitle { color: #94a3b8; font-size: 1.15rem; margin-bottom: 2.5rem; }
.card { background: #1e293b; border: 1px solid #334155; border-radius: 12px; padding: 2rem; text-align: left; }
.primary-btn { background: #38bdf8; color: #0f172a; border: none; padding: 0.75rem 1.5rem; font-weight: 600; border-radius: 8px; cursor: pointer; transition: all 0.2s; }
.primary-btn:hover { background: #7dd3fc; }
pre { background: #0f172a; padding: 1rem; border-radius: 8px; margin-top: 1.5rem; color: #38bdf8; font-family: monospace; overflow-x: auto; font-size: 0.9rem; }`

		mainJS := `document.getElementById('btn-fetch').addEventListener('click', async () => {
  const output = document.getElementById('output');
  output.textContent = 'Calling RPC...';
  try {
    const res = await fetch('/user.v1.UserService/GetProfile', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-Slug': 'demo_tenant',
        'Connect-Protocol-Version': '1'
      },
      body: JSON.stringify({ user_id: 'usr_frontend_demo' })
    });
    const data = await res.json();
    output.textContent = JSON.stringify(data, null, 2);
  } catch (err) {
    output.textContent = 'Error: ' + err.message;
  }
});`

		_ = os.WriteFile(filepath.Join(frontendDir, "index.html"), []byte(indexHTML), 0644)
		_ = os.WriteFile(filepath.Join(frontendDir, "style.css"), []byte(styleCSS), 0644)
		_ = os.WriteFile(filepath.Join(frontendDir, "main.js"), []byte(mainJS), 0644)

		log.Println("[Craft] Generated frontend starter in ./frontend/")
	},
}

// Helpers
func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "")
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func init() {
	MakeModelCmd.Flags().BoolVarP(&makeModelWithMigration, "migration", "m", false, "Generate a migration file for the model")
	MakeCmd.AddCommand(MakeModelCmd)
	MakeCmd.AddCommand(MakeMigrationCmd)
	MakeCmd.AddCommand(MakeServiceCmd)
	MakeCmd.AddCommand(MakeHandlerCmd)
	MakeCmd.AddCommand(MakeFrontendCmd)
}
