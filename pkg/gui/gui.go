// Package gui provides a lightweight web-based GUI for the SonicScrewdriver
// container management system. It serves a dashboard for monitoring containers,
// managing deployments, and viewing system health.
package gui

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/uDosGo/SonicScrewdriver/pkg/container"
)

//go:embed static/*
var staticFiles embed.FS

// Server is the SonicScrewdriver GUI HTTP server
type Server struct {
	runtime container.Runtime
	port    int
	server  *http.Server
}

// NewServer creates a new GUI server
func NewServer(runtime container.Runtime, port int) *Server {
	return &Server{
		runtime: runtime,
		port:    port,
	}
}

// Start starts the GUI server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/containers", s.handleContainers)
	mux.HandleFunc("/api/containers/", s.handleContainerAction)
	mux.HandleFunc("/api/health", s.handleHealth)

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to get static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	// Try to find an available port
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		// Try a random port
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		s.port = listener.Addr().(*net.TCPAddr).Port
	}

	log.Printf("SonicScrewdriver GUI starting on http://localhost:%d", s.port)
	go s.server.Serve(listener)
	return nil
}

// Port returns the port the server is running on
func (s *Server) Port() int {
	return s.port
}

// Stop stops the GUI server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func (s *Server) handleContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	containers, err := s.runtime.List()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Get health for each container
	type containerInfo struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	var info []containerInfo
	for _, name := range containers {
		health, err := s.runtime.CheckContainerHealth(name)
		status := "unknown"
		if err == nil && health != nil {
			status = health.Status
		}
		info = append(info, containerInfo{Name: name, Status: status})
	}

	json.NewEncoder(w).Encode(info)
}

func (s *Server) handleContainerAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse path: /api/containers/{name}/{action}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/containers/"), "/")
	if len(parts) < 2 {
		http.Error(w, `{"error":"invalid path"}`, http.StatusBadRequest)
		return
	}
	name, action := parts[0], parts[1]

	var err error
	switch action {
	case "start":
		err = s.runtime.Start(name)
	case "stop":
		err = s.runtime.Stop(name)
	case "restart":
		err = s.runtime.RestartContainer(name)
	case "remove":
		err = s.runtime.Remove(name)
	default:
		http.Error(w, fmt.Sprintf(`{"error":"unknown action: %s"}`, action), http.StatusBadRequest)
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	healthStatuses, err := s.runtime.GetAllContainerHealth()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(healthStatuses)
}

// EnsureStaticDir creates the static directory and writes the default HTML
func EnsureStaticDir(basePath string) error {
	staticDir := filepath.Join(basePath, "static")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return err
	}

	indexHTML := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SonicScrewdriver Console</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0d1117; color: #c9d1d9; }
.container { max-width: 1200px; margin: 0 auto; padding: 20px; }
header { padding: 20px 0; border-bottom: 1px solid #30363d; margin-bottom: 20px; }
h1 { font-size: 24px; color: #58a6ff; }
h1 small { font-size: 14px; color: #8b949e; margin-left: 10px; }
.controls { margin-bottom: 20px; }
.controls button { background: #21262d; color: #c9d1d9; border: 1px solid #30363d; padding: 8px 16px; border-radius: 6px; cursor: pointer; margin-right: 8px; }
.controls button:hover { background: #30363d; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px; text-align: left; border-bottom: 1px solid #21262d; }
th { color: #8b949e; font-weight: 600; text-transform: uppercase; font-size: 12px; }
.status-running { color: #3fb950; }
.status-stopped { color: #f85149; }
.status-not_found { color: #8b949e; }
.actions button { background: none; border: 1px solid #30363d; color: #c9d1d9; padding: 4px 8px; border-radius: 4px; cursor: pointer; margin-right: 4px; font-size: 12px; }
.actions button:hover { background: #30363d; }
.loading { text-align: center; padding: 40px; color: #8b949e; }
.error { color: #f85149; padding: 20px; text-align: center; }
</style>
</head>
<body>
<div class="container">
<header><h1>SonicScrewdriver <small>Container Console</small></h1></header>
<div class="controls">
<button onclick="refreshContainers()">🔄 Refresh</button>
</div>
<div id="content"><div class="loading">Loading containers...</div></div>
</div>
<script>
async function refreshContainers() {
const content = document.getElementById('content');
content.innerHTML = '<div class="loading">Loading containers...</div>';
try {
const res = await fetch('/api/containers');
if (!res.ok) throw new Error('Failed to fetch');
const containers = await res.json();
if (containers.length === 0) {
content.innerHTML = '<p style="text-align:center;padding:40px;color:#8b949e;">No containers found. Start one with <code>sonic container start <name></code></p>';
return;
}
let html = '<table><thead><tr><th>Name</th><th>Status</th><th>Actions</th></tr></thead><tbody>';
for (const c of containers) {
const statusClass = 'status-' + c.Status.toLowerCase().replace(/\s+/g, '_');
html += '<tr>';
html += '<td>' + c.Name + '</td>';
html += '<td class="' + statusClass + '">' + c.Status + '</td>';
html += '<td class="actions">';
html += '<button onclick="actionContainer(\'' + c.Name + '\',\'start\')">Start</button>';
html += '<button onclick="actionContainer(\'' + c.Name + '\',\'stop\')">Stop</button>';
html += '<button onclick="actionContainer(\'' + c.Name + '\',\'restart\')">Restart</button>';
html += '<button onclick="actionContainer(\'' + c.Name + '\',\'remove\')">Remove</button>';
html += '</td></tr>';
}
html += '</tbody></table>';
content.innerHTML = html;
} catch (e) {
content.innerHTML = '<div class="error">Error: ' + e.message + '</div>';
}
}
async function actionContainer(name, action) {
try {
const res = await fetch('/api/containers/' + name + '/' + action, { method: 'POST' });
const data = await res.json();
if (data.status === 'error') alert('Error: ' + data.message);
refreshContainers();
} catch (e) {
alert('Error: ' + e.message);
}
}
refreshContainers();
</script>
</body>
</html>`

	return os.WriteFile(filepath.Join(staticDir, "index.html"), []byte(indexHTML), 0644)
}

// Ensure we reference template for compilation
var _ = template.HTMLEscapeString
