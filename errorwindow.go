package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// ErrorWindow a Window that is displayed in case of error covering the entire terminal
type ErrorWindow struct {
	message string
}

// Layout displays the ErrorWindow
func (w *ErrorWindow) Layout(g *gocui.Gui) error {
	g.Cursor = false
	g.Mouse = false

	maxX, maxY := g.Size()

	v, err := g.SetView("Error", 0, 0, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetCurrentView("Error"); err != nil {
			return err
		}

		if _, err := g.SetViewOnTop("Error"); err != nil {
			return err
		}

		v.Title = "Error"
	}

	posY := maxY / 2
	posX := maxX / 2
	msgLen := len(w.message)

	v, err = g.SetView("ErrorMessage", posX-msgLen/2-1, posY-1, posX+msgLen/2, posY+1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop("ErrorMessage"); err != nil {
			return err
		}

		v.Frame = false
		v.Clear()

		fmt.Fprintln(v, w.message)
	}

	return nil
}
