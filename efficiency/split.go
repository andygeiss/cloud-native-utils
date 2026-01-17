package efficiency

// Split takes an input channel `in` and an integer `num`.
// It returns a slice of `num` output channels that distribute items from
// the input channel. Each input item is sent to exactly one output channel,
// enabling concurrent processing by multiple consumers (fan-out pattern).
func Split[T any](in <-chan T, num int) []<-chan T {
	// Iterate `num` times to create the specified number of channels.
	out := make([]<-chan T, 0)
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
