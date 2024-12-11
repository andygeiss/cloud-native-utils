package concurrency

// Split takes an input channel `in` and splits its values into multiple output channels.
// The number of output channels is specified by `num`.
// Each output channel receives all the values from the input channel.
func Split[T any](in <-chan T, num int) (out []chan T) {
	out = make([]chan T, 0)
	for i := 0; i < num; i++ {
		ch := make(chan T)
		out = append(out, ch)
		// Launch a goroutine for each output channel.
		go func(i int) {
			defer close(ch)
			for val := range in {
				out[i] <- val
			}
		}(i)
	}
	return out
}
