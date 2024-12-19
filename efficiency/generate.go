package efficiency

// Generate takes a variadic input of any type T and returns a read-only channel of type T.
// It sends each input value into the returned channel in a separate goroutine.
func Generate[T any](in ...T) <-chan T {
	out := make(chan T)
	// Start a goroutine to send the input values into the channel.
	go func() {
		defer close(out)
		for _, val := range in {
			out <- val
		}
	}()
	return out
}
