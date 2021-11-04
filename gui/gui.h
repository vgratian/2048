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
static unsigned long _find_window(unsigned long, const char*);
static void _fix_coordinates();

