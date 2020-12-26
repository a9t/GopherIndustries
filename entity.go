package main

import (
	"fmt"
	"unicode/utf8"
)

// Direction indicates the movement direction
type Direction = uint8

const (
	// DirectionDown indicates downward movement
	DirectionDown Direction = iota
	// DirectionLeft indicates leftward movement
	DirectionLeft
	// DirectionRight indicates rightward movement
	DirectionRight
	// DirectionUp idnicates upward movement
	DirectionUp
)

// Transfer is a transfer point for the input/output of a Structure
type Transfer struct {
	x int
	y int
	d Direction
}

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
	UnderlyingResource() *RawResource
	SetUnderlyingResource(*RawResource)
	CopyStructureTile() StructureTile
}

// Structure is a Tile containing a player created entity
type Structure interface {
	Tiles() [][]StructureTile
	RotateLeft()
	RotateRight()
	CopyStructure() Structure
	Outputs() []Transfer
	Inputs() []Transfer
}

// BaseStructureTile one Tile component of a Structure
type BaseStructureTile struct {
	rotationPosition int
	maxRotations     int
	symbolID         string

	structure Structure
	resource  *RawResource
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
func (b *BaseStructureTile) UnderlyingResource() *RawResource {
	return b.resource
}

// SetUnderlyingResource sets the resource beneath the BaseStructureTile
func (b *BaseStructureTile) SetUnderlyingResource(r *RawResource) {
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
	tiles   [][]StructureTile
	inputs  []Transfer
	outputs []Transfer
}

// Inputs return Transfer point outside the Structure where the inputs are expected to come from
func (s *BaseStructure) Inputs() []Transfer {
	return s.inputs
}

// Outputs return Transfer point inside the Structure where the outputs are generated
func (s *BaseStructure) Outputs() []Transfer {
	return s.outputs
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

	for _, input := range s.inputs {
		input.x, input.y = height-1-input.y, input.x
		input.d = (input.d + 1) % 4
	}

	for _, output := range s.outputs {
		output.x, output.y = height-1-output.y, output.x
		output.d = (output.d + 1) % 4
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

	for _, input := range s.inputs {
		input.x, input.y = input.y, width-1-input.x
		input.d = (input.d + 3) % 4
	}

	for _, output := range s.outputs {
		output.x, output.y = output.y, width-1-output.x
		output.d = (output.d + 3) % 4
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

	structure.inputs = make([]Transfer, len(s.inputs))
	for i, input := range s.inputs {
		structure.inputs[i] = input
	}

	structure.outputs = make([]Transfer, len(s.outputs))
	for i, output := range s.outputs {
		structure.outputs[i] = output
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

// FillerMidTile replacement for FillerMidTile
type FillerMidTile struct {
	BaseStructureTile
}

// NewFillerMidTile creates a new *FillerMidTile
func NewFillerMidTile(pos int) *FillerMidTile {
	tile := FillerMidTile{BaseStructureTile{pos % 4, 4, "fillerMid", nil, nil}}
	return &tile
}

// FillerCenterTile replacement for FillerCenterTile
type FillerCenterTile struct {
	BaseStructureTile
}

// NewFillerCenterTile creates a new *FillerMidTile
func NewFillerCenterTile(pos int) *FillerCenterTile {
	tile := FillerCenterTile{BaseStructureTile{pos % 4, 4, "fillerCenter", nil, nil}}
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

	block.inputs = make([]Transfer, 0)
	block.outputs = make([]Transfer, 0)

	return block
}

// ThreeXThreeBlock replacement for ThreeXThreeBlock
type ThreeXThreeBlock struct {
	BaseStructure
}

// NewThreeXThreeBlock creates a new *ThreeXThreeBlock
func NewThreeXThreeBlock() *ThreeXThreeBlock {
	block := new(ThreeXThreeBlock)
	block.tiles = [][]StructureTile{
		{NewFillerCornerTile(0), NewFillerMidTile(0), NewFillerCornerTile(1)},
		{NewFillerMidTile(3), NewFillerCenterTile(0), NewFillerMidTile(1)},
		{NewFillerCornerTile(3), NewFillerMidTile(2), NewFillerCornerTile(2)},
	}

	block.inputs = make([]Transfer, 0)
	block.outputs = make([]Transfer, 0)

	return block
}

// BeltTile is the map representation of a conveyor belt
type BeltTile struct {
	BaseStructureTile
	product *Product
}

// Display create a string to be diplayed for the BellTile
func (b *BeltTile) Display(mode DisplayMode) string {
	if b.product != nil {
		return string([]rune{b.product.representation})
	}

	return b.BaseStructureTile.Display(mode)
}

// NewBeltTile creates a new *BeltTile
func NewBeltTile() *BeltTile {
	return &BeltTile{BaseStructureTile{0, 12, "belt", nil, nil}, nil}
}

// Belt is the structure representation of a conveyor belt
type Belt struct {
	BaseStructure
	rotationPosition int
}

// NewBelt creates a new *Belt
func NewBelt() *Belt {
	block := new(Belt)
	block.tiles = [][]StructureTile{
		{NewBeltTile()},
	}

	block.inputs = make([]Transfer, 1)
	block.outputs = make([]Transfer, 1)

	block.outputs[0] = Transfer{0, 0, 0}
	block.inputs[0] = Transfer{0, -1, 0}

	return block
}

// CopyStructure creates a copy of the Belt
func (b *Belt) CopyStructure() Structure {
	belt := new(Belt)
	belt.rotationPosition = b.rotationPosition

	structure := b.BaseStructure.CopyStructure()
	if baseStructure, ok := structure.(*BaseStructure); ok {
		belt.BaseStructure = *baseStructure
	}

	return belt
}

// RotateRight gets the next rotation of a Belt
func (b *Belt) RotateRight() {
	b.rotationPosition++
	b.rotationPosition %= 12

	b.tiles[0][0].RotateRight()

	b.setTransfers()
}

// RotateLeft gets the next rotation of a Belt
func (b *Belt) RotateLeft() {
	b.rotationPosition += 11
	b.rotationPosition %= 12

	b.tiles[0][0].RotateRight()

	b.setTransfers()
}

func (b *Belt) setTransfers() {
	b.inputs[0].d = uint8(b.rotationPosition) / 3
	b.outputs[0].d = uint8(b.rotationPosition) % 3

	switch b.inputs[0].d {
	case DirectionDown:
		b.inputs[0].x = 0
		b.inputs[0].y = -1
	case DirectionLeft:
		b.inputs[0].x = 1
		b.inputs[0].y = 0
	case DirectionUp:
		b.inputs[0].x = 0
		b.inputs[0].y = 1
	case DirectionRight:
		b.inputs[0].x = -1
		b.inputs[0].y = 0
	}
}

// RawResource is a Tile containing natural resources
type RawResource struct {
	amount   int
	resource int
}

// Display displays the Resource tile
func (t *RawResource) Display(mode DisplayMode) string {
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

	symbolColor := 33
	colorMode := 4
	if mode == DisplayModeMapSelected {
		symbolColor = 37
		colorMode = 7
	}

	return fmt.Sprintf("\033[%d;%dm%c\033[0m", symbolColor, colorMode, symbol)
}
