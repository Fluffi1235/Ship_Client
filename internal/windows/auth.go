package windows

import (
	ssov1 "diplom/gen"
	"diplom/internal/clients/grpc"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v2/errors"
)

const (
	PlaceHolderLogin    = "Логин"
	PlaceHolderPassword = "Пароль"
	PlaceHolderRegister = "Создать учетную запись"
	PlaceHolderEnter    = "Вход"
)

func GetLoginContent(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	usernameEntry := newWidget(PlaceHolderLogin, fyne.NewSize(600, 40))

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = PlaceHolderPassword
	passwordEntry.Resize(fyne.NewSize(600, 40))
	passwordEntry.Password = true

	registerButton := widget.NewButton(PlaceHolderRegister, func() {
		window.SetContent(RegisterUser(window, ssoClient))
	})
	registerButton.Resize(fyne.NewSize(200, 70))

	loginButton := widget.NewButtonWithIcon(PlaceHolderEnter, theme.LoginIcon(), func() {
		if usernameEntry.Text == "" || passwordEntry.Text == "" {
			dialog.ShowError(errors.New("Заполните все поля"), window)
			return
		}

		ssoClient.Password = passwordEntry.Text

		clientData := &ssov1.LoginRequest{
			UserName: usernameEntry.Text,
			Password: []byte(passwordEntry.Text),
		}

		resp, err := ssoClient.Login(context.Background(), clientData)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		ssoClient.UserName = usernameEntry.Text
		ssoClient.Token = resp.Token

		if ssoClient.Token != "" {
			window.SetContent(GetMenu(window, ssoClient))
		} else {
			dialog.ShowError(errors.New("Неверные данные пользователя"), window)
		}

	})
	loginButton.Resize(fyne.NewSize(200, 70))

	spaceBetweenFields := container.New(layout.NewGridWrapLayout(fyne.NewSize(0, spaceBetweenFieldsSize)), layout.NewSpacer())
	spaceBetweenButtons := container.New(layout.NewGridWrapLayout(fyne.NewSize(spaceBetweenButtonsWidth+10, 0)), layout.NewSpacer())

	return container.NewCenter(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewWithoutLayout(usernameEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(passwordEntry),
			spaceBetweenFields,
			container.NewHBox(
				container.NewWithoutLayout(loginButton),
				spaceBetweenButtons,
				container.NewWithoutLayout(registerButton),
			),
			layout.NewSpacer(),
		),
	)
}
