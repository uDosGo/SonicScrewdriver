---
title: "USX Design Token System — npm Structure & Reference"
status: active
last_updated: 2026-07-08
category: architecture
tags: [usx, design-tokens, npm, uDos, reference]
---

# USX Design Token System — Reference for uDos Projects

## Overview

**USX v3.0.0** is the canonical **uDos design token system** published to npm as **`@udos/usx-tokens`**.

- **Source home:** `uCore/packages/usx-tokens/` (canonical)
- **Published:** npm registry (public, `@udos` scope)
- **Primary users:** uCore (Vue), HomeNest (Vue), Groovebox (React)
- **SonicScrewdriver role:** Reference documentation; used if/when CLI tools gain UI components

**Full documentation:** See [uCore docs/USX-NPM-STRUCTURE.md](../../../uCore/docs/USX-NPM-STRUCTURE.md)

---

## What is USX?

USX tokens are **shared CSS custom properties** that define a unified design language across all uDos UI surfaces:

- **Color system** — Primary, secondary, surfaces, borders, accessibility-focused
- **Spacing scale** — xs (2px) through 2xl (32px) for consistent rhythm
- **Typography** — Font sizes, weights, families for readable hierarchies
- **Touch targets** — 48px minimum, 56px comfortable (for remotes/gamepads)
- **Components** — Radii, card padding, tab borders, UI micro-measurements

### Example Token Usage

```css
.my-component {
  background: var(--surface-container);           /* Elevated surface */
  color: var(--usx-color-on-surface);             /* Text on surface */
  padding: var(--usx-spacing-md);                 /* 8px gutters */
  border-radius: var(--usx-radius-md);            /* Card corners */
  min-height: var(--usx-touch-min);               /* Touch target: 48px */
  font-size: var(--usx-font-size-body);           /* 16px body text */
}
```

---

## Package Structure

```
@udos/usx-tokens/
├── tokens/                 ← Canonical CSS variables
│   ├── tokens-color.css
│   ├── tokens-components.css
│   ├── tokens-spacing.css
│   ├── tokens-touch.css
│   └── tokens-typography.css
├── themes/                 ← Theme overrides (dark, light, c64, teletext, etc.)
├── home-nest/              ← HomeNest 10-foot console specifics
└── usx-standard.css        ← Base layout primitives
```

---

## Key Principles

1. **Single source of truth** — All projects import from `@udos/usx-tokens`, never fork.
2. **Variables only** — Tokens define CSS custom property values; projects define layout & behavior.
3. **Themes as value swaps** — A theme (dark, light, c64) only overrides variable values.
4. **Extensions, not forks** — New UI needs? Add to `home-nest/` or create project-specific exports.
5. **PicoCSS as base** — USX layers on top of PicoCSS for HTML5 semantic defaults.

---

## Token Reference

### Color Variables

```css
--usx-color-primary              /* Brand primary */
--usx-color-on-primary           /* Text on primary */
--usx-color-secondary            /* Secondary accent */
--surface                        /* Default surface background */
--surface-container              /* Elevated surface (cards, panels) */
--surface-container-high         /* Very elevated (modals) */
--usx-color-on-surface           /* Text on surface */
--border-subtle                  /* Subtle borders */
--border-strong                  /* Emphasis borders */
```

### Spacing Variables

```css
--usx-spacing-xs      2px     /* Tight gutters */
--usx-spacing-sm      4px     /* Small spacing */
--usx-spacing-md      8px     /* Default gutters */
--usx-spacing-lg      16px    /* Section margins */
--usx-spacing-xl      24px    /* Large sections */
--usx-spacing-2xl     32px    /* Major sections */
```

### Touch Target Variables

```css
--usx-touch-min                  /* 48px (accessibility minimum) */
--usx-touch-comfortable          /* 56px (comfortable for remote/gamepad) */
```

### Typography Variables

```css
--usx-font-size-caption          /* 12px */
--usx-font-size-body             /* 16px (default) */
--usx-font-size-headline         /* 24px */
--usx-font-family-sans           /* System sans-serif */
--usx-font-weight-regular        /* 400 */
--usx-font-weight-bold           /* 700 */
```

---

## NPM Installation

For projects that need UI components:

```bash
npm install @udos/usx-tokens
# or
npm install @udos/usx-tokens --save
```

Then import in your CSS or JS:

```css
@import '@udos/usx-tokens/usx-standard.css';
```

or 

```javascript
import '@udos/usx-tokens/usx-standard.css';
```

---

## Publishing & Version Management

### uCore Maintainers Only

When updating tokens in `uCore/packages/usx-tokens/`:

```bash
cd uCore/packages/usx-tokens
npm version patch    # or minor / major
npm publish --access public
```

**Version bumping:**
- `patch` (3.0.0 → 3.0.1) — Bug fixes, minor tweaks
- `minor` (3.0.0 → 3.1.0) — New tokens, new themes, new features
- `major` (3.0.0 → 4.0.0) — Breaking changes (renamed tokens, restructure)

### All Consumers (HomeNest, Groovebox, etc.)

After npm publish:

```bash
npm update @udos/usx-tokens
npm run test     # or verify in dev
git add package.json package-lock.json
git commit -m "chore: bump @udos/usx-tokens to 3.0.1"
```

---

## For SonicScrewdriver

**Current state:** SonicScrewdriver is a Python/Linux bootloader and CLI tool. No direct USX integration is needed.

**If future plans include:**
- CLI UI components (e.g., TUI dashboard)
- Web-based installer UI
- Monitoring dashboard

Then use the USX token system for consistency with uCore, HomeNest, and Groovebox.

**Integration path:**
1. Add web UI component (Vite + Vue/React)
2. Add `@udos/usx-tokens` to `package.json`
3. Import tokens in CSS
4. Reference this guide and [uCore docs/USX-NPM-STRUCTURE.md](../../../uCore/docs/USX-NPM-STRUCTURE.md)

---

## Related Documentation

- [uCore Architecture](../../../uCore/docs/architecture.md)
- [HomeNest Console Design](../../../HomeNest/docs/architecture.md)
- [Full USX npm Structure Guide](../../../uCore/docs/USX-NPM-STRUCTURE.md)
- [PicoCSS Reference](https://picocss.com/)
- [uCore GitHub Issues (label: usx-tokens)](https://github.com/fredporter/uCore/issues)
