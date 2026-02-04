package main

/*
app.go - Wails application bindings for character management

Exposes to frontend:
- GetCharacters: List available characters from Frames directory
- SpawnCharacter: Create new SDL character window in separate OS thread
- DestroyCharacter: Close specific character window
- GetActiveWindows: List currently spawned windows
- SetCharacterScale: Adjust scale of specific window
*/

import (
	"boccho-ui/AnimationEngine"
	"boccho-ui/PackManagement"
	"boccho-ui/Window"
	"boccho-ui/config"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx           context.Context
	activeWindows map[string]*Window.CharacterWindow
	mu            sync.RWMutex
	framesPath    string
	cfg           config.Config
}

type CharacterWindowInfo struct {
	ID            string  `json:"id"`
	CharacterName string  `json:"characterName"`
	IsRunning     bool    `json:"isRunning"`
	Scale         float64 `json:"scale"`
}

func NewApp() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v, using defaults\n", err)
		cfg = config.GetDefaultConfig()
	}

	return &App{
		activeWindows: make(map[string]*Window.CharacterWindow),
		framesPath:    cfg.FramesPath,
		cfg:           cfg,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if err := config.EnsureFramesDir(a.cfg); err != nil {
		fmt.Printf("Error ensuring Frames directory: %v\n", err)
	}

	fmt.Printf("Frames path: %s\n", a.framesPath)

	go a.cleanupDeadWindows()
}

func (a *App) cleanupDeadWindows() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.mu.Lock()
			for id, cw := range a.activeWindows {
				if !cw.IsRunning() {
					delete(a.activeWindows, id)
					fmt.Printf("Cleaned up window: %s\n", id)
				}
			}
			a.mu.Unlock()
		}
	}
}

func (a *App) GetCharacters() []AnimationEngine.CharacterInfo {
	characters, err := AnimationEngine.ScanCharacters(a.framesPath)
	if err != nil {
		fmt.Printf("Error scanning characters: %v\n", err)
		return []AnimationEngine.CharacterInfo{}
	}
	return characters
}

func (a *App) SpawnCharacter(characterName string) CharacterWindowInfo {
	charPath := AnimationEngine.GetCharacterFramesPath(a.framesPath, characterName)

	if _, err := os.Stat(charPath); os.IsNotExist(err) {
		fmt.Printf("Character path not found: %s\n", charPath)
		return CharacterWindowInfo{}
	}

	id := uuid.New().String()[:8]

	charWindow := Window.NewCharacterWindow(id, characterName, charPath)

	a.mu.Lock()
	a.activeWindows[id] = charWindow
	a.mu.Unlock()

	charWindow.Start()

	time.Sleep(50 * time.Millisecond)

	return CharacterWindowInfo{
		ID:            id,
		CharacterName: characterName,
		IsRunning:     charWindow.IsRunning(),
		Scale:         charWindow.GetScale(),
	}
}

func (a *App) DestroyCharacter(windowId string) bool {
	a.mu.RLock()
	charWindow, exists := a.activeWindows[windowId]
	a.mu.RUnlock()

	if !exists {
		return false
	}

	charWindow.Close()

	a.mu.Lock()
	delete(a.activeWindows, windowId)
	a.mu.Unlock()

	return true
}

func (a *App) GetActiveWindows() []CharacterWindowInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()

	windows := make([]CharacterWindowInfo, 0, len(a.activeWindows))
	for id, cw := range a.activeWindows {
		if cw.IsRunning() {
			windows = append(windows, CharacterWindowInfo{
				ID:            id,
				CharacterName: cw.GetCharacterName(),
				IsRunning:     true,
				Scale:         cw.GetScale(),
			})
		}
	}
	return windows
}

func (a *App) DestroyAllCharacters() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, cw := range a.activeWindows {
		cw.Close()
	}
	a.activeWindows = make(map[string]*Window.CharacterWindow)
}

func (a *App) SetCharacterScale(windowId string, scale float64) bool {
	a.mu.RLock()
	charWindow, exists := a.activeWindows[windowId]
	a.mu.RUnlock()

	if !exists {
		return false
	}

	charWindow.SetScale(scale)
	return true
}

func (a *App) GetPreviewImageBase64(characterName string) string {
	previewPath, err := AnimationEngine.GetPreviewImage(a.framesPath, characterName)
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(previewPath)
	if err != nil {
		return ""
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
}

func (a *App) GetPreviewFrames(characterName string, maxFrames int) []string {
	charPath := AnimationEngine.GetCharacterFramesPath(a.framesPath, characterName)

	entries, err := os.ReadDir(charPath)
	if err != nil {
		return []string{}
	}

	var frames []string
	count := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := ""
		if len(name) > 4 {
			ext = name[len(name)-4:]
		}

		if ext != ".png" && ext != ".jpg" && ext != ".PNG" && ext != ".JPG" {
			continue
		}

		framePath := charPath + "/" + name
		data, err := os.ReadFile(framePath)
		if err != nil {
			continue
		}

		frames = append(frames, "data:image/png;base64,"+base64.StdEncoding.EncodeToString(data))
		count++

		if maxFrames > 0 && count >= maxFrames {
			break
		}
	}

	return frames
}

func (a *App) OpenFramesDir() error {
	framesPath := a.cfg.FramesPath

	if err := config.EnsureFramesDir(a.cfg); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", framesPath)
	case "darwin":
		cmd = exec.Command("open", framesPath)
	default:
		cmd = exec.Command("xdg-open", framesPath)
	}

	return cmd.Start()
}

func (a *App) OpenConfig() error {
	configPath := config.GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := config.SaveConfig(a.cfg); err != nil {
			return err
		}
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad", configPath)
	case "darwin":
		cmd = exec.Command("open", "-t", configPath)
	default:
		editors := []string{"xdg-open", "nano", "vi"}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				cmd = exec.Command(editor, configPath)
				break
			}
		}
		if cmd == nil {
			cmd = exec.Command("xdg-open", configPath)
		}
	}

	return cmd.Start()
}

func (a *App) GetFramesPath() string {
	return a.cfg.FramesPath
}

func (a *App) GetConfigPath() string {
	return config.GetConfigPath()
}

func (a *App) BrowseBfkFile() string {
	filePath, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Boccho Frame Pack",
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "Boccho Frame Pack (*.bfk)",
				Pattern:     "*.bfk",
			},
		},
	})
	if err != nil {
		fmt.Printf("Error opening file dialog: %v\n", err)
		return ""
	}
	return filePath
}

func (a *App) GetBfkPackInfo(filePath string) PackManagement.PackInfo {
	return PackManagement.GetPackInfo(filePath)
}

func (a *App) InstallBfkPack(filePath string) error {
	return PackManagement.InstallPack(filePath, a.cfg.FramesPath)
}
