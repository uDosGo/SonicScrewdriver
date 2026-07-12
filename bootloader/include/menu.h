#ifndef MENU_H
#define MENU_H

#include <stdint.h>
#include "teletext.h"

/* Maximum menu entries */
#define MENU_MAX_ENTRIES    20
#define MENU_MAX_TITLE      64
#define MENU_MAX_LABEL      32
#define MENU_MAX_HELP_TEXT  128

/* Boot entry types */
typedef enum {
    BOOT_TYPE_CHAINLOAD,       /* Chainload another bootloader */
    BOOT_TYPE_LINUX_KERNEL,    /* Direct Linux kernel boot */
    BOOT_TYPE_EFI_STUB,        /* EFI application execution */
    BOOT_TYPE_REBOOT,          /* Reboot system */
    BOOT_TYPE_SHUTDOWN,        /* Shutdown system */
    BOOT_TYPE_FIRMWARE_SETUP,  /* Enter firmware setup */
    BOOT_TYPE_MEMTEST,         /* Memory test */
    BOOT_TYPE_SUBMENU          /* Nested submenu */
} boot_entry_type_t;

/* Boot entry */
typedef struct {
    char label[MENU_MAX_LABEL];         /* Display label */
    boot_entry_type_t type;              /* Entry type */
    char path[256];                      /* Path to boot file or kernel */
    char initrd[256];                    /* Initrd path (for Linux kernel) */
    char cmdline[512];                   /* Kernel command line */
    uint8_t default_entry : 1;           /* Is this the default? */
    uint8_t requires_mac   : 1;          /* Mac only */
    uint8_t requires_pc    : 1;          /* PC only */
    uint8_t requires_bios  : 1;          /* Legacy BIOS only */
    uint8_t requires_uefi  : 1;          /* UEFI only */
} menu_entry_t;

/* Platform detection result */
typedef enum {
    PLATFORM_UNKNOWN = 0,
    PLATFORM_MAC_UEFI,       /* Mac with UEFI (Intel + Apple Silicon) */
    PLATFORM_PC_UEFI,        /* PC with UEFI */
    PLATFORM_PC_BIOS,        /* PC with legacy BIOS */
    PLATFORM_VIRTUAL,        /* Virtual machine */
    PLATFORM_ARM64           /* ARM64 (Apple Silicon, Raspberry Pi) */
} platform_type_t;

/* Menu configuration */
typedef struct {
    char title[MENU_MAX_TITLE];          /* Menu title */
    char help_text[MENU_MAX_HELP_TEXT];  /* Help footer text */
    uint8_t entry_count;                 /* Number of entries */
    menu_entry_t entries[MENU_MAX_ENTRIES];
    uint8_t default_index;               /* Default selected entry */
    uint8_t timeout_seconds;             /* Auto-boot timeout (0 = no timeout) */
    uint8_t fg_color;                    /* Foreground color */
    uint8_t bg_color;                    /* Background color */
    uint8_t highlight_fg;               /* Highlight foreground */
    uint8_t highlight_bg;               /* Highlight background */
    uint8_t title_fg;                   /* Title foreground */
    uint8_t title_bg;                   /* Title background */
} menu_config_t;

/* Menu state */
typedef struct {
    menu_config_t config;
    uint8_t selected_index;
    uint8_t scroll_offset;
    uint8_t visible_entries;
    platform_type_t platform;
    int8_t timeout_remaining;            /* -1 = no timeout */
} menu_state_t;

/* Function prototypes */

/* Initialize a menu from config */
void menu_init(menu_state_t *state, const menu_config_t *config);

/* Detect the current platform */
platform_type_t menu_detect_platform(void);

/* Filter entries based on platform detection */
uint8_t menu_get_visible_count(menu_state_t *state);

/* Get the nth visible entry (accounting for platform filtering) */
menu_entry_t *menu_get_visible_entry(menu_state_t *state, uint8_t index);

/* Navigate the menu */
void menu_move_up(menu_state_t *state);
void menu_move_down(menu_state_t *state);
void menu_move_first(menu_state_t *state);
void menu_move_last(menu_state_t *state);

/* Select the current entry (returns the entry to boot) */
menu_entry_t *menu_select_current(menu_state_t *state);

/* Render the menu to a teletext screen */
void menu_render(teletext_screen_t *screen, menu_state_t *state);

/* Handle a keypress, returns 1 if menu should update, 0 if boot selected */
int menu_handle_key(menu_state_t *state, uint8_t key);

/* Tick the timeout counter, returns 1 if timeout expired (auto-boot) */
int menu_tick_timeout(menu_state_t *state);

/* Get the help text for function keys */
const char *menu_get_help_text(menu_state_t *state);

#endif /* MENU_H */
