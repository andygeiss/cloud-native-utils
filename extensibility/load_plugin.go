package extensibility

import "plugin"

// LoadPlugin loads a plugin from the given path and returns the symbol with the given name.
func LoadPlugin[T any](path string, symName string) (res T, err error) {
	// Open the plugin.
	plugin, err := plugin.Open(path)
	if err != nil {
		return res, err
	}

	// Lookup the symbol (an exported function or variable).
	symbol, err := plugin.Lookup(symName)
	if err != nil {
		return res, err
	}

	// Assert that loaded symbol is of a correct type.
	res = symbol.(T)
	return res, err
}
