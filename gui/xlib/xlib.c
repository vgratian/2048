/* 
 * This script uses methods from XLIB (libx11) to find and interact with
 * an X window. Most of xlib functions are pretty straigtforward, except
 * when you need to work with EWMH properties (e.g. GTK_FRAME_EXTENTS).
 * I looked at the source codes of "wmctrl" and "xprop" for help.
 *
 * Most of the public methods return 0 on success and a positive integer
 * on failure (1 usually means failed to connect to X).
 *
 * TODO: implement verbose/debug mode
 *   for now, you can remove all "// DEBUG: " comments
 *   to get see what's happening.
*/

#include "xlib.h"

// position of target window: x, y, width and height
static struct {
    int x;
    int y;
    int w;
    int h;
} board;

// screen size
static struct {
    unsigned int w;
    unsigned int h;
} screen;

// coordinates of the tiles in the grid
static int **tiles;

// size of the grid
static int N;

// Find window and fix its XID and position. We expect that this window contains
// size x size tiles and we will try to approximate their coordinates.
//
// This method should be called once before all the others.
int find_window(const char* name, int size) {

    // fix expected board size
    N = size;

    // connect to X and fetch some basic info
    Display *dpy;
    Window root, target;
    int screen_num;

    if ( (dpy = XOpenDisplay(NULL)) == NULL ) {
        // DEBUG: printf("failed to connect to X server\n");
        return 1;
    }

    // DEBUG: printf("connected to the X server\n");

    root = XDefaultRootWindow(dpy);

    screen_num = XDefaultScreen(dpy);
    screen.w = XDisplayWidth(dpy, screen_num);
    screen.h = XDisplayHeight(dpy, screen_num);

    // DEBUG: printf("fixed screen size [%d x %d]\n", screen.w, screen.h);

    // find XID of the target window
    if ( (target = _find_window(dpy, root, name)) == 0 ) {
        // DEBUG: printf("target window [%s] not found\n", name);
        return 2;
    }

    // DEBUG: printf("found target window with XID [0x%x]\n", target);

    // fix position of the target window
    if (get_coordinates(dpy, target) != 0) {
        // DEBUG: printf("coordinates not found\n");
        return 3;
    }

    // DEBUG: printf("initial window position: x=%d y=%d (top-left pixel); w=%d h=%d\n",
    // DEBUG:         board.x, board.y, board.w, board.h);

    // the width and height values are sometimes larger than the actual
    // window, so try to find the value for GTK frame extents to normalize
    // them. If we can't get frames, we just assume they don't exist.
    int frames[4];
    frames[0] = frames[1] = frames[2] = frames[3] = 0;
    get_gtk_frame_extents(dpy, target, frames);

    // DEBUG: printf("GTK frames: left=%d, right=%d, up=%d, down=%d\n",
    // DEBUG:        frames[0], frames[1], frames[2], frames[3]);

    normalize_coordinates(frames);

    // DEBUG:  printf("normalized position: x=%d, y=%d (bottm-right pixel); w=%d, h=%d\n",
    // DEBUG:        board.x, board.y, board.w, board.h);

    // DEBUG: printf("coordinates of %dx%d tiles:\n", N, N);
    // DEBUG: for (int i=0; i<N*N; i++) {
    // DEBUG:     printf(" (%d, %d)", tiles[i][0], tiles[i][1]);
    // DEBUG:     if ( (i+1) % N == 0 ) {
    // DEBUG:         printf("\n");
    // DEBUG:     }
    // DEBUG: }
    // DEBUG: printf("\n");

    XCloseDisplay(dpy);
    return 0;
}

// Get the pixel values of the tiles in the grid
unsigned long* scan_pixels() {

    Display *dpy;
    if ( (dpy = XOpenDisplay(NULL)) == NULL ) {
        // DEBUG: printf("failed to connect to X server\n");
        return NULL;
    }

    // DEBUG: printf("requesting image for x=%d y=%d, w=%d h=%d\n",
    // DEBUG:        board.x, board.y, board.w, board.h);

    XImage *image = XGetImage(dpy, XDefaultRootWindow(dpy), 0, 0,
            screen.w, screen.h, AllPlanes, ZPixmap);
    if (image == NULL) {
        // DEBUG: printf("failed to get image\n");
        return NULL;
    }

    int i;
    int size = N*N;
    unsigned long *pixels;

    pixels = (unsigned long *) malloc(N*N*sizeof(unsigned long));

    // DEBUG: printf("scanning %d pixels (%d x %d)....\n", size, N, N);

    for (i=0; i<size; i++) {
        // scan pixels at the (x,y) coordinate of each tile
        pixels[i] = XGetPixel(image, tiles[i][0], tiles[i][1]);
        // DEBUG: printf(" %2d (%3d, %3d): %lu\n", i, tiles[i][0], tiles[i][1], pixels[i]);
    }

    XCloseDisplay(dpy);
    XDestroyImage(image);

    return pixels;
}

int fake_key_event(const char* key, unsigned long delay) {
    // delay is in milliseconds

    Display *dpy;
    KeyCode kc;
    KeySym ks;

    if ( (dpy = XOpenDisplay(NULL)) == NULL ) {
        return 1;
    }

    if ( (ks = XStringToKeysym(key)) == NoSymbol ) {
        return 2;
    }

    if ( (kc = XKeysymToKeycode(dpy, ks)) == 0 ) {
        return 2;
    }

    // key is released
    XTestFakeKeyEvent(dpy, kc, False, delay);
    XFlush(dpy);

    // key is pressed
    XTestFakeKeyEvent(dpy, kc, True, delay);
    XFlush(dpy);

    // key is released
    XTestFakeKeyEvent(dpy, kc, False, delay);
    XFlush(dpy);

    XCloseDisplay(dpy);

    return 0;
}

static Window _find_window(Display *dpy, Window win, const char* name) {
    char *window_name;
    XFetchName(dpy, win, &window_name);

    if (window_name != NULL && strcmp(window_name, name) == 0) {
        return win;
    }

    Window root_return, parent_return, child_match;
    Window *children;
    int i, n;

    XQueryTree(dpy, win, &root_return, &parent_return, &children, &n);

    if (children != NULL) {
        for (i=0; i<n; i++) {
            child_match = _find_window(dpy, children[i], name);
            if ( child_match ) {
                break;
            }
        }
        XFree(children);
        if (i < n) {
            return child_match;
        }
    }
    
    return 0;
}

static int get_coordinates(Display *dpy, Window target) {
    XWindowAttributes a;

    if ( XGetWindowAttributes(dpy, target, &a ) != 1 ) {
        return 1;
    }

    board.x = a.x;
    board.y = a.y;
    board.w = a.width;
    board.h = a.height;

    return 0;
}

static void normalize_coordinates(int frames[4] /* left right up down */ ) {

    int x, y, width, tile_width, margin;
    // first normalize the size and position of the window
    // by removing the frame extents
    board.x += frames[0];
    board.y += frames[2];

    board.w -= (frames[0] + frames[1]);
    board.h -= (frames[2] + frames[3]);

    // DEBUG: printf("attrs: x=%d, y=%d, width=%d, height=%d\n",
    // DEBUG:        board.x, board.y, board.w, board.h);

    // approximate the coordinates of the grid
    // it should be roughly a square, so we strip the top bar
    width = (board.w < board.h) ? board.w : board.h;
    // height == width, no need to store it

    // fix x,y as the bottom-rigth pixel
    x = board.x + board.w;
    y = board.y + board.h;

    // a gross approximation of the margins between tiles
    margin = width / 80;
    // strip the outer margins
    width -= 2 * margin;
    // now we have an estimate of tile width/height
    tile_width = width / N;

    // finally, fix coordinates of the NxN tiles
    tiles = malloc(N * N * sizeof(int *));

    // ensure that when scanning pixel colors, we never hit the digit
    // (in the center) or the margin/border between tiles.
    int offset = 3 * margin;

    int k, m;
    // here I was tired, so didn't come up with a more efficient iteration

    m = 0;
    for (int i=N*N-1; i>0; i-=N) {
        // DEBUG: printf(" i=%2d\n", i);
        k = 0;
        for (int j=i; j>i-N; j--) {
            tiles[j] = malloc(2 * sizeof(int));
            // DEBUG: printf("   i=%2d j=%2d", i, j);
            tiles[j][0] = x - ( k * tile_width ) - offset;
            tiles[j][1] = y - ( m * tile_width ) - offset;
            // DEBUG: printf(" ==> (%3d %3d)\n", tiles[j][0], tiles[j][1]);
            k++;
        }
        m++;
    }
}

// This function and extract_value() are partially copied from:
// xprop 1.2.5: Get_Window_Property_Data_And_Type()
// wmctrl 1.07: get_property()
static int get_gtk_frame_extents(Display *dpy, Window target, int frames[4]) {

    const char* buffer;
    unsigned char *data;
    int size;
    unsigned long nitems, bytes_after;
    Atom xa_property, xa_type;

    // what to do with unused variables?
    // - xa_type
    // - nitems, bytes_after

    if ((xa_property = XInternAtom(dpy, "_GTK_FRAME_EXTENTS", False)) == None) {
        return 1;
    }

    if ( XGetWindowProperty(dpy, target, xa_property, 0, 1024,
            False, XA_CARDINAL, &xa_type, &size, &nitems,
            &bytes_after, &data) != Success) {
        return 2;
    }

    // expecting 4 longs
    if (size != 32 ) {
        return 3;
    }

    buffer = (const char*) data;
    int length = (int) size;
    int i = 0;

    while (length >= size/8) {
        frames[i] = (int) extract_value(&buffer, &length);
        i++;
    }

    return 0;
}

static long extract_value(const char **ptr, int *length) {
    long value;
    value = * (const unsigned short *) *ptr & 0xffffffff;
    *ptr += sizeof(long);
    *length -= sizeof(long);
    return value;
}
