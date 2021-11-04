#include <gui.h>

static Display *dpy;
static Window root, target;

static int pos[4]; // position of target window: x, y, width and height
static int *coords[2]; // coordinates of the tiles in the grid
static int n; // size of the grid

// Open connection to the X server and fix the root window
// Returns 0 on success.
int open() {
    if ( (dpy = XOpenDisplay(NULL) != NULL ) ) {
        root = XDefaultRootWindow(dpy);
        return 0;
    }
    return 1;
}

// Close connection to the X server
void close() {
    XCloseDisplay(dpy);
}

// Find window and fix its XID and position. We expect that this window contains
// size x size tiles and we will try to approximate their coordinates.
int find_window(const char* name, int size) {

    n = size*size;

    // find XID of the window
    if ( (target = find_window_recursively(name) == 0 ) ) {
        return 1;
    }

    // find position of the window, populates static array *pos[4]*
    // with values for x, y (top-left pixel), width, height
    if (get_coordinates() != 0) {
        return 1;
    }

    // the width and heigth values are sometimes larger than the actual
    // window, so try to find the value for GTK frame extents to normalize
    // them. If we can't get frames, we just assume they don't exist.
    int frames[4];
    frames = get_gtk_frame_extents();

    normalize_coordinates(frames);
    return 0;
}

// Get the pixel values of the tiles in the grid
unsigned long* scan_pixels() {

    unsigned long* pixels;
    int i;
    XImage *image;

    image = XGetImage(dpy, target, pos[0], pos[1], pos[2], pos[3], AllPlanes, XYPixmap);
    if (image == NULL) {
        return NULL;
    }

    pixels = (unsigned long*) malloc(n);

    for (i=0; i<n; i++) {
        pixels[i] = XGetPixel(image, coords[0], coords[1]); // x, y coordinate of the grid
    }

    return pixels;
}

void fake_key_event(const char* key) {

    unsigned long delay = 100;

    KeyCode keycode = XKeysymToKeycode(dpy, XStringToKeysym(key));

    XTestFakeKeyEvent(dpy, keycode, False, delay); // key is is not pressed
    XFlush(dpy);
    XTestFakeKeyEvent(dpy, keycode, True, delay);  // key is pressed
    XFlush(dpy);
    XTestFakeKeyEvent(dpy, keycode, False, delay); // key is released
    XFlush(dpy);

}
