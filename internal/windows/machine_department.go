package windows

import (
	"context"
	ssov1 "diplom/gen"
	"diplom/internal/clients/grpc"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/errgo.v2/fmt/errors"
	"image/color"
	"log"
)

const (
	space                    = "       "
	engineOffImage           = "pic/machine_dep/engine_off.png"
	engineWorkingImage       = "pic/machine_dep/engine_working.png"
	engineDamagedImage       = "pic/machine_dep/engine_damaged.png"
	tempNormalImage          = "pic/machine_dep/temp_normal.png"
	tempIncreasedImage       = "pic/machine_dep/temp_increased.png"
	tempHighIncreasedImage   = "pic/machine_dep/temp_high_increased.png"
	tempCriticalImage        = "pic/machine_dep/temp_critical.png"
	coolingOffImage          = "pic/machine_dep/cooling_off.png"
	coolingWorkingImage      = "pic/machine_dep/cooling_working.png"
	coolingDamagedImage      = "pic/machine_dep/cooling_damaged.png"
	electricOffImage         = "pic/machine_dep/electric_off.png"
	electricLowPowerImage    = "pic/machine_dep/electric_low_power.png"
	electricMediumPowerImage = "pic/machine_dep/electric_medium_power.png"
	electricHighPowerImage   = "pic/machine_dep/electric_hight_power.png"
	electricFullPowerImage   = "pic/machine_dep/electric_full_power.png"
	electricDamagedImage     = "pic/machine_dep/electric_damaged.png"
	fuelOffImage             = "pic/machine_dep/fuel_off.png"
	fuelWorkingImage         = "pic/machine_dep/fuel_working.png"
	fuelDamagedImage         = "pic/machine_dep/fuel_damaged.png"
)

func MachineDepartment(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	// Создаем контекст с функцией отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Создаем кнопку "Назад"
	backButton := widget.NewButton("Назад", func() {
		// Отменяем контекст для закрытия соединения
		cancel()
		window.SetContent(GetMenu(window, ssoClient)) // Возвращаемся в меню при нажатии кнопки
	})
	backButton.Resize(fyne.NewSize(100, 25))
	usernameLabel := widget.NewLabel(ssoClient.UserName)

	mainHeader := canvas.NewText("Информация по машинному отделу", color.Black)
	mainHeader.TextStyle = fyne.TextStyle{Bold: true}
	mainHeader.TextSize = 40

	// Создаем заголовок
	header1 := canvas.NewText("Двигатели", color.Black)
	header1.TextStyle = fyne.TextStyle{Bold: true}
	header1.TextSize = 35
	header2 := canvas.NewText("Системы охлаждения", color.Black)
	header2.TextStyle = fyne.TextStyle{Bold: true}
	header2.TextSize = 35
	header3 := canvas.NewText("Электрогенераторы", color.Black)
	header3.TextStyle = fyne.TextStyle{Bold: true}
	header3.TextSize = 35
	header4 := canvas.NewText("Топливные системы", color.Black)
	header4.TextStyle = fyne.TextStyle{Bold: true}
	header4.TextSize = 35

	statusContainer1 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer2 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer3 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer4 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer5 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer6 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer7 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	statusContainer8 := container.NewVBox(
		widget.NewLabel(""),
		container.NewHBox(),
	)

	card0 := NewCustomCard(
		container.NewCenter(mainHeader),
		1920, 50,
	)

	// Создаем кастомные карточки с рамками
	card1 := NewCustomCard(
		container.NewVBox(
			container.NewCenter(header1),
			container.NewHBox(statusContainer1, widget.NewLabel(space), statusContainer2),
		),
		870, 450,
	)

	card2 := NewCustomCard(
		container.NewVBox(
			container.NewCenter(header2),
			container.NewHBox(statusContainer3, widget.NewLabel(space), statusContainer4),
		),
		1050, 450,
	)

	card3 := NewCustomCard(
		container.NewVBox(
			container.NewCenter(header3),
			container.NewHBox(statusContainer5, widget.NewLabel(space), statusContainer6),
		),
		870, 450,
	)

	card4 := NewCustomCard(
		container.NewVBox(
			container.NewCenter(header4),
			container.NewHBox(statusContainer7, widget.NewLabel(space), statusContainer8),
		),
		1050, 450,
	)

	// Создаем контейнер для заголовка и statusContainer
	content := container.NewVBox(
		container.NewHBox(container.NewWithoutLayout(backButton), layout.NewSpacer(), usernameLabel),
		card0,
		container.NewVBox(
			container.NewCenter(
				container.NewHBox(card1, card2),
			),
			container.NewCenter(
				container.NewHBox(card3, card4),
			),
		),
	)

	// Создаем горутину для чтения потока
	go func() {
		stream, err := ssoClient.GetMachineInfo(ctx, &ssov1.UserName{UserName: ssoClient.UserName})
		if err != nil {
			dialog.ShowError(errors.New("Failed to establish connection"), window)
			window.SetContent(GetLoginContent(window, ssoClient))
			return
		}
		defer stream.CloseSend()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.Recv()
				if err != nil {
					// Получаем статус ошибки
					st, ok := status.FromError(err)
					if ok && st.Code() == codes.Canceled {
						// Если ошибка из-за отмены контекста, просто выходим из цикла
						return
					}
					// Обрабатываем другие ошибки
					log.Printf("Error receiving message: %v", err)
					dialog.ShowError(errors.New("Failed to establish connection"), window)
					window.SetContent(GetLoginContent(window, ssoClient))
					return
				}

				// Обновляем информацию в зависимости от полученного сообщения
				statusContainer1 = BuildEngineContainer(msg.Engine1, statusContainer1, "Основной двигатель")
				statusContainer2 = BuildEngineContainer(msg.Engine2, statusContainer2, "Резервный двигатель")
				statusContainer3 = BuildCoolingSystemContainer(msg.CoolingSystem1, statusContainer3, "Основная система охлождения")
				statusContainer4 = BuildCoolingSystemContainer(msg.CoolingSystem2, statusContainer4, "Резервная система охлождения")
				statusContainer5 = BuildGeneratorCon(msg.Generator1, statusContainer5, "Основной генератор")
				statusContainer6 = BuildGeneratorCon(msg.Generator2, statusContainer6, "Резервный генератор")
				statusContainer7 = BuildFuelSystemCon(msg.FuelSystem1, statusContainer7, "Основная система подачи топлива")
				statusContainer8 = BuildFuelSystemCon(msg.FuelSystem2, statusContainer8, "Резезвная система подачи топлива")

				window.Content().(fyne.CanvasObject).Refresh()
			}
		}
	}()

	// Возвращаем контейнер с заголовком и центральным statusContainer
	return content
}

func BuildEngineContainer(engine *ssov1.Engine, con *fyne.Container, headerText string) *fyne.Container {
	header := canvas.NewText(headerText, color.Black)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.TextSize = 20

	imagePath := engineOffImage
	statusLabel := widget.NewLabel("Двигатель Выкл.")
	revLabel := widget.NewLabel("Обороты: - об/мин")
	tempLabel := widget.NewLabel("Температура: - °C")
	voltLabel := widget.NewLabel("Нарпяжение: - В")

	imageTempPath := "pic/temp_dont_info.png"
	imageTemp := canvas.NewImageFromFile(imageTempPath)
	imageTemp.FillMode = canvas.ImageFillContain
	imageTemp.SetMinSize(fyne.NewSize(200, 200))

	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(30, 30))

	if engine.Status == "working" {
		imagePath = engineWorkingImage
		newImage := canvas.NewImageFromFile(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		newImage.SetMinSize(fyne.NewSize(200, 200))
		// Заменяем старое изображение новым на экране

		if engine.Temperature < 40 {
			imageTempPath = tempNormalImage
		} else if engine.Temperature < 60 {
			imageTempPath = tempIncreasedImage
		} else if engine.Temperature < 80 {
			imageTempPath = tempHighIncreasedImage
		} else {
			imageTempPath = tempCriticalImage
		}
		newImageTemp := canvas.NewImageFromFile(imageTempPath)
		newImageTemp.FillMode = canvas.ImageFillContain
		newImageTemp.SetMinSize(fyne.NewSize(30, 30))
		// Заменяем старое изображение новым на экране

		con.Objects[1] = container.NewVBox(
			container.NewCenter(header),
			container.NewHBox(
				newImage,
				container.NewVBox(
					statusLabel,
					revLabel,
					container.NewHBox(
						tempLabel,
						newImageTemp,
					),
					voltLabel,
				),
			),
		)

		statusLabel.SetText("Двигатель Вкл.")
		revLabel.SetText("Обороты: " + fmt.Sprint(engine.Rpm) + " об/мин")
		tempLabel.SetText("Температура: " + fmt.Sprint(engine.Temperature) + " °C")
		voltLabel.SetText("Нарпяжение: " + fmt.Sprint(engine.Voltage) + " В")
	} else {
		if engine.Status == "dont working" {
			imagePath = engineOffImage
			statusLabel.SetText("Двигатель Выкл.")
		} else {
			imagePath = engineDamagedImage
			statusLabel.SetText("Двигатель Поврежден")
		}

		newImage := canvas.NewImageFromFile(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		newImage.SetMinSize(fyne.NewSize(200, 200))

		// Заменяем старое изображение новым на экране
		con.Objects[1] = container.NewVBox(
			container.NewCenter(header),
			container.NewHBox(
				newImage,
				container.NewVBox(
					statusLabel,
				),
			),
		)
	}

	return con
}

func BuildCoolingSystemContainer(cool *ssov1.CoolingSystem, con *fyne.Container, headerText string) *fyne.Container {
	header := canvas.NewText(headerText, color.Black)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.TextSize = 20

	imagePath := coolingOffImage
	statusLabel := widget.NewLabel("Система Выкл.")
	tempLabel := widget.NewLabel("Температура жидкости: - °C")
	pressureLabel := widget.NewLabel("Давление: - бар")

	imageTempPath := "pic/temp_dont_info.png"
	imageTemp := canvas.NewImageFromFile(imageTempPath)
	imageTemp.FillMode = canvas.ImageFillContain
	imageTemp.SetMinSize(fyne.NewSize(200, 200))

	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(30, 30))

	tempLabel.SetText("Температура жидкости: " + fmt.Sprint(cool.CoolantTemperature) + " °C")
	pressureLabel.SetText("Давление в системе: " + fmt.Sprint(cool.SystemPressure) + " бар")

	if cool.Status == "working" {
		imagePath = coolingWorkingImage
		newImage := canvas.NewImageFromFile(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		newImage.SetMinSize(fyne.NewSize(200, 200))
		// Заменяем старое изображение новым на экране

		if cool.CoolantTemperature < 40 {
			imageTempPath = tempNormalImage
		} else if cool.CoolantTemperature < 60 {
			imageTempPath = tempIncreasedImage
		} else if cool.CoolantTemperature < 80 {
			imageTempPath = tempHighIncreasedImage
		} else {
			imageTempPath = tempCriticalImage
		}
		newImageTemp := canvas.NewImageFromFile(imageTempPath)
		newImageTemp.FillMode = canvas.ImageFillContain
		newImageTemp.SetMinSize(fyne.NewSize(30, 30))
		// Заменяем старое изображение новым на экране

		con.Objects[1] = container.NewVBox(
			container.NewCenter(header),
			container.NewHBox(
				newImage,
				container.NewVBox(
					statusLabel,
					container.NewHBox(
						tempLabel,
						newImageTemp,
					),
					pressureLabel,
				),
			),
		)

		statusLabel.SetText("Система Вкл.")
	} else {
		if cool.Status == "dont working" {
			imagePath = coolingOffImage
			statusLabel.SetText("Система Выкл.")
		} else {
			imagePath = coolingDamagedImage
			statusLabel.SetText("Система Повреждена")
		}

		newImage := canvas.NewImageFromFile(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		newImage.SetMinSize(fyne.NewSize(200, 200))

		// Заменяем старое изображение новым на экране
		con.Objects[1] = container.NewVBox(
			container.NewCenter(header),
			container.NewHBox(
				newImage,
				container.NewVBox(
					statusLabel,
					container.NewHBox(
						tempLabel,
					),
					pressureLabel,
				),
			),
		)
	}

	return con
}

func BuildGeneratorCon(gen *ssov1.Generator, con *fyne.Container, headerText string) *fyne.Container {
	header := canvas.NewText(headerText, color.Black)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.TextSize = 20

	statusLabel := widget.NewLabel("Система Выкл.")
	powerLabel := widget.NewLabel("Мощность: - кВт")
	fuelLabel := widget.NewLabel("Тип топлива: -")
	voltLabel := widget.NewLabel("Нарпяжение: - В")

	imagePath := electricOffImage
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(30, 30))

	if gen.Status == "working" {
		if gen.Power < 20 {
			imagePath = electricLowPowerImage
		} else if gen.Power < 100 {
			imagePath = electricMediumPowerImage
		} else if gen.Power < 500 {
			imagePath = electricHighPowerImage
		} else {
			imagePath = electricFullPowerImage
		}

		newImage := canvas.NewImageFromFile(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		newImage.SetMinSize(fyne.NewSize(200, 200))
		// Заменяем старое изображение новым на экране

		con.Objects[1] = container.NewVBox(
			container.NewCenter(header),
			container.NewHBox(
				newImage,
				container.NewVBox(
					statusLabel,
					powerLabel,
					fuelLabel,
					voltLabel,
				),
			),
		)

		statusLabel.SetText("Генератор Вкл.")
		powerLabel.SetText("Мощность: " + fmt.Sprint(gen.Power) + " кВт")
		fuelLabel.SetText("Тип топлева: " + fmt.Sprint(gen.Fuel))
		voltLabel.SetText("Напряжение: " + fmt.Sprint(gen.Voltage) + " В")
	} else {
		if gen.Status == "dont working" {
			imagePath = electricOffImage
			statusLabel.SetText("Генератор Выкл.")
		} else {
			imagePath = electricDamagedImage
			statusLabel.SetText("Генератор Поврежден")
		}

		newImage := canvas.NewImageFromFile(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		newImage.SetMinSize(fyne.NewSize(200, 200))

		// Заменяем старое изображение новым на экране
		con.Objects[1] = container.NewVBox(
			container.NewCenter(header),
			container.NewHBox(
				newImage,
				container.NewVBox(
					statusLabel,
				),
			),
		)
	}

	return con
}

func BuildFuelSystemCon(system *ssov1.FuelSystem, con *fyne.Container, headerText string) *fyne.Container {
	header := canvas.NewText(headerText, color.Black)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.TextSize = 20

	statusLabel := widget.NewLabel(fmt.Sprintf("Статус: %s", system.Status))
	fuelLevelLabel := widget.NewLabel(fmt.Sprintf("Уровень топлива: %.2f%%", system.FuelLevel))
	contaminantsLevelLabel := widget.NewLabel(fmt.Sprintf("Уровень загрязнителей: %.2f ppm", system.ContaminantsLevel))
	fuelFilterStatusLabel := widget.NewLabel(fmt.Sprintf("Состояние фильтра: %s", system.FuelFilterStatus))
	flowRateLabel := widget.NewLabel(fmt.Sprintf("Расход топлива: %.2f л/ч", system.FlowRate))
	leakDetectionLabel := widget.NewLabel(fmt.Sprintf("Обнаружение утечек: %t", system.LeakDetection))
	fuelPumpStatusLabel := widget.NewLabel(fmt.Sprintf("Статус насоса: %s", system.FuelPumpStatus))

	imagePath := fuelOffImage // Укажите путь к вашему изображению
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(200, 200))

	if system.Status == "working" {
		imagePath = fuelWorkingImage // Укажите путь к изображению рабочего состояния
		statusLabel = widget.NewLabel(fmt.Sprintf("Cистема Вкл."))
	} else if system.Status == "dont working" {
		statusLabel = widget.NewLabel(fmt.Sprintf("Cистема Выкл."))
		imagePath = fuelOffImage // Укажите путь к изображению нерабочего состояния
	} else {
		statusLabel = widget.NewLabel(fmt.Sprintf("Cистема Повреждена"))
		imagePath = fuelDamagedImage // Укажите путь к изображению состояния "выведен из строя"
	}

	newImage := canvas.NewImageFromFile(imagePath)
	newImage.FillMode = canvas.ImageFillContain
	newImage.SetMinSize(fyne.NewSize(200, 200))

	con.Objects[1] = container.NewVBox(
		container.NewCenter(header),
		container.NewHBox(
			newImage,
			container.NewVBox(
				statusLabel,
				fuelLevelLabel,
				contaminantsLevelLabel,
				fuelFilterStatusLabel,
				flowRateLabel,
				leakDetectionLabel,
				fuelPumpStatusLabel,
			),
		),
	)

	return con
}

func NewCustomCard(content fyne.CanvasObject, width, height float32) fyne.CanvasObject {
	border := canvas.NewRectangle(color.RGBA{R: 96, G: 96, B: 96, A: 170}) // Темно-зеленый цвет рамки
	border.SetMinSize(fyne.NewSize(width, height))                         // Устанавливаем минимальный размер рамки
	inner := container.NewPadded(content)                                  // Добавляем отступы внутри карточки

	return container.New(&borderLayout{border: border, width: width, height: height}, border, inner)
}

// Кастомный лэйаут для рамки
type borderLayout struct {
	border *canvas.Rectangle
	width  float32
	height float32
}

func (b *borderLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(b.width, b.height) // Используем заданные размеры для рамки
}

func (b *borderLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	borderSize := fyne.NewSize(b.width, b.height)
	innerSize := fyne.NewSize(b.width-20, b.height-20) // Отнимаем отступы для внутреннего контента
	objects[0].Resize(borderSize)
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Resize(innerSize)
	objects[1].Move(fyne.NewPos(10, 10)) // Отступ для внутреннего контента
}
