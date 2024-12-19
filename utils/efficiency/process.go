package efficiency

import (
	"context"
	"runtime"
	"sync"

	"github.com/andygeiss/cloud-native/utils/service"
)

// Process concurrently processes items from the input channel using the provided function `fn`.
// It spawns a number of worker goroutines equal to the number of available CPU cores.
func Process[IN, OUT any](in <-chan IN, fn service.Function[IN, OUT]) (<-chan OUT, <-chan error) {
	out := make(chan OUT)
	errCh := make(chan error)
	ctx := context.Background()
	// Launch `num` worker goroutines.
	num := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(num)
	for range num {
		go func() {
			defer wg.Done()
			// Process items from the input channel.
			for val := range in {
				// Call the processing function `fn` with the current value.
				res, err := fn(ctx, val)
				// If an error occurs, send it to the error channel and stop processing.
				if err != nil {
					errCh <- err
					return
				}
				// Send the processed result to the output channel.
				out <- res
			}
		}()
	}
	// Start a goroutine to close the output channel after all workers finish.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out, errCh
}
