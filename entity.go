package main

import "fmt"

// DisplayMode indicates an entity's display mode
type DisplayMode = uint8

const (
	// DisplayModeMap default representation on the map
	DisplayModeMap DisplayMode = iota

	// DisplayModePreview representation in the preview section
	DisplayModePreview

	// DisplayModeGhostValid valid representation when placing on the map
	DisplayModeGhostValid

	// DisplayModeGhostInvalid invalid representation when placing on the map
	DisplayModeGhostInvalid
)

// Tile is a component of the game map
type Tile interface {
	Display() string
}

// StructureTile one Tile component of a Structure
type StructureTile interface {
	Tile
	Group() Structure
	UnderlyingResource() *Resource
	SetUnderlyingResource(*Resource)
}

// Structure is a Tile containing a player created entity
type Structure interface {
	Tiles() [][]StructureTile
	RotateLeft()
	RotateRight()
	Copy() Structure
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
func (b *Belt) Display() string {
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

	return fmt.Sprintf("\033[33;40m%c\033[0m", symbol)
}

// Group returns the Structure the Belt is associated with, itself
func (b *Belt) Group() Structure {
	return b
}

// Tiles return the Tiles associated with the Belt, an array of itself
func (b *Belt) Tiles() [][]StructureTile {
	return [][]StructureTile{{b}}
}

// Copy creates a deep copy of the current Belt
func (b *Belt) Copy() Structure {
	return &Belt{b.index, nil}
}

// Resource is a Tile containing natural resources
type Resource struct {
	amount   int
	resource int
}

// Display displays the Resource tile
func (t *Resource) Display() string {
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

	return fmt.Sprintf("\033[33;40m%c\033[0m", symbol)
}
