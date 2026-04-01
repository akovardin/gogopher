package main

import (
	"gopher/gopher"

	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

// Server configuration - change these to match your test server
const (
	server = "192.168.0.108"
	port   = 7070
)

var (
	client     *gopher.Client
	controller *Controller
	view       *View
)

func initGame() {
	font, err := pd.Graphics.LoadFont("assets/fonts/JetBrainsMono-ExtraLight")
	if err != nil {
		pd.System.LogToConsole("Error on load font: " + err.Error())
	} else {
		pd.Graphics.SetFont(font)
	}

	pd.System.LogToConsole("Networking Demo: Starting")

	client = gopher.New(pd, server, port)
	view = NewView(pd.Graphics, pd.System)
	controller = NewController(view, client, pd.System)
}

// update is called every frame
func update() int {
	// Clear screen
	pd.Graphics.Clear(pdgo.SolidWhite)

	// render all data
	view.Render()

	// handle input
	_, pushed, _ := pd.System.GetButtonState()
	// Fire on A or B press
	if pushed&pdgo.ButtonB != 0 {
		pd.System.LogToConsole("Press button A")

		// connect to server
		controller.PressB()

	}

	if pushed&pdgo.ButtonA != 0 {
		pd.System.LogToConsole("Press button B")

		// query by selector
		controller.PressA(view.Cursor)
	}

	// Draw FPS
	pd.System.DrawFPS(0, 0)

	return 1
}

// нужно отдельно рендерить интерфейс и отдельно обрабатывать нажатия

func main() {}
