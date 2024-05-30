package main

import (
	"diplom/internal/clients/grpc"
	"diplom/internal/config"
	"diplom/internal/windows"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/net/context"
	"log"
)

func main() {
	cfg, err := config.LoadConfigFromYaml()
	if err != nil {
		log.Fatalf("Error load config app: %w", err)
	}

	myApp := app.New()
	window := myApp.NewWindow(cfg.AppSetting.Name)
	myApp.Settings().SetTheme(theme.DarkTheme())
	window.Resize(fyne.NewSize(cfg.AppSetting.Width, cfg.AppSetting.Height))

	ssoClient, err := grpc.New(
		context.Background(),
		cfg.AddrConnection,
		cfg.TimeoutConnection,
		cfg.RetriesCount,
		myApp,
	)
	if err != nil {
		log.Fatalf("Error initializing gRPC client: %w", err)
	}

	window.SetMainMenu(fyne.NewMainMenu(SettingMenuApp(ssoClient)))
	loginContent := windows.GetLoginContent(window, ssoClient)

	window.CenterOnScreen()
	window.SetContent(loginContent)
	window.ShowAndRun()
}

func SettingMenuApp(ssoClient *grpc.Client) *fyne.Menu {
	themeApp := fyne.NewMenuItem("Выбрать тему", nil)

	themeApp.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Темная", func() {
			ssoClient.App.Settings().SetTheme(theme.DarkTheme())
		}),
		fyne.NewMenuItem("Светлая", func() {
			ssoClient.App.Settings().SetTheme(theme.DefaultTheme())
		}),
	)

	settingMenu := fyne.NewMenu("Настройки", themeApp)

	return settingMenu
}
