package windows

import (
	"context"
	"io"
	"net/http"
	"sort"

	ssov1 "diplom/gen"
	"diplom/internal/clients/grpc"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/errgo.v2/errors"
)

const defaultAvatarPath = "pic/user/avatar.jpg"

func updateImageUsers(url string, image *canvas.Image) {
	go func() {
		respPhoto, err := http.Get(url)
		if err != nil {
			return
		}
		defer respPhoto.Body.Close()

		imageData, err := io.ReadAll(respPhoto.Body)
		if err != nil {
			return
		}

		image.Resource = fyne.NewStaticResource("userAvatar", imageData)
		image.Refresh()
	}()
}

func createUserRow(user *ssov1.UserInfo) []fyne.CanvasObject {
	image := canvas.NewImageFromFile(defaultAvatarPath)
	image.SetMinSize(fyne.NewSize(50, 50))
	image.FillMode = canvas.ImageFillContain
	updateImageUsers(user.PhotoUrl, image)

	return []fyne.CanvasObject{
		image,
		widget.NewLabelWithStyle(user.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(user.Surname, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(user.Department, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(user.Rank, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(user.Email, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	}
}

func sortUsers(users []*ssov1.UserInfo, column string, asc bool) {
	sort.Slice(users, func(i, j int) bool {
		switch column {
		case "name":
			if asc {
				return users[i].Name < users[j].Name
			}
			return users[i].Name > users[j].Name
		case "surname":
			if asc {
				return users[i].Surname < users[j].Surname
			}
			return users[i].Surname > users[j].Surname
		case "department":
			if asc {
				return users[i].Department < users[j].Department
			}
			return users[i].Department > users[j].Department
		case "rank":
			if asc {
				return users[i].Rank < users[j].Rank
			}
			return users[i].Rank > users[j].Rank
		case "email":
			if asc {
				return users[i].Email < users[j].Email
			}
			return users[i].Email > users[j].Email
		}
		return false
	})
}

func updateUserList(userList *fyne.Container, users []*ssov1.UserInfo) {
	userList.Objects = nil
	for _, user := range users {
		row := createUserRow(user)
		userRow := container.NewGridWithColumns(6, row...)
		userList.Add(userRow)
	}
	userList.Refresh()
}

func UserList(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	// Загружаем список пользователей
	resp, err := ssoClient.GetAllUsers(context.Background(), &ssov1.UserName{UserName: ssoClient.UserName})
	if err != nil {
		dialog.ShowError(errors.New("Ошибка входа"), window)
		return GetLoginContent(window, ssoClient)
	}

	users := resp.Users
	userList := container.NewVBox()

	// Создаем заголовки таблицы с кнопками для сортировки
	headers := []fyne.CanvasObject{
		widget.NewLabelWithStyle("Фотография", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButtonWithIcon("Имя", theme.ListIcon(), func() {
			sortUsers(users, "name", true)
			updateUserList(userList, users)
		}),
		widget.NewButtonWithIcon("Фамилия", theme.ListIcon(), func() {
			sortUsers(users, "surname", true)
			updateUserList(userList, users)
		}),
		widget.NewButtonWithIcon("Боевая часть", theme.ListIcon(), func() {
			sortUsers(users, "department", true)
			updateUserList(userList, users)
		}),
		widget.NewButtonWithIcon("Звание", theme.ListIcon(), func() {
			sortUsers(users, "rank", true)
			updateUserList(userList, users)
		}),
		widget.NewButtonWithIcon("Почта", theme.ListIcon(), func() {
			sortUsers(users, "email", true)
			updateUserList(userList, users)
		}),
	}

	headerContainer := container.NewGridWithColumns(6, headers...)

	// Создаем контейнер для списка пользователей
	for _, user := range users {
		row := createUserRow(user)
		userRow := container.NewGridWithColumns(6, row...)
		userList.Add(userRow)
	}

	scrollContainer := container.NewScroll(userList)
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))

	// Основной контейнер, включающий заголовки и прокручиваемую часть
	content := container.NewVBox(
		headerContainer,
		scrollContainer,
	)

	return content
}
