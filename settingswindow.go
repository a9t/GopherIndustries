package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// SettingsWindow a Window that handles the game settings
type SettingsWindow struct {
	manager WindowManager
	widgets []Widget
}

// NewSettingsWindow creates a new SettingsWindow
func NewSettingsWindow(manager WindowManager) *SettingsWindow {
	var w SettingsWindow
	w.manager = manager

	colorMenuWidget := &ColorMenuWidget{"ColorMenu", 0, nil, true}
	symbolMenuWidget := &SymbolMenuWidget{"SymbolMenu", 0, nil, false}
	sampleWidget := newSampleWidget()

	colorMenuWidget.other = symbolMenuWidget
	symbolMenuWidget.other = sampleWidget
	sampleWidget.other = colorMenuWidget

	w.widgets = append(w.widgets, colorMenuWidget)
	w.widgets = append(w.widgets, symbolMenuWidget)
	w.widgets = append(w.widgets, sampleWidget)

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
	other *SampleWidget
	focus bool
}

// Layout displays the SymbolMenuWidget
func (w *SymbolMenuWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, 1, 5, 12, 8)

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
	name   string
	sel    int
	focus  bool
	r      []*RawResource
	s      [][]Structure
	sNames []string

	other *ColorMenuWidget
}

func newSampleWidget() *SampleWidget {
	w := new(SampleWidget)
	w.name = "SampleWidget"
	w.sel = -1

	w.r = []*RawResource{&RawResource{0, 1}, &RawResource{1, 1}, &RawResource{2, 1}, &RawResource{3, 1}}

	w.s = make([][]Structure, 3)
	w.s[0] = make([]Structure, 12)
	for i := range w.s[0] {
		w.s[0][i] = NewBelt()
	}

	for i, b := range w.s[0] {
		for j := i; j < len(w.s[0]); j++ {
			b.RotateRight()
		}
	}

	w.s[1] = make([]Structure, 1)
	w.s[1][0] = NewTwoXTwoBlock()

	w.s[2] = make([]Structure, 1)
	w.s[2][0] = NewExtractor()

	w.sNames = []string{"Belt", "TwoXTwo", "Extractor"}

	return w
}

// Layout displays the SampleWidget
func (w *SampleWidget) Layout(g *gocui.Gui) error {
	resurceModes := []DisplayMode{DisplayModeMap, DisplayModeMapSelected}
	structureModes := []DisplayMode{DisplayModeMap, DisplayModeMapSelected, DisplayModeGhostValid, DisplayModeGhostInvalid}

	v, err := g.SetView(w.name, 20, 1, 40, 20)

	if err == gocui.ErrUnknownView {
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				w.sel--
				if w.sel < -1 {
					w.sel = -1
				}

				return nil
			}); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				size := len(w.s)
				w.sel++
				if w.sel > size-1 {
					w.sel = size - 1
				}

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
		if _, err := g.SetViewOnTop(w.name); err != nil {
			return err
		}

		if w.focus {
			if _, err := g.SetCurrentView(w.name); err != nil {
				return err
			}
		}

		if w.sel == -1 {
			colorPrefix := ""
			colorSuffix := ""
			if w.focus {
				colorPrefix = "\033[31;1m"
				colorSuffix = "\033[0m"
			}

			fmt.Fprintf(v, "%s>%s Resources\n", colorPrefix, colorSuffix)

			for _, mode := range resurceModes {
				fmt.Fprintf(v, "  ")
				for _, resource := range w.r {
					fmt.Fprintf(v, "%s", resource.Display(mode))
				}
				fmt.Fprintf(v, "\n")
			}
			fmt.Fprintf(v, "\n")
		}

		for i, structures := range w.s {
			if i < w.sel {
				continue
			}

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

			fmt.Fprintf(v, "%s%c%s %s\n", colorPrefix, prefix, colorSuffix, w.sNames[i])

			for _, mode := range structureModes {
				tileMatrix := w.getTiles(structures)
				for _, tiles := range tileMatrix {
					for _, tile := range tiles {
						var msg string
						if tile != nil {
							msg = tile.Display(mode)
						} else {
							msg = " "
						}

						fmt.Fprintf(v, "%s", msg)
					}
					fmt.Fprintf(v, "\n")
				}
			}
			fmt.Fprintf(v, "\n")
		}
	}

	return nil
}

func (w *SampleWidget) getTiles(structures []Structure) [][]Tile {
	totalW := 0
	maxH := 0
	for _, s := range structures {
		tiles := s.Tiles()
		h := len(tiles)
		w := len(tiles[0])

		if h > maxH {
			maxH = h
		}

		totalW += w
	}

	responseTiles := make([][]Tile, maxH)
	for i := range responseTiles {
		responseTiles[i] = make([]Tile, totalW)
	}

	crtW := 0
	for _, structure := range structures {
		tileMatrix := structure.Tiles()
		for i, tiles := range tileMatrix {
			for j, tile := range tiles {
				responseTiles[i][crtW+j] = tile
			}
		}
		crtW += len(tileMatrix[0])
	}

	return responseTiles
}
