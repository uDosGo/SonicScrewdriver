#include "menu.h"
#include "detect.h"
#include <string.h>

void menu_init(menu_state_t *state, const menu_config_t *config) {
    memcpy(&state->config, config, sizeof(menu_config_t));
    state->selected_index = config->default_index;
    state->scroll_offset = 0;
    state->platform = detect_platform();
    state->visible_entries = menu_get_visible_count(state);
    state->timeout_remaining = (config->timeout_seconds > 0)
        ? (int8_t)config->timeout_seconds
        : -1;
}

platform_type_t menu_detect_platform(void) {
    return detect_platform();
}

uint8_t menu_get_visible_count(menu_state_t *state) {
    uint8_t count = 0;
    for (uint8_t i = 0; i < state->config.entry_count; i++) {
        menu_entry_t *entry = &state->config.entries[i];
        /* Skip entries that don't match current platform */
        if (entry->requires_mac && state->platform != PLATFORM_MAC_UEFI) continue;
        if (entry->requires_pc && state->platform != PLATFORM_PC_UEFI
            && state->platform != PLATFORM_PC_BIOS) continue;
        if (entry->requires_bios && state->platform != PLATFORM_PC_BIOS) continue;
        if (entry->requires_uefi && state->platform != PLATFORM_MAC_UEFI
            && state->platform != PLATFORM_PC_UEFI) continue;
        count++;
    }
    return count;
}

menu_entry_t *menu_get_visible_entry(menu_state_t *state, uint8_t index) {
    uint8_t vis_idx = 0;
    for (uint8_t i = 0; i < state->config.entry_count; i++) {
        menu_entry_t *entry = &state->config.entries[i];
        if (entry->requires_mac && state->platform != PLATFORM_MAC_UEFI) continue;
        if (entry->requires_pc && state->platform != PLATFORM_PC_UEFI
            && state->platform != PLATFORM_PC_BIOS) continue;
        if (entry->requires_bios && state->platform != PLATFORM_PC_BIOS) continue;
        if (entry->requires_uefi && state->platform != PLATFORM_MAC_UEFI
            && state->platform != PLATFORM_PC_UEFI) continue;
        if (vis_idx == index) return entry;
        vis_idx++;
    }
    return NULL;
}

void menu_move_up(menu_state_t *state) {
    if (state->selected_index > 0) {
        state->selected_index--;
    }
    /* Scroll if needed */
    if (state->selected_index < state->scroll_offset) {
        state->scroll_offset = state->selected_index;
    }
    state->timeout_remaining = -1; /* Cancel timeout on user input */
}

void menu_move_down(menu_state_t *state) {
    if (state->selected_index < state->visible_entries - 1) {
        state->selected_index++;
    }
    /* Scroll if needed */
    uint8_t max_visible = TELETEXT_ROWS - 6; /* Leave room for header/footer */
    if (state->selected_index >= state->scroll_offset + max_visible) {
        state->scroll_offset = state->selected_index - max_visible + 1;
    }
    state->timeout_remaining = -1;
}

void menu_move_first(menu_state_t *state) {
    state->selected_index = 0;
    state->scroll_offset = 0;
    state->timeout_remaining = -1;
}

void menu_move_last(menu_state_t *state) {
    state->selected_index = state->visible_entries - 1;
    state->timeout_remaining = -1;
}

menu_entry_t *menu_select_current(menu_state_t *state) {
    return menu_get_visible_entry(state, state->selected_index);
}

void menu_render(teletext_screen_t *screen, menu_state_t *state) {
    teletext_clear(screen);

    uint8_t fg = state->config.fg_color;
    uint8_t bg = state->config.bg_color;
    uint8_t hf = state->config.highlight_fg;
    uint8_t hb = state->config.highlight_bg;
    uint8_t tf = state->config.title_fg;
    uint8_t tb = state->config.title_bg;

    /* Draw title bar */
    teletext_fill(screen, 0, 0, 0, TELETEXT_COLS - 1, ' ', tf, tb);
    teletext_puts(screen, 0, 2, state->config.title, tf, tb);

    /* Draw platform info on the right side of title bar */
    const char *platform_name = detect_platform_name(state->platform);
    uint8_t plat_col = TELETEXT_COLS - 2 - strlen(platform_name);
    teletext_puts(screen, 0, plat_col, "[", tf, tb);
    teletext_puts(screen, 0, plat_col + 1, platform_name, TELETEXT_YELLOW, tb);
    teletext_puts(screen, 0, plat_col + 1 + strlen(platform_name), "]", tf, tb);

    /* Draw separator line below title */
    teletext_draw_hline(screen, 1, 0, TELETEXT_COLS - 1, fg);

    /* Draw menu entries */
    uint8_t max_visible = TELETEXT_ROWS - 6;
    uint8_t start_row = 3;
    uint8_t entry_row = start_row;

    for (uint8_t i = state->scroll_offset;
         i < state->visible_entries && i < state->scroll_offset + max_visible;
         i++) {
        menu_entry_t *entry = menu_get_visible_entry(state, i);
        if (!entry) continue;

        uint8_t is_selected = (i == state->selected_index);
        uint8_t ef = is_selected ? hf : fg;
        uint8_t eb = is_selected ? hb : bg;

        /* Selection indicator */
        if (is_selected) {
            teletext_puts(screen, entry_row, 2, "▸", TELETEXT_YELLOW, eb);
        } else {
            teletext_puts(screen, entry_row, 2, " ", fg, eb);
        }

        /* Entry number */
        char num_str[4];
        num_str[0] = '[';
        num_str[1] = '0' + (i + 1) % 10;
        if (i + 1 >= 10) num_str[0] = '0' + (i + 1) / 10;
        num_str[2] = ']';
        num_str[3] = '\0';
        teletext_puts(screen, entry_row, 4, num_str, ef, eb);

        /* Entry label */
        teletext_puts(screen, entry_row, 8, entry->label, ef, eb);

        /* Entry type indicator */
        const char *type_str = "";
        switch (entry->type) {
            case BOOT_TYPE_CHAINLOAD:    type_str = " [chainload]"; break;
            case BOOT_TYPE_LINUX_KERNEL: type_str = " [kernel]";   break;
            case BOOT_TYPE_EFI_STUB:     type_str = " [efi]";      break;
            case BOOT_TYPE_MEMTEST:      type_str = " [memtest]";  break;
            case BOOT_TYPE_SUBMENU:      type_str = " >";          break;
            default: break;
        }
        if (type_str[0]) {
            teletext_puts(screen, entry_row, TELETEXT_COLS - 2 - strlen(type_str),
                          type_str, TELETEXT_CYAN, eb);
        }

        entry_row++;
    }

    /* Draw separator line */
    teletext_draw_hline(screen, TELETEXT_ROWS - 3, 0, TELETEXT_COLS - 1, fg);

    /* Draw help text */
    teletext_puts(screen, TELETEXT_ROWS - 2, 2, state->config.help_text, TELETEXT_CYAN, bg);

    /* Draw timeout if active */
    if (state->timeout_remaining >= 0) {
        char timeout_str[20];
        int len = snprintf(timeout_str, sizeof(timeout_str),
                          "Auto-boot in %ds...", state->timeout_remaining);
        teletext_puts(screen, TELETEXT_ROWS - 2,
                     TELETEXT_COLS - 2 - len, timeout_str, TELETEXT_YELLOW, bg);
    }

    /* Draw function key bar */
    teletext_fill(screen, TELETEXT_ROWS - 1, 0, TELETEXT_ROWS - 1, TELETEXT_COLS - 1,
                  ' ', TELETEXT_BLACK, TELETEXT_WHITE);
    teletext_puts(screen, TELETEXT_ROWS - 1, 1,
                  "[F1] Help  [F2] CHASIS  [F3] Security  [Tab] Boot Order",
                  TELETEXT_BLACK, TELETEXT_WHITE);
}

int menu_handle_key(menu_state_t *state, uint8_t key) {
    switch (key) {
        case 0x48: /* Up arrow */
        case 'k':
            menu_move_up(state);
            return 1;
        case 0x50: /* Down arrow */
        case 'j':
            menu_move_down(state);
            return 1;
        case 0x47: /* Home */
        case 'g':
            menu_move_first(state);
            return 1;
        case 0x4F: /* End */
        case 'G':
            menu_move_last(state);
            return 1;
        case 0x1B: /* Escape */
            return 0; /* Boot default */
        case '\r':  /* Enter */
        case ' ':
            return 0; /* Boot selected */
        default:
            /* Number keys for direct selection */
            if (key >= '1' && key <= '9') {
                uint8_t idx = key - '1';
                if (idx < state->visible_entries) {
                    state->selected_index = idx;
                    return 0; /* Boot selected */
                }
            }
            return 1; /* Unknown key, keep menu */
    }
}

int menu_tick_timeout(menu_state_t *state) {
    if (state->timeout_remaining < 0) return 0;
    if (state->timeout_remaining == 0) return 1; /* Timeout! */
    state->timeout_remaining--;
    return 0;
}

const char *menu_get_help_text(menu_state_t *state) {
    (void)state;
    return "SonicScrewdriver v2.0.0 — Universal USB Bootloader";
}
