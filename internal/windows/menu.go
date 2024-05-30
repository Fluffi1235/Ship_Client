package windows

import (
	"diplom/internal/clients/grpc"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func GetMenu(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	usernameLabel := widget.NewLabel(ssoClient.UserName)

	logoutButton := widget.NewButtonWithIcon("Выход", theme.LogoutIcon(), func() {
		loginContent := GetLoginContent(window, ssoClient)
		window.SetContent(loginContent)
	})

	profile := widget.NewButtonWithIcon("Учетная карточка", theme.AccountIcon(), func() {
		window.SetContent(Profile(window, ssoClient))
	})

	machine := widget.NewButton("Корабельные системы", func() {
		window.SetContent(MachineDepartment(window, ssoClient))
	})
	machine.Resize(fyne.NewSize(960, 140))

	radar := widget.NewButton("Система РЛС", func() {
		window.SetContent(Radar(window, ssoClient))
	})
	radar.Resize(fyne.NewSize(960, 140))

	tabs := container.NewAppTabs(
		&container.TabItem{Text: "Сотрудники", Content: UserList(window, ssoClient)},
	)

	return container.NewVBox(
		container.NewHBox(layout.NewSpacer(), usernameLabel, profile),
		container.NewHBox(layout.NewSpacer(), logoutButton),
		container.NewGridWithRows(0, container.NewWithoutLayout(machine), container.NewWithoutLayout(radar)),
		layout.NewSpacer(),
		tabs,
		layout.NewSpacer(),
	)
}
