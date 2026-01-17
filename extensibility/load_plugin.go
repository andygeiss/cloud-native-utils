package extensibility

import (
	"errors"
	"fmt"
	"plugin"
)

// LoadPlugin loads a plugin from the given path and returns the symbol with the given name.
func LoadPlugin[T any](path string, symName string) (T, error) {
	var zero T
	// Open the plugin.
	p, err := plugin.Open(path)
	if err != nil {
		return zero, err
	}
	if p == nil {
		return zero, errors.New("plugin open returned nil plugin")
	}

	// Lookup the symbol (an exported function or variable).
	symbol, err := p.Lookup(symName)
	if err != nil {
		return zero, err
	}

	// Assert that loaded symbol is of a correct type.
	cast, ok := symbol.(T)
	if !ok {
		return zero, fmt.Errorf("symbol %q has unexpected type", symName)
	}
	return cast, nil
}
