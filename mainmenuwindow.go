package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// PrimaryMenuWindow a Window that handles the main game menu
type PrimaryMenuWindow struct {
	manager WindowManager
	tick    int
	widgets []Widget
}

// NewPrimaryMenuWindow creates a new MainMenuWindow
func NewPrimaryMenuWindow(manager WindowManager, gw *GameWindow, sw *SettingsWindow) *PrimaryMenuWindow {
	var w PrimaryMenuWindow
	w.manager = manager

	w.widgets = append(w.widgets, newMascotWidget("Mascot", 1, 1))
	w.widgets = append(w.widgets, newConveyorBeltWidget("ConveyorBelt", 24, 19))
	w.widgets = append(w.widgets, newPrimaryMenuWidget("MainMenu", 24, 14, manager, gw, sw))

	return &w
}

// Layout displays the PrimaryMenuWindow
func (w *PrimaryMenuWindow) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = false

	v, err := g.SetView("PrimaryMenuWindow", 0, 0, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop("PrimaryMenuWindow"); err != nil {
			return err
		}

		v.Title = "Main Menu"
	}

	for _, widget := range w.widgets {
		widget.Layout(g)
	}

	return nil
}

// MascotWidget a Widget that displays an image of the Golang Gopher wearing a hat
type MascotWidget struct {
	name string
	x, y int
}

func newMascotWidget(name string, x, y int) *MascotWidget {
	return &MascotWidget{name: name, x: x, y: y}
}

// Layout displays the MascotWidget
func (w *MascotWidget) Layout(g *gocui.Gui) error {
	mat := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 0, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0},
		{0, 0, 0, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 0, 0, 0},
		{0, 0, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 0, 0},
		{0, 0, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 0, 0},
		{0, 0, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 0, 0},
		{0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 0},
		{0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0},
		{6, 6, 6, 6, 7, 7, 7, 7, 6, 6, 6, 6, 6, 7, 7, 7, 7, 6, 6, 6, 6},
		{6, 0, 6, 7, 7, 7, 7, 7, 7, 6, 6, 6, 7, 7, 7, 7, 7, 7, 6, 0, 6},
		{6, 6, 6, 7, 7, 7, 7, 0, 0, 6, 6, 6, 7, 7, 7, 7, 0, 0, 6, 6, 6},
		{0, 0, 6, 7, 7, 7, 0, 0, 0, 6, 6, 6, 7, 7, 7, 0, 0, 0, 6, 0, 0},
		{0, 0, 6, 7, 7, 7, 0, 0, 0, 6, 6, 6, 7, 7, 7, 0, 0, 0, 6, 0, 0},
		{0, 0, 6, 6, 7, 7, 7, 7, 6, 0, 0, 0, 6, 7, 7, 7, 7, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 1, 0, 0, 0, 1, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 1, 1, 1, 1, 1, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 7, 7, 0, 7, 7, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 7, 7, 0, 7, 7, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 7, 7, 0, 7, 7, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0},
		{0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0},
	}

	v, err := g.SetView(w.name, w.x, w.y, w.x+len(mat[0])+2, w.y+len(mat)+2)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		v.Frame = false
		v.Clear()

		for _, line := range mat {
			for _, cell := range line {
				fmt.Fprintf(v, "\033[3%d;7m \033[0m", cell)
			}
			fmt.Fprintln(v)
		}
	} else {
		return err
	}
	return nil
}

// ConveyorBeltWidget a Widget that displays an animated coveyor belt
type ConveyorBeltWidget struct {
	name string
	x, y int
	tick int
}

func newConveyorBeltWidget(name string, x, y int) *ConveyorBeltWidget {
	return &ConveyorBeltWidget{name: name, x: x, y: y, tick: 0}
}

// Layout displays the ConveyorBeltWidget
func (w *ConveyorBeltWidget) Layout(g *gocui.Gui) error {
	mat := [][]int{
		{1, 1, 0, 3, 3, 3, 3, 0, 4, 4, 4, 4, 0, 1, 1},
		{1, 1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1},
		{1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1, 1},
		{3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1, 1, 1},
		{1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1, 1},
		{1, 1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1},
		{1, 1, 0, 3, 3, 3, 3, 0, 4, 4, 4, 4, 0, 1, 1},
	}
	maxX, _ := g.Size()

	v, err := g.SetView(w.name, w.x, w.y, maxX-1, w.y+len(mat)+1)
	if err != nil || err != gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		v.Frame = false
		v.Clear()

		w.tick++
		offset := (w.tick / 10) % len(mat[0])
		bx, _ := v.Size()
		bx -= offset
		repeat := bx/len(mat[0]) + 2

		for _, line := range mat {
			for _, cell := range line[offset:] {
				fmt.Fprintf(v, "\033[3%d;7m \033[0m", cell)
			}

			for i := 0; i < repeat; i++ {
				for _, cell := range line {
					fmt.Fprintf(v, "\033[3%d;7m \033[0m", cell)
				}
			}

			fmt.Fprintln(v)
		}
	} else {
		return err
	}

	return nil
}

// PrimaryMenuWidget a Widget that display the main menu
type PrimaryMenuWidget struct {
	name           string
	x, y           int
	selection      int
	manager        WindowManager
	gameWindow     *GameWindow
	settingsWindow *SettingsWindow
}

func newPrimaryMenuWidget(name string, x, y int, manager WindowManager, gameWindow *GameWindow, settingsWindow *SettingsWindow) *PrimaryMenuWidget {
	return &PrimaryMenuWidget{name: name, x: x, y: y, selection: 0, gameWindow: gameWindow, manager: manager, settingsWindow: settingsWindow}
}

// Layout displays the PrimaryMenuWidget
func (w *PrimaryMenuWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+16, w.y+5)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}

		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		if err == gocui.ErrUnknownView {
			v.Frame = false

			if err := g.SetKeybinding(w.name, 'n', gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					game := GenerateGame(120, 100)
					w.gameWindow.SetGame(game)
					w.manager.SetTopWindow(w.gameWindow)

					return nil
				}); err != nil {
				return err
			}
			if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					if w.gameWindow.HasGame() {
						w.manager.SetTopWindow(w.gameWindow)
					}
					return nil
				}); err != nil {
				return err
			}
			if err := g.SetKeybinding(w.name, 's', gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					w.manager.SetTopWindow(w.settingsWindow)

					return nil
				}); err != nil {
				return err
			}
			if err := g.SetKeybinding(w.name, 'q', gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					return gocui.ErrQuit
				}); err != nil {
				return err
			}
		}

		v.Clear()

		fmt.Fprintf(v, "[N]ew game\n")
		if w.gameWindow.HasGame() {
			fmt.Fprintf(v, "[C]ontinue game\n")
		}
		fmt.Fprintf(v, "[S]ettings\n")
		fmt.Fprintf(v, "[Q]uit\n")
	}

	return nil
}
