package Window

/*
CharacterWindow.go - SDL3 transparent borderless window for displaying animated characters

SDL requires its event loop to run on a locked OS thread for proper event handling on Windows.
This implementation uses runtime.LockOSThread() and channels for thread-safe communication.

Functions:
- NewCharacterWindow: Create new character window instance
- (CharacterWindow) Start: Launch window in dedicated OS thread
- (CharacterWindow) Close: Signal window to close via channel
- (CharacterWindow) SetScale: Thread-safe scale adjustment via channel
- (CharacterWindow) IsRunning: Check if window is still active
- (CharacterWindow) GetID: Get unique window identifier
*/

import (
	"boccho-ui/AnimationEngine"
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type CharacterWindow struct {
	id            string
	characterName string
	framesPath    string
	running       atomic.Bool
	closeChan     chan struct{}
	doneChan      chan struct{}
	scaleChan     chan float64
	currentScale  atomic.Value
}

func NewCharacterWindow(id, characterName, framesPath string) *CharacterWindow {
	cw := &CharacterWindow{
		id:            id,
		characterName: characterName,
		framesPath:    framesPath,
		closeChan:     make(chan struct{}),
		doneChan:      make(chan struct{}),
		scaleChan:     make(chan float64, 10),
	}
	cw.currentScale.Store(0.51)
	return cw
}

func (cw *CharacterWindow) Start() {
	go cw.runInOSThread()
}

func (cw *CharacterWindow) runInOSThread() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(cw.doneChan)

	cw.running.Store(true)
	defer cw.running.Store(false)

	title := fmt.Sprintf("Boccho - %s", cw.characterName)
	winW, winH := int32(400), int32(400)

	flags := sdl.WindowTransparent | sdl.WindowAlwaysOnTop | sdl.WindowBorderless

	window := sdl.CreateWindow(title, winW, winH, flags)
	if window == nil {
		fmt.Printf("[%s] Failed to create window: %s\n", cw.id, sdl.GetError())
		return
	}
	defer sdl.DestroyWindow(window)

	if !sdl.SetWindowHitTest(window, hitTestCallback, nil) {
		fmt.Printf("[%s] Warning: Could not set hit test callback: %s\n", cw.id, sdl.GetError())
	}

	renderer := sdl.CreateRenderer(window, "")
	if renderer == nil {
		fmt.Printf("[%s] Failed to create renderer: %s\n", cw.id, sdl.GetError())
		return
	}
	defer sdl.DestroyRenderer(renderer)

	sdl.SetRenderDrawBlendMode(renderer, sdl.BlendModeBlend)

	animation := AnimationEngine.NewAnimationPlayer(cw.framesPath)
	if err := animation.LoadFrames(renderer); err != nil {
		fmt.Printf("[%s] Failed to load frames: %v\n", cw.id, err)
		return
	}
	defer animation.Cleanup()

	// Register render context for event watcher (fixes animation freeze during window drag on Windows)
	windowID := sdl.GetWindowID(window)
	RegisterRenderContext(windowID, &RenderContext{
		Renderer:  renderer,
		Window:    window,
		Animation: animation,
		WindowID:  windowID,
	})
	defer UnregisterRenderContext(windowID)

	fmt.Printf("[%s] Character window started\n", cw.id)
	fmt.Println("  Controls: Arrow Up/Down = Scale, Escape = Close")

	var event sdl.Event
	for {
		select {
		case <-cw.closeChan:
			fmt.Printf("[%s] Received close signal\n", cw.id)
			return
		case newScale := <-cw.scaleChan:
			animation.SetScale(newScale)
			cw.currentScale.Store(newScale)
			fmt.Printf("[%s] Scale set to: %.2f\n", cw.id, newScale)
		default:
		}

		for sdl.PollEvent(&event) {
			eventType := event.Type()

			switch eventType {
			case sdl.EventQuit:
				return
			case sdl.EventKeyDown:
				key := event.Key().Key
				if key == sdl.KeycodeEscape {
					return
				} else if key == sdl.KeycodeUp {
					animation.ScaleUp()
					cw.currentScale.Store(animation.GetScale())
				} else if key == sdl.KeycodeDown {
					animation.ScaleDown()
					cw.currentScale.Store(animation.GetScale())
				}
			}
		}

		animation.Update()

		sdl.SetRenderDrawColor(renderer, 0, 0, 0, 0)
		sdl.RenderClear(renderer)
		animation.Render(renderer, window)
		sdl.RenderPresent(renderer)

		sdl.DelayNS(16 * 1000000)
	}
}

func (cw *CharacterWindow) Close() {
	select {
	case <-cw.closeChan:
	default:
		close(cw.closeChan)
	}
}

func (cw *CharacterWindow) SetScale(scale float64) {
	select {
	case cw.scaleChan <- scale:
	default:
	}
}

func (cw *CharacterWindow) GetScale() float64 {
	if v := cw.currentScale.Load(); v != nil {
		return v.(float64)
	}
	return 0.51
}

func (cw *CharacterWindow) Wait() {
	<-cw.doneChan
}

func (cw *CharacterWindow) IsRunning() bool {
	return cw.running.Load()
}

func (cw *CharacterWindow) GetID() string {
	return cw.id
}

func (cw *CharacterWindow) GetCharacterName() string {
	return cw.characterName
}
