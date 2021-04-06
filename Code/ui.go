package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var App = app.New()

var logLabel = widget.NewLabel("[LOG] OK")
var progressLabel = widget.NewLabel("(0/0)")

var dogNumber = widget.NewEntry()
var content = container.New(layout.NewMaxLayout())

var waitTime = 10

// Buttons
var github = widget.NewButton("Github", openGithub)
var removeImg = widget.NewButton("Remove Images", removeImages)
var begin = widget.NewButton("Begin", beginFront)
var buttonGrid = container.NewGridWithColumns(3)

var gridAvailable = true // Set to false for a [Please wait] to replace the button grid. Set to true to cancel this.

func gridLoop(){
	lastGrid := true
	for true {
		if gridAvailable != lastGrid {
			if gridAvailable == true {
				buttonGrid.Objects = []fyne.CanvasObject{begin, removeImg, github}
			} else {
				progressText := widget.NewLabel("[Please Wait]")
				progressText.TextStyle = fyne.TextStyle{Monospace: true}
				buttonGrid.Objects = []fyne.CanvasObject{layout.NewSpacer(), progressText, layout.NewSpacer()}
			}
			lastGrid = gridAvailable
		}
	}
}

func themeChange(dark bool) {
	if dark == true {
		App.Settings().SetTheme(theme.DarkTheme())
	} else {
		App.Settings().SetTheme(theme.LightTheme())
	}
}

func speedChange(speed string){
	switch speed {
	case "Moderate":
		waitTime = 20
	case "Fast":
		waitTime = 10
	case "Blitz":
		waitTime = 0
	}
}

func beginFront(){ // Does front-end modifications.
	err := dogNumber.Validate()
	if err != nil {
		setLog("Invalid number of dogs.", "0", "0")
		return
	}

	gridAvailable = false
	beginBack()
	gridAvailable = true
}

func main() {
	window := App.NewWindow("Dog Downloader")

	dogNumber.Validator = validation.NewRegexp("^[0-9]+$", "Must be a number.")
	dogNumber.Text = "100"

	themeCheck := widget.NewCheck("Dark Theme", themeChange)
	themeCheck.Checked = true

	speedWid := widget.NewSelect([]string{"Moderate", "Fast", "Blitz"}, speedChange)
	speedWid.Selected = "Fast"

	title := widget.NewLabel("Welcome to the Dog Downloader!")
	title.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	logLabel.TextStyle = fyne.TextStyle{Monospace: true}
	progressLabel.TextStyle = fyne.TextStyle{Monospace: true}

	buttonGrid.Objects = []fyne.CanvasObject{begin, removeImg, github}

	content = container.NewVBox(
		container.NewHBox(title, layout.NewSpacer(), themeCheck),
		container.New(layout.NewFormLayout(),
			widget.NewLabel("How many dog pictures do you want?"),
			dogNumber,
			widget.NewLabel("Choose a speed:"),
			speedWid,
			),
		layout.NewSpacer(),
		buttonGrid,
		container.NewHBox(logLabel, layout.NewSpacer(), progressLabel),
	)

	window.Resize(fyne.NewSize(40, 240))
	window.SetContent(content)

	go gridLoop()
	window.ShowAndRun()
}
