# SonicScrewloader — Bootloader Status Taxonomy

**Created:** 2026-07-09
**Status:** Foundation deliverable — `task.sonic.bootloader-ux.001`
**Purpose:** Map every bootloader lifecycle state to uCore-compatible status codes
and severity levels for diagnostics and spool event emission.

---

## Lifecycle States

| # | State | uCore Status | Severity | Description |
|---|-------|-------------|----------|-------------|
| 1 | `init` | `success` | `INFO` | Bootloader entry point invoked. Stack and BSS initialised. |
| 2 | `detect` | `success` | `INFO` | Hardware detection running (SMBIOS, ACPI, device-tree). |
| 3 | `detect_complete` | `success` | `INFO` | Platform identified: `mac`, `pc-uefi`, `pc-bios`, `arm64`. |
| 4 | `framebuffer_init` | `success` | `INFO` | Framebuffer (UEFI GOP) or VGA text mode (BIOS) initialised. |
| 5 | `fb_fallback` | `warn` | `WARNING` | GOP unavailable, falling back to serial/text-only output. |
| 6 | `config_load` | `success` | `INFO` | YAML menu config parsed and loaded into C structs. |
| 7 | `config_empty` | `warn` | `WARNING` | No menu config found — bootloader will show minimal prompt. |
| 8 | `theme_load` | `success` | `INFO` | Teletext theme loaded (16-color palette, block glyphs). |
| 9 | `theme_default` | `warn` | `INFO` | Theme missing, using built-in fallback (monochrome). |
| 10 | `menu_render` | `success` | `INFO` | Menu grid rendered to framebuffer — user can interact. |
| 11 | `menu_timeout` | `warn` | `INFO` | Menu timeout reached, auto-booting default entry. |
| 12 | `chainload_begin` | `success` | `INFO` | Chainloading target OS bootloader (GRUB, rEFInd, EFI stub). |
| 13 | `chainload_fail` | `error` | `ERROR` | Failed to locate or launch chainload target. |
| 14 | `chainload_success` | `success` | `INFO` | Control transferred to target bootloader. |
| 15 | `panic_oob` | `error` | `CRITICAL` | Unrecoverable error — out of bounds memory, triple fault, etc. |
| 16 | `panic_corrupt` | `error` | `CRITICAL` | Corrupt config, missing binary, or invalid EFI image. |
| 17 | `halt` | `error` | `ERROR` | Bootloader intentionally halted (user requested, diagnostic). |

---

## State Transitions

```
init
  └─► detect
         └─► detect_complete
                ├─► framebuffer_init ─► (fb_fallback if GOP missing)
                │
                ├─► config_load ─► (config_empty if no YAML found)
                │
                ├─► theme_load ─► (theme_default if theme missing)
                │
                └─► menu_render
                       ├─► menu_timeout ─► chainload_begin
                       │                     ├─► chainload_success
                       │                     └─► chainload_fail ─► menu_render
                       │
                       └─► (user selection) ─► chainload_begin
                      
  ANY STATE ─► panic_oob | panic_corrupt | halt
```

---

## Spool Event Shape

Each lifecycle state should emit a spool event:

```python
{
    "timestamp": "2026-07-09T22:30:00+08:00",
    "module": "sonic.bootloader",
    "level": "INFO",   # ← from severity column above
    "message": "Bootloader: detect_complete — platform=pc-uefi",
    "tags": ["bootloader", "detect", "uefi"]
}
```

---

## Architecture Mapping

| SonicScrewloader Component | Source File | Lifecycle States Covered |
|---|---|---|
| Entry / init | `bootloader/src/main.c` | `init` |
| Hardware detection | `bootloader/src/detect.c` | `detect`, `detect_complete` |
| Framebuffer setup | `bootloader/src/framebuffer.c` | `framebuffer_init`, `fb_fallback` |
| Menu engine | `bootloader/src/menu.c` | `config_load`, `config_empty`, `theme_load`, `theme_default`, `menu_render`, `menu_timeout` |
| Chainloading | `bootloader/src/chainload.c` | `chainload_begin`, `chainload_fail`, `chainload_success` |
| Teletext renderer | `bootloader/src/teletext.c` | `panic_oob` (renderer watchdog), `halt` |

---

## Alignment Notes

- **uCore compatibility:** Status codes map to `EnvelopeStatus` enum (`success`/`warn`/`error`).
- **Spool severity:** Maps to `EventLevel` (`INFO`/`WARNING`/`ERROR`/`CRITICAL`).
- **Module:** All bootloader events use `module: "sonic.bootloader"`.
- **Tags:** Lowercase kebab-case — `bootloader`, `detect`, `uefi`, `bios`, `arm64`, `chainload`, `menu`, `framebuffer`, `panic`.