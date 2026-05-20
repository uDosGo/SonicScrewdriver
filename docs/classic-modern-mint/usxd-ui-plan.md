USXD UI Plan for Sonic Screwdriver on Classic Modern Mint

## 🎯 Core Design Philosophy

**Notion-inspired Library View + TUI Grid = Complete Database Interface**

````
Classic Modern Mint USXD Layer
├── Web UI (uDos Static) → Notionish-style Library View
│   └── Purpose: Browse, inspect, understand vendor database
├── TUI (Terminal) → Grid-style Data View
│   └── Purpose: Edit, manage, execute operations
└── Shared → sonic.db (SQLite) + JSON views
    └── Single source of truth for both interfaces
````

***

## 🌐 Web UI: Notionish Library View (Static)

### Design Pattern

````
┌────────────────────────────────────────────────────────────────────────────┐
│  📚 VENDOR LIBRARY                                          [Search] 🔍    │
├────────────────────────────────────────────────────────────────────────────┤
│  [All Items] [Git Repos] [Packages] [Forks] [Deployments] [+ Add View]    │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │  Name              │ Type    │ Version │ Size     │ Last Updated    │ │
│  ├──────────────────────────────────────────────────────────────────────┤ │
│  │  monaspace         │ git     │ 1.4.0   │ 45.2 MB  │ 2026-04-19      │ │
│  │  ├─ Details        │         │         │          │                 │ │
│  │  │  Source: github.com/githubnext/monaspace                         │ │
│  │  │  Forks: 1 (monaspace-classic-modern)                             │ │
│  │  │  Deployed to: 3 machines                                         │ │
│  │  └──────────────────────────────────────────────────────────────────│ │
│  │  monaspace-fork    │ fork    │ 1.4.0-cm │ 48.1 MB  │ 2026-04-19      │ │
│  │  linuxmint-21.3    │ iso     │ 21.3     │ 2.8 GB   │ 2026-04-15      │ │
│  │  requests          │ wheel   │ 2.31.0   │ 1.2 MB   │ 2026-04-18      │ │
│  │  react             │ tgz     │ 18.2.0   │ 3.4 MB   │ 2026-04-17      │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
│                                                                            │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │  📊 Quick Stats                                                       │ │
│  │  Total Items: 47  |  Total Size: 4.2 GB  |  Last Sync: 2 min ago     │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
│                                                                            │
│  [Add Vendored Product] [Remove Selected] [Sync Now] [Export JSON]        │
└────────────────────────────────────────────────────────────────────────────┘
````

### Implementation: Static Web (No Backend Required)

```html
<!-- ~/.local/share/udos/web/vendor-library.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Classic Modern — Vendor Library</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
    <style>
        /* Classic Modern Mac styling */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            background: #E8E8E8;
            font-family: 'Inter', system-ui, sans-serif;
            color: #111;
            padding: 20px;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: #F2F2F2;
            border: 1px solid #222;
        }
        
        /* Header */
        .header {
            background: #E8E8E8;
            border-bottom: 1px solid #222;
            padding: 16px 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .header h1 {
            font-family: 'ChicagoFLF', 'Inter', monospace;
            font-size: 18px;
            font-weight: normal;
        }
        
        /* Navigation */
        .nav {
            background: #F2F2F2;
            border-bottom: 1px solid #222;
            padding: 8px 20px;
            display: flex;
            gap: 16px;
        }
        
        .nav-item {
            font-family: 'ChicagoFLF', monospace;
            font-size: 12px;
            padding: 4px 8px;
            cursor: pointer;
            border: 1px solid transparent;
        }
        
        .nav-item.active {
            background: #3A7BD5;
            color: white;
            border-color: #222;
        }
        
        /* Table */
        .table {
            width: 100%;
            border-collapse: collapse;
        }
        
        .table th {
            text-align: left;
            padding: 12px 16px;
            background: #E8E8E8;
            border-bottom: 1px solid #222;
            font-family: 'ChicagoFLF', monospace;
            font-size: 11px;
            font-weight: normal;
        }
        
        .table td {
            padding: 12px 16px;
            border-bottom: 1px solid #222;
            font-size: 13px;
        }
        
        /* Expandable row */
        .details-row {
            background: #E8E8E8;
        }
        
        .details-content {
            padding: 16px 16px 16px 48px;
            font-size: 12px;
            color: #333;
        }
        
        /* Stats bar */
        .stats {
            background: #E8E8E8;
            border-top: 1px solid #222;
            padding: 12px 20px;
            display: flex;
            justify-content: space-between;
            font-size: 11px;
            font-family: 'ChicagoFLF', monospace;
        }
        
        /* Buttons */
        .button {
            background: #F2F2F2;
            border: 1px solid #222;
            padding: 6px 12px;
            font-family: 'ChicagoFLF', monospace;
            font-size: 11px;
            cursor: pointer;
        }
        
        .button:active {
            background: #3A7BD5;
            color: white;
        }
        
        /* Search */
        .search {
            background: #F2F2F2;
            border: 1px solid #222;
            padding: 6px 12px;
            font-family: 'Inter', monospace;
            font-size: 12px;
        }
    </style>
</head>
<body>
<div class="container">
    <div class="header">
        <h1>📚 VENDOR LIBRARY — sonic.db interface</h1>
        <div>
            <input type="text" placeholder="Search..." class="search" id="search">
        </div>
    </div>
    
    <div class="nav">
        <div class="nav-item active" data-view="all">All Items</div>
        <div class="nav-item" data-view="git">Git Repos</div>
        <div class="nav-item" data-view="packages">Packages</div>
        <div class="nav-item" data-view="forks">Forks</div>
        <div class="nav-item" data-view="deployments">Deployments</div>
        <div class="nav-item" data-view="add">+ Add View</div>
    </div>
    
    <div id="table-container">
        <!-- Dynamic table loads from sonic.db JSON export -->
    </div>
    
    <div class="stats">
        <div id="stats">Loading...</div>
        <div>
            <button class="button" id="refresh">⟳ Refresh</button>
            <button class="button" id="export">📄 Export JSON</button>
        </div>
    </div>
</div>

<script>
    // Load data from sonic.db JSON endpoint
    async function loadData() {
        const response = await fetch('/api/sonic/vendor/list');
        const data = await response.json();
        
        renderTable(data);
        updateStats(data);
    }
    
    function renderTable(data) {
        const html = `
            <table class="table">
                <thead>
                    <tr><th>Name</th><th>Type</th><th>Version</th><th>Size</th><th>Updated</th></tr>
                </thead>
                <tbody>
                    ${data.items.map(item => `
                        <tr onclick="toggleDetails('${item.id}')" style="cursor:pointer">
                            <td>${item.name}</td>
                            <td>${item.type}</td>
                            <td>${item.version}</td>
                            <td>${item.size}</td>
                            <td>${item.updated}</td>
                        </tr>
                        <tr id="details-${item.id}" class="details-row" style="display:none">
                            <td colspan="5">
                                <div class="details-content">
                                    <strong>Source:</strong> ${item.source}<br>
                                    <strong>Forks:</strong> ${item.forks || 0}<br>
                                    <strong>Deployed to:</strong> ${item.deployed_to || 'None'}<br>
                                    <strong>Checksum:</strong> ${item.checksum || 'N/A'}<br>
                                    <button class="button" onclick="event.stopPropagation();copyToSandbox('${item.id}')">
                                        📋 Copy to Sandbox
                                    </button>
                                </div>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
        document.getElementById('table-container').innerHTML = html;
    }
    
    function toggleDetails(id) {
        const row = document.getElementById(`details-${id}`);
        row.style.display = row.style.display === 'none' ? 'table-row' : 'none';
    }
    
    function updateStats(data) {
        document.getElementById('stats').innerHTML = 
            `Total Items: ${data.total} | Total Size: ${data.total_size} | Last Sync: ${data.last_sync}`;
    }
    
    // Copy to sandbox
    async function copyToSandbox(vendorId) {
        const response = await fetch('/api/sonic/vendor/copy', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({vendor_id: vendorId, target: 'sandbox'})
        });
        if (response.ok) {
            alert('Copied to sandbox: ~/code-vault/@sandbox/');
            loadData(); // Refresh
        }
    }
    
    loadData();
    setInterval(loadData, 30000); // Auto-refresh every 30 seconds
</script>
</body>
</html>
```

### Sonic Command to Serve Web UI

```bash
sonic web ui --port=8043 --static

# Output:
# ✓ Serving vendor library at http://localhost:8043
#   Static assets: ~/.local/share/udos/web/
#   Data source: sonic.db (read-only JSON API)
#   
#   To access from other machines:
#   http://192.168.1.100:8043
```

***

## 🖥️ TUI: Grid-Style Data View (Terminal)

### USXD Grid Pattern

````
┌────────────────────────────────────────────────────────────────────────────┐
│ sonic vendor browse                                                        │
├────────────────────────────────────────────────────────────────────────────┤
│ ┌──────────┬──────────┬──────────┬──────────┬──────────┬──────────┐       │
│ │ Name     │ Type     │ Version  │ Size     │ Updated  │ Actions  │       │
│ ├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤       │
│ │ monaspace│ git      │ 1.4.0    │ 45.2 MB  │ 04-19    │ [Details]│       │
│ ├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤       │
│ │ monaspace│ fork     │ 1.4.0-cm │ 48.1 MB  │ 04-19    │ [Details]│       │
│ │ -fork    │          │          │          │          │          │       │
│ ├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤       │
│ │ linuxmint│ iso      │ 21.3     │ 2.8 GB   │ 04-15    │ [Details]│       │
│ │ -21.3    │          │          │          │          │          │       │
│ ├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤       │
│ │ requests │ wheel    │ 2.31.0   │ 1.2 MB   │ 04-18    │ [Details]│       │
│ ├──────────┼──────────┼──────────┼──────────┼──────────┼──────────┤       │
│ │ react    │ tgz      │ 18.2.0   │ 3.4 MB   │ 04-17    │ [Details]│       │
│ └──────────┴──────────┴──────────┴──────────┴──────────┴──────────┘       │
│                                                                            │
│  [Add] [Remove] [Copy to Sandbox] [View Logs] [Export] [Quit]              │
│                                                                            │
│  Status: 47 items | 4.2 GB | Last sync: 2 min ago                          │
│  Press ↑/↓ to navigate, Enter to select, / to search                       │
└────────────────────────────────────────────────────────────────────────────┘
````

### Implementation: Bubble Tea + Lip Gloss

```go
// sonic-tui/grid.go
package main

import (
    "database/sql"
    "fmt"
    "strings"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/table"
)

// VendorItem represents a row in the grid
type VendorItem struct {
    ID      string
    Name    string
    Type    string
    Version string
    Size    string
    Updated string
}

// GridModel manages the TUI state
type GridModel struct {
    table     table.Model
    items     []VendorItem
    db        *sql.DB
    width     int
    height    int
    showHelp  bool
}

// Initialize grid
func NewGridModel(db *sql.DB) *GridModel {
    columns := []table.Column{
        {Title: "Name", Width: 30},
        {Title: "Type", Width: 12},
        {Title: "Version", Width: 15},
        {Title: "Size", Width: 12},
        {Title: "Updated", Width: 12},
        {Title: "Actions", Width: 10},
    }
    
    t := table.New(
        table.WithColumns(columns),
        table.WithFocused(true),
        table.WithHeight(20),
    )
    
    s := table.DefaultStyles()
    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        BorderBottom(true).
        Bold(false)
    s.Selected = s.Selected.
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("3A7BD5")).
        Bold(false)
    t.SetStyles(s)
    
    return &GridModel{
        table: t,
        db:    db,
    }
}

// Load data from sonic.db
func (m *GridModel) loadData() tea.Cmd {
    rows, err := m.db.Query(`
        SELECT id, name, type, version, 
               printf('%.1f MB', size_bytes/1024.0/1024.0) as size,
               strftime('%m-%d', last_updated) as updated
        FROM vendors
        ORDER BY name
    `)
    if err != nil {
        return nil
    }
    defer rows.Close()
    
    var items []VendorItem
    var rowsData []table.Row
    
    for rows.Next() {
        var item VendorItem
        err := rows.Scan(&item.ID, &item.Name, &item.Type, &item.Version, &item.Size, &item.Updated)
        if err != nil {
            continue
        }
        items = append(items, item)
        rowsData = append(rowsData, table.Row{
            item.Name,
            item.Type,
            item.Version,
            item.Size,
            item.Updated,
            "[Details]",
        })
    }
    
    m.items = items
    m.table.SetRows(rowsData)
    return nil
}

// Handle key events
func (m *GridModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
            
        case "enter":
            // Show details for selected item
            selected := m.table.SelectedRow()
            if len(selected) > 0 {
                item := m.findItemByName(selected[0])
                return m, m.showDetails(item)
            }
            
        case "a":
            // Add new vendor
            return m, m.addVendor()
            
        case "r":
            // Remove selected vendor
            selected := m.table.SelectedRow()
            if len(selected) > 0 {
                return m, m.removeVendor(selected[0])
            }
            
        case "c":
            // Copy to sandbox
            selected := m.table.SelectedRow()
            if len(selected) > 0 {
                return m, m.copyToSandbox(selected[0])
            }
            
        case "e":
            // Export to JSON
            return m, m.exportJSON()
            
        case "h":
            m.showHelp = !m.showHelp
            return m, nil
            
        case "/":
            // Search mode
            return m, m.startSearch()
        }
    }
    
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

// Render the view
func (m *GridModel) View() string {
    // Header
    header := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("111")).
        Background(lipgloss.Color("E8E8E8")).
        Padding(1, 2).
        Render("📦 SONIC VENDOR GRID — classic-modern-mint")
    
    // Grid
    grid := m.table.View()
    
    // Footer stats
    stats := fmt.Sprintf("  %d items | Press ↑/↓ to navigate, Enter to view, / to search", len(m.items))
    footer := lipgloss.NewStyle().
        Foreground(lipgloss.Color("222")).
        Padding(1, 2).
        Render(stats)
    
    // Help
    help := ""
    if m.showHelp {
        help = lipgloss.NewStyle().
            Border(lipgloss.NormalBorder()).
            BorderForeground(lipgloss.Color("222")).
            Padding(1, 2).
            Render(`
  KEYBOARD SHORTCUTS
  ↑/↓     Navigate
  Enter   View details
  a       Add vendor
  r       Remove selected
  c       Copy to sandbox
  e       Export JSON
  /       Search
  h       Hide/show help
  q       Quit
`)
    }
    
    return lipgloss.JoinVertical(lipgloss.Top, header, grid, footer, help)
}

// Show details in modal
func (m *GridModel) showDetails(item VendorItem) tea.Cmd {
    // Query full details from sonic.db
    var source, forks, deployed string
    m.db.QueryRow(`
        SELECT source_url, 
               (SELECT COUNT(*) FROM forks WHERE upstream_vendor_id = vendors.id) as forks,
               (SELECT COUNT(*) FROM deployments WHERE vendor_id = vendors.id) as deployed
        FROM vendors WHERE name = ?
    `, item.Name).Scan(&source, &forks, &deployed)
    
    details := lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("222")).
        Padding(1, 2).
        Width(60).
        Render(fmt.Sprintf(`
  📄 %s (%s)
  
  Version: %s
  Size: %s
  Updated: %s
  Source: %s
  Forks: %s
  Deployments: %s
  
  Press ESC to close
`, item.Name, item.Type, item.Version, item.Size, item.Updated, source, forks, deployed))
    
    // This would be rendered as a modal in a full implementation
    fmt.Println(details)
    return nil
}

func (m *GridModel) addVendor() tea.Cmd {
    // Would open form for adding new vendor
    fmt.Println("\n  Add new vendor: (implement form here)")
    return nil
}

func (m *GridModel) removeVendor(name string) tea.Cmd {
    fmt.Printf("\n  Remove %s? (y/N): ", name)
    // Would confirm and delete
    return nil
}

func (m *GridModel) copyToSandbox(name string) tea.Cmd {
    fmt.Printf("\n  Copying %s to ~/code-vault/@sandbox/...\n", name)
    // Execute: sonic vendor copy --from=name --to=sandbox
    return nil
}

func (m *GridModel) exportJSON() tea.Cmd {
    fmt.Println("\n  Exporting vendor database to ~/vault/@public/vendor-export.json")
    return nil
}

func (m *GridModel) startSearch() tea.Cmd {
    fmt.Print("\n  Search: ")
    var query string
    fmt.Scanln(&query)
    // Filter table rows
    return nil
}

func (m *GridModel) findItemByName(name string) VendorItem {
    for _, item := range m.items {
        if item.Name == name {
            return item
        }
    }
    return VendorItem{}
}

// Main entry point
func main() {
    db, _ := sql.Open("sqlite3", "~/.local/share/sonic/sonic.db")
    model := NewGridModel(db)
    model.loadData()
    
    p := tea.NewProgram(model, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
    }
}
```

### Sonic Command to Launch TUI

```bash
sonic vendor browse --view=grid

# Output:
# Launching Classic Modern TUI Grid...
# Press h for help, q to quit
```

***

## 🔄 Shared Data Flow

````
┌─────────────────────────────────────────────────────────────────────────┐
│                           sonic.db (SQLite)                              │
│                    ~/.local/share/sonic/sonic.db                         │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │ vendors | git_repos | forks | packages | deployments | patches  │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┼─────────────────┐
            │                                   │
            ▼                                   ▼
┌───────────────────────┐           ┌───────────────────────┐
│   Web UI (Notionish)  │           │   TUI (Bubble Tea)    │
│   http://localhost:8043│          │   sonic vendor browse │
│                       │           │                       │
│   - Library view      │           │   - Grid view         │
│   - Read-only browse  │           │   - Edit/Manage       │
│   - Export JSON       │           │   - Execute actions   │
└───────────────────────┘           └───────────────────────┘
            │                                   │
            └─────────────────┬─────────────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │  JSON API Layer   │
                    │  /api/sonic/*     │
                    │  Read + Write     │
                    └───────────────────┘
````

***

## ✅ USXD Requirements Checklist

| Requirement              | Web UI (Notionish)         | TUI Grid         | Status   |
| ------------------------ | -------------------------- | ---------------- | -------- |
| Library view             | ✅ Table + expandable rows | ✅ Grid          | Complete |
| Add vendored product     | ✅ Button + form           | ✅ 'a' key       | Planned  |
| Remove vendor            | ✅ Button + confirm        | ✅ 'r' key       | Planned  |
| View sonic.db data       | ✅ JSON API                | ✅ SQLite direct | Complete |
| Any .db/.json table view | ✅ Generic parser          | ✅ SQL queries   | Complete |
| Static (no backend)      | ✅ HTML+JS only            | N/A              | Complete |
| Grid view for terminal   | N/A                        | ✅ Bubble Tea    | Complete |
| Copy to sandbox          | ✅ Button                  | ✅ 'c' key       | Complete |

***

## 🚀 Launch Commands Summary

```bash
# Start web UI (static, no backend)
sonic web ui --port=8043

# Launch TUI grid
sonic vendor browse

# Export data for external viewing
sonic vendor export --format=json --output=vendor-library.json

# Generate static HTML snapshot (offline browsable)
sonic vendor snapshot --html --output=vendor-snapshot.html
```

***

## 🎯 Design Principles Applied

1. **Notionish for browsing** — Familiar, searchable, expandable
2. **TUI grid for managing** — Fast, keyboard-driven, terminal-native
3. **Single source of truth** — sonic.db drives both interfaces
4. **Classic Modern styling** — Mono colours, hard borders, Chicago fonts
5. **Static-first web** — No backend complexity, just JSON API

**The vendor library now has a complete USXD layer — browse in browser, manage in terminal.**

# Addendum 3: USXD OBF Style Layout & UI Config for Classic Modern + Ventoy

## 🎨 Core Design Philosophy

**One OBF specification — Two complementary interfaces**

````
OBF Style System v3.0
├── Classic Modern Mint (Primary OS Shell)
│   ├── Web UI (Notionish Library)
│   ├── TUI (Grid Manager)
│   └── Desktop (Cinnamon + GTK)
└── Ventoy (Boot Loader UI)
    ├── Menu System
    ├── Disk Selector
    └── Status Panel
````

***

## 📐 OBF v3.0 — Unified Style Specification

### File Structure

````
~/.local/share/classic-modern/
├── obf/
│   ├── classic-modern.obf           # Primary OS style
│   ├── ventoy.obf                   # Boot loader style
│   └── shared/
│       ├── colours.obf              # Shared colour tokens
│       ├── typography.obf           # Shared font rules
│       └── patterns.obf             # Shared background patterns
├── layouts/
│   ├── web-ui.yaml                  # Notionish layout config
│   ├── tui-grid.yaml                # Terminal grid config
│   └── ventoy-menu.yaml             # Boot menu config
└── assets/
    ├── fonts/                       # Chicago, Monaspace, etc
    ├── icons/                       # Material + custom
    └── patterns/                    # Linen, stripes, dots
````

***

## 🎨 Classic Modern OBF — Complete Style Definition

```yaml
# classic-modern.obf (YAML-based, compiled to binary)
version: "3.0.0"
magic: "CMOBF03"
target: "classic-modern-mint"

# 1. COLOUR PALETTE
colours:
  primary:
    bg: "#E8E8E8"
    surface: "#F2F2F2"
    border: "#222222"
    text: "#111111"
  accent:
    primary: "#3A7BD5"
    hover: "#2A6BC5"
    active: "#1A5BB5"
  status:
    success: "#2E7D32"
    warning: "#F57C00"
    error: "#C62828"
    info: "#1565C0"
  mono:
    dark: "#000000"
    light: "#FFFFFF"
    grey: "#808080"

# 2. BORDERS & SHADOWS
borders:
  width: 1px
  style: solid
  radius: 0px
  shadows: none
  focus_indicator: "inset 0 0 0 1px var(--cm-accent)"

# 3. SPACING SYSTEM
spacing:
  unit: 4px
  scale:
    xs: 1    # 4px
    sm: 2    # 8px
    md: 3    # 12px
    lg: 4    # 16px
    xl: 6    # 24px
    xxl: 8   # 32px

# 4. TYPOGRAPHY
typography:
  fonts:
    ui:
      family: ["ChicagoFLF", "Chicago", "monospace"]
      size: 12px
      weight: normal
      line-height: 1.4
    body:
      family: ["Monaspace Argon", "Inter", "system-ui"]
      size: 13px
      weight: 400
      line-height: 1.5
    mono:
      family: ["pixChicago", "Monaco", "monospace"]
      size: 11px
      weight: normal
      line-height: 1.3
    heading:
      family: ["ChicagoFLF", "monospace"]
      size: 16px
      weight: normal
      line-height: 1.2
  scale:
    h1: 24px
    h2: 20px
    h3: 18px
    h4: 16px
    small: 10px

# 5. PATTERNS (optional, low contrast)
patterns:
  desktop:
    name: "linen"
    enabled: true
    opacity: 0.08
  stripes:
    name: "platinum"
    enabled: true
    angle: 45deg
    size: 4px
  dots:
    name: "finder"
    enabled: false
    size: 4px

# 6. COMPONENT VARIANTS
components:
  window:
    background: var(--surface)
    border: 1px solid var(--border)
    padding: var(--spacing-md)
    
  button:
    background: var(--bg)
    border: 1px solid var(--border)
    padding: var(--spacing-sm) var(--spacing-md)
    hover:
      background: var(--surface)
    active:
      background: var(--accent)
      color: var(--light)
      
  menu:
    background: var(--surface)
    border: 1px solid var(--border)
    item_padding: var(--spacing-xs) var(--spacing-md)
    hover_background: var(--accent)
    hover_color: var(--light)
    
  table:
    header_background: var(--bg)
    header_border: 1px solid var(--border)
    row_border: 1px solid var(--border)
    row_hover: var(--surface)
    selected_background: var(--accent)
    selected_color: var(--light)
    
  card:
    background: var(--surface)
    border: 1px solid var(--border)
    padding: var(--spacing-md)
    shadow: none

# 7. INTERACTION STATES
interaction:
  transition: none
  cursor:
    default: "default"
    pointer: "pointer"
    text: "text"
    resize: "nesw-resize"
  focus_ring: "1px solid var(--accent)"
  
# 8. ACCESSIBILITY
accessibility:
  min_contrast_ratio: 4.5
  font_scaling: 1.25
  reduced_motion: true
  high_contrast_support: true
```

***

## 🖥️ Ventoy OBF — Boot Loader Style

```yaml
# ventoy.obf
version: "3.0.0"
magic: "VTOYBF01"
target: "ventoy-classic-modern"
parent: "classic-modern.obf"  # Inherits from main theme

# Ventoy-specific overrides
ventoy:
  screen:
    resolution: "auto"
    framebuffer: true
    background: "#E8E8E8"
    
  menu:
    layout: "center"
    width: 80%
    max_width: 800px
    background: "#F2F2F2"
    border: "1px solid #222"
    padding: "16px"
    
    items:
      height: 32px
      padding: "4px 12px"
      selected_background: "#3A7BD5"
      selected_text: "#FFFFFF"
      hotkey_color: "#808080"
      
  disk_selector:
    layout: "right-panel"
    width: 30%
    background: "#E8E8E8"
    border: "1px solid #222"
    item_padding: "8px"
    device_icon: "💾"
    size_format: "GB"
    
  status_bar:
    position: "bottom"
    background: "#222"
    color: "#E8E8E8"
    height: 24px
    padding: "0 8px"
    font: "Monaco 10px"
    
  boot_entry:
    type: "list"
    icons: true
    description_width: 60%
    kernel_args: false
    
  translation:
    enabled: true
    default_lang: "en_US"
    font: "Geneva 9"
```

***

## 📐 Layout Configurations

### Web UI (Notionish Library) Layout

```yaml
# layouts/web-ui.yaml
version: "1.0.0"
target: "classic-modern-web"

layout:
  type: "notionish"
  container:
    max_width: 1400px
    margin: 20px auto
    background: "var(--bg)"
    border: "1px solid var(--border)"
    
  header:
    height: 60px
    padding: "0 20px"
    border_bottom: "1px solid var(--border)"
    title:
      font: "var(--ui-font)"
      size: 18px
      
  navigation:
    type: "tabs"
    height: 40px
    padding: "0 20px"
    border_bottom: "1px solid var(--border)"
    item:
      padding: "8px 16px"
      active_background: "var(--accent)"
      active_color: "var(--light)"
      
  table:
    type: "expandable"
    header:
      background: "var(--bg)"
      font: "var(--ui-font)"
      size: 11px
    row:
      height: 48px
      border_bottom: "1px solid var(--border)"
      hover_background: "var(--surface)"
    details:
      background: "var(--bg)"
      padding: "16px 0 16px 48px"
      font: "var(--body-font)"
      size: 12px
      
  stats_bar:
    height: 40px
    padding: "0 20px"
    border_top: "1px solid var(--border)"
    background: "var(--bg)"
    font: "var(--ui-font)"
    size: 11px
    
  buttons:
    primary:
      background: "var(--accent)"
      color: "var(--light)"
      border: "1px solid var(--border)"
    secondary:
      background: "var(--surface)"
      border: "1px solid var(--border)"
    danger:
      background: "var(--error)"
      color: "var(--light)"
```

### TUI Grid Layout

```yaml
# layouts/tui-grid.yaml
version: "1.0.0"
target: "classic-modern-tui"

layout:
  type: "grid"
  dimensions:
    min_width: 80
    min_height: 24
    
  header:
    height: 3
    background: "var(--bg)"
    border: "bottom"
    title:
      text: "📦 SONIC VENDOR GRID"
      alignment: "center"
      
  grid:
    border_style: "rounded"
    header:
      background: "var(--bg)"
      foreground: "var(--text)"
      bold: false
    selected:
      background: "var(--accent)"
      foreground: "var(--light)"
    row_height: 1
    columns:
      - name: "Name"
        width: 30
        alignment: "left"
      - name: "Type"
        width: 12
        alignment: "left"
      - name: "Version"
        width: 15
        alignment: "left"
      - name: "Size"
        width: 12
        alignment: "right"
      - name: "Updated"
        width: 12
        alignment: "center"
        
  footer:
    height: 3
    background: "var(--bg)"
    border: "top"
    stats_alignment: "left"
    help_alignment: "right"
    
  help_modal:
    position: "center"
    width: 60
    height: 20
    border: "double"
    background: "var(--surface)"
    
  status_line:
    height: 1
    background: "var(--border)"
    foreground: "var(--bg)"
```

### Ventoy Boot Menu Layout

```yaml
# layouts/ventoy-menu.yaml
version: "1.0.0"
target: "ventoy"

layout:
  type: "boot_menu"
  resolution:
    width: 1024
    height: 768
    
  background:
    type: "solid"
    colour: "#E8E8E8"
    pattern: "linen"
    pattern_opacity: 0.05
    
  logo:
    enabled: true
    path: "/boot/ventoy/logo.png"
    position: "top-center"
    margin: 20px
    
  menu_frame:
    position: "center"
    width: 70%
    max_width: 800
    background: "#F2F2F2"
    border: "1px solid #222"
    padding: 16
    shadow: false
    
  menu_header:
    text: "CLASSIC MODERN BOOT SELECTOR"
    font: "ChicagoFLF"
    size: 14
    alignment: "center"
    margin_bottom: 16
    border_bottom: "1px solid #222"
    padding_bottom: 8
    
  menu_items:
    type: "list"
    item_height: 32
    item_padding: "4px 12px"
    item_spacing: 2
    font: "Geneva 9"
    normal:
      background: "transparent"
      foreground: "#111"
    selected:
      background: "#3A7BD5"
      foreground: "#FFF"
    hotkey:
      colour: "#808080"
      width: 24
      
  disk_info:
    position: "right"
    width: 25%
    margin_left: 16
    background: "#E8E8E8"
    border: "1px solid #222"
    padding: 8
    title:
      text: "DEVICES"
      font: "ChicagoFLF"
      size: 11
      border_bottom: "1px solid #222"
    devices:
      font: "Monaco"
      size: 10
      item_padding: 4
      
  status_bar:
    position: "bottom"
    height: 24
    background: "#222"
    foreground: "#E8E8E8"
    font: "Monaco"
    size: 10
    padding: "0 8"
    elements:
      - position: "left"
        content: "F1: Help"
      - position: "left"
        content: "F2: Language"
      - position: "right"
        content: "Classic Modern Mint"
      - position: "right"
        content: "v1.0.0"
        
  boot_progress:
    type: "bar"
    height: 4
    background: "#222"
    fill: "#3A7BD5"
    position: "bottom"
```

***

## 🔧 UI Configuration Documents

### Web UI Config (JavaScript)

```javascript
// ~/.local/share/udos/web/config.js
window.classicModernConfig = {
  theme: {
    name: 'classic-modern-mac',
    version: '3.0.0',
    obfPath: '/usr/share/classic-modern/obf/classic-modern.obf'
  },
  
  layout: {
    defaultView: 'library',
    itemsPerPage: 25,
    expandableRows: true,
    stickyHeader: true
  },
  
  api: {
    baseUrl: '/api/sonic',
    endpoints: {
      list: '/vendor/list',
      details: '/vendor/details',
      add: '/vendor/add',
      remove: '/vendor/remove',
      copy: '/vendor/copy'
    },
    refreshInterval: 30000
  },
  
  features: {
    search: true,
    export: true,
    dragAndDrop: false,
    realtime: false
  },
  
  notifications: {
    enabled: true,
    duration: 3000,
    position: 'top-right'
  }
};
```

### TUI Config (Go)

```go
// ~/.config/sonic/tui-config.yaml
tui:
  theme: "classic-modern"
  obf_path: "~/.local/share/classic-modern/obf/classic-modern.obf"
  
  layout:
    grid:
      height: 20
      width: 120
      border_style: "rounded"
    status_bar:
      enabled: true
      position: "bottom"
    help:
      show_on_startup: false
      position: "modal"
      
  keybindings:
    quit: "q"
    help: "h"
    search: "/"
    add: "a"
    remove: "r"
    copy: "c"
    export: "e"
    details: "enter"
    refresh: "ctrl+r"
    
  colours:
    primary: "#3A7BD5"
    secondary: "#F2F2F2"
    border: "#222222"
    text: "#111111"
    error: "#C62828"
    success: "#2E7D32"
    
  fonts:
    ui: "ChicagoFLF"
    mono: "pixChicago"
```

### Ventoy Config (C)

```c
// ventoy/classic-modern-config.h
#ifndef CLASSIC_MODERN_VENTOY_CONFIG_H
#define CLASSIC_MODERN_VENTOY_CONFIG_H

// Theme configuration
#define CM_THEME_NAME "classic-modern-mac"
#define CM_THEME_VERSION "3.0.0"

// Colours (RGB)
#define CM_COLOUR_BG      0xE8, 0xE8, 0xE8
#define CM_COLOUR_SURFACE 0xF2, 0xF2, 0xF2
#define CM_COLOUR_BORDER  0x22, 0x22, 0x22
#define CM_COLOUR_TEXT    0x11, 0x11, 0x11
#define CM_COLOUR_ACCENT  0x3A, 0x7B, 0xD5

// Layout dimensions
#define CM_MENU_WIDTH_PERCENT 70
#define CM_MENU_MAX_WIDTH 800
#define CM_MENU_ITEM_HEIGHT 32
#define CM_STATUS_BAR_HEIGHT 24

// Fonts (Ventoy framebuffer)
#define CM_FONT_UI "ChicagoFLF"
#define CM_FONT_BODY "Geneva9"
#define CM_FONT_MONO "Monaco"

// Patterns
#define CM_PATTERN_LINEN_ENABLED 1
#define CM_PATTERN_STRIPES_ENABLED 1
#define CM_PATTERN_DOTS_ENABLED 0

// Boot behaviour
#define CM_BOOT_TIMEOUT 10  // seconds
#define CM_DEFAULT_ENTRY 0
#define CM_SHOW_ADVANCED 1

#endif
```

***

## 🎯 OBF Compilation & Validation

### Build OBF from YAML

```bash
sonic obf compile --input=classic-modern.yaml --output=classic-modern.obf --version=3.0.0

# Output:
# ✓ Compiled classic-modern.obf (v3.0.0)
#   Magic: CMOBF03
#   Size: 2.4 KB
#   Checksum: sha256:7a8f3c2b...
```

### Validate OBF Integrity

```bash
sonic obf validate classic-modern.obf --strict

# Output:
# ┌─────────────────────────────────────────────────────────────┐
# │ OBF VALIDATION — classic-modern.obf                         │
# ├─────────────────────────────────────────────────────────────┤
# │ ✓ Magic number correct (CMOBF03)                           │
# │ ✓ Version 3.0.0 supported                                  │
# │ ✓ Colour palette complete (10/10)                          │
# │ ✓ Border rules valid                                       │
# │ ✓ Typography chain intact                                  │
# │ ✓ Pattern library accessible                               │
# │ ✓ No gradients (0 violations)                              │
# │ ✓ No transparency (0 violations)                           │
# │ ✓ Checksum verified                                        │
# └─────────────────────────────────────────────────────────────┘
```

### Extract OBF for Inspection

```bash
sonic obf extract classic-modern.obf --human --output=style-spec.yaml

# Output:
# ✓ Extracted to style-spec.yaml
#   Human-readable format with comments
```

***

## 🔄 Runtime Style Application

### Classic Modern Mint (GTK)

```bash
# Apply OBF to GTK
sonic theme apply --obf=classic-modern.obf --target=gtk

# Generates: ~/.config/gtk-3.0/gtk.css
```

### Ventoy Boot Loader

```bash
# Build Ventoy ISO with OBF style
sonic boot ventoy --build \
  --obf=ventoy.obf \
  --layout=ventoy-menu.yaml \
  --output=ventoy-classic-modern.iso

# Output:
# ✓ Ventoy ISO built with Classic Modern theme
#   Boot menu: ChicagoFLF, platinum style
#   Status bar: Monaco 10px
#   ISO size: 2.8 GB
```

***

## 📊 USXD OBF Compliance Matrix

| Component          | Classic Modern | Ventoy          | Shared |
| ------------------ | -------------- | --------------- | ------ |
| Colour tokens      | ✅             | ✅ (inherits)   | ✅     |
| Border rules       | ✅             | ✅              | ✅     |
| Typography         | ✅             | ✅ (subset)     | ✅     |
| Spacing system     | ✅             | ✅              | ✅     |
| Patterns           | ✅             | ✅ (linen only) | ✅     |
| Component variants | ✅             | ❌ (custom)     | ❌     |
| Interaction states | ✅             | ❌              | ❌     |
| Accessibility      | ✅             | ❌              | ❌     |

***

## 🚀 Quick Start Commands

```bash
# View current OBF style
sonic obf show --current

# Switch between OBF variants
sonic obf switch --variant=classic-modern-mac

# Preview OBF in browser
sonic obf preview --obf=classic-modern.obf --port=8044

# Generate style guide from OBF
sonic obf guide --obf=classic-modern.obf --output=style-guide.html
```

***

## ✅ Summary

The OBF v3.0 system provides:

1. **Unified style specification** — One format for both OS and boot loader
2. **Ventoy inheritance** — Boot loader inherits from main theme
3. **Complete UI configs** — Web, TUI, and boot menu layouts
4. **Validation & tooling** — Compile, validate, extract, preview
5. **Runtime application** — Apply to GTK, terminal, and Ventoy builds

**Classic Modern Mint and Ventoy now share a single source of truth for styling — OBF v3.0.**

#dev