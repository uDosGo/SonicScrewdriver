# USX Token Reference — SonicScrewdriver Surfaces

**Created:** 2026-07-10
**Status:** sprint.2026-07-10 — `task.sonic.usx.001`
**Source:** `@udos/usx-tokens` v3.0.0 (canonical in `uCore/packages/usx-tokens/`)
**Purpose:** Subset of USX tokens applicable to Sonic CLI output and documentation surfaces.

---

## Status Colors (Semantic)

These are the critical tokens for Sonic's snackbar envelope→visual mapping:

| Token | USX Variable | Hex (Approx) | Rich Style | Usage |
|---|---|---|---|---|
| Success | `--usx-color-success` | `#4CAF50` | `bold green` | ✅ Completed operations |
| Warning | `--usx-color-warning` | `#FF9800` | `bold yellow` | ⚠ Destructive confirmations, fallbacks |
| Error | `--usx-color-error` | `#F44336` | `bold red` | ❌ Failures, panics, schema drift |
| Info | `--usx-color-info` | `#2196F3` | `bold cyan` | 🔍 Diagnostics, status, lookups |
| Primary | `--usx-color-primary` | `#1976D2` | `bold blue` | 🔨 Build, create, install |
| On-Surface | `--usx-color-on-surface` | `#FFFFFF`/`#121212` | `white` | Body text on surfaces |
| Surface | `--surface` | `#1E1E1E` | `on #1E1E1E` | Panel backgrounds |
| Surface-Container | `--surface-container` | `#2D2D2D` | `on #2D2D2D` | Card/table backgrounds |

---

## Status Icon Convention

Sonic CLI uses Rich markup emoji consistently:

| Snackbar Status | Icon | Rich Markup | Spool Level |
|---|---|---|---|
| `success` | ✅ | `[bold green]✅[/bold green]` | `INFO` |
| `warn` | ⚠ | `[bold yellow]⚠[/bold yellow]` | `WARNING` |
| `error` | 💥 | `[bold red]💥[/bold red]` | `ERROR` / `CRITICAL` |
| `info` | 🔍 | `[bold cyan]🔍[/bold cyan]` | `INFO` |

---

## Action Icons

| Action | Icon | Rich Markup |
|---|---|---|
| Create / Build | 🔨 | `[bold green]🔨[/bold green]` |
| Install | 💿 | `[bold]💿[/bold]` |
| Destroy / Wipe | 💥 | `[bold red]💥[/bold red]` |
| Remove | 🗑 | `[bold red]🗑[/bold red]` |
| Add | ➕ | `[bold]➕[/bold]` |
| Export | 📤 | `[bold]📤[/bold]` |
| Import | 📥 | `[bold]📥[/bold]` |
| Verify / Check | 🔍 | `[bold]🔍[/bold]` |
| List / Scan | 🔍 | `[bold]🔍[/bold]` |
| Flash | ⚡ | `[bold]⚡[/bold]` |
| Security | 🔐 | `[bold]🔐[/bold]` |
| Keygen | 🔑 | `[bold]🔑[/bold]` |
| Repurpose | 🔄 | `[bold]🔄[/bold]` |
| Build ISO | 🔧 | `[bold]🔧[/bold]` |
| Mesh | 🌐 | `[bold]🌐[/bold]` |
| Game / Launch | 🚀 | `[bold]🚀[/bold]` |
| Diagnostics | 📊 | `[bold blue]📊[/bold blue]` |

---

## Result Panel Convention

All commands use `rich.panel.Panel.fit()` with these conventions:

```python
Panel.fit(
    f"[status_icon] Message\n\n"
    f"  Detail 1:  value\n"
    f"  Detail 2:  value\n"
    f"  Next step text.",
    title="Command — Result",
    border_style=border,  # green=success, yellow=warn, red=error, blue=info
)
```

---

## Spacing Scale (For Rich Layouts)

When building table padding or panel spacing in Rich:

| USX Token | Value | Rich Equivalent |
|---|---|---|
| `--usx-spacing-xs` | 2px | N/A (Rich uses cell padding) |
| `--usx-spacing-sm` | 4px | `padding=(0, 1)` |
| `--usx-spacing-md` | 8px | `padding=(0, 2)` |
| `--usx-spacing-lg` | 16px | `padding=(1, 2)` |
| `--usx-spacing-xl` | 24px | `padding=(1, 4)` |
| `--usx-spacing-2xl` | 32px | `padding=(2, 4)` |

---

## Typography (For Docs)

When authoring operator-facing markdown documentation:

| USX Token | Value | Markdown Equivalent |
|---|---|---|
| `--usx-font-size-caption` | 12px | `> small text` or footnotes |
| `--usx-font-size-body` | 16px | Default paragraph text |
| `--usx-font-size-headline` | 24px | `## Heading 2` |
| `--usx-font-weight-regular` | 400 | Normal body |
| `--usx-font-weight-bold` | 700 | `**bold**` or `### Heading 3` |

---

## Surface Conventions for Docs

| Pattern | USX Alignment | Markdown |
|---|---|---|
| Success block | `--usx-color-success` bg | Not applicable in plain markdown — use `✅` prefix |
| Warning block | `--usx-color-warning` bg | `> ⚠ WARNING: ...` blockquote |
| Error block | `--usx-color-error` bg | `> ❌ ERROR: ...` blockquote |
| Info panel | `--surface-container` bg | Standard markdown table or fenced block |
| Card | `--surface-container` + `--usx-radius-md` | Markdown table |