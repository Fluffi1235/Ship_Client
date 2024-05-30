package windows

import (
	"context"
	ssov1 "diplom/gen"
	"diplom/internal/clients/grpc"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/errgo.v2/errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type Object struct {
	X, Y int
	Type string
}

type PointInfo struct {
	X         float32
	Y         float32
	Angel     float32
	Color     color.RGBA
	Container *fyne.Container
	Line      *canvas.Line
}

var count int

func Radar(window fyne.Window, ssoClient *grpc.Client) fyne.CanvasObject {
	ctx, cancel := context.WithCancel(context.Background())

	backButton := widget.NewButton("Назад", func() {
		cancel()
		window.SetContent(GetMenu(window, ssoClient))
	})
	backButton.Resize(fyne.NewSize(100, 25))
	usernameLabel := widget.NewLabel(ssoClient.UserName)

	imagePath := "pic/radar/img.png"
	originalImage := canvas.NewImageFromFile(imagePath)
	originalImage.FillMode = canvas.ImageFillStretch
	originalImage.SetMinSize(fyne.NewSize(3628, 2052))

	imageWithPoints := container.NewStack(originalImage)
	scrollContainer := container.NewScroll(imageWithPoints)
	scrollContainer.SetMinSize(fyne.NewSize(1919, 600))

	pointsMap := make(map[string]*PointInfo)
	pointsMutex := sync.Mutex{}

	saveButton := widget.NewButton("Сохранить как PNG", func() {
		saveImageWithPoint("saved_image"+strconv.Itoa(count)+".png", imagePath, pointsMap)
		dialog.ShowInformation("Сохранено", "Изображение сохранено как saved_image.png", window)
		count++
	})

	angleStep := 0
	speedStep := 0.0

	angleStepLabel := widget.NewLabel(strconv.Itoa(angleStep) + " град.")
	speedStepLabel := widget.NewLabel(strconv.FormatFloat(speedStep, 'f', 1, 32) + " узел")

	increaseAngleStepButton := widget.NewButton("+", func() {
		if angleStep < 15 {
			angleStep += 1
		}
		angleStepLabel.Text = strconv.Itoa(angleStep) + " град."
		angleStepLabel.Refresh()
	})
	decreaseAngleStepButton := widget.NewButton("-", func() {
		if angleStep > 0 {
			angleStep -= 1
		}
		angleStepLabel.Text = strconv.Itoa(angleStep) + " град."
		angleStepLabel.Refresh()
	})

	increaseSpeedStepButton := widget.NewButton("+", func() {
		if speedStep < 15 {
			speedStep += 0.05
		}
		speedStepLabel.Text = strconv.FormatFloat(speedStep*2, 'f', 1, 32) + " узел"
		speedStepLabel.Refresh()
	})
	decreaseSpeedStepButton := widget.NewButton("-", func() {
		if speedStep > 0 {
			speedStep -= 0.05
		}
		speedStepLabel.Text = strconv.FormatFloat(speedStep*2, 'f', 1, 32) + " узел"
		speedStepLabel.Refresh()
	})

	upButton := widget.NewButton("Увеличить скорость", func() {
		_, err := ssoClient.ChangeShipParameters(context.Background(), &ssov1.UpdateShipParameters{TypeParam: "speed", Value: float32(speedStep)})
		if err != nil {
			log.Printf("Error changing ship parameters: %v", err)
		}
	})
	downButton := widget.NewButton("Уменьшить скорость", func() {
		_, err := ssoClient.ChangeShipParameters(context.Background(), &ssov1.UpdateShipParameters{TypeParam: "speed", Value: float32(0 - speedStep)})
		if err != nil {
			log.Printf("Error changing ship parameters: %v", err)
		}
	})
	rightButton := widget.NewButton("Поворот направо", func() {
		_, err := ssoClient.ChangeShipParameters(context.Background(), &ssov1.UpdateShipParameters{TypeParam: "angel", Value: float32(angleStep)})
		if err != nil {
			log.Printf("Error changing ship parameters: %v", err)
		}
	})
	leftButton := widget.NewButton("Поворот налево", func() {
		_, err := ssoClient.ChangeShipParameters(context.Background(), &ssov1.UpdateShipParameters{TypeParam: "angel", Value: float32(0 - angleStep)})
		if err != nil {
			log.Printf("Error changing ship parameters: %v", err)
		}
	})

	buttonsContainer := container.NewHBox(
		leftButton,
		container.NewVBox(
			upButton,
			downButton,
		),
		rightButton,
	)

	angleStepContainer := container.NewVBox(
		widget.NewLabel("Шаг изменения угла"),
		container.NewCenter(container.NewHBox(decreaseAngleStepButton, angleStepLabel, increaseAngleStepButton)),
	)

	speedStepContainer := container.NewVBox(
		widget.NewLabel("Шаг изменения скорости"),
		container.NewCenter(container.NewHBox(decreaseSpeedStepButton, speedStepLabel, increaseSpeedStepButton)),
	)

	mainObjectLabels := map[string]*widget.Label{
		"coords": widget.NewLabel("Координаты: "),
		"angle":  widget.NewLabel("Курс корабля: "),
		"speed":  widget.NewLabel("Скорость корабля: "),
	}

	mainObjectContainer := container.NewHBox(
		mainObjectLabels["coords"],
		mainObjectLabels["angle"],
		mainObjectLabels["speed"],
	)

	infoShip := container.NewVBox(
		container.NewCenter(widget.NewLabel("Данные по движению коробля")),
		mainObjectContainer,
	)

	infoBox := container.NewVBox()
	infoScrollContainer := container.NewScroll(infoBox)
	infoScrollContainer.SetMinSize(fyne.NewSize(1920, 150))

	content := container.NewVBox(
		container.NewHBox(backButton, layout.NewSpacer(), usernameLabel),
		container.NewHBox(
			container.NewVBox(
				container.NewCenter(widget.NewLabel("Внешняя обстановка")),
				scrollContainer,
				container.NewHBox(saveButton, layout.NewSpacer(), infoShip, layout.NewSpacer(),
					container.NewVBox(
						container.NewCenter(widget.NewLabel("Управление кораблем")),
						container.NewHBox(angleStepContainer, speedStepContainer, buttonsContainer),
					),
				),
				container.NewVBox(
					container.NewCenter(widget.NewLabel("Цели в радиусе")),
					infoScrollContainer,
				),
			),
		),
		layout.NewSpacer(),
	)

	go func() {
		stream, err := ssoClient.GetRadarObjects(ctx, &ssov1.UserName{UserName: ssoClient.UserName})
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
				object, err := stream.Recv()
				if err != nil {
					if ctx.Err() == context.Canceled {
						return
					}
					log.Printf("Error receiving message: %v", err)
					dialog.ShowError(errors.New("Failed to establish connection"), window)
					window.SetContent(GetLoginContent(window, ssoClient))
					return
				}

				go func(object *ssov1.Object) {
					t := time.Now()

					if object.Type == "main" {
						updateMainObjectLabels(object, mainObjectLabels)
						mainObjectContainer.Refresh()
					} else {
						updateInfoBox(object, infoBox)
					}
					updatePoints(object, imageWithPoints, pointsMap, &pointsMutex)

					end := time.Now()
					duration := end.Sub(t)
					log.Printf("Time taken: %d ms", duration.Milliseconds())
				}(object)
			}
		}
	}()

	return content
}

func updateMainObjectLabels(object *ssov1.Object, labels map[string]*widget.Label) {
	labels["coords"].SetText("Координаты: (X: " + strconv.FormatFloat(float64(object.X), 'f', 2, 32) + ", Y: " + strconv.FormatFloat(float64(object.Y), 'f', 2, 32) + ")")
	labels["angle"].SetText("Курс корабля: " + strconv.FormatFloat(float64(object.Angel), 'f', 2, 32))
	labels["speed"].SetText("Скорость корабля: " + strconv.FormatFloat(float64(object.Speed*2), 'f', 2, 32))
}

func updateInfoBox(object *ssov1.Object, infoBox *fyne.Container) {
	for i, item := range infoBox.Objects {
		if box, ok := item.(*fyne.Container); ok {
			if label, ok := box.Objects[0].(*widget.Label); ok && label.Text == "Объект: "+object.Name {
				if !object.InRangeShip {
					infoBox.Objects = append(infoBox.Objects[:i], infoBox.Objects[i+1:]...)
					infoBox.Refresh()
					return
				}

				labels := box.Objects
				labels[0].(*widget.Label).SetText("Объект: " + object.Name)
				labels[1].(*widget.Label).SetText("Тип: " + object.Type)
				labels[2].(*widget.Label).SetText("Координаты: (X:" + strconv.FormatFloat(float64(object.X), 'f', 2, 32) + ", Y: " + strconv.FormatFloat(float64(object.Y), 'f', 2, 32) + ")")
				labels[3].(*widget.Label).SetText("Курс корабля: " + strconv.FormatFloat(float64(object.Angel), 'f', 2, 32))
				labels[4].(*widget.Label).SetText("Скорость корабля: " + strconv.FormatFloat(float64(object.Speed), 'f', 2, 32))
				infoBox.Refresh()
				return
			}
		}
	}

	if object.InRangeShip {
		newBox := container.NewGridWithRows(1,
			widget.NewLabel("Объект: "+object.Name),
			widget.NewLabel("Тип: "+object.Type),
			widget.NewLabel("Координаты: (X:"+strconv.FormatFloat(float64(object.X), 'f', 2, 32)+", Y: "+strconv.FormatFloat(float64(object.Y), 'f', 2, 32)+")"),
			widget.NewLabel("Курс корабля: "+strconv.FormatFloat(float64(object.Angel), 'f', 2, 32)),
			widget.NewLabel("Скорость корабля: "+strconv.FormatFloat(float64(object.Speed), 'f', 2, 32)),
		)
		infoBox.Add(newBox)
		infoBox.Refresh()
	}
}

func updatePoints(object *ssov1.Object, imageWithPoints *fyne.Container, pointsMap map[string]*PointInfo, pointsMutex *sync.Mutex) {
	var originalImage *canvas.Image

	pointColor := color.RGBA{0, 0, 0, 255}
	switch object.Type {
	case "Враг":
		pointColor = color.RGBA{153, 0, 0, 255}
	case "Союзник":
		pointColor = color.RGBA{0, 102, 0, 255}
	case "main":
		imagePath := "pic/radar/okrug.png"
		originalImage = canvas.NewImageFromFile(imagePath)
		originalImage.FillMode = canvas.ImageFillContain
		originalImage.SetMinSize(fyne.NewSize(40, 40))
		pointColor = color.RGBA{0, 0, 0, 255}
	}

	var pointCon *fyne.Container
	var line *canvas.Line

	if pi, exists := pointsMap[object.Name]; exists {
		pi.X = object.X
		pi.Y = object.Y
		pi.Angel = object.Angel
		pointCon = pi.Container
		line = pi.Line
	} else {
		point := canvas.NewCircle(pointColor)
		point.Resize(fyne.NewSize(20, 20))
		point.Move(fyne.NewPos(object.X, object.Y))

		line = canvas.NewLine(color.Black)
		line.Position1 = fyne.NewPos(object.X, object.Y)
		line.Position2 = fyne.NewPos(object.X+50*float32(math.Sin(float64(object.Angel)*math.Pi/180)), object.Y-50*float32(math.Cos(float64(object.Angel)*math.Pi/180)))
		line.StrokeWidth = 2

		if object.Type == "main" {
			pointCon = container.NewStack(
				container.NewCenter(
					container.New(layout.NewGridWrapLayout(fyne.NewSize(20, 20)), point),
					container.New(layout.NewGridWrapLayout(fyne.NewSize(350, 350)), originalImage),
				),
			)
		} else {
			pointCon = container.NewStack(
				container.NewCenter(
					container.New(layout.NewGridWrapLayout(fyne.NewSize(20, 20)), point),
				),
			)
		}

		pointsMap[object.Name] = &PointInfo{
			X: object.X, Y: object.Y, Color: pointColor, Container: pointCon, Line: line, Angel: object.Angel,
		}

		imageWithPoints.Add(container.NewWithoutLayout(pointCon, line))
	}

	go func() {
		pointsMutex.Lock()
		defer pointsMutex.Unlock()

		pointCon.Move(fyne.NewPos(object.X, object.Y))
		line.Position1 = fyne.NewPos(object.X, object.Y)
		line.Position2 = fyne.NewPos(object.X+50*float32(math.Sin(float64(object.Angel)*math.Pi/180)), object.Y-50*float32(math.Cos(float64(object.Angel)*math.Pi/180)))
		imageWithPoints.Refresh()
	}()
}

func saveImageWithPoint(outputPath string, imagePath string, pointsMap map[string]*PointInfo) {
	file, err := os.Open(imagePath)
	if err != nil {
		fyne.LogError("Failed to load image", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fyne.LogError("Failed to decode image", err)
		return
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, image.Point{}, draw.Over)

	for _, info := range pointsMap {
		lineEndX := info.X + 50*float32(math.Sin(float64(info.Angel)*math.Pi/180))
		lineEndY := info.Y - 50*float32(math.Cos(float64(info.Angel)*math.Pi/180))
		drawLine(newImg, int(info.X), int(info.Y), int(lineEndX), int(lineEndY), color.Black)

		for dx := -5; dx <= 5; dx++ {
			for dy := -5; dy <= 5; dy++ {
				if dx*dx+dy*dy <= 25 {
					newImg.Set(int(info.X)+dx, int(info.Y)+dy, info.Color)
				}
			}
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		fyne.LogError("Failed to create file", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, newImg)
	if err != nil {
		fyne.LogError("Failed to save image", err)
	}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := x2 - x1
	dy := y2 - y1

	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}

	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)

	x := float64(x1)
	y := float64(y1)

	for i := 0; i <= steps; i++ {
		img.Set(int(x), int(y), col)
		x += xInc
		y += yInc
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
