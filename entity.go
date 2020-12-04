package main

import (
	"fmt"
	"unicode/utf8"
)

// DisplayMode indicates an entity's display mode
type DisplayMode = uint8

const (
	// DisplayModeMap default representation on the map
	DisplayModeMap DisplayMode = iota

	// DisplayModeMapSelected representation of selected item on the map
	DisplayModeMapSelected

	// DisplayModeGhostValid valid representation when placing on the map
	DisplayModeGhostValid

	// DisplayModeGhostInvalid invalid representation when placing on the map
	DisplayModeGhostInvalid
)

// Tile is a component of the game map
type Tile interface {
	Display(DisplayMode) string
}

// StructureTile one Tile component of a Structure
type StructureTile interface {
	Tile
	RotateLeft()
	RotateRight()
	Group() Structure
	SetGroup(Structure)
	UnderlyingResource() *Resource
	SetUnderlyingResource(*Resource)
	CopyStructureTile() StructureTile
}

// Structure is a Tile containing a player created entity
type Structure interface {
	Tiles() [][]StructureTile
	RotateLeft()
	RotateRight()
	CopyStructure() Structure
}

// BaseStructureTile one Tile component of a Structure
type BaseStructureTile struct {
	rotationPosition int
	maxRotations     int
	symbolID         string

	structure Structure
	resource  *Resource
}

// RotateRight gets the next rotation of a BaseStructureTile
func (b *BaseStructureTile) RotateRight() {
	b.rotationPosition++
	b.rotationPosition %= b.maxRotations
}

// RotateLeft gets the next rotation of a BaseStructureTile
func (b *BaseStructureTile) RotateLeft() {
	if b.rotationPosition == 0 {
		b.rotationPosition = b.maxRotations - 1
	} else {
		b.rotationPosition--
	}
}

// UnderlyingResource provides the resource beneath the BaseStructureTile
func (b *BaseStructureTile) UnderlyingResource() *Resource {
	return b.resource
}

// SetUnderlyingResource sets the resource beneath the BaseStructureTile
func (b *BaseStructureTile) SetUnderlyingResource(r *Resource) {
	b.resource = r
}

// Group returns the Structure the BaseStructureTile
func (b *BaseStructureTile) Group() Structure {
	return b.structure
}

// SetGroup specifies the Structure for BaseStructureTile
func (b *BaseStructureTile) SetGroup(s Structure) {
	b.structure = s
}

// CopyStructureTile creates a deep copy of the current BaseStructureTile
func (b *BaseStructureTile) CopyStructureTile() StructureTile {
	return &BaseStructureTile{b.rotationPosition, b.maxRotations, b.symbolID, nil, nil}
}

// Display create a string to be diplayed BaseStructureTile
func (b *BaseStructureTile) Display(mode DisplayMode) string {
	symbolConfig := GlobalDisplayConfigManager.GetSymbolConfig()
	symbols := symbolConfig.Types[b.symbolID]

	var symbol rune
	var width int
	repeat := b.rotationPosition

	for i, w := 0, 0; i < len(symbols); i += w {
		symbol, width = utf8.DecodeRuneInString(symbols[i:])
		w = width

		if repeat == 0 {
			break
		}
		repeat--
	}

	symbolColors := GlobalDisplayConfigManager.GetColorConfig().StructureColors[mode]

	return fmt.Sprintf("\033[%d;%dm%c\033[0m", symbolColors[0], symbolColors[1], symbol)
}

// BaseStructure is a basic implementation of Structure
type BaseStructure struct {
	tiles [][]StructureTile
}

// Tiles return the Tiles associated with the BaseStructure
func (s *BaseStructure) Tiles() [][]StructureTile {
	return s.tiles
}

// RotateRight gets the next rotation of a BaseStructure
func (s *BaseStructure) RotateRight() {
	height := len(s.tiles)
	width := len(s.tiles[0])
	newTiles := make([][]StructureTile, width)
	for i := 0; i < len(s.tiles[0]); i++ {
		newTiles[i] = make([]StructureTile, height)
	}

	for i, tiles := range s.tiles {
		for j, tile := range tiles {
			newTiles[j][height-1-i] = tile
			if tile != nil {
				tile.RotateRight()
			}
		}
	}

	s.tiles = newTiles
}

// RotateLeft gets the next rotation of a BaseStructure
func (s *BaseStructure) RotateLeft() {
	height := len(s.tiles)
	width := len(s.tiles[0])
	newTiles := make([][]StructureTile, width)
	for i := 0; i < len(s.tiles[0]); i++ {
		newTiles[i] = make([]StructureTile, height)
	}

	for i, tiles := range s.tiles {
		for j, tile := range tiles {
			newTiles[width-1-j][i] = tile
			if tile != nil {
				tile.RotateLeft()
			}
		}
	}

	s.tiles = newTiles
}

// CopyStructure creates a deep copy of the current BaseStructure
func (s *BaseStructure) CopyStructure() Structure {
	bTiles := s.Tiles()

	structure := new(BaseStructure)
	structure.tiles = make([][]StructureTile, len(bTiles))
	for i, tiles := range bTiles {
		structure.tiles[i] = make([]StructureTile, len(bTiles[0]))
		for j, tile := range tiles {
			if tile == nil {
				structure.tiles[i][j] = nil
			} else {
				structure.tiles[i][j] = tile.CopyStructureTile()
				structure.tiles[i][j].SetGroup(structure)
			}
		}
	}

	return structure
}

// FillerCornerTile replacement for FillerCornerTile
type FillerCornerTile struct {
	BaseStructureTile
}

// NewFillerCornerTile creates a new *FillerCornerTile
func NewFillerCornerTile(pos int) *FillerCornerTile {
	tile := FillerCornerTile{BaseStructureTile{pos % 4, 4, "fillerCorner", nil, nil}}
	return &tile
}

// TwoXTwoBlock replacement for TwoXTwoBlock
type TwoXTwoBlock struct {
	BaseStructure
}

// NewTwoXTwoBlock creates a new *TwoXTwoBlock
func NewTwoXTwoBlock() *TwoXTwoBlock {
	block := new(TwoXTwoBlock)
	block.tiles = [][]StructureTile{
		{NewFillerCornerTile(0), NewFillerCornerTile(1)},
		{NewFillerCornerTile(3), NewFillerCornerTile(2)},
	}

	return block
}

// BeltTile is the map representation of a conveyor belt
type BeltTile struct {
	BaseStructureTile
}

// NewBeltTile creates a new *BeltTile
func NewBeltTile() *BeltTile {
	return &BeltTile{BaseStructureTile{0, 12, "belt", nil, nil}}
}

// Belt is the structure representation of a conveyor belt
type Belt struct {
	BaseStructure
}

// NewBelt creates a new *Belt
func NewBelt() *Belt {
	block := new(Belt)
	block.tiles = [][]StructureTile{
		{NewBeltTile()},
	}

	return block
}

// Resource is a Tile containing natural resources
type Resource struct {
	amount   int
	resource int
}

// Display displays the Resource tile
func (t *Resource) Display(mode DisplayMode) string {
	symbolConfig := GlobalDisplayConfigManager.GetSymbolConfig()
	symbols := symbolConfig.Types["resource"]

	var symbol rune
	var width int
	repeat := t.amount

	for i, w := 0, 0; i < len(symbols); i += w {
		symbol, width = utf8.DecodeRuneInString(symbols[i:])
		w = width

		if repeat == 0 {
			break
		}
		repeat--
	}

	// var symbol rune
	// switch t.amount {
	// case 0:
	// 	symbol = ' '
	// case 1:
	// 	symbol = '\u2591'
	// case 2:
	// 	symbol = '\u2592'
	// case 3:
	// 	symbol = '\u2593'
	// }

	symbolColor := 33
	colorMode := 4
	if mode == DisplayModeMapSelected {
		symbolColor = 37
		colorMode = 7
	}

	return fmt.Sprintf("\033[%d;%dm%c\033[0m", symbolColor, colorMode, symbol)
}
