#ifndef TELETEXT_H
#define TELETEXT_H

#include <stdint.h>
#include <stddef.h>

/* Teletext screen dimensions */
#define TELETEXT_COLS        80
#define TELETEXT_ROWS        25
#define TELETEXT_CELLS       (TELETEXT_COLS * TELETEXT_ROWS)

/* Teletext colors (16-color palette) */
typedef enum {
    TELETEXT_BLACK       = 0,
    TELETEXT_RED         = 1,
    TELETEXT_GREEN       = 2,
    TELETEXT_YELLOW      = 3,
    TELETEXT_BLUE        = 4,
    TELETEXT_MAGENTA     = 5,
    TELETEXT_CYAN        = 6,
    TELETEXT_WHITE       = 7,
    TELETEXT_FLASH       = 8,  /* Flash mask */
    TELETEXT_STEADY      = 9,  /* Steady mask */
    TELETEXT_BLACK_BG    = 10,
    TELETEXT_NEW_BG      = 11,
    TELETEXT_HOLD_MOSAIC = 12,
    TELETEXT_CONCEAL     = 13,
    TELETEXT_CONTIGUOUS  = 14,
    TELETEXT_SEPARATED   = 15
} teletext_color_t;

/* Teletext character cell */
typedef struct {
    uint8_t character;      /* ASCII or block graphics character */
    uint8_t foreground : 4; /* Foreground color (0-15) */
    uint8_t background : 4; /* Background color (0-15) */
    uint8_t bold       : 1; /* Bold flag */
    uint8_t flash      : 1; /* Flash flag */
    uint8_t underline  : 1; /* Underline flag */
    uint8_t graphics   : 1; /* Block graphics mode */
    uint8_t reserved   : 4;
} __attribute__((packed)) teletext_cell_t;

/* Teletext screen buffer */
typedef struct {
    teletext_cell_t cells[TELETEXT_ROWS][TELETEXT_COLS];
    uint8_t cursor_row;
    uint8_t cursor_col;
    uint8_t default_fg;
    uint8_t default_bg;
} teletext_screen_t;

/* Block graphics characters (CP437 / Teletext semigraphics) */
#define TELETEXT_BLOCK_UPPER_HALF    0xDF
#define TELETEXT_BLOCK_LOWER_HALF    0xDC
#define TELETEXT_BLOCK_FULL          0xDB
#define TELETEXT_BLOCK_LEFT_HALF     0xDE
#define TELETEXT_BLOCK_RIGHT_HALF    0xDD
#define TELETEXT_BLOCK_DARK_SHADE    0xB1
#define TELETEXT_BLOCK_MED_SHADE     0xB2
#define TELETEXT_BLOCK_LIGHT_SHADE   0xB0

/* Box drawing characters */
#define TELETEXT_BOX_HORIZONTAL      0xC4
#define TELETEXT_BOX_VERTICAL        0xB3
#define TELETEXT_BOX_TOP_LEFT        0xDA
#define TELETEXT_BOX_TOP_RIGHT       0xBF
#define TELETEXT_BOX_BOTTOM_LEFT     0xC0
#define TELETEXT_BOX_BOTTOM_RIGHT    0xD9
#define TELETEXT_BOX_CROSS           0xC5
#define TELETEXT_BOX_T_DOWN          0xC2
#define TELETEXT_BOX_T_UP            0xC1
#define TELETEXT_BOX_T_RIGHT         0xC3
#define TELETEXT_BOX_T_LEFT          0xB4

/* Double-line box drawing */
#define TELETEXT_BOX_DBL_HORIZ       0xCD
#define TELETEXT_BOX_DBL_VERT        0xBA
#define TELETEXT_BOX_DBL_TL          0xC9
#define TELETEXT_BOX_DBL_TR          0xBB
#define TELETEXT_BOX_DBL_BL          0xC8
#define TELETEXT_BOX_DBL_BR          0xBC

/* Arrow characters */
#define TELETEXT_ARROW_UP            0x18
#define TELETEXT_ARROW_DOWN          0x19
#define TELETEXT_ARROW_LEFT          0x1B
#define TELETEXT_ARROW_RIGHT         0x1A
#define TELETEXT_TRIANGLE_RIGHT      0x10
#define TELETEXT_TRIANGLE_LEFT       0x11

/* Special characters */
#define TELETEXT_MUSIC_NOTE          0x0D
#define TELETEXT_SKULL               0x0E
#define TELETEXT_HEART               0x03
#define TELETEXT_CHECKMARK           0xFB
#define TELETEXT_CROSSMARK           0xFD
#define TELETEXT_DOT_OPERATOR        0xF9
#define TELETEXT_DEGREE_SIGN         0xF8
#define TELETEXT_PLUS_MINUS          0xF1

/* Function prototypes */

/* Initialize a teletext screen with default colors */
void teletext_init(teletext_screen_t *screen);

/* Clear the screen (fill with spaces and default colors) */
void teletext_clear(teletext_screen_t *screen);

/* Set a character at (row, col) with specified attributes */
void teletext_putchar(teletext_screen_t *screen, uint8_t row, uint8_t col,
                      uint8_t ch, uint8_t fg, uint8_t bg);

/* Write a null-terminated string at (row, col) */
void teletext_puts(teletext_screen_t *screen, uint8_t row, uint8_t col,
                   const char *str, uint8_t fg, uint8_t bg);

/* Write a string with length limit at (row, col) */
void teletext_putsn(teletext_screen_t *screen, uint8_t row, uint8_t col,
                    const char *str, size_t n, uint8_t fg, uint8_t bg);

/* Draw a horizontal line at (row, col_start) to (row, col_end) */
void teletext_draw_hline(teletext_screen_t *screen, uint8_t row,
                         uint8_t col_start, uint8_t col_end, uint8_t fg);

/* Draw a vertical line at (col, row_start) to (col, row_end) */
void teletext_draw_vline(teletext_screen_t *screen, uint8_t col,
                         uint8_t row_start, uint8_t row_end, uint8_t fg);

/* Draw a bordered box from (row1, col1) to (row2, col2) */
void teletext_draw_box(teletext_screen_t *screen, uint8_t row1, uint8_t col1,
                       uint8_t row2, uint8_t col2, uint8_t fg, uint8_t bg);

/* Draw a double-line bordered box */
void teletext_draw_dbl_box(teletext_screen_t *screen, uint8_t row1, uint8_t col1,
                           uint8_t row2, uint8_t col2, uint8_t fg, uint8_t bg);

/* Set cursor position */
void teletext_set_cursor(teletext_screen_t *screen, uint8_t row, uint8_t col);

/* Scroll the screen up by one row */
void teletext_scroll(teletext_screen_t *screen);

/* Fill a region with a character and attributes */
void teletext_fill(teletext_screen_t *screen, uint8_t row1, uint8_t col1,
                   uint8_t row2, uint8_t col2, uint8_t ch, uint8_t fg, uint8_t bg);

/* Set default foreground/background colors */
void teletext_set_defaults(teletext_screen_t *screen, uint8_t fg, uint8_t bg);

/* Render the screen to a framebuffer (platform-specific callback) */
typedef void (*teletext_render_cb)(teletext_cell_t *cells, size_t count, void *context);
void teletext_render(teletext_screen_t *screen, teletext_render_cb callback, void *context);

#endif /* TELETEXT_H */
