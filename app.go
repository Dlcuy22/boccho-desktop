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
	"boccho-ui/Window"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type App struct {
	ctx           context.Context
	activeWindows map[string]*Window.CharacterWindow
	mu            sync.RWMutex
	framesPath    string
}

type CharacterWindowInfo struct {
	ID            string  `json:"id"`
	CharacterName string  `json:"characterName"`
	IsRunning     bool    `json:"isRunning"`
	Scale         float64 `json:"scale"`
}

func NewApp() *App {
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	return &App{
		activeWindows: make(map[string]*Window.CharacterWindow),
		framesPath:    filepath.Join(execDir, "Frames"),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if _, err := os.Stat(a.framesPath); os.IsNotExist(err) {
		workDir, _ := os.Getwd()
		a.framesPath = filepath.Join(workDir, "Frames")
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
