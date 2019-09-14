package gui

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func StartGUI() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}

func testDrawPianoRoll(canvas Canvas) {

}
