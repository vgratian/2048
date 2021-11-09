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
int find_window(const char*, int);
unsigned long* scan_pixels();
int fake_key_event(const char*, unsigned long);

// private functions
static Window _find_window(Display*, Window, const char*);
static int get_coordinates(Display*, Window);
static int get_gtk_frame_extents(Display*, Window, int[4]);
static long extract_value(const char**, int*);
static void normalize_coordinates(int[4]);

