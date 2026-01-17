package service

import (
	"context"
	"os/signal"
	"syscall"
)

// Context creates a new context with a cancel function that listens for
// SIGTERM, SIGINT, SIGQUIT, and SIGKILL signals.
func Context() (context.Context, context.CancelFunc) {
	// Create a new context with a cancel function.
	return signal.NotifyContext(
		context.Background(),
		// SIGTERM is sent by Kubernetes to gracefully stop a container.
		syscall.SIGTERM,
		// SIGINT is sent by a user terminal to interrupt a running process.
		syscall.SIGINT,
		// SIGQUIT is sent by a user terminal to make a core dump.
		syscall.SIGQUIT,
		// SIGKILL is sent by a user terminal to kill a process immediately.
		syscall.SIGKILL,
	)
}
