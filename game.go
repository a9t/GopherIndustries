package main

import (
	"math/rand"
)

// Game implementation
type Game struct {
	WorldMap [][]Tile
}

// PlaceBuilding puts a Building at the specified location on the map
func (g *Game) PlaceBuilding(y, x int, b Structure) bool {
	if y < 0 || y >= len(g.WorldMap) {
		return false
	}

	if x < 0 || x >= len(g.WorldMap[0]) {
		return false
	}

	for i, tiles := range b.Tiles() {
		for j, tile := range tiles {
			t := g.WorldMap[y+i][x+j]
			switch t.(type) {
			case *Resource:
				tile.SetUnderlyingResource(t.(*Resource))
				g.WorldMap[y][x] = tile
				return true
			}
		}
	}

	return false
}

// Tick update the state of the game
func (g *Game) Tick() {
	//
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
				worldMap[i][j] = &Resource{0, 1}
			case r < 0.9:
				worldMap[i][j] = &Resource{1, 1}
			case r < 0.98:
				worldMap[i][j] = &Resource{2, 1}
			default:
				worldMap[i][j] = &Resource{3, 1}
			}
		}
	}

	return &Game{worldMap}
}
