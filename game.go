package main

import (
	"math"
	"math/rand"
)

type position struct {
	x int
	y int
}

// Game implementation
type Game struct {
	WorldMap  [][]Tile
	roots     map[Structure]position
	splitters map[*Splitter]position
	cursor    position
	inventory *Storage
}

// WithinBounds indicates if the position is within the map limits
func (g *Game) WithinBounds(x, y int) bool {
	if y < 0 || y >= len(g.WorldMap) {
		return false
	}

	if x < 0 || x >= len(g.WorldMap[0]) {
		return false
	}

	return true
}

// GetCursor returns the cursor position in the game world
func (g *Game) GetCursor() (int, int) {
	return g.cursor.x, g.cursor.y
}

// MoveCursor updates the cursor position with the specified offset
func (g *Game) MoveCursor(offsetX, offsetY int) bool {
	x := g.cursor.x + offsetX
	y := g.cursor.y + offsetY

	if !g.WithinBounds(x, y) {
		return false
	}

	g.cursor.x, g.cursor.y = x, y
	return true
}

// GetStructureAt returns the Structure that covers the location and its top right corner
func (g *Game) GetStructureAt(y, x int) (Structure, int, int) {
	if !g.WithinBounds(x, y) {
		return nil, -1, -1
	}

	tile := g.WorldMap[y][x]

	switch t := tile.(type) {
	case StructureTile:
		s := t.Group()

		for yy, tiles := range s.Tiles() {
			for xx, tile := range tiles {
				if t == tile {
					return s, y - yy, x - xx
				}
			}
		}
	}
	return nil, -1, -1
}

// Tick advances the internal state of the game
func (g *Game) Tick() {
	// handle splitters first
	for s := range g.splitters {
		s.Tick()
	}

	for s, p := range g.splitters {
		for _, input := range s.Inputs() {
			if !s.CanAcceptProduct(nil) {
				// splitted cannot accept products anymore
				break
			}

			x := p.x + input.x
			y := p.y + input.y

			neighbour, _, _ := g.GetNeighbour(y, x, input.d, true)
			if neighbour == nil {
				continue
			}

			retrieved, _ := neighbour.RetrieveProduct()
			if retrieved == nil {
				continue
			}

			s.AcceptProduct(retrieved)
		}
	}

	inProgress := make(map[Structure]position)
	for s, p := range g.roots {
		inProgress[s] = p
	}

	for len(inProgress) != 0 {
		var crt Structure
		var p position

		for crt, p = range inProgress {
			break
		}
		delete(inProgress, crt)

		crt.Tick()
		for _, input := range crt.Inputs() {
			x := p.x + input.x
			y := p.y + input.y

			neighbour, ny, nx := g.GetNeighbour(y, x, input.d, true)
			if neighbour == nil {
				continue
			}

			// the Structure is a valid input provider
			switch neighbour.(type) {
			case *Splitter:
				// do nothing, splitters are handled separately
			default:
				inProgress[neighbour] = position{x: nx, y: ny}
			}

			_, hasProduct := crt.CanRetrieveProduct()
			if hasProduct {
				// current structure has a product, so it cannot consume input
				continue
			}

			retrieved, _ := neighbour.CanRetrieveProduct()
			if retrieved == nil {
				// input provider has no item, so no need to look into it anymore
				continue
			}

			// check and see if the current structure can consume Product
			if !crt.CanAcceptProduct(retrieved) {
				continue
			}

			retrieved, _ = neighbour.RetrieveProduct()
			crt.AcceptProduct(retrieved)
			continue
		}
	}
}

// GetNeighbour returns a Structure that is linked to the Transfer point (in or out)
func (g *Game) GetNeighbour(y int, x int, d Direction, in bool) (Structure, int, int) {
	correction := 1
	if !in {
		correction = -1
	}

	switch d {
	case DirectionDown:
		y -= correction
	case DirectionLeft:
		x += correction
	case DirectionUp:
		y += correction
	case DirectionRight:
		x -= correction
	}

	neighbour, ny, nx := g.GetStructureAt(y, x)
	if neighbour == nil {
		return nil, -1, -1
	}

	var transfers []Transfer
	if in {
		transfers = neighbour.Outputs()
	} else {
		transfers = neighbour.Inputs()
	}

	for _, transfer := range transfers {
		tx := nx + transfer.x
		if x != tx {
			continue
		}

		ty := ny + transfer.y
		if y != ty {
			continue
		}

		if transfer.d == d {
			return neighbour, ny, nx
		}
	}

	return nil, -1, -1
}

// RemoveStructure remove the structure at the specified position from the map
func (g *Game) RemoveStructure(y, x int) Structure {
	s, y, x := g.GetStructureAt(y, x)
	if s == nil {
		return nil
	}

	structureTiles := s.Tiles()
	for yy, tiles := range structureTiles {
		for xx, tile := range tiles {
			g.WorldMap[y+yy][x+xx] = tile.UnderlyingResource()
		}
	}

	switch ss := s.(type) {
	case *Splitter:
		delete(g.splitters, ss)
		return s
	}

	for _, input := range s.Inputs() {
		tx := input.x + x
		ty := input.y + y

		neighbour, nx, ny := g.GetNeighbour(ty, tx, input.d, true)
		if neighbour == nil {
			continue
		}

		switch neighbour.(type) {
		case *Splitter:
			// splitters cannot be roots, because they are handled separately
			continue
		}

		pos := position{x: nx, y: ny}
		g.roots[neighbour] = pos
	}

	delete(g.roots, s)

	return s
}

// PlaceStructure puts a Building at the specified location on the map
func (g *Game) PlaceStructure(y, x int, s Structure) bool {
	tilesMatrix := s.Tiles()

	if y < 0 || y+len(tilesMatrix) >= len(g.WorldMap) {
		return false
	}

	if x < 0 || x+len(tilesMatrix[0]) >= len(g.WorldMap[0]) {
		return false
	}

	for i, tiles := range tilesMatrix {
		for j, tile := range tiles {
			if tile == nil {
				continue
			}

			t := g.WorldMap[y+i][x+j]
			switch t.(type) {
			case StructureTile:
				return false
			}
		}
	}

	for yy, tiles := range tilesMatrix {
		for xx, tile := range tiles {
			if tile == nil {
				continue
			}

			t := g.WorldMap[y+yy][x+xx]
			switch t.(type) {
			case *RawResource:
				tile.SetUnderlyingResource(t.(*RawResource))
				g.WorldMap[y+yy][x+xx] = tile
			}
		}
	}

	switch ss := s.(type) {
	case *Splitter:
		g.splitters[ss] = position{x: x, y: y}
		return true
	}

	for _, input := range s.Inputs() {
		sx := input.x + x
		sy := input.y + y

		neighbour, _, _ := g.GetNeighbour(sy, sx, input.d, true)
		if neighbour == nil {
			continue
		}

		switch s.(type) {
		case *Splitter:
			// special handling for splitters, ignore them as roots
			continue
		}

		// current structure is a consumer, so the neighbour cannot be a root
		delete(g.roots, neighbour)
	}

	isRoot := true
	for _, output := range s.Outputs() {
		sx := output.x + x
		sy := output.y + y

		neighbour, _, _ := g.GetNeighbour(sy, sx, output.d, false)
		if neighbour == nil {
			continue
		}

		switch neighbour.(type) {
		case *Splitter:
			// special handling for splitters, ignore them as roots
			continue
		}

		isRoot = false
		break
	}

	if isRoot {
		g.roots[s] = position{x: x, y: y}
	}

	return true
}

// GenerateGame creates a new game instance
func GenerateGame(height int, width int) *Game {
	if width <= 0 || height <= 0 {
		return nil
	}

	g := new(Game)
	g.WorldMap = generateMap(height, width)
	g.roots = make(map[Structure]position)
	g.splitters = make(map[*Splitter]position)

	g.inventory = NewStorage(100)
	g.inventory.Add(GlobalProductFactory.GetProduct(ProductStructureBelt), 50)
	g.inventory.Add(GlobalProductFactory.GetProduct(ProductStructureChest), 3)
	g.inventory.Add(GlobalProductFactory.GetProduct(ProductStructureExtractor), 3)
	g.inventory.Add(GlobalProductFactory.GetProduct(ProductStructureSplitter), 6)
	g.inventory.Add(GlobalProductFactory.GetProduct(ProductStructureFactory), 8)

	return g
}

func distance(x, y int, xx, yy int) float64 {
	dx := x - xx
	if dx < 0 {
		dx = -dx
	}

	dy := y - yy
	if dy < 0 {
		dy = -dy
	}

	return math.Sqrt(float64(dx*dx + dy*dy))
}

func generateMap(height, width int) [][]Tile {
	worldMap := make([][]Tile, height)
	for i := 0; i < height; i++ {
		worldMap[i] = make([]Tile, width)
		for j := 0; j < width; j++ {
			worldMap[i][j] = &RawResource{0, -1}
		}
	}

	minDistance := 20.
	refMaxRay := 10.

	size := 10 + rand.Int()%4
	centers := make([][]int, size)
	for index := range centers {
		centers[index] = make([]int, 2)

		var x int
		var y int

		for {
			x = rand.Int() % width
			y = rand.Int() % height

			valid := true
			for i := 0; i < index; i++ {
				d := distance(x, y, centers[i][0], centers[i][1])
				if d <= minDistance {
					valid = false
					break
				}
			}

			if valid {
				break
			}
		}

		centers[index][0] = x
		centers[index][1] = y

		otherX := x - 4 + rand.Int()%9
		otherY := y - 4 + rand.Int()%9

		var minX, maxX int
		if x < otherX {
			minX, maxX = x, otherX
		} else {
			minX, maxX = otherX, x
		}
		minX -= 10
		maxX += 10

		var minY, maxY int
		if y < otherY {
			minY, maxY = y, otherY
		} else {
			minY, maxY = otherY, y
		}
		minY -= 10
		maxY += 10

		maxRay := refMaxRay + rand.Float64()*6
		for i := minY; i <= maxY; i++ {
			for j := minX; j <= maxX; j++ {
				if i < 0 || j < 0 || i >= height || j >= width {
					continue
				}

				d1 := distance(x, y, j, i)
				d2 := distance(otherX, otherY, j, i)

				if d1+d2 > maxRay {
					continue
				}

				var amount int
				switch r := rand.Float32(); {
				case r <= 0.7:
					amount = rand.Int() % 100
				case r <= 0.9:
					amount = 100 + rand.Int()%100
				default:
					amount = 200 + rand.Int()%100
				}

				worldMap[i][j] = &RawResource{amount, index % 3}
			}

		}
	}

	return worldMap
}
