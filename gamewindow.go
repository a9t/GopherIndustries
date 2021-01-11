package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	stateNavigate int = iota
	stateStructureSelect
	stateStructureGhost
	stateMoveFromInventory
	stateMoveFromStructure
	stateSetRecipe
)

type state struct {
	state int
	ghost Structure
	st    []*Storage
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
	s.st = make([]*Storage, 2)

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

	structureSelectorWidget := newStructureSelectorWidget()
	structureSelectorWidget.name = "Structure"
	structureSelectorWidget.width = 20
	structureSelectorWidget.height = 7
	structureSelectorWidget.offsetY = infoWidget.height + 1
	structureSelectorWidget.s = s

	recipeSelectorWidget := newRecipeSelectorWidget()
	recipeSelectorWidget.name = "Recipes"
	recipeSelectorWidget.width = 20
	recipeSelectorWidget.height = 7
	recipeSelectorWidget.offsetY = infoWidget.height + 1
	recipeSelectorWidget.s = s

	inventoryWidget := newInventoryWidget()
	inventoryWidget.name = "Inventory"
	inventoryWidget.width = 20
	inventoryWidget.height = 6
	inventoryWidget.offsetY = infoWidget.height + 1
	inventoryWidget.activeState = stateMoveFromInventory
	inventoryWidget.s = s

	chestInventoryWidget := newInventoryWidget()
	chestInventoryWidget.name = "Chest"
	chestInventoryWidget.width = 20
	chestInventoryWidget.height = 6
	chestInventoryWidget.offsetY = infoWidget.height + 1 + inventoryWidget.height + 1
	chestInventoryWidget.activeState = stateMoveFromStructure
	chestInventoryWidget.storageIndex = 1
	chestInventoryWidget.s = s

	w.widgets = append(w.widgets, &gameMapWidget)
	w.widgets = append(w.widgets, &infoWidget)
	w.widgets = append(w.widgets, structureSelectorWidget)
	w.widgets = append(w.widgets, inventoryWidget)
	w.widgets = append(w.widgets, chestInventoryWidget)
	w.widgets = append(w.widgets, recipeSelectorWidget)

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

	if w.s.state == stateStructureSelect {
		fmt.Fprintf(v, "Choose structure\n")
		fmt.Fprint(v, "navigate: ↑↓\n")
		fmt.Fprint(v, "cancel  : c\n")
		fmt.Fprint(v, "select  : ˽")

		return nil
	}

	if w.s.state == stateStructureGhost {
		var structureName string
		switch w.s.ghost.(type) {
		case *Belt:
			structureName = "belt"
		case *Chest:
			structureName = "chest"
		case *Extractor:
			structureName = "extractor"
		case *Splitter:
			structureName = "splitter"
		default:
			structureName = "unknown"
		}

		fmt.Fprintf(v, "Placing: %s\n", structureName)
		fmt.Fprint(v, "rotate: qr\n")
		fmt.Fprint(v, "cancel: c\n")
		fmt.Fprint(v, "place : ˽\n")
		fmt.Fprint(v, "move  : ↑←↓→\n")

		return nil
	}

	if w.s.state == stateMoveFromInventory || w.s.state == stateMoveFromStructure {
		if w.s.state == stateMoveFromInventory {
			fmt.Fprintf(v, "Inventory → chest\n")
		} else {
			fmt.Fprintf(v, "Chest → inventory\n")
		}

		fmt.Fprint(v, "navigate: ↑↓\n")
		fmt.Fprint(v, "cancel  : c\n")
		fmt.Fprint(v, "switch  : t\n")
		fmt.Fprint(v, "delete  : d\n")
		fmt.Fprint(v, "transfer: ˽\n")

		return nil
	}

	structure, _, _ := w.game.GetStructureAt(cursorY, cursorX)
	if structure == nil {
		switch r := w.game.WorldMap[cursorY][cursorX].(type) {
		case *RawResource:
			if r.amount > 0 {
				fmt.Fprintf(v, "Resource %d\n", r.amount)
			} else {
				fmt.Fprint(v, "Empty tile\n")
			}
		}

		fmt.Fprint(v, "navigate: ↑←↓→\n")
		fmt.Fprint(v, "add     : a\n")

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
	case *Splitter:
		structureName = "splitter"
	default:
		structureName = "unknown"
	}
	fmt.Fprintf(v, "Structure: %s\n", structureName)
	fmt.Fprint(v, "navigate: ↑←↓→\n")
	fmt.Fprint(v, "delete  : d\n")
	fmt.Fprint(v, "add     : a\n")
	if structureName == "chest" {
		fmt.Fprint(v, "transfer: t\n")
	}

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

	if w.s.state == stateNavigate || w.s.state == stateStructureGhost {
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
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
				j < cursorX+ghostWidth &&
				ghost[i-cursorY][j-cursorX] != nil {

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
	if err := g.SetKeybinding(w.name, 'a', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.state == stateNavigate {
				w.s.state = stateStructureSelect
			}

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 't', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.state != stateNavigate {
				return nil
			}

			x, y := w.game.GetCursor()
			structure, _, _ := w.game.GetStructureAt(y, x)
			switch c := structure.(type) {
			case *Chest:
				w.s.st[0] = w.game.inventory
				w.s.st[1] = c.s
				w.s.state = stateMoveFromInventory
			}

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.state != stateStructureGhost {
				return nil
			}

			w.s.state = stateNavigate
			w.s.ghost = nil

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'd', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.state != stateNavigate {
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
			if w.s.state != stateStructureGhost {
				return nil
			}

			if w.s.ghost != nil {
				w.s.ghost.RotateRight()
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'q', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.state != stateStructureGhost {
				return nil
			}

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

				switch w.s.ghost.(type) {
				case *Factory:
					w.s.ghost = nil
					w.s.state = stateSetRecipe
				}
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

// StructureSelectorWidget a GameWidget that displays available Structures to build
type StructureSelectorWidget struct {
	name    string
	offsetY int
	width   int
	height  int

	game *Game
	s    *state

	position int
	products []*Product
}

func newStructureSelectorWidget() *StructureSelectorWidget {
	w := new(StructureSelectorWidget)
	w.products = make([]*Product, 6)
	w.products[0] = GlobalProductFactory.GetProduct(ProductStructureBelt)
	w.products[1] = GlobalProductFactory.GetProduct(ProductStructureChest)
	w.products[2] = GlobalProductFactory.GetProduct(ProductStructureExtractor)
	w.products[3] = GlobalProductFactory.GetProduct(ProductStructureSplitter)
	w.products[4] = GlobalProductFactory.GetProduct(ProductStructureFactory)
	w.products[5] = GlobalProductFactory.GetProduct(ProductStructureUnderground)

	return w
}

// SetGame sets the Game associated with StructureSelectorWidget
func (w *StructureSelectorWidget) SetGame(game *Game) {
	w.game = game
}

// Layout displays the StructureSelectorWidget
func (w *StructureSelectorWidget) Layout(g *gocui.Gui) error {
	if w.s.state != stateStructureSelect {
		return nil
	}

	maxX, _ := g.Size()

	v, err := g.SetView(w.name, maxX-w.width, w.offsetY, maxX-1, w.offsetY+w.height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		err = w.initBindings(g)
		if err == nil {
			return err
		}
	}

	if _, err := g.SetViewOnTop(w.name); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(w.name); err != nil {
		return err
	}

	v.Title = "Structures"

	v.Clear()

	for i, product := range w.products {
		var prefix string
		if i == w.position {
			prefix = ">"
		} else {
			prefix = " "
		}

		fmt.Fprintf(v, "%s %s\n", prefix, product.name)
	}

	return nil
}

func (w *StructureSelectorWidget) initBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
		w.move(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
		w.move(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			w.s.ghost = w.products[w.position].structure.CopyStructure()
			w.s.state = stateStructureGhost

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			w.s.ghost = nil
			w.s.state = stateNavigate

			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (w *StructureSelectorWidget) move(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		w.position += d + len(w.products)
		w.position %= len(w.products)

		return nil
	}
}

// InventoryWidget a GameWidget that displays the full player inventory for transfer purposes
type InventoryWidget struct {
	name    string
	offsetY int
	width   int
	height  int

	game *Game
	s    *state

	storageIndex int
	activeState  int

	position int
}

func newInventoryWidget() *InventoryWidget {
	w := new(InventoryWidget)

	return w
}

// SetGame sets the Game associated with InventoryWidget
func (w *InventoryWidget) SetGame(game *Game) {
	w.game = game
}

// Layout displays the InventoryWidget
func (w *InventoryWidget) Layout(g *gocui.Gui) error {
	if w.s.state != stateMoveFromInventory && w.s.state != stateMoveFromStructure {
		return nil
	}

	maxX, _ := g.Size()

	v, err := g.SetView(w.name, maxX-w.width, w.offsetY, maxX-1, w.offsetY+w.height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		err = w.initBindings(g)
		if err == nil {
			return err
		}
	}

	if _, err := g.SetViewOnTop(w.name); err != nil {
		return err
	}

	if w.activeState == w.s.state {
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
	}

	v.Title = w.name

	v.Clear()

	storage := w.s.st[w.storageIndex]

	if storage.Size() == 0 {
		fmt.Fprint(v, "Empty inventory\n")
		return nil
	}

	if w.position >= storage.Size() {
		w.position = storage.Size() - 1
	}

	var start int
	uniqueCount := storage.UniqueObjects()
	if uniqueCount-5 > w.position {
		start = w.position
	} else {
		start = uniqueCount - 5
	}

	index := -1
	for _, product := range GlobalProductFactory.cannonicalOrder {
		count, present := storage.objects[product]
		if !present {
			continue
		}

		index++
		if index < start {
			continue
		}

		var prefix string
		if index == w.position {
			if w.activeState == w.s.state {
				prefix = ">"
			} else {
				prefix = " "
			}
		} else {
			prefix = " "
		}

		fmt.Fprintf(v, "%s %3d x %s\n", prefix, count, product.name)
	}

	return nil
}

func (w *InventoryWidget) initBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
		w.move(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
		w.move(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			w.position = 0
			w.s.state = stateNavigate

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 't', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if w.s.state == stateMoveFromInventory {
				w.s.state = stateMoveFromStructure
			} else {
				w.s.state = stateMoveFromInventory
			}

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'd', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			product := w.getProduct()
			w.s.st[w.storageIndex].Remove(product, 1)

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			product := w.getProduct()

			storage := w.s.st[w.storageIndex]
			otherStorage := w.s.st[(w.storageIndex+1)%2]

			added := otherStorage.Add(product, 1)
			if added == 1 {
				removed := storage.Remove(product, 1)
				if removed == 1 {
					return nil
				}

				otherStorage.Remove(product, 1)
			}

			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (w *InventoryWidget) getProduct() *Product {
	storage := w.s.st[w.storageIndex]

	index := -1
	for _, product := range GlobalProductFactory.cannonicalOrder {
		_, present := storage.objects[product]
		if !present {
			continue
		}

		index++
		if index < w.position {
			continue
		}

		if index > w.position {
			break
		}

		return product
	}

	return nil
}

func (w *InventoryWidget) move(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if w.activeState != w.s.state {
			return nil
		}

		size := w.s.st[w.storageIndex].UniqueObjects()
		if size == 0 {
			return nil
		}

		newPosition := w.position + d
		if newPosition >= 0 && newPosition < size {
			w.position = newPosition
		}

		return nil
	}
}

// RecipeSelectorWidget a GameWidget that displays available recipes
type RecipeSelectorWidget struct {
	name    string
	offsetY int
	width   int
	height  int

	game *Game
	s    *state

	position int
}

func newRecipeSelectorWidget() *RecipeSelectorWidget {
	w := new(RecipeSelectorWidget)

	return w
}

// SetGame sets the Game associated with RecipeSelectorWidget
func (w *RecipeSelectorWidget) SetGame(game *Game) {
	w.game = game
}

// Layout displays the RecipeSelectorWidget
func (w *RecipeSelectorWidget) Layout(g *gocui.Gui) error {
	if w.s.state != stateSetRecipe {
		return nil
	}

	maxX, _ := g.Size()

	v, err := g.SetView(w.name, maxX-w.width, w.offsetY, maxX-1, w.offsetY+w.height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		err = w.initBindings(g)
		if err == nil {
			return err
		}
	}

	if _, err := g.SetViewOnTop(w.name); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(w.name); err != nil {
		return err
	}

	v.Title = w.name
	v.Clear()

	maxPrintLines := w.height - 2
	printedLines := 0

	for index, recipe := range GlobalRecipeFactory.Assembly {
		if index < w.position {
			continue
		}

		if printedLines > maxPrintLines {
			break
		}

		var prefix string
		if index == w.position {
			prefix = ">"
		} else {
			prefix = " "
		}

		fmt.Fprintf(v, "%s %s\n", prefix, recipe.output.name)
		printedLines++

		for _, product := range recipe.inputOrder {
			fmt.Fprintf(v, "    %2d x %s\n", recipe.input[product], product.name)
			printedLines++
		}

	}

	return nil
}

func (w *RecipeSelectorWidget) initBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
		w.move(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
		w.move(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			w.position = 0
			w.s.state = stateNavigate

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			x, y := w.game.GetCursor()
			s, _, _ := w.game.GetStructureAt(y, x)
			if s == nil {
				// this should not happen, this should only display if on a structure
				return nil
			}

			switch f := s.(type) {
			case *Factory:
				f.SetRecipe(GlobalRecipeFactory.Assembly[w.position])
				w.s.state = stateNavigate
			}

			return nil
		}); err != nil {
		return err
	}

	return nil
}

func (w *RecipeSelectorWidget) move(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		size := len(GlobalRecipeFactory.Assembly)

		newPosition := w.position + d
		if newPosition >= 0 && newPosition < size {
			w.position = newPosition
		}

		return nil
	}
}
