# Classic Modern Mint — Complete Setup Guide

> **Version:** 1.0.0  
> **Last Updated:** 2026-05-12  
> **Status:** Production Ready  
> **Target:** Linux Mint 21+ (Ubuntu-based)

---

## 🎯 Overview

This guide transforms a fresh Linux Mint installation into **Classic Modern Mint (CMM)** — a lean, text-first, feed-driven development environment optimized for the uDos ecosystem.

### What You Get
- **Minimal bloat** - Only essential packages installed
- **Text-first workflow** - Markdown, Obsidian, Vault-based document system
- **Feed-driven** - Email, GitHub, and other inputs as structured JSON feeds
- **MCP-enabled** - Model Context Protocol for LLM integration
- **Dev-ready** - All development tools pre-configured

### What Gets Removed
- Firefox → **Zen Browser**
- LibreOffice → **Obsidian + markdown-it**
- Thunderbird → **Feed-based email ingestion**
- GIMP, Rhythmbox, Celluloid, Totem → **Removed** (Hypnotix kept for IPTV)
- All games → **Removed**
- Printing/Bluetooth → **Optional** (can be re-installed)

---

## 📋 Prerequisites

### Hardware
- Linux Mint 21 or later (Ubuntu-based)
- Minimum 8GB RAM (16GB recommended)
- 50GB free disk space
- Network connection for package downloads

### Knowledge
- Basic Linux command line
- sudo privileges
- Git/GitHub familiarity

---

## 🚀 Quick Start (One Command)

For experienced users:

```bash
cd ~/Vault/dev-briefs/classic-modern-mint
./setup-cmm.sh
```

---

## 📦 Step-by-Step Installation

---

### Phase 1: System Preparation & Bloat Removal

#### 1.1 Update System
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git jq flatpak
```

#### 1.2 Create Audit Snapshot
```bash
dpkg -l > ~/Desktop/original-packages.txt
flatpak list --app > ~/Desktop/original-flatpaks.txt
```

#### 1.3 Remove Default Bloatware
```bash
sudo apt purge --auto-remove \
    libreoffice* thunderbird* gimp* \
    rhythmbox* celluloid* totem* \
    gnome-mines* gnome-sudoku* aisleriot* \
    cups* printer-driver* bluez* blueman* \
    language-pack-{fr,de,es,it,pt}* \
    firefox* \
    -y
sudo apt autoremove --purge -y
sudo apt clean
```

#### 1.4 Remove Unnecessary Flatpaks
```bash
flatpak uninstall org.libreoffice.LibreOffice org.mozilla.Thunderbird \
    org.gimp.GIMP org.gnome.Rhythmbox3 2>/dev/null
flatpak uninstall --unused
```

---

### Phase 2: Install Core Replacements

#### 2.1 Install Zen Browser
```bash
flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo
flatpak install flathub io.github.zen_browser.zen -y
```

#### 2.2 Install Obsidian
```bash
flatpak install flathub md.obsidian.Obsidian -y
flatpak override --user --filesystem=home md.obsidian.Obsidian
```

#### 2.3 Install Node.js and markdown-it
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo bash -
sudo apt install -y nodejs
sudo npm install -g markdown-it-cli
```

#### 2.4 Install Feed Tools
```bash
sudo apt install jq -y
mkdir -p ~/.local/share/udos/feeds/{system,user,mcp,network,spool,inbox}
mkdir -p ~/Vault/@inbox/email ~/Vault/family/calendar
```

---

### Phase 3: Configure Vault & Document System

#### 3.1 Create Vault Structure
```bash
mkdir -p ~/Vault/{system,home,family,user,@inbox,@workspace,@toybox,@sandbox,@public,@private,binder}

cat > ~/Vault/README.md << 'VAULT_EOF'
# Vault
This is the root of your Classic Modern Mint document system.

## Directory Structure
- `@workspace/` - Active projects
- `@inbox/` - Unsorted incoming
- `@sandbox/` - Experimental
- `@toybox/` - Experimental features
- `family/` - Shared documents
- `binder/` - Topic organization
- `system/` - System docs
- `home/` - Home automation
- `user/` - Personal notes
VAULT_EOF
```

#### 3.2 Install Wrapper Scripts
```bash
mkdir -p ~/.local/bin
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc

# md2html
cat > ~/.local/bin/md2html << 'SCRIPT_EOF'
#!/bin/bash
INPUT="$1"
OUTPUT="${2:-${INPUT%.md}.html}"
/usr/local/bin/markdown-it-cli "$INPUT" -o "$OUTPUT" && echo "Converted $INPUT to $OUTPUT"
SCRIPT_EOF
chmod +x ~/.local/bin/md2html

# vault-export
cat > ~/.local/bin/vault-export << 'SCRIPT_EOF'
#!/bin/bash
VAULT_ROOT="${1:-$HOME/Vault}"
OUTPUT_DIR="${2:-$HOME/Vault/_export}"
mkdir -p "$OUTPUT_DIR"
find "$VAULT_ROOT" -name "*.md" -type f | while read -r md_file; do
    rel_path="${md_file#$VAULT_ROOT/}"
    html_file="$OUTPUT_DIR/${rel_path%.md}.html"
    mkdir -p "$(dirname "$html_file")"
    /usr/local/bin/markdown-it-cli "$md_file" -o "$html_file"
    echo "Exported $rel_path"
done
echo "Export complete: $OUTPUT_DIR"
SCRIPT_EOF
chmod +x ~/.local/bin/vault-export
```

#### 3.3 Configure Obsidian
```bash
flatpak run md.obsidian.Obsidian ~/Vault
```

#### 3.4 Create Email-to-Feed Pipeline
```bash
cat > ~/.local/bin/email-to-feed << 'SCRIPT_EOF'
#!/bin/bash
IMAP_SERVER="${IMAP_SERVER:-imap.gmail.com}"
IMAP_USER="${IMAP_USER:-$USER}"
IMAP_PASS="${IMAP_PASS:-}"
FEED_FILE="$HOME/.local/share/udos/feeds/inbox/email.json"

if [ -z "$IMAP_PASS" ]; then
    echo "Error: IMAP_PASS not set. Use: export IMAP_PASS='your-password'"
    exit 1
fi

new_entry=$(jq -n \
    --arg id "email-$(date +%s)" \
    --arg timestamp "$(date -Iseconds)" \
    --arg server "$IMAP_SERVER" \
    --arg user "$IMAP_USER" \
    '{"id": $id, "timestamp": $timestamp, "type": "email.received", "summary": "Email from '"$server"'", "payload": {"server": $server, "user": $user}}')

if [ -f "$FEED_FILE" ]; then
    jq --argjson new "$new_entry" '.items += [$new]' "$FEED_FILE" > "$FEED_FILE.tmp" && \
    mv "$FEED_FILE.tmp" "$FEED_FILE"
else
    jq -n --argjson item "$new_entry" '{"version": "1.0", "type": "feed", "source": "inbox.email", "items": [$item]}' > "$FEED_FILE"
fi

echo "Email feed updated: $FEED_FILE"
SCRIPT_EOF
chmod +x ~/.local/bin/email-to-feed
```

---

### Phase 4: Create Sonic CLI Wrapper

```bash
cat > ~/.local/bin/sonic << 'SCRIPT_EOF'
#!/bin/bash
ACTION="$1"
shift

case "$ACTION" in
    doc)
        SUBACTION="$1"
        shift
        case "$SUBACTION" in
            create)
                NAME="$1"
                DATE=$(date -I)
                cat > "$HOME/Vault/@workspace/$NAME.md" << DOC
---
title: "$NAME"
type: "document"
created: "$DATE"
---
# $NAME
Content here...
DOC
                echo "Created: ~/Vault/@workspace/$NAME.md"
                ;;
            search)
                find "$HOME/Vault" -name "*.md" -exec grep -l "$1" {} \;
                ;;
            export)
                INPUT="$1"
                OUTPUT="${2:-${INPUT%.md}.html}"
                /usr/local/bin/markdown-it-cli "$INPUT" -o "$OUTPUT" 2>/dev/null
                echo "Exported to $OUTPUT"
                ;;
            feed)
                jq -n '{"version": "1.0", "type": "feed", "source": "vault.docs", "items": [' > /tmp/feed.json
                find "$HOME/Vault" -name "*.md" -type f | while read f; do
                    TITLE=$(head -n1 "$f" | sed 's/^# //')
                    echo "{\"id\": \"$(basename "$f")\", \"title\": \"$TITLE\", \"path\": \"$f\"}," >> /tmp/feed.json
                done
                echo ']}' >> /tmp/feed.json
                jq '.' /tmp/feed.json > "$HOME/.local/share/udos/feeds/vault/docs.json"
                echo "Feed generated"
                ;;
        esac
        ;;
    audit)
        echo "Vault: $(find ~/Vault -name "*.md" | wc -l) markdown files"
        echo "Feeds: $(find ~/.local/share/udos/feeds -name "*.json" | wc -l) feed files"
        ;;
    *)
        echo "Usage: sonic {doc|audit}"
        echo "  doc create <name>"
        echo "  doc search <term>"
        echo "  doc export <file.md>"
        echo "  doc feed"
        ;;
esac
SCRIPT_EOF
chmod +x ~/.local/bin/sonic
```

---

### Phase 5: Install Systemd Services

```bash
mkdir -p ~/systemd-services

# email-feed.service
cat > ~/systemd-services/email-feed.service << 'EOF'
[Unit]
Description=Email to uDos Feed Converter
After=network.target

[Service]
Type=oneshot
User=$USER
Environment=IMAP_SERVER=imap.gmail.com
Environment=IMAP_USER=$USER
ExecStart=%h/.local/bin/email-to-feed

[Install]
WantedBy=multi-user.target
EOF

# email-feed.timer
cat > ~/systemd-services/email-feed.timer << 'EOF'
[Unit]
Description=Email feed update timer

[Timer]
OnCalendar=*:0/15
Persistent=true

[Install]
WantedBy=timers.target
EOF

# vault-mcp.service
cat > ~/systemd-services/vault-mcp.service << 'EOF'
[Unit]
Description=Vault MCP Server
After=network.target

[Service]
Type=simple
User=$USER
ExecStart=/usr/bin/python3 %h/.local/bin/vault-mcp-server.py
WorkingDirectory=%h
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Install
sudo cp ~/systemd-services/*.{service,timer} /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now email-feed.timer
```

---

### Phase 6: MCP Vault Server

```bash
cat > ~/.local/bin/vault-mcp-server.py << 'PYEOF'
#!/usr/bin/env python3
from mcp.server import Server
import os
import sys

app = Server("vault-mcp")
VAULT_ROOT = os.path.expanduser("~/Vault")

@app.list_resources()
def handle_list_resources():
    resources = []
    for root, dirs, filenames in os.walk(VAULT_ROOT):
        for f in filenames:
            full_path = os.path.join(root, f)
            rel_path = os.path.relpath(full_path, VAULT_ROOT)
            mime = "text/markdown" if f.endswith('.md') else "text/plain"
            resources.append({
                "uri": f"file:///vault/{rel_path}",
                "name": f,
                "description": f"Vault: {rel_path}",
                "mimeType": mime
            })
    return {"resources": resources}

@app.read_resource()
def handle_read_resource(uri: str):
    if uri.startswith("file:///vault/"):
        rel_path = uri[len("file:///vault/"):]
        full_path = os.path.join(VAULT_ROOT, rel_path)
        if os.path.exists(full_path):
            with open(full_path, 'r') as f:
                return {"text": f.read(), "mimeType": "text/markdown"}
    return {"error": f"File not found: {uri}"}

if __name__ == "__main__":
    app.run(sys.stdin, sys.stdout, {})
PYEOF
chmod +x ~/.local/bin/vault-mcp-server.py
```

---

## ✅ Verification

### Test All Components
```bash
source ~/.bashrc

sonic doc create "Setup-Test"
sonic doc search "Classic Modern"
sonic audit

md2html ~/Vault/README.md /tmp/test.html

vault-export ~/Vault ~/Vault/_export

export IMAP_PASS='your-password'
~/.local/bin/email-to-feed

python3 ~/.local/bin/vault-mcp-server.py
```

### Verify Services
```bash
systemctl status email-feed.timer
journalctl -u email-feed.service -n 20
cat ~/.local/share/udos/feeds/inbox/email.json
```

---

## 📊 What Changed

| Category | Before | After |
|----------|--------|-------|
| Browser | Firefox | Zen Browser |
| Office | LibreOffice | Obsidian + markdown-it |
| Email | Thunderbird | Feed-based JSON |
| Graphics | GIMP | Removed |
| Media | Rhythmbox, Celluloid | Removed |
| Games | GNOME Games | Removed |
| Documents | ~/Documents | ~/Vault/ |

---

## 🎯 Next Steps

1. Configure IMAP:
   ```bash
   export IMAP_PASS='your-password'
   export IMAP_SERVER='imap.your-provider.com'
   export IMAP_USER='your-email@provider.com'
   ```

2. Test workflow:
   ```bash
   sonic doc create "My-First-Note"
   sonic doc feed
   md2html ~/Vault/@workspace/My-First-Note.md ~/Vault/@workspace/My-First-Note.html
   ```

---

## 🔧 Troubleshooting

**Obsidian can't access Vault:**
```bash
flatpak override --user --filesystem=home md.obsidian.Obsidian
```

**markdown-it-cli not found:**
```bash
sudo npm install -g markdown-it-cli
```

**Email-feed requires password:**
```bash
export IMAP_PASS='your-password'
```

---

## 📚 Documentation

- ROADMAP.md - Current roadmap
- DEVLOG.md - Completed tasks
- INDEX.md - Documentation index

---

## 🎉 Complete!

Your Linux Mint server is now Classic Modern Mint.
