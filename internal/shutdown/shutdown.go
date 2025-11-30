package shutdown

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var ErrGracefullyShutdown = errors.New("error graceful shutdown")

// Receive will get termination signals and gracefuly exit application
func Receive(ctx context.Context) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-stop:
			return ErrGracefullyShutdown
		}
	}
}
