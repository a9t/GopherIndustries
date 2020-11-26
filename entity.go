package main

import "fmt"

// ShowMode indicates an entity's display mode
type ShowMode = uint8

const (
	// ShowModeMap default representation on the map
	ShowModeMap ShowMode = iota

	// ShowModePreview representation in the preview section
	ShowModePreview

	// ShowModeGhostValid valid representation when placing on the map
	ShowModeGhostValid

	// ShowModeGhostInvalid invalid representation when placing on the map
	ShowModeGhostInvalid
)

// Tile is a component of the game map
type Tile interface {
	Show() string
}

// Structure is a Tile containing a player created entity
type Structure interface {
	Tile
	RotateLeft()
	RotateRight()
	UnderlyingResource() *Resource
	SetUnderlyingResource(*Resource)
	Copy() Structure
}

// Resource is a Tile containing natural resources
type Resource struct {
	amount   int
	resource int
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

// Show displays the belt
func (b *Belt) Show() string {
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
		symbol = '\u257B'
	case 10:
		symbol = '\u2512'
	case 11:
		symbol = '\u251A'
	}

	return fmt.Sprintf("\033[33;40m%c\033[0m", symbol)
}

// Copy creates a deep copy of the current Belt
func (b *Belt) Copy() Structure {
	return &Belt{b.index, nil}
}

// Belt is a structure that transports resources
type Belt struct {
	index    int
	Resource *Resource
}

// Show displays the Resource tile
func (t *Resource) Show() string {
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
