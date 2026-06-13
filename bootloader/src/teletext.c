#include "teletext.h"
#include <string.h>

void teletext_init(teletext_screen_t *screen) {
    screen->cursor_row = 0;
    screen->cursor_col = 0;
    screen->default_fg = TELETEXT_WHITE;
    screen->default_bg = TELETEXT_BLACK;
    teletext_clear(screen);
}

void teletext_clear(teletext_screen_t *screen) {
    for (int r = 0; r < TELETEXT_ROWS; r++) {
        for (int c = 0; c < TELETEXT_COLS; c++) {
            screen->cells[r][c].character = ' ';
            screen->cells[r][c].foreground = screen->default_fg;
            screen->cells[r][c].background = screen->default_bg;
            screen->cells[r][c].bold = 0;
            screen->cells[r][c].flash = 0;
            screen->cells[r][c].underline = 0;
            screen->cells[r][c].graphics = 0;
        }
    }
}

void teletext_putchar(teletext_screen_t *screen, uint8_t row, uint8_t col,
                      uint8_t ch, uint8_t fg, uint8_t bg) {
    if (row >= TELETEXT_ROWS || col >= TELETEXT_COLS) return;
    screen->cells[row][col].character = ch;
    screen->cells[row][col].foreground = fg;
    screen->cells[row][col].background = bg;
}

void teletext_puts(teletext_screen_t *screen, uint8_t row, uint8_t col,
                   const char *str, uint8_t fg, uint8_t bg) {
    if (row >= TELETEXT_ROWS) return;
    size_t len = strlen(str);
    size_t max = TELETEXT_COLS - col;
    if (len > max) len = max;
    for (size_t i = 0; i < len; i++) {
        screen->cells[row][col + i].character = (uint8_t)str[i];
        screen->cells[row][col + i].foreground = fg;
        screen->cells[row][col + i].background = bg;
    }
}

void teletext_putsn(teletext_screen_t *screen, uint8_t row, uint8_t col,
                    const char *str, size_t n, uint8_t fg, uint8_t bg) {
    if (row >= TELETEXT_ROWS) return;
    size_t max = TELETEXT_COLS - col;
    if (n > max) n = max;
    for (size_t i = 0; i < n; i++) {
        screen->cells[row][col + i].character = (uint8_t)str[i];
        screen->cells[row][col + i].foreground = fg;
        screen->cells[row][col + i].background = bg;
    }
}

void teletext_draw_hline(teletext_screen_t *screen, uint8_t row,
                         uint8_t col_start, uint8_t col_end, uint8_t fg) {
    if (row >= TELETEXT_ROWS) return;
    if (col_start > col_end) {
        uint8_t tmp = col_start;
        col_start = col_end;
        col_end = tmp;
    }
    if (col_end >= TELETEXT_COLS) col_end = TELETEXT_COLS - 1;
    for (uint8_t c = col_start; c <= col_end; c++) {
        screen->cells[row][c].character = TELETEXT_BOX_HORIZONTAL;
        screen->cells[row][c].foreground = fg;
    }
}

void teletext_draw_vline(teletext_screen_t *screen, uint8_t col,
                         uint8_t row_start, uint8_t row_end, uint8_t fg) {
    if (col >= TELETEXT_COLS) return;
    if (row_start > row_end) {
        uint8_t tmp = row_start;
        row_start = row_end;
        row_end = tmp;
    }
    if (row_end >= TELETEXT_ROWS) row_end = TELETEXT_ROWS - 1;
    for (uint8_t r = row_start; r <= row_end; r++) {
        screen->cells[r][col].character = TELETEXT_BOX_VERTICAL;
        screen->cells[r][col].foreground = fg;
    }
}

void teletext_draw_box(teletext_screen_t *screen, uint8_t row1, uint8_t col1,
                       uint8_t row2, uint8_t col2, uint8_t fg, uint8_t bg) {
    if (row1 > row2) { uint8_t t = row1; row1 = row2; row2 = t; }
    if (col1 > col2) { uint8_t t = col1; col1 = col2; col2 = t; }
    if (row2 >= TELETEXT_ROWS) row2 = TELETEXT_ROWS - 1;
    if (col2 >= TELETEXT_COLS) col2 = TELETEXT_COLS - 1;

    /* Fill interior */
    for (uint8_t r = row1 + 1; r < row2; r++) {
        for (uint8_t c = col1 + 1; c < col2; c++) {
            screen->cells[r][c].character = ' ';
            screen->cells[r][c].foreground = fg;
            screen->cells[r][c].background = bg;
        }
    }

    /* Draw corners */
    screen->cells[row1][col1].character = TELETEXT_BOX_TOP_LEFT;
    screen->cells[row1][col2].character = TELETEXT_BOX_TOP_RIGHT;
    screen->cells[row2][col1].character = TELETEXT_BOX_BOTTOM_LEFT;
    screen->cells[row2][col2].character = TELETEXT_BOX_BOTTOM_RIGHT;

    /* Draw edges */
    for (uint8_t c = col1 + 1; c < col2; c++) {
        screen->cells[row1][c].character = TELETEXT_BOX_HORIZONTAL;
        screen->cells[row2][c].character = TELETEXT_BOX_HORIZONTAL;
    }
    for (uint8_t r = row1 + 1; r < row2; r++) {
        screen->cells[r][col1].character = TELETEXT_BOX_VERTICAL;
        screen->cells[r][col2].character = TELETEXT_BOX_VERTICAL;
    }

    /* Set corner colors */
    for (uint8_t r = row1; r <= row2; r++) {
        for (uint8_t c = col1; c <= col2; c++) {
            screen->cells[r][c].foreground = fg;
            screen->cells[r][c].background = bg;
        }
    }
}

void teletext_draw_dbl_box(teletext_screen_t *screen, uint8_t row1, uint8_t col1,
                           uint8_t row2, uint8_t col2, uint8_t fg, uint8_t bg) {
    if (row1 > row2) { uint8_t t = row1; row1 = row2; row2 = t; }
    if (col1 > col2) { uint8_t t = col1; col1 = col2; col2 = t; }
    if (row2 >= TELETEXT_ROWS) row2 = TELETEXT_ROWS - 1;
    if (col2 >= TELETEXT_COLS) col2 = TELETEXT_COLS - 1;

    /* Fill interior */
    for (uint8_t r = row1 + 1; r < row2; r++) {
        for (uint8_t c = col1 + 1; c < col2; c++) {
            screen->cells[r][c].character = ' ';
            screen->cells[r][c].foreground = fg;
            screen->cells[r][c].background = bg;
        }
    }

    /* Draw corners */
    screen->cells[row1][col1].character = TELETEXT_BOX_DBL_TL;
    screen->cells[row1][col2].character = TELETEXT_BOX_DBL_TR;
    screen->cells[row2][col1].character = TELETEXT_BOX_DBL_BL;
    screen->cells[row2][col2].character = TELETEXT_BOX_DBL_BR;

    /* Draw edges */
    for (uint8_t c = col1 + 1; c < col2; c++) {
        screen->cells[row1][c].character = TELETEXT_BOX_DBL_HORIZ;
        screen->cells[row2][c].character = TELETEXT_BOX_DBL_HORIZ;
    }
    for (uint8_t r = row1 + 1; r < row2; r++) {
        screen->cells[r][col1].character = TELETEXT_BOX_DBL_VERT;
        screen->cells[r][col2].character = TELETEXT_BOX_DBL_VERT;
    }

    for (uint8_t r = row1; r <= row2; r++) {
        for (uint8_t c = col1; c <= col2; c++) {
            screen->cells[r][c].foreground = fg;
            screen->cells[r][c].background = bg;
        }
    }
}

void teletext_set_cursor(teletext_screen_t *screen, uint8_t row, uint8_t col) {
    if (row < TELETEXT_ROWS) screen->cursor_row = row;
    if (col < TELETEXT_COLS) screen->cursor_col = col;
}

void teletext_scroll(teletext_screen_t *screen) {
    /* Move all rows up by one */
    for (int r = 1; r < TELETEXT_ROWS; r++) {
        for (int c = 0; c < TELETEXT_COLS; c++) {
            screen->cells[r - 1][c] = screen->cells[r][c];
        }
    }
    /* Clear last row */
    for (int c = 0; c < TELETEXT_COLS; c++) {
        screen->cells[TELETEXT_ROWS - 1][c].character = ' ';
        screen->cells[TELETEXT_ROWS - 1][c].foreground = screen->default_fg;
        screen->cells[TELETEXT_ROWS - 1][c].background = screen->default_bg;
    }
}

void teletext_fill(teletext_screen_t *screen, uint8_t row1, uint8_t col1,
                   uint8_t row2, uint8_t col2, uint8_t ch, uint8_t fg, uint8_t bg) {
    if (row1 > row2) { uint8_t t = row1; row1 = row2; row2 = t; }
    if (col1 > col2) { uint8_t t = col1; col1 = col2; col2 = t; }
    if (row2 >= TELETEXT_ROWS) row2 = TELETEXT_ROWS - 1;
    if (col2 >= TELETEXT_COLS) col2 = TELETEXT_COLS - 1;
    for (uint8_t r = row1; r <= row2; r++) {
        for (uint8_t c = col1; c <= col2; c++) {
            screen->cells[r][c].character = ch;
            screen->cells[r][c].foreground = fg;
            screen->cells[r][c].background = bg;
        }
    }
}

void teletext_set_defaults(teletext_screen_t *screen, uint8_t fg, uint8_t bg) {
    screen->default_fg = fg;
    screen->default_bg = bg;
}

void teletext_render(teletext_screen_t *screen, teletext_render_cb callback, void *context) {
    callback(&screen->cells[0][0], TELETEXT_CELLS, context);
}
