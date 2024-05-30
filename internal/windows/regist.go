package windows

import (
	"context"
	ssov1 "diplom/gen"
	"diplom/internal/clients/grpc"
	"diplom/internal/windows/auth"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/errgo.v2/errors"
	"image/color"
)

const (
	PlaceHolderName = "Имя"

	PlaceHolderSurname    = "Фамилия"
	PlaceHolderDepartment = "Отдел"
	PlaceHolderRank       = "Ранг"

	spaceTitleSize           = 20
	spaceBetweenFieldsSize   = 10
	spaceBetweenButtonsWidth = 305
)

var (
	departments = []string{"Навигационный", "Машинный", "Радиосвязи", "Административный", "Безопасности"}
	ranks       = []string{"Лейтенант", "Старший лейтенант", "Капитан-лейтенант", "Капитан 3-го ранга",
		"Капитан 2-го ранга", "Капитан 1-го ранга", "Контр-адмирал", "Вице-адмирал", "Адмирал", "Адмирал флота"}
)

func RegisterUser(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	var user ssov1.RegisterRequest

	titleText := canvas.NewText("Регистрация пользователя", color.NRGBA{R: 164, G: 151, B: 151, A: 255})
	titleText.TextSize = 45

	loginEntry := newWidget(PlaceHolderLogin, fyne.NewSize(600, 40))

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.PlaceHolder = PlaceHolderPassword
	passwordEntry.Resize(fyne.NewSize(600, 40))
	passwordEntry.Password = true

	nameEntry := newWidget(PlaceHolderName, fyne.NewSize(600, 40))
	surnameEntry := newWidget(PlaceHolderSurname, fyne.NewSize(600, 40))

	departmentSelect := widget.NewSelect(departments, func(selected string) {
		user.Department = selected
	})
	departmentSelect.PlaceHolder = PlaceHolderDepartment
	departmentSelect.Resize(fyne.NewSize(600, 40))

	rankEntry := widget.NewSelect(ranks, func(selected string) {
		user.Rank = selected
	})
	rankEntry.PlaceHolder = PlaceHolderRank
	rankEntry.Resize(fyne.NewSize(600, 40))

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Электронная почта")
	emailEntry.Resize(fyne.NewSize(600, 40))

	registerButton := widget.NewButtonWithIcon("Регистрация", theme.ContentAddIcon(), func() {
		if len(passwordEntry.Text) < 8 {
			dialog.ShowError(errors.New("Слишком короткий пароль, минимум 8 символов"), window)
			return
		}

		password, err := auth.PasswordHash(passwordEntry.Text)
		if err != nil {
			dialog.ShowError(errors.New("Ошибка сервера"), window)
			return
		}

		user.UserName = loginEntry.Text
		user.Password = password
		user.Name = nameEntry.Text
		user.Surname = surnameEntry.Text
		user.Email = emailEntry.Text

		if user.UserName == "" || password == nil || user.Name == "" || user.Surname == "" || user.Department == "" ||
			user.Rank == "" || user.Email == "" {
			dialog.ShowError(errors.New("Заполните все поля"), window)
			return
		}

		resp, err := ssoClient.Register(context.Background(), &user)
		if err != nil {
			dialog.ShowError(errors.New("Ошибка при регистрации пользователя"), window)
			return
		}

		w := dialog.NewInformation("Успех", fmt.Sprintf("Пользователь %s зарегистрирован", resp.UserName), window)
		w.SetOnClosed(func() { window.SetContent(GetLoginContent(window, ssoClient)) })
		w.Show()
	})
	registerButton.Resize(fyne.NewSize(200, 70))

	backButton := widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {
		window.SetContent(GetLoginContent(window, ssoClient))
	})
	backButton.Resize(fyne.NewSize(200, 70))

	spaceTitle := container.New(layout.NewGridWrapLayout(fyne.NewSize(0, spaceTitleSize)), layout.NewSpacer())
	spaceBetweenFields := container.New(layout.NewGridWrapLayout(fyne.NewSize(0, spaceBetweenFieldsSize)), layout.NewSpacer())
	spaceBetweenButtons := container.New(layout.NewGridWrapLayout(fyne.NewSize(spaceBetweenButtonsWidth, 0)), layout.NewSpacer())

	return container.NewCenter(
		container.NewVBox(
			container.NewCenter(titleText),
			spaceTitle,
			container.NewWithoutLayout(loginEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(passwordEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(nameEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(surnameEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(emailEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(rankEntry),
			spaceBetweenFields,
			container.NewWithoutLayout(departmentSelect),
			spaceBetweenFields,
			container.NewHBox(
				container.NewWithoutLayout(backButton),
				spaceBetweenButtons,
				container.NewWithoutLayout(registerButton),
			),
		),
	)
}
