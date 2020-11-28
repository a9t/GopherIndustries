package main

import (
	"github.com/jroimartin/gocui"
)

// Widget abstractization of a UI component
type Widget interface {
	Layout(g *gocui.Gui) error
}

// GameWidget a Widget with knowledge about the Game
type GameWidget interface {
	Widget
	SetGame(*Game)
}

// Window a top level UI compoment that fills the entire terminal
type Window interface {
	Widget
}

// WindowManager manager of Window instances of an application
type WindowManager interface {
	SetTopWindow(Window)
}
