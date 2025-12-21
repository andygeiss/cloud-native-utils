package extensibility

import (
	"fmt"
	"plugin"
)

// LoadPlugin loads a plugin from the given path and returns the symbol with the given name.
func LoadPlugin[T any](path string, symName string) (res T, err error) {
	// Open the plugin.
	p, err := plugin.Open(path)
	if err != nil {
		return res, err
	}
	if p == nil {
		return res, fmt.Errorf("plugin open returned nil plugin")
	}

	// Lookup the symbol (an exported function or variable).
	symbol, err := p.Lookup(symName)
	if err != nil {
		return res, err
	}

	// Assert that loaded symbol is of a correct type.
	cast, ok := symbol.(T)
	if !ok {
		return res, fmt.Errorf("symbol %q has unexpected type", symName)
	}
	return cast, nil
}
