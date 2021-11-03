package gui

/*
#cgo CFLAGS: -I${SRCDIR}/gui
#include <stdlib.h>
#include "gui.h"
*/
import "C"

import (
    "errors"
    "unsafe"
)

var (
    pseudoColors = make(map[int64]uint8)
    nextValue uint8 = 0
    boardSize int
)

func FindWindow(name string, size int) error {
    boardSize = size
    ptr := C.CString(name)
    defer C.free(unsafe.Pointer(ptr))

    if C.find_window(ptr) == 0 {
        return nil
    }
    return errors.New("window [" + name + "] not found"
}

func InitBoard() ([]uint8, error) {
    // only two tiles should contain "2", everything else is white/"0"
    // so we just need to fix the pixel value for white

    var (
        aPixel, bPixel uint64
        aCount, bCount int
    )

    size, pixels := scanPixels()

    if size != boardSize {
        return nil, errors.New(fmt.Sprintf("failed to scan %d pixels (%d)", boardSize, size))
    }


    for _, pixel := range pixels {

        if pixel == aPixel {
            aCount++
        } else if pixel == bPixel {
            bPixel++
        }

        if aCount == 0 {
            aPixel = pixel
            aCount = 1
        } else if bCount == 0 {
            bPixel = pixel
            bCount = 1
        }

        if aCount == 3 {
            pseudoColors[aPixel] = 0
            nextValue = 1
            break
        }

        if bCount == 3 {
            pseudoColors[bPixel] = 0
            nextValue = 1
            break
        }
    }

    return ScanBoard()

}

func ScanBoard() ([]uint8, error) {

    size, pixels := scanPixels()

    if size != boardSize {
        return nil, errors.New(fmt.Sprintf("failed to scan %d pixels (%d)", boardSize, size))
    }

    slice := make([]uint8, boardSize)
 
    for i, pixel := pixels {
        if value, has := pseuoColors[pixel] {
            slice[i] = value
        } else {
            pseudoColors[pixel] = nextValue
            slice[i] = nextValue
            nextValue++
        }
    }
}

func scanPixels() (int, []int64) {
    size := C.int(0)
    arr := C.scan_pixels(&size)

    slice := make([]int64, size)

    // trick to convert pointer into an array
    // https://stackoverflow.com/questions/28925179/cgo-how-to-pass-struct-array-from-c-to-go
    for i, pixel := (*[1<<30]C.long(unsafe.Pointer(arr))[:size:size] {
        slice[i] = int64(pixel)
    }
    return size, slice
}

func EmulateKeyPress(key string) {
    ptr := C.CString(key)
    defer C.free(unsafe.Pointer(ptr))
    C.fake_key_event(ptr)
}

