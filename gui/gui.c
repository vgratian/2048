#include <gui.h>

static Display *dpy;
static Window root, target;

static XWindowAttributes attrs; // position of target window: x, y, width and height
static int *coords[2]; // coordinates of the tiles in the grid
static int N; // size of the grid

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
    XFree(attrs);
    XCloseDisplay(dpy);
}

// Find window and fix its XID and position. We expect that this window contains
// size x size tiles and we will try to approximate their coordinates.
int find_window(const char* name, int size) {

    N = size;

    // find XID of the window
    if ( (target = find_window_recursively(root, name) == 0 ) ) {
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

    normalize_coordinates(frames[0], frames[1], frames[2], frames[3]);
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

    pixels = (unsigned long*) malloc(N*N);

    for (i=0; i<n; i++) {
        pixels[i] = XGetPixel(image, coords[i][0], coords[i][1]); // x, y coordinate of the grid
    }

    return pixels;
}

void fake_key_event(const char* key) {

    unsigned long delay = 100; // 0.1 second

    KeyCode keycode = XKeysymToKeycode(dpy, XStringToKeysym(key));

    XTestFakeKeyEvent(dpy, keycode, False, delay); // key is not pressed
    XFlush(dpy);
    XTestFakeKeyEvent(dpy, keycode, True, delay);  // key is pressed
    XFlush(dpy);
    XTestFakeKeyEvent(dpy, keycode, False, delay); // key is released
    XFlush(dpy);
}

static Window find_window_recursively(Window window, const char* target_name) {
    char *name;
    XFetchName(dpy, window, &name);
    if (name != NULL && strcmp(name, target_name) == 0) {
        return window;
    }

    Window root_return, parent_return, child_match;
    Window *children;
    int i, num;

    XQueryTree(dpy, window, &root_return, &parent_return, &children, &num);
    if (children != NULL) {
        for (i=0; i<num; i++) {
            child_match = find_window_recursively(children[i], target_name);
            if ( child_match ) {
                break;
            }
        }
        XFree(children);
        if (i < num) {
            return child_match;
        }
    }
    
    return 0;
}

static int get_coordinates() {
    if (XGetWindowAttributes(dpy, target, &attrs) == Success) {
        return 0;
    }
    return 1;
}

static void normalize_coordinates(int frame_left, frame_right, frame_up, frame_down) {

    int x, y, width, tile_width, margin;
    // first normalize the size and position of the window
    // by removing the frame extents
    attrs.x += frame_left;
    attrs.y += frame_up;

    attrs.width -= (frame_left + frame_right);
    attrs.height -= (frame_up + frame_down);

    // approximate the coordinates of the grid
    // it should be roughly a square, so we strip the top bar
    width = (attrs.width < attrs.height) ? attrs.width : attrs.heigth;
    // heigth == width, no need to store it

    // fix x,y as the bottom-rigth pixel
    x = attrs.x + attrs.width;
    y = attrs.y + attrs.height;

    // a gross approximation of the margins between tiles
    margin = width / 80;
    // strip the outer margins
    width -= 2 * margin;
    // now we have an estimate of tile width/height
    tile_width = width / N;

    // finally, fix coordinates of the NxN tiles
    int (*coords)[2] = malloc(sizeof(int[N*N][2]));

    // ensure that when scanning pixel colors, we never hit the digit
    // (in the center) or the margin/border between tiles.
    offset = 3 * margin;

    int i, j, k, m;
    // here I was tired, so didn't come up with a more efficient iteration
    m = 0;
    for (i=N*N-1; i>0; i-=N) {
        k = 0;
        for (j=i; j>i-N; j--) {
            coords[j][0] = x - ( k * tile_width ) - offset;
            coords[j][1] = y - ( m * tile_width ) - offset;
            k++;
        }
        m++;
    }
}
