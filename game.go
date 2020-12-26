package main

import (
	"math/rand"
)

// Game implementation
type Game struct {
	WorldMap [][]Tile
}

// PlaceBuilding puts a Building at the specified location on the map
func (g *Game) PlaceBuilding(y, x int, s Structure) bool {
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
				worldMap[i][j] = &RawResource{0, 1}
			case r < 0.9:
				worldMap[i][j] = &RawResource{1, 1}
			case r < 0.98:
				worldMap[i][j] = &RawResource{2, 1}
			default:
				worldMap[i][j] = &RawResource{3, 1}
			}
		}
	}

	return &Game{worldMap}
}
