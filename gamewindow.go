package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type state struct {
	ghost Structure
}

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
	s := new(state)

	var w GameWindow
	w.manager = manager

	var infoWidget InfoWidget
	infoWidget.name = "Info"
	infoWidget.width = 20
	infoWidget.height = 10
	infoWidget.s = s

	var gameMapWidget GameMapWidget
	gameMapWidget.name = "GameMap"
	gameMapWidget.maxViewX = 80
	gameMapWidget.maxViewY = 80
	gameMapWidget.reservedX = infoWidget.width
	gameMapWidget.s = s

	w.widgets = append(w.widgets, &gameMapWidget)
	w.widgets = append(w.widgets, &infoWidget)

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

// InfoWidget a GameWidget that displays general information about the cursor
type InfoWidget struct {
	name   string
	width  int
	height int

	game *Game
	s    *state
}

// Layout displays the InfoWidget
func (w *InfoWidget) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	cursorX, cursorY := w.game.GetCursor()
	displayY := len(w.game.WorldMap) - cursorY

	v, err := g.SetView(w.name, maxX-w.width, 0, maxX-1, w.height)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		v.Title = fmt.Sprintf("%s - %d:%d", w.name, cursorX+1, displayY)
	}

	v.Clear()

	if w.s.ghost != nil {
		var structureName string
		switch w.s.ghost.(type) {
		case *Belt:
			structureName = "belt"
		case *Chest:
			structureName = "chest"
		case *Extractor:
			structureName = "extractor"
		default:
			structureName = "unknown"
		}

		fmt.Fprintf(v, "Placing: %s\n", structureName)
		fmt.Fprint(v, "q or r - rotate\n")
		fmt.Fprint(v, "d      - cancel\n")
		fmt.Fprint(v, "SPACE  - place\n")
		fmt.Fprint(v, "\n↑←↓→ - move\n")

		return nil
	}

	structure, _, _ := w.game.GetStructureAt(cursorY, cursorX)
	if structure == nil {
		switch r := w.game.WorldMap[cursorY][cursorX].(type) {
		case *RawResource:
			if r.amount > 0 {
				fmt.Fprintf(v, "Resource %d\n\n", r.amount)
			} else {
				fmt.Fprint(v, "Empty tile\n\n")
			}
		}

		fmt.Fprint(v, "r - extractor\n")
		fmt.Fprint(v, "b - belt\n")
		fmt.Fprint(v, "c - chest\n")
		fmt.Fprint(v, "\n↑←↓→ - move\n")

		return nil
	}

	var structureName string
	switch structure.(type) {
	case *Belt:
		structureName = "belt"
	case *Chest:
		structureName = "chest"
	case *Extractor:
		structureName = "extractor"
	default:
		structureName = "unknown"
	}
	fmt.Fprintf(v, "Structure: %s\n", structureName)
	fmt.Fprint(v, "\n↑←↓→ - move\n")

	return nil
}

// SetGame sets the Game associated with InfoWidget
func (w *InfoWidget) SetGame(game *Game) {
	w.game = game
}

// GameMapWidget a GameWidget that displays the game map
type GameMapWidget struct {
	name                 string
	maxViewX, maxViewY   int
	reservedX, reservedY int

	game *Game

	offsetX, offsetY int
	s                *state
}

// SetGame sets the Game associated with the GameMapWidget
func (w *GameMapWidget) SetGame(game *Game) {
	w.s.ghost = nil
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
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		err = w.initBindings(g)
		if err == nil {
			return err
		}
	}

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

	cursorX, cursorY := w.game.GetCursor()

	if w.s.ghost != nil {
		ghost = w.s.ghost.Tiles()

		ghostHeight = len(ghost)
		ghostWidth = len(ghost[0])

		maxGhostX, maxGhostY := cursorX+ghostWidth, cursorY+ghostHeight

		if !w.game.WithinBounds(maxGhostX, maxGhostY) {
			mode = DisplayModeGhostInvalid
		} else {
			// check overlap
			for i, tileStructures := range ghost {
				for j, tileStructure := range tileStructures {
					if tileStructure == nil {
						continue
					}
					ghostMap[tileStructure] = tileStructure

					switch w.game.WorldMap[i+cursorY][j+cursorX].(type) {
					case StructureTile:
						mode = DisplayModeGhostInvalid
						break
					}
				}
				if mode == DisplayModeGhostInvalid {
					break
				}
			}
		}
	} else {
		switch selectTile := w.game.WorldMap[cursorY][cursorX].(type) {
		case StructureTile:
			structure := selectTile.Group()

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
			if w.s.ghost != nil &&
				cursorY <= i &&
				i < cursorY+ghostHeight &&
				cursorX <= j &&
				j < cursorX+ghostWidth {

				fmt.Fprintf(v, "%s", ghost[i-cursorY][j-cursorX].Display(mode))
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
		func(g *gocui.Gui, v *gocui.View) error { w.s.ghost = NewBelt(); return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'r', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.s.ghost = NewExtractor(); return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error { w.s.ghost = NewChest(); return nil }); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'd', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.ghost != nil {
				w.s.ghost = nil
				return nil
			}

			x, y := w.game.GetCursor()
			w.game.RemoveStructure(y, x)
			return nil

		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'e', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.ghost != nil {
				w.s.ghost.RotateRight()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'q', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.ghost != nil {
				w.s.ghost.RotateLeft()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.ghost != nil {
				copy := w.s.ghost.CopyStructure()
				x, y := w.game.GetCursor()
				w.game.PlaceStructure(y, x, copy)
			}
			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (w *GameMapWidget) move(dx, dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cx, cy := w.game.GetCursor()
		cx += dx
		cy += dy

		if !w.game.WithinBounds(cx, cy) {
			// do nothing if would be moved outside of bounds
			return nil
		}

		maxX, maxY := v.Size()
		if cx < w.offsetX {
			w.offsetX = cx
		} else if cx >= w.offsetX+maxX {
			w.offsetX = cx - maxX + 1
		}

		if cy < w.offsetY {
			w.offsetY = cy
		} else if cy >= w.offsetY+maxY {
			w.offsetY = cy - maxY + 1
		}

		w.game.MoveCursor(dx, dy)

		return nil
	}
}
