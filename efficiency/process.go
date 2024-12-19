package efficiency

import (
	"context"
	"runtime"
	"sync"

	"github.com/andygeiss/cloud-native-utils/service"
)

// Process processes items from `in` using `fn`, running workers equal to CPU cores.
// Errors are sent to the returned error channel.
func Process[IN any](in <-chan IN, fn service.Function[IN]) <-chan error {
	errCh := make(chan error)
	ctx := context.Background()
	num := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(num)

	// Run a gouroutine for each CPU.
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for val := range in {
				if err := fn(ctx, val); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	// Wait until everything is done.
	go func() {
		wg.Wait()
		close(errCh)
	}()
	return errCh
}

// Process2 processes items from `in` using `fn`, running workers equal to CPU cores.
// Results are sent to `out`, and errors to `errCh`.
func Process2[IN, OUT any](in <-chan IN, fn service.Function2[IN, OUT]) (<-chan OUT, <-chan error) {
	out := make(chan OUT)
	errCh := make(chan error)
	ctx := context.Background()
	num := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(num)

	// Run a gouroutine for each CPU.
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for val := range in {
				res, err := fn(ctx, val)
				if err != nil {
					errCh <- err
					return
				}
				out <- res
			}
		}()
	}

	// Wait until everything is done.
	go func() {
		wg.Wait()
		close(out)
		close(errCh)
	}()
	return out, errCh
}
