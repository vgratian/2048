#include "xlib.h"

Display *dpy;
Window root, target;

static struct {
    int x;
    int y;
    int width;
    int height;
    int display_w;
    int display_h;
} attrs; // position of target window: x, y, width and height
static int **coords; // coordinates of the tiles in the grid
static int N; // size of the grid

// Open connection to the X server and fix the root window
// Returns 0 on success.
int open_x_socket() {
    dpy = XOpenDisplay(NULL);
    //if ( (dpy = XOpenDisplay(NULL) != NULL ) ) {
    if ( dpy == NULL ) {
        return 1;
    }
    printf("connected to the X server\n");
    root = XDefaultRootWindow(dpy);
    int screen_num = XDefaultScreen(dpy);
    attrs.display_w = XDisplayWidth(dpy, screen_num);
    attrs.display_h = XDisplayHeight(dpy, screen_num);
    printf("fixed screen size [%d x %d]\n", attrs.display_w, attrs.display_h);
    return 0;
}

// Close connection to the X server
void close_x_socket() {
    //XFree(attrs);
    XCloseDisplay(dpy);
}

int find_something_else(char* name, int size) {
    printf("Okay okay ... [%s] (%d)\n", name, size);

    return find_x_window(name, size);
}

// Find window and fix its XID and position. We expect that this window contains
// size x size tiles and we will try to approximate their coordinates.
int find_x_window(const char* name, int size) {

    N = size;

    // find XID of the window
    if ( (target = find_window_recursively(root, name)) == 0 ) {
        printf("recursively not found\n");
        return 1;
    }

    printf("found target window [0x%x]\n", target);

    // find position of the window, populates static array *pos[4]*
    // with values for x, y (top-left pixel), width, height
    if (get_coordinates() != 0) {
        printf("coordinates not found\n");
        return 2;
    }

    /*
    if (attrs == NULL) {
        printf("fuck, null pointer!\n");
        return 3;
    }
    */

    printf("x=%d y=%d width=%d height=%d\n", attrs.x, attrs.y, attrs.width, attrs.height);

    // the width and height values are sometimes larger than the actual
    // window, so try to find the value for GTK frame extents to normalize
    // them. If we can't get frames, we just assume they don't exist.
    int frames[4];
    printf("get_gtk_frame_extents\n");
    get_gtk_frame_extents(&frames[0], &frames[1], &frames[2], &frames[3]);
    printf("frames = %d, %d, %d, %d\n", frames[0], frames[1], frames[2], frames[3]);

    printf("normalize_coordinates\n");
    normalize_coordinates(frames[0], frames[1], frames[2], frames[3]);

    printf("attrs: x=%d, y=%d, width=%d, height=%d\n", attrs.x, attrs.y, attrs.width, attrs.height);

    for (int i=0; i<N*N; i++) {
        printf(" %d - (%d, %d)\n", i, coords[i][0], coords[i][1]);
    }
    printf("\n\n");
    printf("find_x_window: ALL DONE\n");
    return 0;
}

// Get the pixel values of the tiles in the grid
unsigned long* scan_pixels() {

    printf("scan_pixels() :: start\n");

    //unsigned long* pixels;
    int i;
    XImage *image;

    printf("requesting image for x=%d y=%d, width=%d height=%d\n", attrs.x, attrs.y, attrs.width, attrs.height);

    //image = XGetImage(dpy, root, attrs.x, attrs.y, attrs.width, attrs.height, AllPlanes, XYPixmap);
    image = XGetImage(dpy, root, 0, 0, attrs.display_w, attrs.display_h, AllPlanes, XYPixmap);
    if (image == NULL) {
        return NULL;
    }

    int size = N*N;
    //pixels = (unsigned long*) malloc(N*N);
    //int (*coords)[2] = malloc(sizeof(int[N*N][2]));
    unsigned long *pixels;
    unsigned long value;
    pixels = (unsigned long *) malloc(N*N*sizeof(unsigned long));

    /*
    printf("created array:\n");
    for (i=0; i<size; i++)
        printf("%3d: < %d, %d > => %-10lu\n", i, coords[i][0], coords[i][1], pixels[i]);
    */


    printf("checking %d pixels (%d x %d)....\n", size, N, N);

    for (i=0; i<size; i++) {
        printf(" %d => ", i);
        printf("(%d, %d)", coords[i][0], coords[i][1]);
        //value = XGetPixel(image, coords[0][0]+1000, coords[0][1]+1000);
        value = XGetPixel(image, coords[i][0], coords[i][1]); // x, y coordinate of the grid
        printf(" => %lu", value);
        pixels[i] = value;
        printf(" => ok\n");
    }

    printf("all done:\n");

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

    /*
    if (name != NULL) {
        printf(" -> [%s]\n", name);
    }
    */

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
    
    return None;
}

static int get_coordinates() {
    printf("fetching window [0x%x] attributes....\n", target);
    XWindowAttributes a;
    int result;

    result = XGetWindowAttributes(dpy, target, &a);

    printf("x=%d, y=%d, width=%d, height=%d\n", a.x, a.y, a.width, a.height);

    if (result == 1) {
        printf("%d - OK!\n", result);
        attrs.x = a.x;
        attrs.y = a.y;
        attrs.width = a.width;
        attrs.height = a.height;
        return 0;
    }
    printf("%d - FAIL\n", result);
    return 1;
}

static void normalize_coordinates(int frame_left, int frame_right, int frame_up, int frame_down) {

    printf("normalize_coordinates() :: start\n");

    int x, y, width, tile_width, margin;
    // first normalize the size and position of the window
    // by removing the frame extents
    attrs.x += frame_left;
    attrs.y += frame_up;

    attrs.width -= (frame_left + frame_right);
    attrs.height -= (frame_up + frame_down);

    printf("attrs: x=%d, y=%d, width=%d, height=%d\n", attrs.x, attrs.y, attrs.width, attrs.height);

    // approximate the coordinates of the grid
    // it should be roughly a square, so we strip the top bar
    width = (attrs.width < attrs.height) ? attrs.width : attrs.height;
    // height == width, no need to store it

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
    //static int (*coords)[2] = malloc(sizeof(int[N*N][2]));
    //**coords = malloc(sizeof(int[N*N][2]));
    coords = malloc(N * N * sizeof(int *));

    // ensure that when scanning pixel colors, we never hit the digit
    // (in the center) or the margin/border between tiles.
    int offset = 3 * margin;

    int k, m;
    // here I was tired, so didn't come up with a more efficient iteration
    printf("generating pixel coordinates...\n");

    m = 0;
    for (int i=N*N-1; i>0; i-=N) {
        printf(" i=2%d\n", i);
        k = 0;
        for (int j=i; j>i-N; j--) {
            coords[j] = malloc(2 * sizeof(int));
            printf("   i=%2d j=%2d", i, j);
            coords[j][0] = x - ( k * tile_width ) - offset;
            coords[j][1] = y - ( m * tile_width ) - offset;
            printf(" ==> (%3d %3d)\n", coords[j][0], coords[j][1]);
            k++;
        }
        m++;
    }
}

/* welcome to hell */
static long extract_value(const char **ptr, int *length) {
    long value;
    value = * (const unsigned short *) *ptr & 0xffffffff;
    *ptr += sizeof(long);
    *length -= sizeof(long);
    return value;
}

static int get_gtk_frame_extents(int *left, int *right, int *up, int *down) {

    //const char* buffer;
    const char* buffer;
    unsigned char *data;
    //long length;
    int size, result;
    int frames[4];
    unsigned long nitems, bytes_after;
    Atom xa_property, xa_type;

    // what to do with unused variables?
    // - xa_type
    // - nitems, bytes_after

    if ((xa_property = XInternAtom(dpy, "_GTK_FRAME_EXTENTS", False)) == None) {
        return 1;
    }

    result = XGetWindowProperty(dpy, target, xa_property, 0, 1024, False,
            XA_CARDINAL, &xa_type, &size, &nitems, &bytes_after, &data);

    if (result != Success) {
        return 1;
    }

    // expecting 4 longs
    if (size != 32 ) {
        return 1;
    }

    buffer = (const char*) data;
    int i = 0;
    int length = (int) size;
    long value;

    while (length >= size/8) {
        //value = * (const unsigned short *) buffer & 0xffffffff;
        //*buffer += sizeof(long);
        //frames[i] = (int) value;
        frames[i] = (int) extract_value(&buffer, &length);
        i++;
    }

    *left = frames[0];
    *right = frames[1];
    *up = frames[2];
    *down = frames[3];

    return 0;
}
