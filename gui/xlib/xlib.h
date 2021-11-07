#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>
#include <X11/keysym.h>
#include <X11/extensions/XTest.h>

// public functions
int open_x_socket();
void close_x_socket();
int find_x_window(const char*, int);
unsigned long* scan_pixels();
void fake_key_event(const char*);
int find_something_else(char*, int);

// private functions
static Window find_window_recursively(Window, const char*);
static int get_coordinates();
static int get_gtk_frame_extents(int*, int*, int*, int*);
static void normalize_coordinates(int, int, int, int);

