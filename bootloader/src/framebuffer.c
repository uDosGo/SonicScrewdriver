#include "teletext.h"
#include <string.h>

/*
 * Framebuffer abstraction layer.
 * Supports both UEFI GOP (Graphics Output Protocol) and BIOS VGA text mode.
 * At build time, only one backend is compiled into the binary.
 */

/* Platform detection flag — set at compile time */
#ifdef UEFI_BUILD
    #define IS_UEFI 1
#else
    #define IS_UEFI 0
#endif

/* UEFI framebuffer state */
typedef struct {
    void *frame_buffer;
    uint32_t horizontal_resolution;
    uint32_t vertical_resolution;
    uint32_t pixels_per_scan_line;
    uint32_t bytes_per_pixel;  /* 4 for BGRA32 */
} uefi_fb_t;

static uefi_fb_t uefi_fb = {0};

/* BIOS VGA text mode state */
#define VGA_TEXT_BASE 0xB8000
#define VGA_TEXT_COLS 80
#define VGA_TEXT_ROWS 25

/* VGA text mode cell (hardware format) */
typedef struct {
    uint8_t character;
    uint8_t attributes;  /* 4-bit bg, 4-bit fg */
} __attribute__((packed)) vga_text_cell_t;

/* CP437 to Teletext block graphics mapping */
static const uint8_t cp437_glyphs[128] = {
    /* 0x00-0x1F: Control characters — map to spaces */
    0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
    0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
    0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
    0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
    /* 0x20-0x3F: Standard ASCII */
    0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
    0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
    0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
    0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
    /* 0x40-0x5F */
    0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
    0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
    0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
    0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
    /* 0x60-0x7F */
    0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
    0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
    0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77,
    0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F
};

/* Convert teletext color to VGA attribute byte */
static uint8_t teletext_to_vga_attr(uint8_t fg, uint8_t bg) {
    /* VGA: bits 0-3 = fg, bit 3 = bright, bits 4-6 = bg, bit 7 = blink */
    uint8_t vga_fg = fg & 0x07;
    uint8_t vga_bright = (fg & 0x08) ? 0x08 : 0x00;
    uint8_t vga_bg = (bg & 0x07) << 4;
    return vga_fg | vga_bright | vga_bg;
}

/* Convert teletext color to BGRA32 pixel */
static uint32_t teletext_to_bgra(uint8_t color) {
    /* Teletext 16-color palette mapped to approximate BGRA32 values */
    static const uint32_t palette[16] = {
        0xFF000000, /* Black   */
        0xFF0000AA, /* Red     */
        0xFF00AA00, /* Green   */
        0xFF00AAAA, /* Yellow  */
        0xFFAA0000, /* Blue    */
        0xFFAA00AA, /* Magenta */
        0xFFAA5500, /* Cyan    */
        0xFFAAAAAA, /* White   */
        0xFF555555, /* Bright Black (Dark Gray) */
        0xFF5555FF, /* Bright Red */
        0xFF55FF55, /* Bright Green */
        0xFF55FFFF, /* Bright Yellow */
        0xFFFF5555, /* Bright Blue */
        0xFFFF55FF, /* Bright Magenta */
        0xFFFFFF55, /* Bright Cyan */
        0xFFFFFFFF, /* Bright White */
    };
    if (color >= 16) color = 7;
    return palette[color];
}

/* Character bitmap cache for UEFI rendering (8x16 pixels per char) */
static const uint8_t font_8x16[256][16] = {{0}}; /* Will be populated from embedded font */

/* Initialize UEFI framebuffer */
void framebuffer_init_uefi(void *gop_buffer,
                           uint32_t hres, uint32_t vres,
                           uint32_t scanline, uint32_t bpp) {
    uefi_fb.frame_buffer = gop_buffer;
    uefi_fb.horizontal_resolution = hres;
    uefi_fb.vertical_resolution = vres;
    uefi_fb.pixels_per_scan_line = scanline;
    uefi_fb.bytes_per_pixel = bpp;
}

/* Check if we're in UEFI mode */
int framebuffer_is_uefi(void) {
    return IS_UEFI;
}

/* Render a single character in UEFI GOP framebuffer */
static void framebuffer_putchar_uefi(uint8_t row, uint8_t col,
                                     teletext_cell_t *cell) {
    if (!uefi_fb.frame_buffer) return;
    if (row >= TELETEXT_ROWS || col >= TELETEXT_COLS) return;

    uint32_t bg_color = teletext_to_bgra(cell->background);
    uint32_t fg_color = teletext_to_bgra(cell->foreground);
    uint8_t ch = cell->character;

    /* Calculate pixel position */
    uint32_t px = col * 8;
    uint32_t py = row * 16;

    /* Draw character bitmap (8x16) */
    for (int y = 0; y < 16 && (py + y) < uefi_fb.vertical_resolution; y++) {
        uint8_t glyph_row = font_8x16[ch][y];
        for (int x = 0; x < 8 && (px + x) < uefi_fb.horizontal_resolution; x++) {
            uint32_t pixel = (glyph_row & (0x80 >> x)) ? fg_color : bg_color;
            uint32_t *fb = (uint32_t *)uefi_fb.frame_buffer;
            fb[(py + y) * uefi_fb.pixels_per_scan_line + (px + x)] = pixel;
        }
    }
}

/* Render the entire teletext screen to UEFI framebuffer */
static void framebuffer_render_uefi(teletext_cell_t *cells, size_t count) {
    (void)count;
    for (uint8_t r = 0; r < TELETEXT_ROWS; r++) {
        for (uint8_t c = 0; c < TELETEXT_COLS; c++) {
            framebuffer_putchar_uefi(r, c, &cells[r * TELETEXT_COLS + c]);
        }
    }
}

/* Render the entire teletext screen to BIOS VGA text mode */
static void framebuffer_render_vga(teletext_cell_t *cells, size_t count) {
    (void)count;
    volatile vga_text_cell_t *vga = (volatile vga_text_cell_t *)VGA_TEXT_BASE;

    for (int i = 0; i < TELETEXT_CELLS; i++) {
        vga[i].character = cells[i].character;
        vga[i].attributes = teletext_to_vga_attr(cells[i].foreground,
                                                  cells[i].background);
    }
}

/* Public render callback — called by teletext_render() */
void framebuffer_render(teletext_cell_t *cells, size_t count, void *context) {
    (void)context;
    if (IS_UEFI) {
        framebuffer_render_uefi(cells, count);
    } else {
        framebuffer_render_vga(cells, count);
    }
}

/* Clear the framebuffer */
void framebuffer_clear(void) {
    if (IS_UEFI && uefi_fb.frame_buffer) {
        uint32_t *fb = (uint32_t *)uefi_fb.frame_buffer;
        uint32_t black = teletext_to_bgra(TELETEXT_BLACK);
        for (uint32_t y = 0; y < uefi_fb.vertical_resolution; y++) {
            for (uint32_t x = 0; x < uefi_fb.horizontal_resolution; x++) {
                fb[y * uefi_fb.pixels_per_scan_line + x] = black;
            }
        }
    } else {
        /* Clear VGA text mode */
        volatile vga_text_cell_t *vga = (volatile vga_text_cell_t *)VGA_TEXT_BASE;
        for (int i = 0; i < VGA_TEXT_COLS * VGA_TEXT_ROWS; i++) {
            vga[i].character = ' ';
            vga[i].attributes = 0x07; /* White on black */
        }
    }
}
