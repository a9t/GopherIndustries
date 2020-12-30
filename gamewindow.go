package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// GameWindow a Window that manages all the GameWidget-s
type GameWindow struct {
	manager WindowManager

	game    *Game
	running bool

	mainWindow Window

	widgets []GameWidget
}

// NewGameWindow creates a new GameWindow
func NewGameWindow(manager WindowManager) *GameWindow {
	var w GameWindow
	w.manager = manager

	var structureWidget StructureWidget
	structureWidget.name = "Structure"
	structureWidget.width = 15
	structureWidget.height = 5

	var gameMapWidget GameMapWidget
	gameMapWidget.name = "GameMap"
	gameMapWidget.maxViewX = 80
	gameMapWidget.maxViewY = 80
	gameMapWidget.reservedX = structureWidget.width
	gameMapWidget.structureWidget = &structureWidget

	w.widgets = append(w.widgets, &gameMapWidget)
	w.widgets = append(w.widgets, &structureWidget)

	return &w
}

// Tick advances the game state
func (w *GameWindow) Tick() {
	if w.running {
		w.game.Tick()
	}
}

// SetRunning indicates if the game is running or not, having ticks pass through or not
func (w *GameWindow) SetRunning(state bool) {
	w.running = state
}

// SetGame sets a game to be displayed in the GameWindow
func (w *GameWindow) SetGame(g *Game) {
	w.game = g
	for _, widget := range w.widgets {
		widget.SetGame(g)
	}
}

// HasGame indicates if there is a Game associated with the GameWindow
func (w *GameWindow) HasGame() bool {
	return w.game != nil
}

// Layout displays the GameWindow
func (w *GameWindow) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = false

	v, err := g.SetView("GameWindow", 0, 0, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop("GameWindow"); err != nil {
			return err
		}

		v.Title = "Gopher Industries"
	} else {
		return err
	}

	for _, widget := range w.widgets {
		widget.Layout(g)
	}

	return nil
}

// StructureWidget a GameWidget that display information about a Structure
type StructureWidget struct {
	name   string
	width  int
	height int
	s      Structure
}

// Layout displays the StructureWidget
func (w *StructureWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	v, err := g.SetView(w.name, maxX-w.width, 0, maxX-1, w.height)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		v.Title = w.name
	}

	v.Clear()

	var structureType string
	switch s := w.s.(type) {
	case *Belt:
		structureType = fmt.Sprintf("B %d,%d,%d:%d,%d,%d", s.inputs[0].x, s.inputs[0].y, s.inputs[0].d, s.outputs[0].x, s.outputs[0].y, s.outputs[0].d)
	case *Extractor:
		structureType = fmt.Sprintf("E %d,%d,%d", s.outputs[0].x, s.outputs[0].y, s.outputs[0].d)
	default:
		structureType = " - "
	}
	fmt.Fprintf(v, "%s\n", structureType)

	if w.s != nil {
		product, _ := w.s.CanRetrieveProduct()
		var productName string
		if product == nil {
			productName = "-"
		} else {
			productName = product.name
		}
		fmt.Fprintf(v, "Prod: %s\n", productName)
	}

	return nil
}

// SetGame sets the Game associated with the GameMapWidget
func (w *StructureWidget) SetGame(game *Game) {
}

// GameMapWidget a GameWidget that displays the game map
type GameMapWidget struct {
	name                 string
	maxViewX, maxViewY   int
	reservedX, reservedY int

	game            *Game
	structureWidget *StructureWidget

	cursorX, cursorY int
	offsetX, offsetY int
	ghost            Structure
}

// SetGame sets the Game associated with the GameMapWidget
func (w *GameMapWidget) SetGame(game *Game) {
	w.game = game
}

// Layout displays the GameMapWidget
func (w *GameMapWidget) Layout(g *gocui.Gui) error {
	g.Cursor = false

	maxX, maxY := g.Size()

	maxX = maxX - w.reservedX
	maxY = maxY - w.reservedY

	offsetX := 0
	if maxX > w.maxViewX {
		offsetX = (maxX - w.maxViewX) / 2
		maxX = w.maxViewX
	}

	worldY := len(w.game.WorldMap)
	worldX := len(w.game.WorldMap[0])

	offsetY := 0
	if maxY > w.maxViewY {
		offsetY = (maxY - w.maxViewY) / 2
		maxY = w.maxViewY
	}

	v, err := g.SetView(w.name, offsetX, offsetY, offsetX+maxX-1, offsetY+maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}

		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		v.Clear()
		v.Frame = false

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

		ghostMap := make(map[Tile]Tile)
		selectedMap := make(map[Tile]Tile)

		ghostHeight := -1
		ghostWidth := -1
		var ghost [][]StructureTile
		mode := DisplayModeGhostValid

		if w.ghost != nil {
			w.structureWidget.s = w.ghost
			ghost = w.ghost.Tiles()

			ghostHeight = len(ghost)
			ghostWidth = len(ghost[0])

			// check overlap
			for i, tileStructures := range ghost {
				for j, tileStructure := range tileStructures {
					if tileStructure == nil {
						continue
					}
					ghostMap[tileStructure] = tileStructure

					if w.offsetY+i+w.cursorY >= worldY {
						mode = DisplayModeGhostInvalid
						break
					}

					if w.offsetX+j+w.cursorX >= worldX {
						mode = DisplayModeGhostInvalid
						break
					}

					switch w.game.WorldMap[w.offsetY+i+w.cursorY][w.offsetX+j+w.cursorX].(type) {
					case StructureTile:
						mode = DisplayModeGhostInvalid
						break
					}
				}
				if mode == DisplayModeGhostInvalid {
					break
				}
			}
		} else {
			switch selectTile := w.game.WorldMap[w.offsetY+w.cursorY][w.offsetX+w.cursorX].(type) {
			case StructureTile:
				structure := selectTile.Group()
				w.structureWidget.s = structure

				for _, tiles := range structure.Tiles() {
					for _, tile := range tiles {
						selectedMap[tile] = tile
					}
				}
			default:
				selectedMap[selectTile] = selectTile
			}
		}

		for i := w.offsetY; i < worldMaxY; i++ {
			for j := w.offsetX; j < worldMaxX; j++ {
				if w.ghost != nil &&
					w.cursorY+w.offsetY <= i &&
					i < w.cursorY+w.offsetY+ghostHeight &&
					w.cursorX+w.offsetX <= j &&
					j < w.cursorX+w.offsetX+ghostWidth {

					fmt.Fprintf(v, "%s", ghost[i-w.cursorY-w.offsetY][j-w.cursorX-w.offsetX].Display(mode))
				} else {
					if _, ok := selectedMap[w.game.WorldMap[i][j]]; ok {
						fmt.Fprintf(v, "%s", w.game.WorldMap[i][j].Display(DisplayModeMapSelected))
					} else {
						fmt.Fprintf(v, "%s", w.game.WorldMap[i][j].Display(DisplayModeMap))
					}
				}

			}
			fmt.Fprintln(v, "")
		}

		if err == gocui.ErrUnknownView {
			err = w.initBindings(g)
			if err == nil {
				return err
			}
		}
	}

	return nil
}

func (w *GameMapWidget) initBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
		w.move(0, 1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
		w.move(0, -1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone,
		w.move(-1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone,
		w.move(1, 0)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'b', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.ghost = NewBelt(); return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'r', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.ghost = NewExtractor(); return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.ghost = NewChest(); return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'd', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.ghost = nil; return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'e', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.ghost != nil {
				w.ghost.RotateRight()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'q', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.ghost != nil {
				w.ghost.RotateLeft()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.ghost != nil {
				copy := w.ghost.CopyStructure()
				w.game.PlaceBuilding(w.offsetY+w.cursorY, w.offsetX+w.cursorX, copy)
			}
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (w *GameMapWidget) move(dx, dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cx, cy := w.cursorX, w.cursorY
		maxX, maxY := v.Size()

		worldY := len(w.game.WorldMap)
		worldX := len(w.game.WorldMap[0])

		newCX := cx + dx
		if newCX < 0 {
			if newCX+w.offsetX >= 0 {
				w.offsetX = newCX + w.offsetX
				newCX = 0
			} else {
				w.offsetX = 0
				newCX = 0
			}
		} else if newCX >= maxX {
			if newCX+w.offsetX < worldX {
				w.offsetX = w.offsetX + newCX - maxX + 1
			} else {
				w.offsetX = worldX - maxX
			}
			newCX = maxX - 1
		}

		newCY := cy + dy
		if newCY < 0 {
			if newCY+w.offsetY >= 0 {
				w.offsetY = newCY + w.offsetY
				newCY = 0
			} else {
				w.offsetY = 0
				newCY = 0
			}
		} else if newCY >= maxY {
			if newCY+w.offsetY < worldY {
				w.offsetY = w.offsetY + newCY - maxY + 1
			} else {
				w.offsetY = worldY - maxY
			}
			newCY = maxY - 1
		}

		w.cursorX = newCX
		w.cursorY = newCY
		return nil
	}
}
