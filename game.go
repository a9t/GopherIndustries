package main

import (
	"math/rand"
)

type position struct {
	x int
	y int
}

// Game implementation
type Game struct {
	WorldMap [][]Tile
	roots    map[Structure]position
	cursor   position
	invetory *Storage
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

		_, hasProduct := crt.CanRetrieveProduct()

		crt.Tick()
		for _, input := range crt.Inputs() {
			x := p.x + input.x
			y := p.y + input.y

			neighbour, ny, nx := g.GetNeighbour(y, x, input.d, true)
			if neighbour == nil {
				continue
			}

			// the Structure is a valid input provider
			inProgress[neighbour] = position{x: nx, y: ny}

			if hasProduct {
				// current structure has a product, so it cannot consume input
				continue
			}

			retrieved, _ := neighbour.CanRetrieveProduct()
			if retrieved == nil {
				// input provider has no item, so no need to look into it anymore
				break
			}

			// check and see if the current structure can consume Product
			if !crt.CanAcceptProduct(retrieved) {
				break
			}

			retrieved, _ = neighbour.RetrieveProduct()
			crt.AcceptProduct(retrieved)
			break
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

	for _, input := range s.Inputs() {
		tx := input.x + x
		ty := input.y + y

		neighbour, nx, ny := g.GetNeighbour(ty, tx, input.d, true)
		if neighbour == nil {
			continue
		}

		pos := position{x: nx, y: ny}
		g.roots[neighbour] = pos
	}

	structureTiles := s.Tiles()
	for xx, tiles := range structureTiles {
		for yy, tile := range tiles {
			g.WorldMap[y+yy][x+xx] = tile.UnderlyingResource()
		}
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
		for j := range tiles {
			t := g.WorldMap[y+i][x+j]
			switch t.(type) {
			case StructureTile:
				return false
			}
		}
	}

	for i, tiles := range tilesMatrix {
		for j, tile := range tiles {
			t := g.WorldMap[y+i][x+j]
			switch t.(type) {
			case *RawResource:
				tile.SetUnderlyingResource(t.(*RawResource))
				g.WorldMap[y+i][x+j] = tile
			}
		}
	}

	for _, input := range s.Inputs() {
		sx := input.x + x
		sy := input.y + y

		neighbour, _, _ := g.GetNeighbour(sy, sx, input.d, true)
		if neighbour == nil {
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
		if neighbour != nil {
			isRoot = false
			break
		}
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

	worldMap := make([][]Tile, height)
	for i := 0; i < height; i++ {
		worldMap[i] = make([]Tile, width)

		for j := 0; j < width; j++ {
			switch r := rand.Float32(); {
			case r < 0.8:
				worldMap[i][j] = &RawResource{0, -1}
			case r < 0.9:
				worldMap[i][j] = &RawResource{rand.Int() % 100, 0}
			case r < 0.98:
				worldMap[i][j] = &RawResource{100 + rand.Int()%100, 0}
			default:
				worldMap[i][j] = &RawResource{200 + rand.Int()%100, 0}
			}
		}
	}

	g := new(Game)
	g.WorldMap = worldMap
	g.roots = make(map[Structure]position)
	g.invetory = NewStorage(100)

	return g
}
