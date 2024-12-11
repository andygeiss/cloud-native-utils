package concurrency

import "sync"

// Merge is a generic function that takes multiple input channels of the same type (T)
// and merges them into a single output channel. It uses goroutines and a WaitGroup to
// ensure that all input channels are processed before closing the output channel.
func Merge[T any](in ...chan T) (out chan T) {
	out = make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(in))
	// Launch a goroutine for each input channel.
	for _, ch := range in {
		go func(c <-chan T) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}
	// Launch a goroutine to wait for all input channel processing to complete.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
