package main

import (
	"fmt"
	"sync"

	"github.com/jroimartin/gocui"
)

// GameWindowManager represents the display window
type GameWindowManager struct {
	maxViewX, maxViewY int
	cursorX, cursorY   int
	offsetX, offsetY   int
	game               *Game
	g                  *gocui.Gui
	mainMenu           bool
	tick               int

	mu    sync.Mutex
	ghost Structure
}

// Update the interface
func (d *GameWindowManager) Update() {
	d.g.Update(d.Layout)
}

// NewGameWindowManager will initialize the game display windows
func NewGameWindowManager(game *Game, g *gocui.Gui) *GameWindowManager {
	var d *GameWindowManager
	d = new(GameWindowManager)
	d.maxViewX = 100
	d.maxViewY = 80
	d.game = game
	d.g = g
	d.mainMenu = false
	d.tick = 0

	mapHeight := len(game.WorldMap)
	if mapHeight < d.maxViewY-2 {
		return nil
	}

	mapWidth := len(game.WorldMap)
	if mapWidth < d.maxViewX-2 {
		return nil
	}

	g.Mouse = false
	g.Cursor = true

	g.SetManagerFunc(d.Layout)
	initKeybindings(d, g)

	return d
}

func (w *GameWindowManager) Layout(g *gocui.Gui) error {

	w.tick++
	if w.mainMenu {
		createMainMenuLayout(w, g)
	} else {
		createGameLayout(w, g)
	}

	return nil

}

func createGameLayout(w *GameWindowManager, g *gocui.Gui) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	g.Cursor = true

	maxX, maxY := g.Size()

	if w.maxViewX < maxX {
		maxX = w.maxViewX
	}

	if w.maxViewY < maxY {
		maxY = w.maxViewY
	}

	worldY := len(w.game.WorldMap)
	worldX := len(w.game.WorldMap[0])

	v, err := g.SetView("Map", 0, 0, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetCurrentView("Map"); err != nil {
			return err
		}

		if _, err := g.SetViewOnTop("Map"); err != nil {
			return err
		}

		v.Clear()
		v.Title = "World"

		worldMaxY := w.offsetY + maxY - 2
		worldMaxX := w.offsetX + maxX - 2

		adjustY := worldY - worldMaxY
		if adjustY < 0 {
			w.offsetY += adjustY
		}

		adjustX := worldX - worldMaxX
		if adjustX < 0 {
			w.offsetX += adjustX
		}

		for i := w.offsetY; i < worldMaxY; i++ {
			for j := w.offsetX; j < worldMaxX; j++ {
				if w.ghost != nil && w.cursorY == i && w.cursorX == j {
					fmt.Fprintf(v, "%s", w.ghost.Show())
				} else {
					fmt.Fprintf(v, "%s", w.game.WorldMap[i][j].Show())
				}

			}
			fmt.Fprintln(v, "")
		}
	} else {
		if w.cursorY > maxY-2 || w.cursorX > maxX-2 {
			w.cursorY = 0
			w.cursorX = 0
		}
	}
	v.SetCursor(w.cursorX, w.cursorY)

	return nil
}

func createMainMenuLayout(w *GameWindowManager, g *gocui.Gui) error {
	g.Cursor = false

	maxX, maxY := g.Size()

	v, err := g.SetView("MainMenu", 0, 0, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop("MainMenu"); err != nil {
			return err
		}

		v.Title = "Main Menu"
		v.Clear()

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

		for _, line := range mat {
			for _, cell := range line {
				fmt.Fprintf(v, "\033[3%d;7m \033[0m", cell)
			}
			fmt.Fprintln(v)
		}
	}

	v, err = g.SetView("SelectionMenu", 24, 14, 40, 18)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetCurrentView("SelectionMenu"); err != nil {
			return err
		}

		if _, err := g.SetViewOnTop("SelectionMenu"); err != nil {
			return err
		}

		v.Frame = false
		v.Clear()

		fmt.Fprintf(v, "[N]ew game\n")
		fmt.Fprintf(v, "[C]ontinue game\n")
		fmt.Fprintf(v, "[S]ettings\n")
	}

	v, err = g.SetView("ConveyorBeltAnimation", 24, 18, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetCurrentView("ConveyorBeltAnimation"); err != nil {
			return err
		}

		if _, err := g.SetViewOnTop("ConveyorBeltAnimation"); err != nil {
			return err
		}

		v.Frame = false
		v.Clear()

		mat := [][]int{
			{1, 1, 0, 3, 3, 3, 3, 0, 4, 4, 4, 4, 0, 1, 1},
			{1, 1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1},
			{1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1, 1},
			{3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1, 1, 1},
			{1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1, 1},
			{1, 1, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 1, 1, 1},
			{1, 1, 0, 3, 3, 3, 3, 0, 4, 4, 4, 4, 0, 1, 1},
		}

		offset := w.tick % len(mat[0])
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
	}

	return nil
}

func initKeybindings(w *GameWindowManager, g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyBackspace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			w.mainMenu = !w.mainMenu
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Map", gocui.KeyArrowDown, gocui.ModNone,
		moveCursor(w, 0, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeyArrowUp, gocui.ModNone,
		moveCursor(w, 0, -1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeyArrowLeft, gocui.ModNone,
		moveCursor(w, -1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeyArrowRight, gocui.ModNone,
		moveCursor(w, 1, 0)); err != nil {
		return err
	}

	if err := g.SetKeybinding("Map", 'b', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.ghost = &Belt{0, nil}; return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", 'd', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.ghost = nil; return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", 'e', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.ghost != nil {
				w.ghost.RotateRight()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", 'q', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.ghost != nil {
				w.ghost.RotateLeft()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.ghost != nil {
				copy := w.ghost.Copy()
				w.game.PlaceBuilding(w.offsetY+w.cursorY, w.offsetX+w.cursorX, copy)
			}
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func moveCursor(d *GameWindowManager, dx, dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cx, cy := d.cursorX, d.cursorY
		maxX, maxY := v.Size()

		worldY := len(d.game.WorldMap)
		worldX := len(d.game.WorldMap[0])

		newCX := cx + dx
		if newCX < 0 {
			if newCX+d.offsetX >= 0 {
				d.offsetX = newCX + d.offsetX
				newCX = 0
			} else {
				d.offsetX = 0
				newCX = 0
			}
		} else if newCX >= maxX {
			if newCX+d.offsetX < worldX {
				d.offsetX = d.offsetX + newCX - maxX + 1
			} else {
				d.offsetX = worldX - maxX
			}
			newCX = maxX - 1
		}

		newCY := cy + dy
		if newCY < 0 {
			if newCY+d.offsetY >= 0 {
				d.offsetY = newCY + d.offsetY
				newCY = 0
			} else {
				d.offsetY = 0
				newCY = 0
			}
		} else if newCY >= maxY {
			if newCY+d.offsetY < worldY {
				d.offsetY = d.offsetY + newCY - maxY + 1
			} else {
				d.offsetY = worldY - maxY
			}
			newCY = maxY - 1
		}

		d.cursorX = newCX
		d.cursorY = newCY
		return nil
	}
}
