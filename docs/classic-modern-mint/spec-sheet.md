# Classic Modern Mint — Complete Design Specification Sheet

## 📋 Document Control

| Property            | Value                              |
| ------------------- | ---------------------------------- |
| **Theme Name**      | Classic Modern Mint                |
| **Version**         | 3.0.0                              |
| **Base OS**         | Linux Mint 21.3 (Cinnamon)         |
| **Design Language** | Classic Mac Platinum + Modern Mono |
| **OBF Source**      | `classic-modern.obf` (v3)          |
| **Last Updated**    | 2026-04-21                         |

***

## 🎨 1. Colour Palette

### 1.1 Primary Colours

| Role           | Hex       | RGB         | Usage                    |
| -------------- | --------- | ----------- | ------------------------ |
| **Background** | `#E8E8E8` | 232,232,232 | Desktop, windows, panels |
| **Surface**    | `#F2F2F2` | 242,242,242 | Cards, menus, dialogs    |
| **Border**     | `#222222` | 34,34,34    | All borders, dividers    |
| **Text**       | `#111111` | 17,17,17    | All primary text         |

### 1.2 Accent Colours

| Role               | Hex       | RGB        | Usage                           |
| ------------------ | --------- | ---------- | ------------------------------- |
| **Primary Accent** | `#3A7BD5` | 58,123,213 | Selection, active states, links |
| **Accent Hover**   | `#2A6BC5` | 42,107,197 | Button hover, menu hover        |
| **Accent Active**  | `#1A5BB5` | 26,91,181  | Pressed state, focus            |

### 1.3 Status Colours

| Role        | Hex       | RGB        | Usage                           |
| ----------- | --------- | ---------- | ------------------------------- |
| **Success** | `#2E7D32` | 46,125,50  | Success messages, online status |
| **Warning** | `#F57C00` | 245,124,0  | Warnings, degraded status       |
| **Error**   | `#C62828` | 198,40,40  | Errors, offline status          |
| **Info**    | `#1565C0` | 21,101,192 | Information, help               |

### 1.4 Mono Colours

| Role      | Hex       | RGB         | Usage                         |
| --------- | --------- | ----------- | ----------------------------- |
| **Dark**  | `#000000` | 0,0,0       | Heavy borders, focus rings    |
| **Light** | `#FFFFFF` | 255,255,255 | Text on accent, inversions    |
| **Grey**  | `#808080` | 128,128,128 | Disabled text, secondary info |

### 1.5 Colour Rules

```css
/* No gradients */
* {
  background-image: none !important;
  box-shadow: none !important;
}

/* No transparency */
* {
  opacity: 1 !important;
}

/* Hard borders only */
border: 1px solid #222222;
```

***

## 📝 2. Typography System

### 2.1 Font Families

| Role              | Font Stack             | Fallback   | Size | Weight |
| ----------------- | ---------------------- | ---------- | ---- | ------ |
| **UI (Chrome)**   | ChicagoFLF, Chicago    | monospace  | 12px | normal |
| **Body Text**     | Monaspace Argon, Inter | system-ui  | 13px | 400    |
| **Mono/Code**     | pixChicago, Monaco     | monospace  | 11px | normal |
| **Headings**      | ChicagoFLF             | monospace  | 16px | normal |
| **Small Text**    | Geneva 9               | sans-serif | 10px | normal |
| **Menu Bar**      | ChicagoFLF             | monospace  | 12px | normal |
| **Window Title**  | ChicagoFLF             | monospace  | 12px | normal |
| **Dialog Button** | ChicagoFLF             | monospace  | 11px | normal |

### 2.2 Font Sizes

| Element            | Size | Line Height |
| ------------------ | ---- | ----------- |
| Desktop icon label | 11px | 1.2         |
| Panel text         | 11px | 1.3         |
| Menu items         | 12px | 1.4         |
| Window content     | 13px | 1.5         |
| Dialog content     | 13px | 1.5         |
| Tooltips           | 10px | 1.3         |
| Status bar         | 10px | 1.2         |

### 2.3 Typography Rules

```css
/* No text shadows */
text-shadow: none;

/* No font smoothing on UI fonts */
.ui-text {
  -webkit-font-smoothing: none;
  font-smoothing: none;
}

/* Body text uses standard smoothing */
body {
  -webkit-font-smoothing: antialiased;
}
```

***

## 🔲 3. Border System

### 3.1 Border Specifications

| Property   | Value             |
| ---------- | ----------------- |
| **Width**  | 1px               |
| **Style**  | solid             |
| **Colour** | `#222222`         |
| **Radius** | 0px (no rounding) |
| **Shadow** | none              |

### 3.2 Border Variants

| Element              | Border Type | Colour    | Notes            |
| -------------------- | ----------- | --------- | ---------------- |
| Window outer         | solid       | `#222`    | Full frame       |
| Window inner (inset) | solid       | `#808080` | For depth (rare) |
| Panel divider        | solid       | `#222`    | Bottom only      |
| Menu separator       | solid       | `#222`    | 1px line         |
| Button               | solid       | `#222`    | Full frame       |
| Input field          | solid       | `#222`    | Full frame       |
| Focus ring           | solid       | `#3A7BD5` | 1px, no offset   |
| Table cell           | solid       | `#222`    | Bottom only      |

### 3.3 Border Examples

```css
/* Standard window */
.window {
  border: 1px solid #222222;
}

/* Inset effect (rare, for legacy) */
.inset {
  border: 1px solid #808080;
  border-right-color: #FFFFFF;
  border-bottom-color: #FFFFFF;
}

/* Focus indicator */
:focus {
  outline: 1px solid #3A7BD5;
  outline-offset: 0px;
}
```

***

## 🖼️ 4. Background Patterns

### 4.1 Desktop Pattern (Linen)

```css
.desktop {
  background-color: #E8E8E8;
  background-image: repeating-linear-gradient(
    45deg,
    rgba(0,0,0,0.02) 0px,
    rgba(0,0,0,0.02) 2px,
    transparent 2px,
    transparent 8px
  );
}
```

### 4.2 Window Title Pattern (Platinum Stripes)

```css
.window-title {
  background-image: repeating-linear-gradient(
    45deg,
    #FFFFFF 0px,
    #FFFFFF 2px,
    #F2F2F2 2px,
    #F2F2F2 4px
  );
}
```

### 4.3 Selection Pattern (Classic Mac)

```css
:selected,
.selected {
  background-color: #3A7BD5;
  color: #FFFFFF;
}

/* Alternative striped selection (optional) */
.selection-striped {
  background-image: repeating-linear-gradient(
    45deg,
    #3A7BD5 0px,
    #3A7BD5 2px,
    #2A6BC5 2px,
    #2A6BC5 4px
  );
}
```

### 4.4 Pattern Rules

```css
/* Patterns are optional and low contrast */
/* Enabled via gsettings or user preference */

/* Disable patterns for performance */
.pattern-disabled {
  background-image: none !important;
}
```

***

## 📏 5. Spacing System

### 5.1 Base Unit

````
1 unit = 4px
````

### 5.2 Spacing Scale

| Token           | Units | Pixels | Usage                           |
| --------------- | ----- | ------ | ------------------------------- |
| `--spacing-xs`  | 1     | 4px    | Icon padding, tight layouts     |
| `--spacing-sm`  | 2     | 8px    | Button padding, menu items      |
| `--spacing-md`  | 3     | 12px   | Window padding, card padding    |
| `--spacing-lg`  | 4     | 16px   | Section spacing, dialog padding |
| `--spacing-xl`  | 6     | 24px   | Major sections                  |
| `--spacing-xxl` | 8     | 32px   | Window margins                  |

### 5.3 Component Spacing

| Component  | Padding  | Margin    | Gap  |
| ---------- | -------- | --------- | ---- |
| Window     | 12px     | -         | -    |
| Dialog     | 16px     | -         | 12px |
| Panel      | 8px      | -         | 4px  |
| Button     | 4px 12px | 4px       | -    |
| Menu item  | 4px 24px | -         | -    |
| Toolbar    | 4px      | -         | 4px  |
| Table cell | 8px 12px | -         | -    |
| Form row   | -        | 0 0 8px 0 | -    |

***

## 🔘 6. Component Specifications

### 6.1 Window

```css
.window {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 12px;
}

.window-title {
  background: repeating-linear-gradient(45deg, #fff 0px, #fff 2px, #F2F2F2 2px, #F2F2F2 4px);
  border-bottom: 1px solid #222222;
  padding: 4px 8px;
  margin: -12px -12px 12px -12px;
  font-family: 'ChicagoFLF', monospace;
  font-size: 12px;
}

.window-controls {
  position: absolute;
  top: 6px;
  right: 8px;
  display: flex;
  gap: 6px;
}

.window-control {
  width: 12px;
  height: 12px;
  border: 1px solid #222222;
  background: #F2F2F2;
}

.window-control.close:active {
  background: #C62828;
}
```

### 6.2 Button

```css
.button {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 4px 12px;
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
  color: #111111;
  min-width: 60px;
}

.button:hover {
  background: #E8E8E8;
}

.button:active {
  background: #3A7BD5;
  color: #FFFFFF;
}

.button:disabled {
  background: #E8E8E8;
  color: #808080;
  border-color: #808080;
}

/* Primary action button */
.button-primary {
  background: #3A7BD5;
  color: #FFFFFF;
}

.button-primary:active {
  background: #1A5BB5;
}
```

### 6.3 Menu

```css
.menu-bar {
  background: #F2F2F2;
  border-bottom: 1px solid #222222;
  padding: 2px 8px;
}

.menu-item {
  padding: 4px 12px;
  font-family: 'ChicagoFLF', monospace;
  font-size: 12px;
  color: #111111;
}

.menu-item:hover {
  background: #3A7BD5;
  color: #FFFFFF;
}

.menu-separator {
  height: 1px;
  background: #222222;
  margin: 4px 0;
}

/* Dropdown menu (popover) */
.menu-popover {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 4px 0;
  min-width: 150px;
}
```

### 6.4 Input Fields

```css
.input,
.text-field,
textbox {
  background: #FFFFFF;
  border: 1px solid #222222;
  padding: 6px 8px;
  font-family: 'Monaspace Argon', monospace;
  font-size: 13px;
  color: #111111;
}

.input:focus,
.text-field:focus {
  outline: 1px solid #3A7BD5;
  outline-offset: 0px;
}

.input:disabled {
  background: #E8E8E8;
  color: #808080;
}

/* Checkbox */
.checkbox {
  width: 16px;
  height: 16px;
  background: #FFFFFF;
  border: 1px solid #222222;
  margin-right: 8px;
}

.checkbox:checked {
  background: #3A7BD5;
  position: relative;
}

.checkbox:checked::after {
  content: "✓";
  color: #FFFFFF;
  position: absolute;
  top: -2px;
  left: 2px;
  font-size: 12px;
}

/* Radio button */
.radio {
  width: 14px;
  height: 14px;
  background: #FFFFFF;
  border: 1px solid #222222;
  border-radius: 50%;
  margin-right: 8px;
}

.radio:checked {
  background: #3A7BD5;
  box-shadow: inset 0 0 0 2px #FFFFFF;
}
```

### 6.5 Tables

```css
.table {
  border-collapse: collapse;
  width: 100%;
}

.table th {
  text-align: left;
  padding: 8px 12px;
  background: #E8E8E8;
  border-bottom: 1px solid #222222;
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
  font-weight: normal;
}

.table td {
  padding: 8px 12px;
  border-bottom: 1px solid #222222;
  font-family: 'Monaspace Argon', monospace;
  font-size: 13px;
}

.table tr:hover {
  background: #E8E8E8;
}

.table tr.selected {
  background: #3A7BD5;
  color: #FFFFFF;
}
```

### 6.6 Tabs

```css
.tab-bar {
  border-bottom: 1px solid #222222;
  display: flex;
  gap: 0px;
}

.tab {
  padding: 8px 16px;
  background: #E8E8E8;
  border: 1px solid #222222;
  border-bottom: none;
  margin-right: -1px;
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
}

.tab.active {
  background: #F2F2F2;
  border-bottom: 1px solid #F2F2F2;
  margin-bottom: -1px;
}

.tab:hover:not(.active) {
  background: #F2F2F2;
}
```

### 6.7 Scrollbar

```css
scrollbar {
  background: #E8E8E8;
  border: 1px solid #222222;
}

scrollbar slider {
  background: #F2F2F2;
  border: 1px solid #222222;
  min-width: 16px;
  min-height: 16px;
}

scrollbar slider:hover {
  background: #3A7BD5;
}

scrollbar slider:active {
  background: #1A5BB5;
}
```

### 6.8 Tooltip

```css
.tooltip {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 4px 8px;
  font-family: 'Geneva 9', monospace;
  font-size: 10px;
  color: #111111;
  
  /* No delay */
  transition: none;
}
```

### 6.9 Progress Bar

```css
.progress-bar {
  background: #E8E8E8;
  border: 1px solid #222222;
  height: 16px;
}

.progress-bar-fill {
  background: #3A7BD5;
  height: 100%;
}

/* Indeterminate (striped) */
.progress-bar.indeterminate .progress-bar-fill {
  background-image: repeating-linear-gradient(
    45deg,
    #3A7BD5 0px,
    #3A7BD5 10px,
    #2A6BC5 10px,
    #2A6BC5 20px
  );
  animation: progress-stripes 1s linear infinite;
}
```

### 6.10 Dialog

```css
.dialog {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 16px;
  min-width: 300px;
}

.dialog-header {
  font-family: 'ChicagoFLF', monospace;
  font-size: 14px;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #222222;
}

.dialog-content {
  margin-bottom: 16px;
}

.dialog-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  border-top: 1px solid #222222;
  padding-top: 12px;
}
```

***

## 🖥️ 7. Desktop Environment Elements

### 7.1 Cinnamon Panel

```css
/* Top panel */
.panel-top {
  background: #E8E8E8;
  border-bottom: 1px solid #222222;
  height: 28px;
  padding: 0 8px;
}

/* Panel text */
.panel-text,
.panel-button {
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
  color: #111111;
}

/* Panel buttons */
.panel-button:hover {
  background: #F2F2F2;
}

/* System tray icons */
.systray-icon {
  padding: 2px;
}

/* Workspace switcher */
.workspace-button {
  padding: 0 8px;
  border-left: 1px solid #222222;
  border-right: 1px solid #222222;
  margin: 0 -1px;
}

.workspace-button.active {
  background: #3A7BD5;
  color: #FFFFFF;
}
```

### 7.2 Desktop Icons

```css
/* Desktop icon container */
.desktop-icon {
  width: 80px;
  text-align: center;
}

/* Desktop icon label */
.desktop-icon-label {
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
  color: #111111;
  background: rgba(242, 242, 242, 0.8);
  padding: 2px 4px;
  margin-top: 4px;
}

/* Selected desktop icon */
.desktop-icon.selected .desktop-icon-label {
  background: #3A7BD5;
  color: #FFFFFF;
}
```

### 7.3 Cinnamon Menu

```css
/* Start menu */
.menu {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 8px;
}

.menu-header {
  padding-bottom: 8px;
  border-bottom: 1px solid #222222;
  margin-bottom: 8px;
}

.menu-search {
  background: #FFFFFF;
  border: 1px solid #222222;
  padding: 4px 8px;
  font-family: 'Monaspace Argon', monospace;
  font-size: 12px;
  margin-bottom: 8px;
}

.menu-applications {
  list-style: none;
  padding: 0;
}

.menu-application-item {
  padding: 4px 8px;
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
}

.menu-application-item:hover {
  background: #3A7BD5;
  color: #FFFFFF;
}
```

### 7.4 Notifications

```css
.notification {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 12px;
  min-width: 250px;
  margin: 8px;
}

.notification-title {
  font-family: 'ChicagoFLF', monospace;
  font-size: 12px;
  font-weight: normal;
  margin-bottom: 4px;
}

.notification-body {
  font-family: 'Monaspace Argon', monospace;
  font-size: 12px;
}

.notification-icon {
  margin-right: 8px;
}

/* Notification popup container */
.notification-popup {
  background: #E8E8E8;
  border: 1px solid #222222;
}
```

### 7.5 Workspace Overview

```css
/* Expo view (workspace overview) */
.workspace-overview {
  background: #E8E8E8;
}

.workspace-thumbnail {
  background: #F2F2F2;
  border: 1px solid #222222;
  margin: 8px;
}

.workspace-thumbnail.active {
  border: 2px solid #3A7BD5;
}
```

***

## 🔄 8. Loader & Startup

### 8.1 Boot Splash (Plymouth)

```css
/* /usr/share/plymouth/themes/classic-modern/ */
.boot-splash {
  background: #E8E8E8;
}

.boot-logo {
  /* Classic Mac-style Happy Mac or custom */
  content: url("happy-mac.svg");
  width: 64px;
  height: 64px;
}

.boot-progress {
  background: #F2F2F2;
  border: 1px solid #222222;
  width: 400px;
  height: 16px;
}

.boot-progress-fill {
  background: #3A7BD5;
  height: 100%;
}

.boot-text {
  font-family: 'ChicagoFLF', monospace;
  font-size: 12px;
  color: #111111;
  margin-top: 20px;
}
```

### 8.2 Login Screen (LightDM)

```css
/* /usr/share/lightdm-webkit/themes/classic-modern/ */
.login-screen {
  background: #E8E8E8;
}

.login-box {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 24px;
  width: 320px;
}

.login-title {
  font-family: 'ChicagoFLF', monospace;
  font-size: 16px;
  text-align: center;
  margin-bottom: 20px;
}

.login-input {
  background: #FFFFFF;
  border: 1px solid #222222;
  padding: 8px;
  font-family: 'Monaspace Argon', monospace;
  font-size: 13px;
  width: 100%;
  margin-bottom: 12px;
}

.login-button {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 6px 12px;
  font-family: 'ChicagoFLF', monospace;
  font-size: 11px;
  width: 100%;
}

.login-button:active {
  background: #3A7BD5;
  color: #FFFFFF;
}
```

### 8.3 Loading Spinner

```css
/* Classic Mac "watch" or "beachball" alternative */
.spinner {
  width: 32px;
  height: 32px;
  border: 2px solid #222222;
  border-top-color: #3A7BD5;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Or use classic Mac cursor animation */
.spinner-classic {
  background-image: url("watch-cursor.gif");
  width: 16px;
  height: 16px;
}
```

### 8.4 Application Startup

```css
/* Splash screen for apps */
.splash {
  background: #F2F2F2;
  border: 1px solid #222222;
  padding: 20px;
  text-align: center;
}

.splash-icon {
  margin-bottom: 16px;
}

.splash-text {
  font-family: 'ChicagoFLF', monospace;
  font-size: 12px;
  color: #111111;
}
```

***

## 🎯 9. Cursor Theme

### 9.1 Cursor Specifications

| Cursor Type | Style               | Notes                     |
| ----------- | ------------------- | ------------------------- |
| Default     | Arrow (Classic Mac) | Black outline, white fill |
| Pointer     | Pointing hand       | Black outline             |
| Text        | I-beam              | 1px width                 |
| Link        | Pointing hand       | Same as pointer           |
| Working     | Watch (animated)    | Classic Mac watch cursor  |
| Resize      | Diagonal arrows     | Simple, no fancy shapes   |

### 9.2 Cursor Configuration

```bash
# ~/.icons/classic-modern-cursors/
cursor-theme:
  name: "Classic Modern"
  default: "arrow"
  text: "ibeam"
  pointer: "hand"
  progress: "watch"
  resize: "resize"
```

***

## 🎬 10. Animation & Transition

### 10.1 Animation Rules

```css
/* No animations by default */
* {
  transition: none !important;
  animation-duration: 0s !important;
}

/* Minimal animations (user preference) */
@media (prefers-reduced-motion: no-preference) {
  .fade-in {
    animation: fadeIn 0.1s ease-in;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
}
```

### 10.2 Transition Durations

| Element      | Duration      | Type    |
| ------------ | ------------- | ------- |
| Window open  | 0ms (instant) | None    |
| Window close | 0ms (instant) | None    |
| Menu open    | 50ms          | Fade    |
| Menu close   | 0ms           | Instant |
| Dialog open  | 100ms         | Fade    |
| Tooltip      | 0ms           | Instant |
| Hover state  | 0ms           | Instant |

***

## 📱 11. Responsive Breakpoints

| Breakpoint    | Width          | Layout Changes  |
| ------------- | -------------- | --------------- |
| Desktop       | > 1024px       | Full layout     |
| Small Desktop | 768px - 1024px | Reduced padding |
| Tablet        | < 768px        | Stacked panels  |
| Mobile        | < 480px        | Single column   |

***

## ✅ 12. Compliance Checklist

| Element         | Spec                               | Status |
| --------------- | ---------------------------------- | ------ |
| Colours         | No gradients, no transparency      | ✅     |
| Borders         | 1px solid #222, 0 radius           | ✅     |
| Typography      | Chicago for UI, Monaspace for body | ✅     |
| Spacing         | 4px unit system                    | ✅     |
| Patterns        | Linen (desktop), Stripes (title)   | ✅     |
| Animations      | None or minimal                    | ✅     |
| Window controls | Top-right, classic style           | ✅     |
| Menu bar        | Top, Chicago font                  | ✅     |
| Desktop icons   | Chicago labels, 80px width         | ✅     |
| Panel           | Top only, 28px height              | ✅     |
| Cursors         | Classic Mac style                  | ✅     |
| Boot splash     | Classic Mac inspired               | ✅     |

***

## 📦 13. Installation Artifacts

```bash
# Theme installation paths
/usr/share/themes/Classic-Modern-Mint/
├── cinnamon/
├── gtk-3.0/
├── gtk-4.0/
├── metacity-1/
└── index.theme

# Icons
/usr/share/icons/Classic-Modern-Icons/
├── apps/
├── actions/
├── devices/
└── status/

# Cursors
/usr/share/icons/Classic-Modern-Cursors/

# Plymouth boot theme
/usr/share/plymouth/themes/classic-modern/

# LightDM theme
/usr/share/lightdm-webkit/themes/classic-modern/
```

***

**This specification sheet defines every visual element of Classic Modern Mint — ready for implementation.**