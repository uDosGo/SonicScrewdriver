/*
 * SonicScrewloader — Universal USB Bootloader
 *
 * Entry point for the SonicScrewdriver v2 bootloader.
 * Compiles for both UEFI (x86_64) and legacy BIOS targets.
 *
 * Build:
 *   make ARCH=x86_64          # UEFI build
 *   make ARCH=x86_64 BIOS=1   # Legacy BIOS build
 */

#include "teletext.h"
#include "menu.h"
#include "detect.h"
#include <string.h>

/* Global teletext screen */
static teletext_screen_t screen;

/* Menu state */
static menu_state_t menu_state;

/* Forward declarations */
extern void framebuffer_render(teletext_cell_t *cells, size_t count, void *context);
extern void framebuffer_clear(void);
extern int chainload_boot_entry(menu_entry_t *entry);

/* Default menu configuration — will be overridden by YAML config at build time */
static menu_config_t default_config = {
    .title = "SonicScrewdriver v2.0.0",
    .help_text = "Select an option to boot. Press Enter to confirm.",
    .entry_count = 7,
    .default_index = 1,
    .timeout_seconds = 10,
    .fg_color = TELETEXT_WHITE,
    .bg_color = TELETEXT_BLACK,
    .highlight_fg = TELETEXT_BLACK,
    .highlight_bg = TELETEXT_WHITE,
    .title_fg = TELETEXT_YELLOW,
    .title_bg = TELETEXT_BLUE,
    .entries = {
        {
            .label = "macOS Recovery",
            .type = BOOT_TYPE_CHAINLOAD,
            .path = "\\System\\Library\\CoreServices\\boot.efi",
            .requires_mac = 1,
            .requires_uefi = 1
        },
        {
            .label = "Linux Mint 22 (uDos Classic Modern)",
            .type = BOOT_TYPE_CHAINLOAD,
            .path = "\\EFI\\grub\\grubx64.efi",
            .default_entry = 1
        },
        {
            .label = "SonicScrewloader (CHASIS Games)",
            .type = BOOT_TYPE_SUBMENU,
            .path = ""
        },
        {
            .label = "Recovery Tools",
            .type = BOOT_TYPE_SUBMENU,
            .path = ""
        },
        {
            .label = "Memory Test (MemTest86)",
            .type = BOOT_TYPE_MEMTEST,
            .path = "\\EFI\\memtest86.efi"
        },
        {
            .label = "Firmware Setup",
            .type = BOOT_TYPE_FIRMWARE_SETUP,
            .path = ""
        },
        {
            .label = "Reboot System",
            .type = BOOT_TYPE_REBOOT,
            .path = ""
        }
    }
};

/* Recovery tools submenu */
static menu_config_t recovery_config = {
    .title = "Sonic Recovery Tools",
    .help_text = "Select a recovery tool to boot.",
    .entry_count = 8,
    .default_index = 0,
    .timeout_seconds = 0,
    .fg_color = TELETEXT_WHITE,
    .bg_color = TELETEXT_BLACK,
    .highlight_fg = TELETEXT_BLACK,
    .highlight_bg = TELETEXT_CYAN,
    .title_fg = TELETEXT_CYAN,
    .title_bg = TELETEXT_RED,
    .entries = {
        { .label = "GParted — Partition Manager",       .type = BOOT_TYPE_EFI_STUB, .path = "\\EFI\\tools\\gparted.efi" },
        { .label = "TestDisk — Partition Recovery",     .type = BOOT_TYPE_EFI_STUB, .path = "\\EFI\\tools\\testdisk.efi" },
        { .label = "Clonezilla — Disk Imaging",         .type = BOOT_TYPE_EFI_STUB, .path = "\\EFI\\tools\\clonezilla.efi" },
        { .label = "MemTest86 — RAM Diagnostics",       .type = BOOT_TYPE_MEMTEST,  .path = "\\EFI\\memtest86.efi" },
        { .label = "Super GRUB Disk — Boot Repair",     .type = BOOT_TYPE_EFI_STUB, .path = "\\EFI\\tools\\supergrub.efi" },
        { .label = "SystemRescue — Full Rescue Env",    .type = BOOT_TYPE_EFI_STUB, .path = "\\EFI\\tools\\systemrescue.efi" },
        { .label = "uDos Rescue — Repair Sonic Install",.type = BOOT_TYPE_EFI_STUB, .path = "\\EFI\\tools\\udos-rescue.efi" },
        { .label = "CHASIS Repair — Fix Game Containers",.type = BOOT_TYPE_EFI_STUB,.path = "\\EFI\\tools\\chasis-repair.efi" },
    }
};

/* CHASIS game submenu */
static menu_config_t chasis_config = {
    .title = "SonicScrewloader — CHASIS Games",
    .help_text = "Select a game to launch. Games are stored on the exFAT data partition.",
    .entry_count = 5,
    .default_index = 0,
    .timeout_seconds = 0,
    .fg_color = TELETEXT_WHITE,
    .bg_color = TELETEXT_BLACK,
    .highlight_fg = TELETEXT_BLACK,
    .highlight_bg = TELETEXT_GREEN,
    .title_fg = TELETEXT_GREEN,
    .title_bg = TELETEXT_BLACK,
    .entries = {
        { .label = "Grid Runner — uCode1 Algebra Game", .type = BOOT_TYPE_EFI_STUB, .path = "\\CHASIS\\grid-runner.efi" },
        { .label = "Teletext Quest — Text Adventure",   .type = BOOT_TYPE_EFI_STUB, .path = "\\CHASIS\\teletext-quest.efi" },
        { .label = "Block Breaker — Retro Arcade",      .type = BOOT_TYPE_EFI_STUB, .path = "\\CHASIS\\block-breaker.efi" },
        { .label = "Mesh Arena — Multiplayer",          .type = BOOT_TYPE_EFI_STUB, .path = "\\CHASIS\\mesh-arena.efi" },
        { .label = "Back to Main Menu",                 .type = BOOT_TYPE_SUBMENU,  .path = "" },
    }
};

/* Boot the selected entry, handling submenus recursively */
static int boot_selected(menu_state_t *state) {
    menu_entry_t *entry = menu_select_current(state);
    if (!entry) return -1;

    /* Handle submenus */
    if (entry->type == BOOT_TYPE_SUBMENU) {
        menu_state_t sub_state;

        /* Determine which submenu to show */
        if (strstr(entry->label, "CHASIS") != NULL) {
            menu_init(&sub_state, &chasis_config);
        } else if (strstr(entry->label, "Recovery") != NULL) {
            menu_init(&sub_state, &recovery_config);
        } else {
            return 0; /* Back to main menu */
        }

        /* Run the submenu */
        return run_menu_loop(&sub_state);
    }

    /* Boot the entry */
    return chainload_boot_entry(entry);
}

/* Main menu loop */
int run_menu_loop(menu_state_t *state) {
    int running = 1;
    int result = 0;

    while (running) {
        /* Render the menu */
        menu_render(&screen, state);
        teletext_render(&screen, framebuffer_render, NULL);

        /* Handle timeout */
        if (menu_tick_timeout(state)) {
            /* Auto-boot default */
            state->selected_index = state->config.default_index;
            result = boot_selected(state);
            running = 0;
            break;
        }

        /* Wait for keypress and handle it */
        uint8_t key = 0; /* Would read from UEFI/BIOS input */
        /* 
         * In a real bootloader, this would:
         * UEFI: Call ST->ConIn->ReadKeyStroke()
         * BIOS: Read from keyboard port (0x60) or INT 16h
         */

        if (key == 0) {
            /* No key pressed — simple delay loop */
            for (volatile int i = 0; i < 1000000; i++);
            continue;
        }

        if (!menu_handle_key(state, key)) {
            /* User selected an entry */
            result = boot_selected(state);
            running = 0;
        }
    }

    return result;
}

/*
 * Main entry point.
 * Called by UEFI firmware (efi_main) or BIOS boot sector.
 */
int main(void) {
    /* Initialize teletext screen */
    teletext_init(&screen);
    teletext_set_defaults(&screen, TELETEXT_WHITE, TELETEXT_BLACK);

    /* Clear the framebuffer */
    framebuffer_clear();

    /* Show splash screen */
    teletext_clear(&screen);
    teletext_puts(&screen, 10, 25,
                  "SonicScrewdriver v2.0.0 — Loading...",
                  TELETEXT_YELLOW, TELETEXT_BLACK);
    teletext_render(&screen, framebuffer_render, NULL);

    /* Detect platform */
    platform_type_t platform = detect_platform();

    /* Initialize menu with default config */
    menu_init(&menu_state, &default_config);

    /* Override platform detection result */
    menu_state.platform = platform;
    menu_state.visible_entries = menu_get_visible_count(&menu_state);

    /* Run the main menu loop */
    int result = run_menu_loop(&menu_state);

    /* If we get here, something went wrong */
    teletext_clear(&screen);
    teletext_puts(&screen, 12, 20,
                  "Boot failed! Please check your USB drive.",
                  TELETEXT_RED, TELETEXT_BLACK);
    teletext_render(&screen, framebuffer_render, NULL);

    return result;
}

/* UEFI entry point */
#ifdef UEFI_BUILD
void efi_main(void *image_handle, void *system_table) {
    (void)image_handle;
    (void)system_table;
    /* Initialize UEFI-specific features */
    /* framebuffer_init_uefi(gop_buffer, hres, vres, scanline, bpp); */
    /* chainload_init(image_handle, system_table); */
    main();
}
#endif
