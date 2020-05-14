package main

// I just copied my code from my testing folder, needs to be optimized to work with Aeacus

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var id string

type enterEntry struct {
	widget.Entry
}

func (e *enterEntry) onEnter() {
	fmt.Println(e.Entry.Text)
	e.Entry.SetText("")
}

func newEnterEntry() *enterEntry {
	entry := &enterEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func idGUI() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Enter ID")
	w.CenterOnScreen()

	entry := newEnterEntry()

	w.SetContent(widget.NewVBox(
		widget.NewLabelWithStyle("Enter ID:", fyne.TextAlignCenter, fyne.TextStyle{}),
		entry,
		widget.NewButton("Save", func() {
			id = entry.Text
			a.Quit()
		}),
	))

	w.Resize(fyne.Size{Width: 400})
	w.ShowAndRun()

	fmt.Println(id)
}
