package process

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
)

type Process interface {
	// run until the context is cancelled
	Run(ctx context.Context) error
}

func Run(p Process) error {
	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// kick off a goroutine to listen for SIGINT and SIGTERM
	go func() {
		// idle until a signal is caught (must be a buffered channel)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// stop the process
		cancel()
	}()

	// run the process
	err := p.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}
