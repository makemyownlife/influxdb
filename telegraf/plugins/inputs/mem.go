package inputs

import (
	"fmt"
)

// MemStats is based on telegraf MemStats.
type MemStats struct {
	baseInput
}

// PluginName is based on telegraf plugin name.
func (m *MemStats) PluginName() string {
	return "mem"
}

// TOML encodes to toml string
func (m *MemStats) TOML() string {
	return fmt.Sprintf(`[[inputs.%s]]
`, m.PluginName())
}
