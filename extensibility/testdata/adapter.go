package main

// TestAdapter is a test adapter.
type TestAdapter struct{}

// FindByID returns a name by id.
func (a *TestAdapter) FindByID(id string) (name string, err error) {
	return "Andy", nil
}

// Export specific symbol.
var Adapter TestAdapter
