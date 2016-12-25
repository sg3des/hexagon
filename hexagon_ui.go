package main

import "github.com/mattn/go-gtk/gtk"

type UI struct {
	window   *gtk.Window
	vbox     *gtk.VBox
	notebook *gtk.Notebook
}

func CreateUI() *UI {
	ui := &UI{}

	ui.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ui.window.SetDefaultSize(600, 500)
	ui.vbox = gtk.NewVBox(false, 0)
	ui.notebook = gtk.NewNotebook()
	ui.vbox.PackStart(ui.notebook, true, true, 0)

	ui.window.Add(ui.vbox)
	ui.window.Connect("destroy", ui.Quit)
	// ui.window.Connect("check-resize", ui.windowResize)

	ui.window.ShowAll()

	return ui
}

func (ui *UI) Quit() {
	gtk.MainQuit()
}
