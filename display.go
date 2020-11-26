package main

import (
	"fmt"
	"sync"

	"github.com/jroimartin/gocui"
)

// Display represents the display window
type Display struct {
	maxViewX, maxViewY int
	cursorX, cursorY   int
	offsetX, offsetY   int
	game               *Game
	g                  *gocui.Gui

	mu    sync.Mutex
	ghost Structure
}

// Update the interface
func (d *Display) Update() {
	d.g.Update(createLayout(d))
}

// GetDisplay will initialize the game display windows
func GetDisplay(game *Game, g *gocui.Gui) *Display {
	var d *Display
	d = new(Display)
	d.maxViewX = 100
	d.maxViewY = 80
	d.game = game
	d.g = g

	mapHeight := len(game.WorldMap)
	if mapHeight < d.maxViewY-2 {
		return nil
	}

	mapWidth := len(game.WorldMap)
	if mapWidth < d.maxViewX-2 {
		return nil
	}

	g.Mouse = true
	g.Cursor = true

	g.SetManagerFunc(createLayout(d))
	initKeybindings(d, g)

	return d
}

func createLayout(d *Display) func(g *gocui.Gui) error {

	return func(g *gocui.Gui) error {
		d.mu.Lock()
		defer d.mu.Unlock()

		maxX, maxY := g.Size()

		if d.maxViewX < maxX {
			maxX = d.maxViewX
		}

		if d.maxViewY < maxY {
			maxY = d.maxViewY
		}

		worldY := len(d.game.WorldMap)
		worldX := len(d.game.WorldMap[0])

		v, err := g.SetView("Map", 0, 0, maxX-1, maxY-1)
		if err == nil || err == gocui.ErrUnknownView {
			if _, err := g.SetCurrentView("Map"); err != nil {
				return err
			}
			v.Clear()
			v.Title = "World"

			worldMaxY := d.offsetY + maxY - 2
			worldMaxX := d.offsetX + maxX - 2

			adjustY := worldY - worldMaxY
			if adjustY < 0 {
				d.offsetY += adjustY
			}

			adjustX := worldX - worldMaxX
			if adjustX < 0 {
				d.offsetX += adjustX
			}

			for i := d.offsetY; i < worldMaxY; i++ {
				for j := d.offsetX; j < worldMaxX; j++ {
					if d.ghost != nil && d.cursorY == i && d.cursorX == j {
						fmt.Fprintf(v, "%s", d.ghost.Show())
					} else {
						fmt.Fprintf(v, "%s", d.game.WorldMap[i][j].Show())
					}

				}
				fmt.Fprintln(v, "")
			}
		} else {
			if d.cursorY > maxY-2 || d.cursorX > maxX-2 {
				d.cursorY = 0
				d.cursorX = 0
			}
		}
		v.SetCursor(d.cursorX, d.cursorY)

		return nil
	}
}

func initKeybindings(d *Display, g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Map", gocui.KeyArrowDown, gocui.ModNone,
		moveCursor(d, 0, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeyArrowUp, gocui.ModNone,
		moveCursor(d, 0, -1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeyArrowLeft, gocui.ModNone,
		moveCursor(d, -1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeyArrowRight, gocui.ModNone,
		moveCursor(d, 1, 0)); err != nil {
		return err
	}

	if err := g.SetKeybinding("Map", 'b', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { d.ghost = &Belt{0, nil}; return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", 'd', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { d.ghost = nil; return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", 'e', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if d.ghost != nil {
				d.ghost.RotateRight()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", 'q', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if d.ghost != nil {
				d.ghost.RotateLeft()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Map", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if d.ghost != nil {
				copy := d.ghost.Copy()
				d.game.PlaceBuilding(d.offsetY+d.cursorY, d.offsetX+d.cursorX, copy)
			}
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func moveCursor(d *Display, dx, dy int) func(g *gocui.Gui, v *gocui.View) error {
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
