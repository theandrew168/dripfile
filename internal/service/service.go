package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
)

type Service interface {
	// start the service, blocking
	Start() error

	// stop the service, unblock Start()
	Stop() error
}

func Run(s Service) error {
	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// kick off a goroutine to listen for SIGINT and SIGTERM
	stopError := make(chan error)
	go func() {
		// idle until a signal is caught
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// stop the service and track any errors
		err := s.Stop()
		if err != nil {
			stopError <- err
		}

		stopError <- nil
	}()

	// start the service
	err := s.Start()
	if err != nil {
		return err
	}

	// check for stop errors
	err = <-stopError
	if err != nil {
		return err
	}

	return nil
}
