package Window

/*
utils.go - Shared utility functions for Window package

Functions:
- hitTestCallback: SDL hit test callback for making windows draggable
*/

import (
	"unsafe"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

func hitTestCallback(window *sdl.Window, point *sdl.Point, data unsafe.Pointer) sdl.HitTestResult {
	return sdl.HitTestDraggable
}
