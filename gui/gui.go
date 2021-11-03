package src

/*
#cgo CFLAGS: -I${SRCDIR}/gui
#include <stdlib.h>
#include "gui.h"
*/
import "C"

import (
    "unsafe"
)
package gui

func StartGame(bin string, size int, ...args string) error {
}

func FindWindow(name string, size int) error {
}

func ScanBoard() ([]uint8, error) {
}

func EmulateKeyPress(key string) {
}

