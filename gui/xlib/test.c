#include "xlib.h"
#include "xlib.c"

#define STR_MAX_SIZE 128

const char *X_WINDOW_NAME = "GNOME 2048";
const char *GAME_NAME = "gnome-2048";
const char *GAME_ARG = "--size=4";
const char *COLOR_BOLD = "\033[1m";
const char *COLOR_END = "\033[0m";

const int GAME_SIZE = 4;

typedef int (*test)(void);

int test_exec_program() {
    printf(" test_exec_program:%s\n", COLOR_END);
    printf("  => checking if program [%s] is available\n", GAME_NAME);

    FILE *fp;
    char path[STR_MAX_SIZE];
    char *args[10];
    char cmd[] = "which ";
    int len;

    if (strcat(cmd, GAME_NAME) == NULL) {
        fprintf(stderr, "  strcat() failed\n");
        return 1;
    }

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

int test_open_x_socket() {
    printf(" test_open_x_socket:%s\n", COLOR_END);
    return open_x_socket();
}

int test_find_x_window() {
    printf(" test_find_x_window:%s\n", COLOR_END);
    return find_x_window(X_WINDOW_NAME, GAME_SIZE);
}

int test_scan_pixels() {
    printf(" test_scan_pixels:%s\n", COLOR_END);
    unsigned long *pixels = scan_pixels();
    if (pixels == NULL) {
        printf("  => scan_pixels returned NULL\n");
        return 1;
    }

    // all but two tiles should be white

    unsigned long a_pixel, b_pixel;
    int a_count, b_count;
    int i, ok;
    int n = GAME_SIZE * GAME_SIZE;
    a_count = b_count = 0;

    printf("  => got pixels:\n");

    for (i=0; i<n; i++) {
        printf(" %lu", pixels[i]);
        if (i % GAME_SIZE == 0) {
            printf("\n");
        }

        if ( a_count == 0 ) {
            a_pixel = pixels[i];
            a_count = 1;
            continue;
        }
        
        if ( b_count == 0 ) {
            b_pixel = pixels[i];
            b_count = 1;
            continue;
        }

        if ( pixels[i] == a_pixel ) {
            a_count++;
            continue;
        }

        if ( pixels[i] = b_pixel ) {
            b_count++;
            continue;
        }
    }

    ok = 0;

    if (a_count+b_count != n) {
        printf("  => pixels don't contain only 2 colors\n");
        ok = 1;
    }

    printf("  => number of white tiles:     %d\n", (a_count < b_count) ? a_count : b_count);
    printf("  => number of non-white tiles: %d\n", (a_count < b_count) ? b_count : a_count);
 
    if (a_count != 2 && b_count != 2) {
        printf("  => expected 2 white tiles\n");
        ok = 1;
    }

    return ok;
}

int test_scan_pixels2() {
    printf(" test_scan_pixels2:%s (skipped)\n", COLOR_END);
    return 0;
}

int test_fake_key_event() {
    printf(" test_fake_key_event:%s\n", COLOR_END);

    const char *moves[4] = {"Up", "Down", "Left", "Right",};
    for (int i=0; i<4; i++) {
        printf("  => %s\n", moves[i]);
        fake_key_event(moves[i]);
    }

    return 0;
}

int num_tests = 6;

test tests[] = {
    &test_exec_program,
    &test_open_x_socket,
    &test_find_x_window,
    &test_scan_pixels,
    &test_fake_key_event,
    &test_scan_pixels2,
};

int main() {

    int i;

    for (i=0; i<num_tests; i++) {
        printf("%sTEST %-3d:", COLOR_BOLD, i);
        if ( tests[i]() != 0 ) {
            return 1;
        }
    }

    printf("\n%s=> all tests passed%s\n", COLOR_BOLD, COLOR_END);
    return 0;
}
