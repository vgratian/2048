#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>

// public functions
int open();
void close();
int find_wondow(const char*, int);
unsigned long* scan_pixels();
void fake_key_event(const char*);

// private functions
static Window find_window_recursively(Window, const char*);
static int get_coordinates();
static int *get_gtk_frame_extents();
static void normalize_coordinates(int, int, int, int);

