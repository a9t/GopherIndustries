package main

// GlobalDisplayConfigManager contains the global display configuration
var GlobalDisplayConfigManager = newDisplayConfigManager()

// DisplayConfigManager stores information about color and symbol configurations
type DisplayConfigManager struct {
	ColorConfigs  []*ColorConfig
	sColorConfig  *ColorConfig
	SymbolConfigs []*SymbolConfig
	sSymbolConfig *SymbolConfig
}

func newDisplayConfigManager() *DisplayConfigManager {
	m := new(DisplayConfigManager)

	eightColorConfig := new(ColorConfig)
	eightColorConfig.Name = "8 color"
	eightColorConfig.StructureColors = [][]int{
		{36, 1},
		{37, 1},
		{36, 1},
		{31, 1},
	}
	eightColorConfig.ResourceColors = []int{33}

	m.ColorConfigs = []*ColorConfig{eightColorConfig}

	asciiSymbolConfig := new(SymbolConfig)
	asciiSymbolConfig.Name = "ascii"
	asciiSymbolConfig.Types = map[string]string{
		"resource":       string([]rune{32, 188, 189, 190}),
		"belt":           string([]rune{226, 224, 225, 234, 232, 233, 238, 236, 237, 244, 242, 243}),
		"fillerCorner":   "/\\/\\",
		"fillerMid":      "~|_|",
		"fillerCenter":   string([]rune{219}),
		"input":          "V<A>",
		"output":         "V<A>",
		"chest":          "+",
		"splitterLeft":   string([]rune{200, 201, 187, 188}),
		"splitterRight":  string([]rune{217, 192, 218, 191}),
		"cornerTriangle": "/\\/\\",
	}

	unicodeSymbolConfig := new(SymbolConfig)
	unicodeSymbolConfig.Name = "FreeMono"
	unicodeSymbolConfig.Types = map[string]string{
		"resource":       " \u2591\u2592\u2593",
		"belt":           "\u2193\u21B2\u21B3\u2190\u2196\u2199\u2191\u21B1\u21B0\u2192\u2198\u2197",
		"fillerCorner":   "\u259B\u259C\u259F\u2599",
		"fillerMid":      "\u2580\u2590\u2584\u258C",
		"fillerCenter":   "\u2588",
		"input":          "\u21A5\u21A6\u21A7\u21A4",
		"output":         "\u21D3\u21D0\u21D1\u21D2",
		"chest":          "\u25A3",
		"splitterLeft":   "\u2558\u2553\u2555\u255C",
		"splitterRight":  "\u255B\u2559\u2552\u2556",
		"cornerTriangle": "\u25E2\u25E3\u25E4\u25E5",
	}

	m.SymbolConfigs = []*SymbolConfig{unicodeSymbolConfig, asciiSymbolConfig}

	m.sColorConfig = m.ColorConfigs[0]
	m.sSymbolConfig = m.SymbolConfigs[0]

	return m
}

// SymbolConfig a configuration containing the symbols used for StructureTile-s
type SymbolConfig struct {
	Name  string
	Types map[string]string
}

// ColorConfig a configuration containing the colors used for StructureTile-s
type ColorConfig struct {
	Name            string
	StructureColors [][]int
	ResourceColors  []int
}

// GetColorConfig returns the current ColorConfig
func (m *DisplayConfigManager) GetColorConfig() *ColorConfig {
	return m.sColorConfig
}

// SetColorConfig sets the current ColorConfig
func (m *DisplayConfigManager) SetColorConfig(index int) {
	size := len(m.ColorConfigs)
	if index < 0 {
		index = 0
	}

	if index >= size {
		index = size - 1
	}

	m.sColorConfig = m.ColorConfigs[index]
}

// GetSymbolConfig returns the current SymbolConfig
func (m *DisplayConfigManager) GetSymbolConfig() *SymbolConfig {
	return m.sSymbolConfig
}

// SetSymbolConfig sets the current SymbolConfig
func (m *DisplayConfigManager) SetSymbolConfig(index int) {
	size := len(m.SymbolConfigs)
	if index < 0 {
		index = 0
	}

	if index >= size {
		index = size - 1
	}

	m.sSymbolConfig = m.SymbolConfigs[index]
}
