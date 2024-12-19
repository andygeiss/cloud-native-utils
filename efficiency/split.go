package efficiency

// Split takes an input channel `in` and an integer `num`.
// It returns a slice of `num` output channels, where each channel will receive
// the same data from the input channel concurrently.
func Split[T any](in <-chan T, num int) []<-chan T {
	out := make([]<-chan T, 0)

	// Iterate `num` times to create the specified number of channels.
	for range num {
		ch := make(chan T)
		out = append(out, ch)

		// Start a goroutine to forward data from the input channel to the current output channel.
		go func(ch chan T) {
			defer close(ch)
			for val := range in {
				ch <- val
			}
		}(ch)
	}
	return out
}
