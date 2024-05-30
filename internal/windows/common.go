package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func newWidget(text string, size fyne.Size) *widget.Entry {
	w := widget.NewEntry()
	w.SetPlaceHolder(text)
	w.Resize(size)

	return w
}
