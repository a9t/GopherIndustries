package main

import "fmt"

// DisplayMode indicates an entity's display mode
type DisplayMode = uint8

const (
	// DisplayModeMap default representation on the map
	DisplayModeMap DisplayMode = iota

	// DisplayModeMapSelected representation of selected item on the map
	DisplayModeMapSelected

	// DisplayModePreview representation in the preview section
	DisplayModePreview

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

// Belt is a structure that transports resources
type Belt struct {
	index    int
	Resource *Resource
}

// RotateRight gets the next rotation of a belt
func (b *Belt) RotateRight() {
	b.index++
	b.index %= 12
}

// RotateLeft gets the next rotation of a belt
func (b *Belt) RotateLeft() {
	if b.index == 0 {
		b.index = 11
	} else {
		b.index--
	}
}

// UnderlyingResource provides the resource beneath the building
func (b *Belt) UnderlyingResource() *Resource {
	return b.Resource
}

// SetUnderlyingResource sets the resource beneath the building
func (b *Belt) SetUnderlyingResource(r *Resource) {
	b.Resource = r
}

// Display displays the belt
func (b *Belt) Display(mode DisplayMode) string {
	symbolColor := 33
	var symbol rune
	switch b.index {
	case 0:
		symbol = '\u257D'
	case 1:
		symbol = '\u2519'
	case 2:
		symbol = '\u2515'
	case 3:
		symbol = '\u257E'
	case 4:
		symbol = '\u2516'
	case 5:
		symbol = '\u250E'
	case 6:
		symbol = '\u257F'
	case 7:
		symbol = '\u250D'
	case 8:
		symbol = '\u2511'
	case 9:
		symbol = '\u257C'
	case 10:
		symbol = '\u2512'
	case 11:
		symbol = '\u251A'
	}

	var colorMode int
	switch mode {
	case DisplayModeMap:
		colorMode = 4
	case DisplayModeGhostValid:
		colorMode = 4
	case DisplayModeGhostInvalid:
		colorMode = 4
		symbolColor = 31
	case DisplayModeMapSelected:
		colorMode = 4
	default:
		colorMode = 4
	}

	return fmt.Sprintf("\033[%d;%dm%c\033[0m", symbolColor, colorMode, symbol)
}

// Group returns the Structure the Belt is associated with, itself
func (b *Belt) Group() Structure {
	return b
}

// Tiles return the Tiles associated with the Belt, an array of itself
func (b *Belt) Tiles() [][]StructureTile {
	return [][]StructureTile{{b}}
}

// SetGroup specifies the Structure for the Belt (does nothing, as Belt also implements Structure)
func (b *Belt) SetGroup(s Structure) {
}

// CopyStructure creates a deep copy of the current Belt
func (b *Belt) CopyStructure() Structure {
	return &Belt{b.index, nil}
}

// CopyStructureTile creates a deep copy of the current Belt
func (b *Belt) CopyStructureTile() StructureTile {
	return &Belt{b.index, nil}
}

// Resource is a Tile containing natural resources
type Resource struct {
	amount   int
	resource int
}

// Display displays the Resource tile
func (t *Resource) Display(mode DisplayMode) string {
	var symbol rune
	switch t.amount {
	case 0:
		symbol = ' '
	case 1:
		symbol = '\u2591'
	case 2:
		symbol = '\u2592'
	case 3:
		symbol = '\u2593'
	}

	symbolColor := 33
	colorMode := 4
	if mode == DisplayModeMapSelected {
		symbolColor = 37
		colorMode = 7
	}

	return fmt.Sprintf("\033[%d;%dm%c\033[0m", symbolColor, colorMode, symbol)
}

// FillerCornerTile a StructureTile that represents a corner
type FillerCornerTile struct {
	index     int
	structure Structure
	Resource  *Resource
}

// Display displays the FillerCornerTile tile
func (t *FillerCornerTile) Display(mode DisplayMode) string {
	var symbol rune

	switch t.index {
	case 0:
		symbol = '\u259B'
	case 1:
		symbol = '\u259C'
	case 2:
		symbol = '\u259F'
	case 3:
		symbol = '\u2599'
	}

	symbolColor := 33
	colorMode := 4
	if mode == DisplayModeMapSelected {
		symbolColor = 37
		colorMode = 4
	} else if mode == DisplayModeGhostInvalid {
		symbolColor = 31
		colorMode = 4
	}

	return fmt.Sprintf("\033[%d;%dm%c\033[0m", symbolColor, colorMode, symbol)
}

// Group returns the Structure the Belt is associated with, itself
func (t *FillerCornerTile) Group() Structure {
	return t.structure
}

// SetGroup sets the Structure the FillerCornerTile belongs to
func (t *FillerCornerTile) SetGroup(s Structure) {
	t.structure = s
}

// UnderlyingResource provides the resource beneath the building
func (t *FillerCornerTile) UnderlyingResource() *Resource {
	return t.Resource
}

// SetUnderlyingResource sets the resource beneath the building
func (t *FillerCornerTile) SetUnderlyingResource(r *Resource) {
	t.Resource = r
}

// RotateRight gets the next rotation of a FillerCornerTile
func (t *FillerCornerTile) RotateRight() {
	t.index++
	t.index %= 4
}

// RotateLeft gets the next rotation of a FillerCornerTile
func (t *FillerCornerTile) RotateLeft() {
	if t.index == 0 {
		t.index = 3
	} else {
		t.index--
	}
}

// CopyStructureTile creates a copy of the current FillerCornerTile
func (t *FillerCornerTile) CopyStructureTile() StructureTile {
	return &FillerCornerTile{t.index, nil, nil}
}

// TwoXTwoBlock is a Structure of 2x2 tiles
type TwoXTwoBlock struct {
	tiles [][]StructureTile
}

// NewTwoXTwoBlock creates a new TwoXTwoBlock Structure
func NewTwoXTwoBlock() *TwoXTwoBlock {
	block := new(TwoXTwoBlock)
	block.tiles = [][]StructureTile{
		{&FillerCornerTile{0, block, nil}, &FillerCornerTile{1, block, nil}},
		{&FillerCornerTile{3, block, nil}, &FillerCornerTile{2, block, nil}},
	}

	return block
}

// Tiles return the Tiles associated with the TwoXTwoBlock
func (b *TwoXTwoBlock) Tiles() [][]StructureTile {
	return b.tiles
}

// RotateRight gets the next rotation of a TwoXTwoBlock
func (b *TwoXTwoBlock) RotateRight() {
	height := len(b.tiles)
	width := len(b.tiles[0])
	newTiles := make([][]StructureTile, width)
	for i := 0; i < len(b.tiles[0]); i++ {
		newTiles[i] = make([]StructureTile, height)
	}

	for i, tiles := range b.tiles {
		for j, tile := range tiles {
			newTiles[j][height-1-i] = tile
			if tile != nil {
				tile.RotateRight()
			}
		}
	}

	b.tiles = newTiles
}

// RotateLeft gets the next rotation of a TwoXTwoBlock
func (b *TwoXTwoBlock) RotateLeft() {
	height := len(b.tiles)
	width := len(b.tiles[0])
	newTiles := make([][]StructureTile, width)
	for i := 0; i < len(b.tiles[0]); i++ {
		newTiles[i] = make([]StructureTile, height)
	}

	for i, tiles := range b.tiles {
		for j, tile := range tiles {
			newTiles[width-1-j][i] = tile
			if tile != nil {
				tile.RotateLeft()
			}
		}
	}

	b.tiles = newTiles
}

// CopyStructure creates a deep copy of the current TwoXTwoBlock
func (b *TwoXTwoBlock) CopyStructure() Structure {
	bTiles := b.Tiles()

	block := new(TwoXTwoBlock)
	block.tiles = make([][]StructureTile, len(bTiles))
	for i, tiles := range bTiles {
		block.tiles[i] = make([]StructureTile, len(bTiles[0]))
		for j, tile := range tiles {
			if tile == nil {
				block.tiles[i][j] = nil
			} else {
				block.tiles[i][j] = tile.CopyStructureTile()
				block.tiles[i][j].SetGroup(block)
			}
		}
	}

	return block
}
