package main

/*
main.go - Application entry point

Initializes SDL and starts Wails application with React frontend.
*/

import (
	"boccho-ui/Window"
	"context"
	"embed"
	"fmt"

	"github.com/jupiterrider/purego-sdl3/sdl"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if !sdl.Init(sdl.InitVideo) {
		fmt.Println("SDL_Init failed:", sdl.GetError())
		return
	}
	defer sdl.Quit()

	// Register event watcher for window move events (fixes animation freeze during drag on Windows)
	eventWatcher := Window.CreateEventWatcher()
	sdl.AddEventWatch(eventWatcher, nil)

	fmt.Println("SDL initialized successfully")

	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Boccho Desktop",
		Width:     550,
		Height:    600,
		MinWidth:  400,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 24, G: 24, B: 27, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			app.DestroyAllCharacters()
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
