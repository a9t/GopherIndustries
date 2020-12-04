package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// SettingsWindow a Window that handles the game settings
type SettingsWindow struct {
	manager    WindowManager
	widgets    []Widget
	focusIndex int
}

// NewSettingsWindow creates a new SettingsWindow
func NewSettingsWindow(manager WindowManager) *SettingsWindow {
	var w SettingsWindow
	w.manager = manager

	colorMenuWidget := &ColorMenuWidget{"ColorMenu", 0, nil, true}
	symbolMenuWidget := &SymbolMenuWidget{"SymbolMenu", 0, nil, false}

	colorMenuWidget.other = symbolMenuWidget
	symbolMenuWidget.other = colorMenuWidget

	w.widgets = append(w.widgets, colorMenuWidget)
	w.widgets = append(w.widgets, symbolMenuWidget)

	return &w
}

// Layout displays the NewSettingsWindow
func (w *SettingsWindow) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = false

	v, err := g.SetView("SettingsWindow", 0, 0, maxX-1, maxY-1)
	if err == nil || err == gocui.ErrUnknownView {
		if _, err := g.SetViewOnTop("SettingsWindow"); err != nil {
			return err
		}

		v.Title = "Settings"
	}

	for _, widget := range w.widgets {
		widget.Layout(g)
	}

	return nil
}

// ColorMenuWidget allows the selection of new color modes
type ColorMenuWidget struct {
	name  string
	sel   int
	other *SymbolMenuWidget
	focus bool
}

// Layout displays the ColorMenuWidget
func (w *ColorMenuWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, 1, 1, 12, 4) //GlobalDisplayConfigManager

	if err == gocui.ErrUnknownView {
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				w.sel--
				if w.sel < 0 {
					w.sel = 0
				}

				GlobalDisplayConfigManager.SetColorConfig(w.sel)
				return nil
			}); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				size := len(GlobalDisplayConfigManager.ColorConfigs)
				w.sel++
				if w.sel > size-1 {
					w.sel = size - 1
				}

				GlobalDisplayConfigManager.SetColorConfig(w.sel)
				return nil
			}); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				w.focus = false
				w.other.focus = true

				return nil
			}); err != nil {
			return err
		}
	}

	if err == nil || err == gocui.ErrUnknownView {
		v.Clear()
		if w.focus {
			if _, err := g.SetCurrentView(w.name); err != nil {
				return err
			}
		}

		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		for i, colorConfig := range GlobalDisplayConfigManager.ColorConfigs {
			prefix := ' '
			if w.sel == i {
				prefix = '>'
			}

			colorPrefix := ""
			colorSuffix := ""
			if w.focus {
				colorPrefix = "\033[31;1m"
				colorSuffix = "\033[0m"
			}

			fmt.Fprintf(v, "%s%c%s %s\n", colorPrefix, prefix, colorSuffix, colorConfig.Name)
		}
	}

	return nil
}

// SymbolMenuWidget allows the selection of new symbol modes
type SymbolMenuWidget struct {
	name  string
	sel   int
	other *ColorMenuWidget
	focus bool
}

// Layout displays the SymbolMenuWidget
func (w *SymbolMenuWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, 20, 1, 30, 4) //GlobalDisplayConfigManager

	if err == gocui.ErrUnknownView {
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				w.sel--
				if w.sel < 0 {
					w.sel = 0
				}

				GlobalDisplayConfigManager.SetSymbolConfig(w.sel)
				return nil
			}); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				size := len(GlobalDisplayConfigManager.SymbolConfigs)
				w.sel++
				if w.sel > size-1 {
					w.sel = size - 1
				}

				GlobalDisplayConfigManager.SetSymbolConfig(w.sel)
				return nil
			}); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				w.focus = false
				w.other.focus = true

				return nil
			}); err != nil {
			return err
		}
	}

	if err == nil || err == gocui.ErrUnknownView {
		v.Clear()
		if w.focus {
			if _, err := g.SetCurrentView(w.name); err != nil {
				return err
			}
		}

		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		for i, colorConfig := range GlobalDisplayConfigManager.SymbolConfigs {
			prefix := ' '
			if w.sel == i {
				prefix = '>'
			}

			colorPrefix := ""
			colorSuffix := ""
			if w.focus {
				colorPrefix = "\033[31;1m"
				colorSuffix = "\033[0m"
			}

			fmt.Fprintf(v, "%s%c%s %s\n", colorPrefix, prefix, colorSuffix, colorConfig.Name)
		}
	}

	return nil
}

// SampleWidget sample display of the current color and symbol selections
type SampleWidget struct {
}
