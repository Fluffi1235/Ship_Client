package windows

import (
	ssov1 "diplom/gen"
	"diplom/internal/clients/grpc"
	"diplom/internal/windows/auth"
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v2/fmt/errors"
	"image/color"
	"io"
	"net/http"
)

const (
	TitleChangeEmail           = "Изменить почту"
	TitleChangePassword        = "Изменить пароль"
	LabelLogin                 = "Логин: "
	LabelName                  = "Имя: "
	LabelSurname               = "Фамилия: "
	LabelDepartment            = "Отдел: "
	LabelRole                  = "Ранг: "
	LabelEmail                 = "Email: "
	ButtonBack                 = "Назад"
	ButtonUploadImage          = "Загрузить картинку"
	ButtonChangeEmail          = "Изменить почту"
	ButtonChangePassword       = "Изменить пароль"
	ButtonSave                 = "Сохранить"
	ButtonCancel               = "Отмена"
	MsgImageUploadSuccess      = "Изображение успешно загружено"
	MsgEmailChangeSuccess      = "Почта успешно изменена"
	MsgPasswordChangeSuccess   = "Пароль успешно изменен"
	PlaceholderCurrentPassword = "Текущий пароль"
	PlaceholderNewEmail        = "Новая почта"
	PlaceholderNewPassword     = "Новый пароль"
	ImageNotAvailable          = "Не удалось загрузить изображение"
	ImagePath                  = "image.jpg"
	IconSize                   = 600
	TextSize                   = 45
	LoginError                 = "Не удалось установить соединение с сервером"
	ButtonWidth                = 1125
	ButtonHeight               = 40
	ProfileTitle               = "Личный кабинет"
)

func Profile(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	backButton := widget.NewButton(ButtonBack, func() {
		window.SetContent(GetMenu(window, ssoClient))
	})
	backButton.Resize(fyne.NewSize(100, 25))
	usernameLabel := widget.NewLabel(ssoClient.UserName)

	spaceBetweenFields := container.New(layout.NewGridWrapLayout(fyne.NewSize(20, 15)), layout.NewSpacer())

	resp, err := ssoClient.GetUserData(context.Background(), &ssov1.GetUserRequest{
		UserName: ssoClient.UserName,
		Token:    ssoClient.Token,
	})
	if err != nil {
		dialog.ShowError(errors.New(LoginError), window)
		return GetLoginContent(window, ssoClient)
	}

	labels := createUserLabels(resp)
	userInfoContainer := container.NewVBox(
		labels[0],
		spaceBetweenFields,
		labels[1],
		spaceBetweenFields,
		labels[2],
		spaceBetweenFields,
		labels[3],
		spaceBetweenFields,
		labels[4],
		spaceBetweenFields,
		labels[5],
	)

	image := canvas.NewImageFromFile("pic/user/avatar.jpg")
	image.SetMinSize(fyne.NewSize(IconSize, IconSize))
	updateImage(image, resp.PhotoUrl, window)

	uploadButton := createUploadButton(window, ssoClient, image)
	changeEmailButton := createChangeEmailButton(window, ssoClient, userInfoContainer)
	changePasswordButton := createChangePasswordButton(window, ssoClient)

	title := canvas.NewText(ProfileTitle, color.Black)
	title.TextSize = TextSize
	title.Alignment = fyne.TextAlignCenter

	cardBackground := canvas.NewRectangle(color.NRGBA{R: 0, G: 102, B: 102, A: 255})
	card := container.NewStack(
		cardBackground,
		container.NewVBox(
			container.NewCenter(title),
			container.NewHBox(
				image,
				layout.NewSpacer(),
				container.NewCenter(userInfoContainer),
			),
			uploadButton,
			changeEmailButton,
			changePasswordButton,
		),
	)
	card.Resize(fyne.NewSize(900, 600))

	return container.NewVBox(
		container.NewHBox(
			container.NewWithoutLayout(backButton),
			layout.NewSpacer(),
			usernameLabel,
		),
		layout.NewSpacer(),
		container.NewCenter(
			container.NewVBox(
				card,
			),
		),
		layout.NewSpacer(),
	)
}

func createUserLabels(resp *ssov1.GetUserResponse) []*canvas.Text {
	loginLabel := canvas.NewText(LabelLogin+resp.UserName, color.Black)
	nameLabel := canvas.NewText(LabelName+resp.Name, color.Black)
	surnameLabel := canvas.NewText(LabelSurname+resp.Surname, color.Black)
	emailLabel := canvas.NewText(LabelEmail+resp.Email, color.Black)
	roleLabel := canvas.NewText(LabelRole+resp.Rank, color.Black)
	departmentLabel := canvas.NewText(LabelDepartment+resp.Department, color.Black)

	labels := []*canvas.Text{loginLabel, nameLabel, surnameLabel, departmentLabel, roleLabel, emailLabel}
	for _, label := range labels {
		label.TextSize = TextSize
	}
	return labels
}

func updateImage(image *canvas.Image, url string, window fyne.Window) {
	respPhoto, err := http.Get(url)
	if err != nil {
		dialog.ShowError(errors.New(ImageNotAvailable), window)
		return
	}
	defer respPhoto.Body.Close()

	imageData, err := io.ReadAll(respPhoto.Body)
	if err != nil {
		dialog.ShowError(errors.New(ImageNotAvailable), window)
		return
	}

	image.Resource = fyne.NewStaticResource(ImagePath, imageData)
	image.Refresh()
}

func createUploadButton(window fyne.Window, ssoClient *grpc.Client, image *canvas.Image) *widget.Button {
	button := widget.NewButtonWithIcon(ButtonUploadImage, theme.UploadIcon(), func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				return // Пользователь нажал "Cancel", ничего не делаем
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			encodedData := base64.StdEncoding.EncodeToString(data)
			respPhoto, err := ssoClient.ChangeProfilePhoto(context.Background(), &ssov1.ChangeProfilePhotoRequest{
				UserName:  ssoClient.UserName,
				PhotoName: file.URI().Name(),
				Photo:     encodedData,
			})
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			dialog.ShowInformation("Информация", MsgImageUploadSuccess, window)
			updateImage(image, respPhoto.UrlPhoto, window)
		}, window)
	})
	return button
}

func createChangeEmailButton(window fyne.Window, ssoClient *grpc.Client, userInfoContainer *fyne.Container) *widget.Button {
	button := widget.NewButton(ButtonChangeEmail, func() {
		emailEntry := widget.NewEntry()
		emailEntry.Resize(fyne.NewSize(ButtonWidth, ButtonHeight))
		currentPasswordEntry := widget.NewPasswordEntry()
		currentPasswordEntry.Resize(fyne.NewSize(ButtonWidth, ButtonHeight))

		form := &widget.Form{
			Items: []*widget.FormItem{
				widget.NewFormItem(PlaceholderCurrentPassword, currentPasswordEntry),
				widget.NewFormItem(PlaceholderNewEmail, emailEntry),
			},
			OnSubmit: func() {
				if currentPasswordEntry.Text != ssoClient.Password {
					dialog.ShowInformation("Ошибка", "Неверный пароль", window)
				}

				resp, err := ssoClient.ChangeEmail(context.Background(), &ssov1.Email{UserName: ssoClient.UserName, Email: emailEntry.Text})
				if err != nil {
					dialog.ShowInformation("Ошибка", "Попробуйте позже", window)
				}

				dialog.ShowInformation("Информация", MsgEmailChangeSuccess, window)

				email := canvas.NewText(LabelEmail+resp.Email, color.Black)
				email.TextSize = TextSize
				userInfoContainer.Objects[10] = email
				userInfoContainer.Objects[10].Refresh()
			},
			CancelText: ButtonCancel,
			SubmitText: ButtonSave,
		}

		d := dialog.NewCustom(TitleChangeEmail, ButtonCancel, form, window)
		d.Resize(fyne.NewSize(ButtonWidth, 300))
		d.Show()
	})
	return button
}

func createChangePasswordButton(window fyne.Window, ssoClient *grpc.Client) *widget.Button {
	button := widget.NewButton(ButtonChangePassword, func() {
		currentPasswordEntry := widget.NewPasswordEntry()
		currentPasswordEntry.Resize(fyne.NewSize(ButtonWidth, ButtonHeight))

		newPasswordEntry := widget.NewPasswordEntry()
		newPasswordEntry.Resize(fyne.NewSize(ButtonWidth, ButtonHeight))

		form := &widget.Form{
			Items: []*widget.FormItem{
				widget.NewFormItem(PlaceholderCurrentPassword, currentPasswordEntry),
				widget.NewFormItem(PlaceholderNewPassword, newPasswordEntry),
			},
			OnSubmit: func() {
				if currentPasswordEntry.Text != ssoClient.Password {
					dialog.ShowInformation("Ошибка", "Неверный пароль", window)
				}

				if len(currentPasswordEntry.Text) < 8 {
					dialog.ShowInformation("Ошибка", "Слишком короткий пароль, минимум 8 символов", window)
					return
				}

				pas, err := auth.PasswordHash(newPasswordEntry.Text)
				if err != nil {
					dialog.ShowInformation("Ошибка", "Попробуйте позже", window)
				}

				err = ssoClient.ChangePassword(context.Background(), &ssov1.NewPassword{UserName: ssoClient.UserName, Pas: pas})
				if err != nil {
					dialog.ShowInformation("Ошибка", "Попробуйте позже", window)
				}

				ssoClient.Password = newPasswordEntry.Text

				dialog.ShowInformation("Информация", MsgPasswordChangeSuccess, window)
			},
			CancelText: ButtonCancel,
			SubmitText: ButtonSave,
		}

		d := dialog.NewCustom(TitleChangePassword, ButtonCancel, form, window)
		d.Resize(fyne.NewSize(ButtonWidth, 300))
		d.Show()
	})

	return button
}
