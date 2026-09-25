package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	newAppMode string
)

var NewCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new TGo web or desktop project with clean git repository",
	Long:  "Scaffolds a complete, high-performance TGo project structure with its own go.mod, clean git history, and ConnectRPC/REST setup.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.TrimSpace(args[0])
		if projectName == "" {
			fmt.Println("❌ Error: Project name cannot be empty.")
			os.Exit(1)
		}

		targetDir := projectName
		if projectName == "." {
			cwd, _ := os.Getwd()
			targetDir = cwd
			projectName = filepath.Base(cwd)
		} else {
			if _, err := os.Stat(targetDir); err == nil {
				fmt.Printf("❌ Error: Directory '%s' already exists.\n", targetDir)
				os.Exit(1)
			}
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				fmt.Printf("❌ Failed to create directory: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("🚀 Scaffolding new TGo project: %s (%s mode)...\n", projectName, newAppMode)

		// 1. Create directory layout
		dirs := []string{
			filepath.Join(targetDir, "app", "handlers"),
			filepath.Join(targetDir, "app", "models"),
			filepath.Join(targetDir, "app", "services"),
			filepath.Join(targetDir, "config"),
			filepath.Join(targetDir, "database", "migrations"),
			filepath.Join(targetDir, "frontend"),
			filepath.Join(targetDir, "proto", "v1"),
		}
		for _, d := range dirs {
			_ = os.MkdirAll(d, 0755)
		}

		// 2. Generate go.mod
		goModContent := fmt.Sprintf(`module %s

go 1.25.3

require (
	github.com/tokalink/tgo latest
)
`, projectName)
		_ = os.WriteFile(filepath.Join(targetDir, "go.mod"), []byte(goModContent), 0644)

		// 3. Generate main.go
		mainContent := fmt.Sprintf(`package main

import (
	"log"
	"net/http"

	"github.com/tokalink/tgo/pkg/app"
	"github.com/tokalink/tgo/pkg/config"
	"github.com/tokalink/tgo/pkg/transport/middleware"
)

func main() {
	// 1. Load Application Configuration
	cfg := config.MustLoad()

	// 2. Initialize TGo App Container
	application := app.New().SetAddr(":" + cfg.App.Port)

	// 3. Register Example REST & RPC Handlers
	application.Server().Register("/api/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := middleware.GetTenant(r.Context())
		if tenant == "" {
			tenant = "default"
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ok\",\"app\":\"%s\",\"tenant\":\"" + tenant + "\"}"))
	}))

	// 4. Run Application Server (ConnectRPC + Native HTTP/2)
	log.Printf("⚡ [%%s] Running on http://localhost:%%s (Multi-Tenant & ConnectRPC Ready)\n", cfg.App.Name, cfg.App.Port)
	if err := application.Run(); err != nil {
		log.Fatalf("Server terminated: %%v", err)
	}
}
`, projectName)
		_ = os.WriteFile(filepath.Join(targetDir, "main.go"), []byte(mainContent), 0644)

		// 4. Generate config/app.yaml
		configYaml := fmt.Sprintf(`app:
  name: "%s"
  env: "development"
  port: "8080"
  debug: true

database:
  driver: "sqlite" # postgres, mysql, sqlite
  database: "data/app.db"
  isolation: "%s" # single or multi_tenant

multitenancy:
  enabled: %t
  default_tenant: "public"
  resolver: "header" # header, jwt, subdomain
`, projectName, newAppMode, newAppMode == "multi")
		_ = os.WriteFile(filepath.Join(targetDir, "config", "app.yaml"), []byte(configYaml), 0644)

		// 5. Generate .env.example
		envExample := `APP_NAME=` + projectName + `
APP_ENV=development
APP_PORT=8080
DB_DRIVER=sqlite
DB_DATABASE=data/app.db
`
		_ = os.WriteFile(filepath.Join(targetDir, ".env.example"), []byte(envExample), 0644)
		_ = os.WriteFile(filepath.Join(targetDir, ".env"), []byte(envExample), 0644)

		// 6. Generate .gitignore
		gitIgnore := `*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/
tmp/
*.db
*.sqlite
.env
.env.local
.idea/
.vscode/
`
		_ = os.WriteFile(filepath.Join(targetDir, ".gitignore"), []byte(gitIgnore), 0644)

		// 7. Generate .air.toml for Live Hot Reload
		airToml := `root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main.exe ./main.go"
  bin = "./tmp/main.exe"
  delay = 500
  exclude_dir = ["tmp", "vendor", "frontend/dist", "data"]
  include_ext = ["go", "yaml", "env", "html", "css", "js"]

[log]
  time = true

[color]
  main = "magenta"
  watcher = "cyan"
  build = "yellow"
  runner = "green"
`
		_ = os.WriteFile(filepath.Join(targetDir, ".air.toml"), []byte(airToml), 0644)

		// 8. Generate Frontend Starter (index.html, style.css, main.js)
		frontendHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>%s • Powered by TGo</title>
  <link rel="stylesheet" href="style.css" />
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;600;700;800&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="badge">⚡ TGo Framework</div>
      <h1>%s</h1>
      <p class="subtitle">Unified ConnectRPC, Native Multi-Tenancy & Desktop IPC Ready</p>
      
      <div class="actions">
        <button id="btn-ping" class="btn-primary">Ping Server</button>
      </div>

      <pre id="output">Waiting for action...</pre>
    </div>
  </div>
  <script src="main.js"></script>
</body>
</html>`, projectName, projectName)
		_ = os.WriteFile(filepath.Join(targetDir, "frontend", "index.html"), []byte(frontendHTML), 0644)

		frontendCSS := `* { box-sizing: border-box; margin: 0; padding: 0; font-family: 'Plus Jakarta Sans', sans-serif; }
body { background: #090d16; color: #f8fafc; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
.container { width: 100%; max-width: 600px; padding: 1.5rem; }
.card { background: rgba(18, 24, 38, 0.75); border: 1px solid rgba(255,255,255,0.08); border-radius: 16px; padding: 2.5rem; text-align: center; backdrop-filter: blur(16px); box-shadow: 0 8px 32px rgba(0,0,0,0.4); }
.badge { display: inline-block; background: rgba(56, 189, 248, 0.1); border: 1px solid rgba(56, 189, 248, 0.25); color: #38bdf8; padding: 4px 12px; border-radius: 9999px; font-size: 0.8rem; font-weight: 700; margin-bottom: 1rem; }
h1 { font-size: 2rem; font-weight: 800; margin-bottom: 0.5rem; }
.subtitle { color: #94a3b8; font-size: 0.95rem; margin-bottom: 2rem; }
.btn-primary { background: linear-gradient(135deg, #38bdf8 0%, #3b82f6 100%); color: #090d16; border: none; padding: 0.75rem 1.5rem; font-weight: 700; border-radius: 8px; cursor: pointer; transition: all 0.2s; }
.btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
pre { background: #060911; border: 1px solid rgba(255,255,255,0.08); padding: 1rem; border-radius: 8px; margin-top: 1.5rem; color: #38bdf8; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; text-align: left; overflow-x: auto; }
`
		_ = os.WriteFile(filepath.Join(targetDir, "frontend", "style.css"), []byte(frontendCSS), 0644)

		frontendJS := `document.getElementById('btn-ping').addEventListener('click', async () => {
  const output = document.getElementById('output');
  output.textContent = 'Calling /api/health...';
  try {
    const res = await fetch('/api/health');
    const data = await res.json();
    output.textContent = JSON.stringify(data, null, 2);
  } catch (err) {
    output.textContent = 'Error: ' + err.message;
  }
});
`
		_ = os.WriteFile(filepath.Join(targetDir, "frontend", "main.js"), []byte(frontendJS), 0644)

		// 9. Generate README.md
		readme := fmt.Sprintf(`# %s

A high-performance modern web & desktop application built with [TGo Framework](https://github.com/tokalink/tgo).

## Getting Started

### Development Mode (with Live Reload)
`+"```bash"+`
air
# or standard go run
go run main.go
`+"```"+`

Open [http://localhost:8080](http://localhost:8080) in your browser.

### Craft CLI Commands
`+"```bash"+`
craft make:model User -m       # Generate model & SQL migration
craft make:handler UserHandler # Generate RPC & REST controller
craft migrate                  # Run database migrations
`+"```"+`
`, projectName)
		_ = os.WriteFile(filepath.Join(targetDir, "README.md"), []byte(readme), 0644)

		// 10. Run git init inside the new project directory
		gitCmd := exec.Command("git", "init")
		gitCmd.Dir = targetDir
		_ = gitCmd.Run()

		// 11. Run go mod tidy inside the new project directory
		fmt.Println("📦 Resolving dependencies with 'go mod tidy'...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = targetDir
		tidyCmd.Env = append(os.Environ(), "GOPROXY=direct")
		_ = tidyCmd.Run()

		fmt.Println("✅ Project scaffolded successfully with a fresh Git repository!")
		fmt.Printf("\nNext steps:\n")
		if projectName != "." {
			fmt.Printf("  cd %s\n", projectName)
		}
		fmt.Println("  go run main.go")
		fmt.Printf("\nHappy coding with ⚡ TGo!\n")
	},
}

func init() {
	NewCmd.Flags().StringVarP(&newAppMode, "mode", "m", "single", "Application mode: 'single' (standard monolithic app) or 'multi' (multi-tenant SaaS)")
}
