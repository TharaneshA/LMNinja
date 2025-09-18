package main

import (
	"embed"
	"lmninja/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	appInstance := app.NewApp()

	err := wails.Run(&options.App{
		Title:              "LMNinja",
		Width:              1000,
		Height:             720,
		MinWidth:           960,
		MinHeight:          640,
		Frameless:          true,
		BackgroundColour:   &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.INFO,
		Assets:             assets,
		OnStartup:          appInstance.Startup,
		OnShutdown:         appInstance.Shutdown,
		Bind: []interface{}{
			appInstance,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false, // Set to false for better performance on Windows
			WindowIsTranslucent:  false,
			DisableFramelessWindowDecorations: false,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false, // Set to false for better performance
			About: &mac.AboutInfo{
				Title:   "LMNinja",
				Message: "The Postman for LLM Security",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}