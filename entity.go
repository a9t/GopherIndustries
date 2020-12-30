package main

import (
	"fmt"
	"unicode/utf8"
)

const (
	// CycleSizeExtractor the cycle size for the extractor
	CycleSizeExtractor int = 40
	// CycleSizeBelt the cycle size for the belt
	CycleSizeBelt int = 20
	// ChestMaxStorage the maximum number of products stored in a chest
	ChestMaxStorage int = 100
)

// Direction indicates the movement direction
type Direction = uint8

const (
	// DirectionDown indicates downward movement
	DirectionDown Direction = iota
	// DirectionLeft indicates leftward movement
	DirectionLeft
	// DirectionUp idnicates upward movement
	DirectionUp
	// DirectionRight indicates rightward movement
	DirectionRight
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
	Tick()
	CanRetrieveProduct() (*Product, bool)
	RetrieveProduct() (*Product, bool)
	CanAcceptProduct(*Product) bool
	AcceptProduct(*Product) bool
}

// BaseStructureTile one Tile component of a Structure
type BaseStructureTile struct {
	rotationPosition int
	maxRotations     int
	symbolID         string

	structure Structure
	resource  *RawResource
	product   *Product
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
	return &BaseStructureTile{b.rotationPosition, b.maxRotations, b.symbolID, nil, nil, nil}
}

// SetProduct sets the Product associated with this BaseStructureTile
func (b *BaseStructureTile) SetProduct(p *Product) {
	b.product = p
}

// Display create a string to be diplayed BaseStructureTile
func (b *BaseStructureTile) Display(mode DisplayMode) string {
	var symbol rune
	if b.product != nil {
		symbol = b.product.representation
	} else {
		symbolConfig := GlobalDisplayConfigManager.GetSymbolConfig()
		symbols := symbolConfig.Types[b.symbolID]

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

// Inputs return Transfer point of the Structure where the inputs are expected to come from
func (s *BaseStructure) Inputs() []Transfer {
	return s.inputs
}

// Outputs return Transfer point of the Structure where the outputs are generated
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

	for i, input := range s.inputs {
		s.inputs[i].x, s.inputs[i].y = height-1-input.y, input.x
		s.inputs[i].d = (input.d + 1) % 4
	}

	for i, output := range s.outputs {
		s.outputs[i].x, s.outputs[i].y = height-1-output.y, output.x
		s.outputs[i].d = (output.d + 1) % 4
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

	for i, input := range s.inputs {
		s.inputs[i].x, s.inputs[i].y = input.y, width-1-input.x
		s.inputs[i].d = (input.d + 3) % 4
	}

	for i, output := range s.outputs {
		s.outputs[i].x, s.outputs[i].y = output.y, width-1-output.x
		s.outputs[i].d = (output.d + 3) % 4
	}

	s.tiles = newTiles
}

func (s *BaseStructure) copyStructure(parent Structure) *BaseStructure {
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
				structure.tiles[i][j].SetGroup(parent)
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
	tile := FillerCornerTile{BaseStructureTile{pos % 4, 4, "fillerCorner", nil, nil, nil}}
	return &tile
}

// FillerMidTile replacement for FillerMidTile
type FillerMidTile struct {
	BaseStructureTile
}

// NewFillerMidTile creates a new *FillerMidTile
func NewFillerMidTile(pos int) *FillerMidTile {
	tile := FillerMidTile{BaseStructureTile{pos % 4, 4, "fillerMid", nil, nil, nil}}
	return &tile
}

// FillerCenterTile replacement for FillerCenterTile
type FillerCenterTile struct {
	BaseStructureTile
}

// NewFillerCenterTile creates a new *FillerMidTile
func NewFillerCenterTile(pos int) *FillerCenterTile {
	tile := FillerCenterTile{BaseStructureTile{pos % 4, 4, "fillerCenter", nil, nil, nil}}
	return &tile
}

// OutputTile tile that indicates an output end of a Structure
type OutputTile struct {
	BaseStructureTile
}

// NewOutputTile creates a new *OutputTile
func NewOutputTile(pos int) *OutputTile {
	tile := OutputTile{BaseStructureTile{pos % 4, 4, "output", nil, nil, nil}}
	return &tile
}

// Extractor Structure that extracts a RawResource from the ground
type Extractor struct {
	BaseStructure
	counter int
	product *Product
}

// NewExtractor creates a new *Extractor
func NewExtractor() *Extractor {
	block := new(Extractor)
	block.tiles = [][]StructureTile{
		{NewFillerCornerTile(0), NewFillerMidTile(0), NewFillerCornerTile(1)},
		{NewFillerMidTile(3), NewFillerCenterTile(0), NewFillerMidTile(1)},
		{NewFillerCornerTile(3), NewOutputTile(0), NewFillerCornerTile(2)},
	}

	block.inputs = make([]Transfer, 0)
	block.outputs = make([]Transfer, 1)

	block.outputs[0] = Transfer{x: 1, y: 2, d: DirectionDown}

	return block
}

// CopyStructure creates a copy of the Extractor
func (e *Extractor) CopyStructure() Structure {
	extractor := new(Extractor)
	extractor.counter = 0

	baseStructure := e.BaseStructure.copyStructure(extractor)
	extractor.BaseStructure = *baseStructure

	return extractor
}

// CanRetrieveProduct indicates if the internal Product can be extracted
func (e *Extractor) CanRetrieveProduct() (*Product, bool) {
	var p *Product
	if e.counter == CycleSizeExtractor-1 {
		p = e.product
	}
	return p, e.product != nil
}

// RetrieveProduct returns the internal Product and resets the internal state
func (e *Extractor) RetrieveProduct() (*Product, bool) {
	if e.product == nil {
		return nil, false
	}

	p := e.product
	e.product = nil
	e.counter = 0

	return p, true
}

// CanAcceptProduct indicates if the Extractor can receive the Product
func (e *Extractor) CanAcceptProduct(*Product) bool {
	return false
}

// AcceptProduct passes the Product to the Extractor
func (e *Extractor) AcceptProduct(p *Product) bool {
	return false
}

// Tick advance the internal state of the Extractor
func (e *Extractor) Tick() {
	product, hasProduct := e.CanRetrieveProduct()
	if hasProduct {
		if product == nil {
			e.counter++
		}
		return
	}

	var rawResource *RawResource
	foundTile := false
	for _, tiles := range e.Tiles() {
		for _, tile := range tiles {
			res := tile.UnderlyingResource()
			if res.amount > 0 {
				rawResource = res
				foundTile = true
				break
			}
		}
		if foundTile {
			break
		}
	}

	if rawResource == nil {
		return
	}

	rawResource.amount--
	e.product = GlobalProductFactory.GetProduct(rawResource.resource)
}

// ChestTile is the map representation of a storage chest
type ChestTile struct {
	BaseStructureTile
}

// NewChestTile creates a new *ChestTile
func NewChestTile() *ChestTile {
	return &ChestTile{BaseStructureTile{0, 1, "chest", nil, nil, nil}}
}

// Chest is the structure representation of a storage chest
type Chest struct {
	BaseStructure
	Products map[*Product]int
	Total    int
}

// NewChest creates a new *Chest
func NewChest() *Chest {
	chest := new(Chest)
	chest.tiles = [][]StructureTile{
		{NewChestTile()},
	}
	chest.Products = make(map[*Product]int)

	chest.inputs = make([]Transfer, 4)
	chest.outputs = make([]Transfer, 0)

	chest.inputs[0] = Transfer{0, 0, DirectionDown}
	chest.inputs[1] = Transfer{0, 0, DirectionLeft}
	chest.inputs[2] = Transfer{0, 0, DirectionUp}
	chest.inputs[3] = Transfer{0, 0, DirectionRight}

	return chest
}

// Tick does nothing for Chest, no internal state to update
func (c *Chest) Tick() {
}

// CanRetrieveProduct does nothing for Chest, it does not allow auto retrieval
func (c *Chest) CanRetrieveProduct() (*Product, bool) {
	return nil, false
}

// RetrieveProduct does nothing for Chest, it does not allow auto retrieval
func (c *Chest) RetrieveProduct() (*Product, bool) {
	return nil, false
}

// CanAcceptProduct indicates if the Chest has space for another Product
func (c *Chest) CanAcceptProduct(*Product) bool {
	return c.Total < ChestMaxStorage
}

// AcceptProduct passes the Product to the Chest
func (c *Chest) AcceptProduct(p *Product) bool {
	if c.Total >= ChestMaxStorage {
		return false
	}

	val, present := c.Products[p]
	if present {
		c.Products[p] = val + 1
	} else {
		c.Products[p] = 1
	}

	c.Total++

	return true
}

// CopyStructure creates a copy of the Chest
func (c *Chest) CopyStructure() Structure {
	chest := new(Chest)
	chest.Products = make(map[*Product]int)

	baseStructure := c.BaseStructure.copyStructure(chest)
	chest.BaseStructure = *baseStructure

	return chest
}

// RotateRight does nothing for a Chest
func (c *Chest) RotateRight() {
}

// RotateLeft does nothing for a Chest
func (c *Chest) RotateLeft() {
}

// BeltTile is the map representation of a conveyor belt
type BeltTile struct {
	BaseStructureTile
}

// NewBeltTile creates a new *BeltTile
func NewBeltTile() *BeltTile {
	return &BeltTile{BaseStructureTile{0, 12, "belt", nil, nil, nil}}
}

// Belt is the structure representation of a conveyor belt
type Belt struct {
	BaseStructure
	RotationPosition int
	Product          *Product
	Counter          int
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
	block.inputs[0] = Transfer{0, 0, 0}

	return block
}

// Tick advance the internal state of the Belt
func (b *Belt) Tick() {
	product, hasProduct := b.CanRetrieveProduct()
	if hasProduct {
		if product == nil {
			b.Counter++
		}
		return
	}
}

// CanRetrieveProduct indicates if the internal Product can be extracted
func (b *Belt) CanRetrieveProduct() (*Product, bool) {
	var p *Product
	if b.Counter == CycleSizeBelt-1 {
		p = b.Product
	}
	return p, b.Product != nil
}

// RetrieveProduct returns the internal Product and resets the internal state
func (b *Belt) RetrieveProduct() (*Product, bool) {
	product, hasProduct := b.CanRetrieveProduct()
	if !hasProduct {
		return nil, false
	}

	if product == nil {
		return nil, true
	}

	b.Product = nil
	b.Counter = 0
	switch t := b.BaseStructure.tiles[0][0].(type) {
	case *BaseStructureTile:
		t.SetProduct(nil)
	}

	return product, hasProduct
}

// CanAcceptProduct indicates if the Belt can receive the Product
func (b *Belt) CanAcceptProduct(*Product) bool {
	_, hasProduct := b.CanRetrieveProduct()
	return !hasProduct
}

// AcceptProduct passes the Product to the BaseStructure
func (b *Belt) AcceptProduct(p *Product) bool {
	_, hasProduct := b.CanRetrieveProduct()
	b.Product = p

	switch t := b.BaseStructure.tiles[0][0].(type) {
	case *BaseStructureTile:
		t.SetProduct(p)
	}

	return !hasProduct
}

// CopyStructure creates a copy of the Belt
func (b *Belt) CopyStructure() Structure {
	belt := new(Belt)
	belt.RotationPosition = b.RotationPosition

	baseStructure := b.BaseStructure.copyStructure(belt)
	belt.BaseStructure = *baseStructure

	return belt
}

// RotateRight gets the next rotation of a Belt
func (b *Belt) RotateRight() {
	b.RotationPosition++
	b.RotationPosition %= 12

	b.tiles[0][0].RotateRight()

	b.setTransfers()
}

// RotateLeft gets the next rotation of a Belt
func (b *Belt) RotateLeft() {
	b.RotationPosition += 11
	b.RotationPosition %= 12

	b.tiles[0][0].RotateLeft()

	b.setTransfers()
}

func (b *Belt) setTransfers() {
	entry := uint8(b.RotationPosition) / 3
	b.inputs[0].d = entry

	exit := uint8(b.RotationPosition) % 3
	switch exit {
	case 0:
		exit = entry
	case 1:
		exit += entry
	case 2:
		exit += entry + 1
	}
	b.outputs[0].d = exit % 4
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
