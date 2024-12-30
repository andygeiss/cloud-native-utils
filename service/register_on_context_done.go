package service

import (
	"context"
	"time"
)

// RegisterOnContextDone waits for the context to be done
// and then calls the function.
func RegisterOnContextDone(ctx context.Context, fn func()) {
	go func() {
		<-ctx.Done()
		// Wait for the readiness check to fail.
		<-time.After(5 * time.Second)
		// Call the function.
		fn()
	}()
}
