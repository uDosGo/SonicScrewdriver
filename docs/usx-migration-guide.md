# USX Migration Guide — SonicScrewdriver Web/Vue Surfaces

**Created:** 2026-07-10
**Status:** sprint.2026-07-10 — `task.sonic.usx.004`
**Purpose:** Document how future Sonic web control surfaces adopt `@udos/usx-tokens`.

---

## Overview

SonicScrewdriver is CLI-first, but future operator surfaces may include:
- A web-based device dashboard (device inventory, USB health)
- A spool/event viewer (chronological event stream)
- A Mint ISO customization wizard

When these are built, they should align with uCore's USX token conventions
rather than inventing new one-off styles.

---

## Step 1: Project Scaffold

```bash
# Create a Vue + Vite project
npm create vite@latest sonic-dashboard -- --template vue-ts
cd sonic-dashboard
npm install @udos/usx-tokens
```

---

## Step 2: Import USX Tokens

In `src/main.ts` or `src/style.css`:

```css
@import '@udos/usx-tokens/usx-standard.css';
```

This brings in:
- All USX CSS custom properties (colors, spacing, typography, touch targets)
- Base reset and layout primitives
- Theme support (dark/light via `data-theme` attribute on `<html>`)

---

## Step 3: Apply USX Variables

### Status Badges

```vue
<template>
  <span
    class="status-badge"
    :class="`status-${status}`"
  >
    {{ label }}
  </span>
</template>

<style scoped>
.status-badge {
  padding: var(--usx-spacing-xs) var(--usx-spacing-sm);
  border-radius: var(--usx-radius-sm);
  font-size: var(--usx-font-size-caption);
  font-weight: var(--usx-font-weight-bold);
}

.status-success {
  background: var(--usx-color-success);
  color: var(--usx-color-on-primary);
}

.status-warn {
  background: var(--usx-color-warning);
  color: var(--usx-color-on-primary);
}

.status-error {
  background: var(--usx-color-error);
  color: var(--usx-color-on-primary);
}
</style>
```

### Card Components

```vue
<template>
  <div class="usx-card">
    <h3 class="usx-card-title">{{ title }}</h3>
    <div class="usx-card-body">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.usx-card {
  background: var(--surface-container);
  border: 1px solid var(--border-subtle);
  border-radius: var(--usx-radius-md);
  padding: var(--usx-spacing-lg);
}

.usx-card-title {
  font-size: var(--usx-font-size-headline);
  font-weight: var(--usx-font-weight-bold);
  color: var(--usx-color-on-surface);
  margin-bottom: var(--usx-spacing-md);
}

.usx-card-body {
  font-size: var(--usx-font-size-body);
  color: var(--usx-color-on-surface);
}
</style>
```

### Device Table (Rich→USX equivalent)

The CLI `rich.table.Table` patterns map directly to HTML tables:

| Rich Pattern | Vue/Vite Equivalent |
|---|---|
| `table.add_column("Device", style="cyan")` | `<th class="col-cyan">Device</th>` |
| `console.print(table)` | Render `<Table>` component with dynamic rows |
| `Panel.fit(...)` | `<div class="usx-card">` wrapper |

---

## Step 4: Match Snackbar Envelope Shape

When building a spool event viewer, match the `SpoolEvent` dataclass shape:

```typescript
interface SpoolEvent {
  timestamp: string;  // ISO 8601
  level: 'INFO' | 'WARNING' | 'ERROR' | 'DEBUG' | 'CRITICAL';
  module: string;     // e.g., 'sonic.usb'
  message: string;
  tags: string[];      // e.g., ['usb', 'create', 'lifecycle']
  metadata?: Record<string, unknown>;
}
```

Feed endpoint: reads `~/.local/share/udos/feeds/spool/sonic-events.jsonl`.

---

## Step 5: Component Library Alignment

Prefer shared uCore components over rebuilding:

| Need | uCore Component | Status |
|---|---|---|
| File explorer | `FilepickerSidebar` | ✅ Built |
| Status badges | `.status-badge` pattern | Use USX tokens |
| Search/filter | `WorkspaceFilter` | ✅ Built |
| Tab navigation | `TabsModule` | ✅ Built |
| Card layout | `.usx-card` pattern | Use USX tokens |
| Data tables | `rich.table.Table` (CLI) → HTML `<table>` | Build new |

---

## Step 6: Theme Support

USX tokens support light/dark themes via `data-theme`:

```html
<html data-theme="dark">
```

```css
/* Auto-detect in JS */
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
document.documentElement.setAttribute(
  'data-theme',
  prefersDark ? 'dark' : 'light'
);
```

---

## Step 7: Safety Gate Pattern

For gated operations (USB create, device flash), use uCore's confirmation modal pattern:

```vue
<template>
  <button
    class="btn-gated"
    :disabled="!confirmed"
    @click="handleGatedAction"
  >
    {{ confirmed ? 'Execute' : 'Confirm First' }}
  </button>
</template>

<script setup lang="ts">
import { ref } from 'vue';
const confirmed = ref(false);

function handleGatedAction() {
  // Call sonic CLI or MCP tool
}
</script>
```

---

## Checklist for New Sonic Web Surfaces

- [ ] `@udos/usx-tokens` installed as dependency (never fork)
- [ ] CSS custom properties used for all colors, spacing, typography
- [ ] Status badges use USX semantic colors (`success`/`warn`/`error`)
- [ ] Cards use `--surface-container` background
- [ ] Tables align with Rich table conventions
- [ ] Gated operations require explicit confirmation
- [ ] Spool event viewer matches `SpoolEvent` schema
- [ ] Theme respects user `prefers-color-scheme`