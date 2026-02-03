package AnimationEngine

/*
Animation.go - Core animation player for character sprites

Functions:
- NewAnimationPlayer: Create new animation player instance
- (AnimationPlayer) LoadFrames: Load PNG frames from directory into textures
- (AnimationPlayer) Update: Advance animation frame based on timing
- (AnimationPlayer) Render: Render current frame to renderer
- (AnimationPlayer) SetScale: Adjust character scale
- (AnimationPlayer) GetScaledSize: Get current scaled dimensions
- (AnimationPlayer) Cleanup: Destroy all textures and free resources
*/

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/jupiterrider/purego-sdl3/img"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

const (
	DefaultFrameDelay = 83 // ~12fps
)

type AnimationPlayer struct {
	textures      []*sdl.Texture
	originalSizes []sdl.Point
	currentFrame  int
	scale         float64
	frameDelay    uint64
	lastFrameTime uint64
	framesPath    string
}

func NewAnimationPlayer(framesPath string) *AnimationPlayer {
	return &AnimationPlayer{
		textures:      make([]*sdl.Texture, 0),
		originalSizes: make([]sdl.Point, 0),
		currentFrame:  0,
		scale:         0.51,
		frameDelay:    DefaultFrameDelay,
		lastFrameTime: 0,
		framesPath:    framesPath,
	}
}

func (ap *AnimationPlayer) LoadFrames(renderer *sdl.Renderer) error {
	imageFiles, err := filepath.Glob(filepath.Join(ap.framesPath, "*.png"))
	if err != nil {
		return fmt.Errorf("error finding images: %w", err)
	}

	if len(imageFiles) == 0 {
		return fmt.Errorf("no PNG images found in %s", ap.framesPath)
	}

	sort.Strings(imageFiles)

	for _, file := range imageFiles {
		surface := img.Load(file)
		if surface == nil {
			fmt.Printf("Failed to load %s: %s\n", file, sdl.GetError())
			continue
		}

		width := int32(surface.W)
		height := int32(surface.H)
		ap.originalSizes = append(ap.originalSizes, sdl.Point{X: width, Y: height})

		texture := sdl.CreateTextureFromSurface(renderer, surface)
		sdl.DestroySurface(surface)

		if texture != nil {
			sdl.SetTextureBlendMode(texture, sdl.BlendModeBlend)
			ap.textures = append(ap.textures, texture)
			fmt.Printf("Loaded: %s (%dx%d)\n", filepath.Base(file), width, height)
		} else {
			fmt.Printf("Failed to create texture for %s: %s\n", filepath.Base(file), sdl.GetError())
			ap.originalSizes = ap.originalSizes[:len(ap.originalSizes)-1]
		}
	}

	if len(ap.textures) == 0 {
		return fmt.Errorf("failed to load any textures")
	}

	fmt.Printf("Total frames loaded: %d\n", len(ap.textures))
	return nil
}

func (ap *AnimationPlayer) Update() {
	if len(ap.textures) == 0 {
		return
	}

	currentTime := sdl.GetTicks()
	if currentTime-ap.lastFrameTime >= ap.frameDelay {
		ap.currentFrame = (ap.currentFrame + 1) % len(ap.textures)
		ap.lastFrameTime = currentTime
	}
}

func (ap *AnimationPlayer) Render(renderer *sdl.Renderer, window *sdl.Window) {
	if len(ap.textures) == 0 {
		return
	}

	texture := ap.textures[ap.currentFrame]
	orig := ap.originalSizes[ap.currentFrame]

	scaledW := float32(float64(orig.X) * ap.scale)
	scaledH := float32(float64(orig.Y) * ap.scale)

	dst := sdl.FRect{X: 0, Y: 0, W: scaledW, H: scaledH}

	sdl.SetWindowSize(window, int32(scaledW), int32(scaledH))
	sdl.RenderTexture(renderer, texture, nil, &dst)
}

func (ap *AnimationPlayer) SetScale(scale float64) {
	if scale < 0.1 {
		scale = 0.1
	}
	ap.scale = scale
}

func (ap *AnimationPlayer) GetScale() float64 {
	return ap.scale
}

func (ap *AnimationPlayer) ScaleUp() {
	ap.scale *= 1.1
	fmt.Printf("Scale: %.2f\n", ap.scale)
}

func (ap *AnimationPlayer) ScaleDown() {
	ap.scale = max(0.1, ap.scale/1.1)
	fmt.Printf("Scale: %.2f\n", ap.scale)
}

func (ap *AnimationPlayer) GetScaledSize() (int32, int32) {
	if len(ap.originalSizes) == 0 {
		return 0, 0
	}
	orig := ap.originalSizes[ap.currentFrame]
	return int32(float64(orig.X) * ap.scale), int32(float64(orig.Y) * ap.scale)
}

func (ap *AnimationPlayer) SetFrameDelay(delay uint64) {
	ap.frameDelay = delay
}

func (ap *AnimationPlayer) FrameCount() int {
	return len(ap.textures)
}

func (ap *AnimationPlayer) Cleanup() {
	for _, t := range ap.textures {
		sdl.DestroyTexture(t)
	}
	ap.textures = nil
	ap.originalSizes = nil
	fmt.Println("Animation resources cleaned up")
}
