package Window

/*
eventwatch.go - Event watcher for window move events (Windows modal loop fix)

On Windows, when a user drags a window, the OS enters a modal loop that blocks
the main SDL event loop, causing animations to freeze. This module implements
an event watcher that continues rendering during window move events.

Functions:
- createEventWatcher: Create event filter callback for window moved events
- windowMoveEventCallback: Callback that triggers render during window drag
*/

import (
	"boccho-ui/AnimationEngine"
	"unsafe"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

// RenderContext holds all resources needed for rendering during window move
type RenderContext struct {
	Renderer  *sdl.Renderer
	Window    *sdl.Window
	Animation *AnimationEngine.AnimationPlayer
	WindowID  sdl.WindowID
}

// Global map to store render contexts for event callback
// Key is SDL window ID since we might have multiple character windows
var renderContexts = make(map[sdl.WindowID]*RenderContext)

// RegisterRenderContext registers a render context for event watching
func RegisterRenderContext(windowID sdl.WindowID, ctx *RenderContext) {
	renderContexts[windowID] = ctx
}

// UnregisterRenderContext removes a render context when window closes
func UnregisterRenderContext(windowID sdl.WindowID) {
	delete(renderContexts, windowID)
}

// CreateEventWatcher creates and registers the event watcher for window move events
func CreateEventWatcher() sdl.EventFilter {
	return sdl.NewEventFilter(func(userdata unsafe.Pointer, event *sdl.Event) bool {
		eventType := event.Type()

		// Handle window events that indicate window is being moved
		if eventType == sdl.EventWindowMoved ||
			eventType == sdl.EventWindowExposed ||
			eventType == sdl.EventWindowResized {

			windowEvent := event.Window()
			windowID := windowEvent.WindowID

			if ctx, exists := renderContexts[windowID]; exists {
				// Update animation timing and render
				ctx.Animation.Update()

				sdl.SetRenderDrawColor(ctx.Renderer, 0, 0, 0, 0)
				sdl.RenderClear(ctx.Renderer)
				ctx.Animation.Render(ctx.Renderer, ctx.Window)
				sdl.RenderPresent(ctx.Renderer)
			}
		}

		// Return true to allow event to be processed normally
		return true
	})
}
