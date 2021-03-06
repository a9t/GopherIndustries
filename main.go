package main

import (
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

// GameWindowManager WindowManager for the game
type GameWindowManager struct {
	minX, minY int

	errorWindow Window
	topWindow   Window
}

// SetTopWindow specifies which window is on top to be displayed
func (m *GameWindowManager) SetTopWindow(w Window) {
	m.topWindow = w
}

// Layout displays the top window or the ErrorWindow if the display area is too small
func (m *GameWindowManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX < m.minX || maxY < m.minY {
		return m.errorWindow.Layout(g)
	}

	return m.topWindow.Layout(g)
}

func newGameWindowManager(g *gocui.Gui) (*GameWindowManager, *GameWindow) {
	var m GameWindowManager

	m.minX = 60
	m.minY = 20

	m.errorWindow = &ErrorWindow{"Window too small to display game"}

	gameWindow := NewGameWindow(&m)
	settingsWindow := NewSettingsWindow(&m)
	mainMenuWindow := NewPrimaryMenuWindow(&m, gameWindow, settingsWindow)

	m.SetTopWindow(mainMenuWindow)
	g.SetManagerFunc(m.Layout)

	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		})
	g.SetKeybinding("", gocui.KeyBackspace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			gameWindow.SetRunning(false)
			m.SetTopWindow(mainMenuWindow)
			return nil
		})
	return &m, gameWindow
}

func main() {
	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	m, gw := newGameWindowManager(g)

	go uiLoop(m, g)
	go logicLoop(gw)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func uiLoop(m *GameWindowManager, g *gocui.Gui) {
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	for {
		<-ticker.C
		g.Update(m.Layout)
	}
}

func logicLoop(w *GameWindow) {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	for {
		<-ticker.C
		w.Tick()
	}
}
