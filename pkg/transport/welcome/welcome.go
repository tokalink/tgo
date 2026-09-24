package welcome

import (
	_ "embed"
	"net/http"
	"os"
	"path/filepath"
)

// Handler returns an HTTP handler that serves frontend files from ./frontend if present, or falls back to the embedded welcome dashboard.
func Handler() http.Handler {
	fileServer := http.FileServer(http.Dir("frontend"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if ./frontend/index.html exists on disk
		if _, err := os.Stat(filepath.Join("frontend", "index.html")); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback: Built-in Welcome Page
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(HTML))
	})
}

// HTML contains the standalone fallback welcome dashboard for TGo Framework.
const HTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>TGo Framework — Next-Gen Go Cloud & Desktop</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Plus+Jakarta+Sans:wght@400;500;600;700;800&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg-dark: #07090e;
      --bg-surface: #0f172a;
      --bg-card: rgba(30, 41, 59, 0.7);
      --border-color: rgba(255, 255, 255, 0.08);
      --border-glow: rgba(56, 189, 248, 0.3);
      --accent-cyan: #38bdf8;
      --accent-indigo: #818cf8;
      --accent-emerald: #34d399;
      --accent-amber: #fbbf24;
      --text-main: #f8fafc;
      --text-muted: #94a3b8;
    }

    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-dark);
      color: var(--text-main);
      font-family: 'Plus Jakarta Sans', -apple-system, sans-serif;
      min-height: 100vh;
      overflow-x: hidden;
      line-height: 1.6;
      background-image: 
        radial-gradient(circle at 15% 15%, rgba(56, 189, 248, 0.08) 0%, transparent 40%),
        radial-gradient(circle at 85% 25%, rgba(129, 140, 248, 0.08) 0%, transparent 45%);
    }

    .container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 0 1.5rem;
    }

    /* Navbar */
    .navbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1.25rem 0;
      border-bottom: 1px solid var(--border-color);
      backdrop-filter: blur(12px);
      position: sticky;
      top: 0;
      z-index: 50;
    }
    .logo-container {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      text-decoration: none;
      color: var(--text-main);
    }
    .logo-icon {
      width: 36px;
      height: 36px;
      background: linear-gradient(135deg, var(--accent-cyan), var(--accent-indigo));
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 800;
      font-size: 1.2rem;
      color: #07090e;
      box-shadow: 0 0 20px rgba(56, 189, 248, 0.4);
    }
    .logo-text { font-size: 1.35rem; font-weight: 800; letter-spacing: -0.5px; }
    .logo-badge {
      font-size: 0.75rem;
      padding: 0.2rem 0.6rem;
      background: rgba(56, 189, 248, 0.12);
      color: var(--accent-cyan);
      border: 1px solid rgba(56, 189, 248, 0.3);
      border-radius: 9999px;
      font-weight: 600;
    }

    .nav-status {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.85rem;
      color: var(--accent-emerald);
      background: rgba(52, 211, 153, 0.1);
      padding: 0.35rem 0.85rem;
      border-radius: 9999px;
      border: 1px solid rgba(52, 211, 153, 0.25);
    }
    .status-dot {
      width: 8px;
      height: 8px;
      background: var(--accent-emerald);
      border-radius: 50%;
      box-shadow: 0 0 10px var(--accent-emerald);
      animation: pulse 2s infinite;
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; transform: scale(1); }
      50% { opacity: 0.5; transform: scale(0.85); }
    }

    /* Hero Section */
    .hero {
      text-align: center;
      padding: 4.5rem 0 3rem;
    }
    .hero-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.85rem;
      background: rgba(129, 140, 248, 0.1);
      color: var(--accent-indigo);
      border: 1px solid rgba(129, 140, 248, 0.25);
      padding: 0.4rem 1rem;
      border-radius: 9999px;
      margin-bottom: 1.5rem;
      font-weight: 600;
    }
    .hero h1 {
      font-size: 3.25rem;
      font-weight: 800;
      letter-spacing: -1.5px;
      line-height: 1.15;
      margin-bottom: 1.25rem;
      background: linear-gradient(180deg, #ffffff 30%, #94a3b8 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
    .hero p {
      font-size: 1.2rem;
      color: var(--text-muted);
      max-width: 680px;
      margin: 0 auto 2.5rem;
    }

    /* Interactive Live RPC Playground */
    .playground-card {
      background: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: 20px;
      padding: 2.25rem;
      backdrop-filter: blur(16px);
      box-shadow: 0 20px 40px rgba(0,0,0,0.4);
      margin-bottom: 4rem;
      position: relative;
      overflow: hidden;
    }
    .playground-card::before {
      content: '';
      position: absolute;
      top: 0; left: 0; right: 0; height: 1px;
      background: linear-gradient(90deg, transparent, var(--accent-cyan), var(--accent-indigo), transparent);
    }
    .playground-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1.5rem;
      flex-wrap: wrap;
      gap: 1rem;
    }
    .playground-title {
      font-size: 1.25rem;
      font-weight: 700;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }
    .endpoint-badge {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.85rem;
      background: rgba(0,0,0,0.5);
      padding: 0.35rem 0.75rem;
      border-radius: 8px;
      border: 1px solid var(--border-color);
      color: var(--accent-cyan);
    }

    .form-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 1.25rem;
      margin-bottom: 1.5rem;
    }
    .form-group label {
      display: block;
      font-size: 0.85rem;
      font-weight: 600;
      color: var(--text-muted);
      margin-bottom: 0.5rem;
    }
    .form-control {
      width: 100%;
      background: rgba(15, 23, 42, 0.8);
      border: 1px solid var(--border-color);
      color: var(--text-main);
      padding: 0.75rem 1rem;
      border-radius: 10px;
      font-family: inherit;
      font-size: 0.95rem;
      transition: all 0.2s;
    }
    .form-control:focus {
      outline: none;
      border-color: var(--accent-cyan);
      box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
    }

    .btn-send {
      background: linear-gradient(135deg, var(--accent-cyan), #0284c7);
      color: #07090e;
      border: none;
      font-weight: 700;
      font-size: 1rem;
      padding: 0.85rem 1.75rem;
      border-radius: 10px;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      transition: all 0.2s;
      box-shadow: 0 4px 15px rgba(56, 189, 248, 0.3);
    }
    .btn-send:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(56, 189, 248, 0.45);
    }
    .btn-send:active { transform: translateY(0); }

    .console-wrapper {
      margin-top: 1.5rem;
      background: #090d16;
      border: 1px solid var(--border-color);
      border-radius: 12px;
      overflow: hidden;
    }
    .console-bar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      background: rgba(255,255,255,0.03);
      padding: 0.6rem 1rem;
      font-size: 0.8rem;
      color: var(--text-muted);
      border-bottom: 1px solid var(--border-color);
    }
    .console-output {
      padding: 1.25rem;
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.88rem;
      color: #38bdf8;
      white-space: pre-wrap;
      max-height: 240px;
      overflow-y: auto;
    }

    /* Grid Features */
    .features-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
      gap: 1.5rem;
      margin-bottom: 4rem;
    }
    .feature-card {
      background: var(--bg-card);
      border: 1px solid var(--border-color);
      border-radius: 16px;
      padding: 1.75rem;
      transition: all 0.25s;
    }
    .feature-card:hover {
      border-color: rgba(56, 189, 248, 0.3);
      transform: translateY(-4px);
    }
    .feature-icon {
      font-size: 1.75rem;
      margin-bottom: 1rem;
      display: inline-block;
    }
    .feature-card h3 {
      font-size: 1.15rem;
      font-weight: 700;
      margin-bottom: 0.5rem;
    }
    .feature-card p {
      font-size: 0.92rem;
      color: var(--text-muted);
    }

    /* CLI Cheatsheet */
    .cli-section {
      background: linear-gradient(180deg, rgba(30, 41, 59, 0.4), rgba(15, 23, 42, 0.6));
      border: 1px solid var(--border-color);
      border-radius: 20px;
      padding: 2.25rem;
      margin-bottom: 4rem;
    }
    .cli-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 1.25rem;
      margin-top: 1.5rem;
    }
    .cli-item {
      background: #090d16;
      border: 1px solid var(--border-color);
      border-radius: 10px;
      padding: 1rem;
      cursor: pointer;
      display: flex;
      justify-content: space-between;
      align-items: center;
      transition: all 0.2s;
    }
    .cli-item:hover {
      border-color: var(--accent-indigo);
      background: #0f172a;
    }
    .cli-code {
      font-family: 'JetBrains Mono', monospace;
      color: var(--accent-indigo);
      font-size: 0.9rem;
    }
    .cli-desc {
      font-size: 0.75rem;
      color: var(--text-muted);
      margin-top: 0.25rem;
    }
    .copy-tag {
      font-size: 0.75rem;
      color: var(--text-muted);
    }

    /* Footer */
    footer {
      border-top: 1px solid var(--border-color);
      padding: 2.5rem 0;
      text-align: center;
      color: var(--text-muted);
      font-size: 0.9rem;
    }

    @media (max-width: 768px) {
      .hero h1 { font-size: 2.25rem; }
      .hero p { font-size: 1rem; }
    }
  </style>
</head>
<body>
  <div class="container">
    <!-- Navbar -->
    <header class="navbar">
      <a href="/" class="logo-container">
        <div class="logo-icon">⚡</div>
        <span class="logo-text">TGo Framework</span>
        <span class="logo-badge">v1.0.0</span>
      </a>
      <div class="nav-status">
        <span class="status-dot"></span>
        <span>Server Ready (HTTP/1.1 & HTTP/2 h2c)</span>
      </div>
    </header>

    <!-- Hero -->
    <section class="hero">
      <div class="hero-badge">⚡ Triple Transport • Native Multi-Tenancy • In-Memory IPC</div>
      <h1>Next-Gen Go Application Framework</h1>
      <p>Modern, high-performance Go application framework for Cloud SaaS, Unified Desktop/Mobile (Wails v3), and Microservices.</p>
    </section>

    <!-- Interactive Live RPC Playground -->
    <section class="playground-card">
      <div class="playground-header">
        <div class="playground-title">
          <span>🧪 Live ConnectRPC Tester</span>
        </div>
        <div class="endpoint-badge">POST /user.v1.UserService/GetProfile</div>
      </div>

      <div class="form-grid">
        <div class="form-group">
          <label for="tenantInput">Tenant Context (Header / Subdomain)</label>
          <select id="tenantInput" class="form-control">
            <option value="acme_corp">acme_corp</option>
            <option value="enterprise_client">enterprise_client</option>
            <option value="demo_tenant">demo_tenant</option>
            <option value="public">public (Default)</option>
          </select>
        </div>
        <div class="form-group">
          <label for="userIdInput">User ID</label>
          <input type="text" id="userIdInput" class="form-control" value="usr_developer_01" placeholder="Enter User ID..." />
        </div>
        <div class="form-group" style="display: flex; align-items: flex-end;">
          <button id="sendBtn" class="btn-send" style="width: 100%;">
            <span>⚡ Send RPC Request</span>
          </button>
        </div>
      </div>

      <div class="console-wrapper">
        <div class="console-bar">
          <span id="responseStatus">Status: Ready to execute</span>
          <span id="latencyBadge">⚡ Latency: -</span>
        </div>
        <pre id="consoleOutput" class="console-output">// Klik 'Send RPC Request' untuk memanggil ConnectRPC handler secara live...</pre>
      </div>
    </section>

    <!-- Architecture Features -->
    <section class="features-grid">
      <div class="feature-card">
        <div class="feature-icon">🚀</div>
        <h3>Triple Mode Transport</h3>
        <p>1 handler melayani ConnectRPC (gRPC binary protobuf), RESTful JSON, dan In-Memory IPC di port yang sama.</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🏢</div>
        <h3>Native Multi-Tenancy</h3>
        <p>Resolusi otomatis via Subdomain, Header, atau JWT Claim dengan PostgreSQL Schema & SQLite File isolator.</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🖥️</div>
        <h3>Wails In-Memory IPC</h3>
        <p>0 network hop (~170,000 req/sec) menghubungkan UI Wails v3 langsung ke backend Go tanpa socket TCP.</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🛠️</div>
        <h3>Craft Developer CLI</h3>
        <p>Scaffolding kode otomatis untuk model, repository, service action, handler, dan migrasi skema database.</p>
      </div>
    </section>

    <!-- CLI Cheatsheet -->
    <section class="cli-section">
      <div class="playground-title">
        <span>🛠️ Quick Craft CLI Commands (Click to Copy)</span>
      </div>
      <div class="cli-grid">
        <div class="cli-item" onclick="copyCode('go run . craft serve')">
          <div>
            <div class="cli-code">go run . craft serve</div>
            <div class="cli-desc">Jalankan application server</div>
          </div>
          <span class="copy-tag">📋 Copy</span>
        </div>
        <div class="cli-item" onclick="copyCode('go run . craft tenant:create acme')">
          <div>
            <div class="cli-code">go run . craft tenant:create acme</div>
            <div class="cli-desc">Provision & migrasi tenant baru</div>
          </div>
          <span class="copy-tag">📋 Copy</span>
        </div>
        <div class="cli-item" onclick="copyCode('go run . craft make:model Product')">
          <div>
            <div class="cli-code">go run . craft make:model Product</div>
            <div class="cli-desc">Generate model struct & repository</div>
          </div>
          <span class="copy-tag">📋 Copy</span>
        </div>
        <div class="cli-item" onclick="copyCode('go run . craft make:service Order')">
          <div>
            <div class="cli-code">go run . craft make:service Order</div>
            <div class="cli-desc">Generate service action logic</div>
          </div>
          <span class="copy-tag">📋 Copy</span>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer>
      <p>⚡ <strong>TGo Framework</strong> — Built for Speed, Multi-Tenancy, and Hybrid Cloud/Desktop Apps.</p>
      <p style="margin-top: 0.5rem; font-size: 0.8rem; color: #64748b;">MIT License © 2026 TGo Platform Core Team.</p>
    </footer>
  </div>

  <script>
    async function sendRequest() {
      const tenant = document.getElementById('tenantInput').value;
      const userId = document.getElementById('userIdInput').value;
      const output = document.getElementById('consoleOutput');
      const statusEl = document.getElementById('responseStatus');
      const latencyEl = document.getElementById('latencyBadge');

      output.textContent = 'Calling RPC procedure /user.v1.UserService/GetProfile...';
      statusEl.textContent = 'Status: Sending...';
      
      const startTime = performance.now();
      try {
        const response = await fetch('/user.v1.UserService/GetProfile', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Tenant-Slug': tenant,
            'Connect-Protocol-Version': '1'
          },
          body: JSON.stringify({ user_id: userId })
        });

        const duration = (performance.now() - startTime).toFixed(1);
        const data = await response.json();
        
        statusEl.textContent = 'Status: ' + response.status + ' ' + response.statusText;
        latencyEl.textContent = '⚡ Latency: ' + duration + ' ms';
        output.textContent = JSON.stringify(data, null, 2);
      } catch (err) {
        statusEl.textContent = 'Status: Error';
        output.textContent = 'Request failed: ' + err.message;
      }
    }

    document.getElementById('sendBtn').addEventListener('click', sendRequest);

    function copyCode(text) {
      navigator.clipboard.writeText(text);
      alert('Copied to clipboard: ' + text);
    }
  </script>
</body>
</html>`
