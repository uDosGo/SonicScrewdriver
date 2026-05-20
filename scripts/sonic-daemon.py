#!/usr/bin/env python3
"""
Sonic Daemon - Provider Probe & Status API
Monitors LLM providers (OpenRouter, Ollama) and exposes a health/status HTTP API.
Designed to run alongside Hivemind on the Ubuntu server.

Usage:
  python3 scripts/sonic-daemon.py              # Run interactively
  python3 scripts/sonic-daemon.py --daemon      # Run as background daemon
  python3 scripts/sonic-daemon.py --status      # Print status and exit
  python3 scripts/sonic-daemon.py --install     # Install as systemd user service

API Endpoints:
  GET  /health          -> {"status": "ok"}
  GET  /status          -> Full provider status JSON
  GET  /providers       -> Provider list with health
  GET  /providers/{id}  -> Single provider details
"""
import argparse
import json
import os
import signal
import subprocess
import sys
import time
import urllib.error
import urllib.request
from datetime import datetime
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path
from threading import Lock, Thread

# ============================================================
# Configuration
# ============================================================
SONIC_PORT = int(os.environ.get("SONIC_PORT", "30002"))
POLL_INTERVAL = int(os.environ.get("SONIC_POLL_INTERVAL", "30"))  # seconds
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

PROVIDERS = {
    "openrouter": {
        "name": "OpenRouter",
        "url": "https://openrouter.ai/api/v1/models",
        "icon": "🌐",
        "timeout": 10,
        "headers": lambda: {
            "Authorization": f"Bearer {os.environ.get('OPENROUTER_API_KEY', '')}",
        },
    },
    "ollama": {
        "name": "Ollama",
        "url": "http://localhost:11434/api/tags",
        "icon": "🤖",
        "timeout": 5,
        "headers": lambda: {},
    },
    "hivemind": {
        "name": "Hivemind",
        "url": "http://localhost:30000/health",
        "icon": "🧠",
        "timeout": 5,
        "headers": lambda: {},
    },
    "secret_server": {
        "name": "Secret Server",
        "url": "http://localhost:30001/health",
        "icon": "🔐",
        "timeout": 5,
        "headers": lambda: {},
    },
    "sse_bridge": {
        "name": "SSE Bridge",
        "url": "http://localhost:30010/health",
        "icon": "🔗",
        "timeout": 5,
        "headers": lambda: {},
    },
}


# ============================================================
# Provider probing
# ============================================================
class ProviderProbe:
    """Probes a provider endpoint and caches results."""

    def __init__(self, provider_id: str, config: dict):
        self.id = provider_id
        self.config = config
        self.last_result = {
            "status": "unknown",
            "latency_ms": None,
            "error": None,
            "last_checked": None,
            "models": None,
        }
        self.lock = Lock()

    def probe(self):
        """Probe the provider and update cached result."""
        start = time.time()
        result = {
            "status": "down",
            "latency_ms": None,
            "error": None,
            "last_checked": datetime.utcnow().isoformat() + "Z",
            "models": None,
        }

        try:
            req = urllib.request.Request(
                self.config["url"],
                headers=self.config["headers"](),
            )
            with urllib.request.urlopen(req, timeout=self.config["timeout"]) as resp:
                latency = (time.time() - start) * 1000
                data = json.loads(resp.read())
                result["status"] = "ok"
                result["latency_ms"] = round(latency, 1)

                # Extract model count from response
                if self.id == "ollama" and "models" in data:
                    result["models"] = [m["name"] for m in data["models"]]
                elif self.id == "openrouter" and "data" in data:
                    result["models"] = [m["id"] for m in data["data"]]
                elif self.id in ("hivemind", "secret_server", "sse_bridge"):
                    result["status"] = "ok" if data.get("status") in ("healthy", "ok") else "warn"

        except urllib.error.HTTPError as e:
            result["status"] = "warn"
            result["error"] = f"HTTP {e.code}: {e.reason}"
            result["latency_ms"] = round((time.time() - start) * 1000, 1)
        except urllib.error.URLError as e:
            result["error"] = str(e.reason)
        except json.JSONDecodeError:
            result["error"] = "Invalid JSON response"
        except OSError as e:
            result["error"] = str(e)
        except Exception as e:
            result["error"] = f"Unexpected: {e}"

        with self.lock:
            self.last_result = result
        return result

    def get_status(self) -> dict:
        with self.lock:
            return {
                "id": self.id,
                "name": self.config["name"],
                "icon": self.config["icon"],
                **self.last_result,
            }


# ============================================================
# Daemon state
# ============================================================
class SonicDaemon:
    """Manages provider probes and HTTP API."""

    def __init__(self):
        self.probes = {
            pid: ProviderProbe(pid, cfg)
            for pid, cfg in PROVIDERS.items()
        }
        self.running = True
        self.start_time = time.time()

    def probe_all(self):
        """Probe all providers."""
        threads = []
        for probe in self.probes.values():
            t = Thread(target=probe.probe, daemon=True)
            t.start()
            threads.append(t)
        for t in threads:
            t.join(timeout=15)

    def get_full_status(self) -> dict:
        """Get full daemon status."""
        providers = {
            pid: probe.get_status()
            for pid, probe in self.probes.items()
        }
        all_ok = all(p["status"] == "ok" for p in providers.values())
        return {
            "status": "ok" if all_ok else "degraded",
            "version": "0.1.0",
            "uptime_seconds": int(time.time() - self.start_time),
            "started_at": datetime.fromtimestamp(self.start_time).isoformat(),
            "providers": providers,
        }

    def stop(self):
        self.running = False


# ============================================================
# HTTP API Server
# ============================================================
class SonicHandler(BaseHTTPRequestHandler):
    """HTTP request handler for Sonic Daemon API."""

    daemon: SonicDaemon = None  # Set externally

    def do_GET(self):
        path = self.path.rstrip("/")

        if path == "/health":
            self._json({"status": "ok"})
        elif path == "/status":
            self._json(self.daemon.get_full_status())
        elif path == "/providers":
            status = self.daemon.get_full_status()
            self._json(status["providers"])
        elif path.startswith("/providers/"):
            pid = path.split("/providers/")[1]
            if pid in self.daemon.probes:
                self._json(self.daemon.probes[pid].get_status())
            else:
                self._json({"error": f"Unknown provider: {pid}"}, 404)
        else:
            self._json({"error": "Not found", "paths": ["/health", "/status", "/providers"]}, 404)

    def _json(self, data: dict, status: int = 200):
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()
        self.wfile.write(json.dumps(data, indent=2).encode())

    def log_message(self, format, *args):
        """Suppress default logging; we log manually."""
        pass


# ============================================================
# Main
# ============================================================
def run_daemon(daemonize: bool = False):
    """Run the Sonic Daemon."""
    daemon = SonicDaemon()
    SonicHandler.daemon = daemon

    # Start HTTP server
    server = HTTPServer(("0.0.0.0", SONIC_PORT), SonicHandler)
    server_thread = Thread(target=server.serve_forever, daemon=True)
    server_thread.start()
    print(f"🔧 Sonic Daemon API running on http://0.0.0.0:{SONIC_PORT}")

    # Handle signals
    def handle_signal(sig, frame):
        print("\n🛑 Shutting down...")
        daemon.stop()
        server.shutdown()
        sys.exit(0)

    signal.signal(signal.SIGINT, handle_signal)
    signal.signal(signal.SIGTERM, handle_signal)

    # Initial probe
    print("🔍 Probing providers...")
    daemon.probe_all()
    status = daemon.get_full_status()
    for pid, p in status["providers"].items():
        emoji = "✅" if p["status"] == "ok" else ("⚠️" if p["status"] == "warn" else "🔴")
        latency = f" ({p['latency_ms']}ms)" if p["latency_ms"] else ""
        print(f"  {emoji} {p['icon']} {p['name']}: {p['status']}{latency}")

    # Polling loop
    while daemon.running:
        time.sleep(POLL_INTERVAL)
        daemon.probe_all()


def print_status():
    """Print one-shot status."""
    daemon = SonicDaemon()
    daemon.probe_all()
    status = daemon.get_full_status()
    print(f"🔧 Sonic Daemon Status")
    print(f"  Overall: {'✅ OK' if status['status'] == 'ok' else '⚠️ Degraded'}")
    print(f"  Uptime:  {status['uptime_seconds']}s")
    print()
    for pid, p in status["providers"].items():
        emoji = "✅" if p["status"] == "ok" else ("⚠️" if p["status"] == "warn" else "🔴")
        latency = f" ({p['latency_ms']}ms)" if p["latency_ms"] else ""
        error = f" - {p['error']}" if p["error"] else ""
        models = f" - {len(p['models'])} models" if p.get("models") else ""
        print(f"  {emoji} {p['icon']} {p['name']}: {p['status']}{latency}{error}{models}")


def install_systemd():
    """Install as systemd user service."""
    service_name = "sonic-daemon"
    service_dir = os.path.expanduser("~/.config/systemd/user")
    os.makedirs(service_dir, exist_ok=True)

    script_path = os.path.abspath(__file__)

    service_content = f"""[Unit]
Description=Sonic Daemon - Provider Probe & Status API
After=network.target

[Service]
Type=simple
ExecStart={sys.executable} {script_path}
Restart=on-failure
RestartSec=5
Environment=SONIC_PORT={SONIC_PORT}
Environment=SONIC_POLL_INTERVAL={POLL_INTERVAL}

[Install]
WantedBy=default.target
"""

    service_path = os.path.join(service_dir, f"{service_name}.service")
    with open(service_path, "w") as f:
        f.write(service_content)

    subprocess.run(["systemctl", "--user", "daemon-reload"], capture_output=True)
    subprocess.run(["systemctl", "--user", "enable", service_name], capture_output=True)

    print(f"✅ Installed systemd user service: {service_name}")
    print(f"   Service file: {service_path}")
    print(f"   Start now: systemctl --user start {service_name}")
    print(f"   Status:    systemctl --user status {service_name}")
    print(f"   Logs:      journalctl --user -u {service_name} -f")


def main():
    parser = argparse.ArgumentParser(description="Sonic Daemon - Provider Probe & Status API")
    parser.add_argument("--daemon", action="store_true", help="Run as background daemon")
    parser.add_argument("--status", action="store_true", help="Print status and exit")
    parser.add_argument("--install", action="store_true", help="Install as systemd user service")
    args = parser.parse_args()

    if args.status:
        print_status()
        return

    if args.install:
        install_systemd()
        return

    if args.daemon:
        print("🔧 Sonic Daemon starting as daemon (background)...")
        if os.fork() > 0:
            return
        os.setsid()
        if os.fork() > 0:
            return
        log_dir = os.path.join(ROOT, "logs")
        os.makedirs(log_dir, exist_ok=True)
        sys.stdout = open(os.path.join(log_dir, "sonic-daemon.log"), "a")
        sys.stderr = sys.stdout

    run_daemon(daemonize=args.daemon)


if __name__ == "__main__":
    main()
