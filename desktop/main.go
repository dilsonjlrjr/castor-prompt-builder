package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "CASTOR Builder",
		Width:            960,
		Height:           640,
		MinWidth:         800,
		MinHeight:        560,
		MaxWidth:         960,
		MaxHeight:        640,
		DisableResize:    true,
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 20, A: 1},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "CASTOR Builder",
				Message: "Construa prompts estruturados para LLMs.",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
