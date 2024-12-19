package efficiency

import "sync"

// Merge combines multiple input channels into a single output channel.
// It starts a goroutine for each input channel to forward its values to the output channel.
// Once all input channels are processed, the output channel is closed.
func Merge[T any](in ...<-chan T) chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(in))

	// Start a goroutine for each input channel.
	for _, ch := range in {

		// Forward all values from the input channel to the output channel.
		go func(ch <-chan T) {
			defer wg.Done()
			for val := range ch {
				out <- val
			}
		}(ch)
	}

	// Wait for all goroutines to finish.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
