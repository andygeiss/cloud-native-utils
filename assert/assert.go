// Package assert provides tools for testing, including utility functions
// to assert value equality and simplify debugging during development.
//
// [That] is the entry point. It takes the *testing.T of the test that calls it,
// so it has no runnable example; it reads like this:
//
//	func Test_Add_With_TwoAndTwo_Should_ReturnFour(t *testing.T) {
//		// Act
//		sum := add(2, 2)
//
//		// Assert
//		assert.That(t, "sum must be 4", sum, 4)
//	}
package assert
