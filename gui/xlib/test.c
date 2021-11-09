/* This file tests/debugs the public methods of xlib.c independent from the 
 * main Go-package. It only makes sense to run this code when gnome-2048
 * is installed on the system.
 */

#include "xlib.h"
#include "xlib.c"

#define STR_MAX_SIZE 128

const char *X_WINDOW_NAME = "GNOME 2048";
const char *GAME_NAME = "gnome-2048";
const char *GAME_ARG = "--size=4";
const char *COLOR_BOLD = "\033[1m";
const char *COLOR_END = "\033[0m";

const int GAME_SIZE = 4;

unsigned long *pixels;

typedef int (*test)(void);


// Check if "gnome-2048" is installed and run in the background.
int test_exec_program() {
    printf(" test_exec_program:%s [%s]\n", COLOR_END, GAME_NAME);

    FILE *fp;
    char path[STR_MAX_SIZE];
    char *args[10];
    char cmd[] = "which ";
    int len;

    if (strcat(cmd, GAME_NAME) == NULL) {
        fprintf(stderr, "  strcat() failed\n");
        return 1;
    }

    // read path of executable from pipe
    if ( (fp = popen(cmd, "r")) == NULL) {
        fprintf(stderr, "  popen() failed\n");
        return 1;
    }

    if (fgets(path, STR_MAX_SIZE, fp) != NULL) {
        // strip trailing newline
        len = strlen(path);
        if (path[len-1] == '\n') {
            path[len-1] = '\0';
        }
    }

    // error here means which exited with non-0 status
    if (pclose(fp) != 0 ) {
        fprintf(stderr, "  program [%s] not found\n");
        return 1;
    }

    printf("  => running [%s] in background\n", path);

    switch (fork()) {
        case -1:
            fprintf(stderr, "  fork() failed\n");
            return 1;
        case 0:
            args[0] = path;
            args[1] = (char*) GAME_ARG;
            args[2] = NULL;
            execv(path, args);
            fprintf(stderr, "  execv() failed\n");
        default:
            printf("  => wait 1 second...\n");
            sleep(1);
            return 0;
    }
}

// Check if we are able to find target window
// and fix its coordinates
int test_find_window() {
    printf(" test_find_x_window:%s\n", COLOR_END);
    return find_window(X_WINDOW_NAME, GAME_SIZE);
}

// Check if we are able to take a screenshot and
// scan pixels of the game board. (Initially the
// board should contain all, but 2 non-white tiles,
// so we can check if the pixels are consistent.
int test_scan_pixels() {
    printf(" test_scan_pixels:%s\n", COLOR_END);
    pixels = scan_pixels();
    if (pixels == NULL) {
        printf("  => scan_pixels returned NULL\n");
        return 1;
    }

    // all but two tiles should be white
    unsigned long a_pixel, b_pixel;
    int a_count, b_count;
    int i, ok;
    int n = GAME_SIZE * GAME_SIZE;
    a_count = 0;
    b_count = 0;

    printf("  => checking %d pixels:\n", n);

    for (i=0; i<n; i++) {
        printf(" %lu%s", pixels[i], ((i+1) % GAME_SIZE == 0) ? "\n" : "");

        if ( a_count != 0 && pixels[i] == a_pixel ) {
            a_count++;
        } else if ( b_count != 0 && pixels[i] == b_pixel ) {
            b_count++;
        } else if ( a_count == 0 ) {
            a_pixel = pixels[i];
            a_count = 1;
        } else if ( b_count == 0 ) {
            b_pixel = pixels[i];
            b_count = 1;
        }
    }
    printf("\n");

    ok = 0;

    if (a_count+b_count != n) {
        printf("  => pixels don't contain only 2 colors\n");
        ok = 1;
    }

    if (a_count > b_count) {
        printf("  => white  [%lu] tiles: %d \n", a_pixel, a_count);
        printf("  => yellow [%lu] tiles: %d \n", b_pixel, b_count);
    } else {
        printf("  => white  [%lu] tiles: %d \n", b_pixel, b_count);
        printf("  => yellow [%lu] tiles: %d \n", a_pixel, a_count);
    }
 
    if (a_count != 2 && b_count != 2) {
        printf("  => expected 2 white tiles\n");
        ok = 1;
    }

    return ok;
}

// This test should run after test_fake_key_event(),
// so we expect that the tiles have changed.
int test_scan_pixels2() {
    printf(" test_scan_pixels2:%s check if board changed after moves\n", COLOR_END);
    int ok = 1;
    int i;
    int n = GAME_SIZE * GAME_SIZE;
    unsigned long *pixels2 = scan_pixels();
    if (pixels2 == NULL) {
        printf("  => scan_pixels returned NULL\n");
        return 1;
    }

    for (i=0; i<n; i++) {
        if ( pixels2[i] != pixels[i] ) {
            ok = 0;
            break;
        }
    }

    return ok;
}

// This is a pseudo-test, since we don't really know
// if the key events really take an effect.
int test_fake_key_event() {
    printf(" test_fake_key_event:%s\n", COLOR_END);

    const char *moves[4] = {"Up", "Down", "Left", "Right",};
    for (int i=0; i<4; i++) {
        printf("  => [%s]\n", moves[i]);
        sleep(1);
        if (fake_key_event(moves[i], 50) != 0) {
            return 1;
        }
    }

    return 0;
}

test tests[] = {
    &test_exec_program,
    &test_find_window,
    &test_scan_pixels,
    &test_fake_key_event,
    &test_scan_pixels2,
    NULL,
};

int main() {

    int i, ok;
    i = 0;
    while (tests[i] != NULL) {
        printf(" =%stest %-2d:", COLOR_BOLD, i);
        if ( (ok = tests[i]()) != 0 ) {
            printf("=FAIL\n");
            break;
        }
        printf("=PASS\n");
        i++;
    }

    return ok;
}
