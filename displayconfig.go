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

	unicodeSymbolConfig := new(SymbolConfig)
	unicodeSymbolConfig.Name = "unicode"
	unicodeSymbolConfig.Types = map[string]string{
		"resource":     " \u2591\u2592\u2593",
		"belt":         "\u257D\u2519\u2515\u257E\u2516\u250E\u257F\u250D\u2511\u257C\u2512\u251A",
		"fillerCorner": "\u259B\u259C\u259F\u2599",
	}

	m.SymbolConfigs = []*SymbolConfig{unicodeSymbolConfig}

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
