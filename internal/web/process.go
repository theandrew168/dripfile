package web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/process"
)

type webProcess struct {
	logger *jsonlog.Logger
	server *http.Server
}

func NewProcess(logger *jsonlog.Logger, addr string, handler http.Handler) process.Process {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	p := webProcess{
		logger: logger,
		server: server,
	}
	return &p
}

func (p *webProcess) Run(ctx context.Context) error {
	p.logger.Info("starting server", map[string]string{
		"addr": p.server.Addr,
	})

	// start a goro to watch for stop signal (context cancelled)
	stopError := make(chan error)
	go func() {
		<-ctx.Done()

		// give the web server 5 seconds to shutdown gracefully
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// disable keepalives and shutdown gracefully
		p.logger.Info("stopping web server", nil)
		p.server.SetKeepAlivesEnabled(false)
		err := p.server.Shutdown(timeout)
		if err != nil {
			stopError <- err
		}

		close(stopError)
	}()

	// listen and serve forever
	// ignore http.ErrServerClosed (expected upon stop)
	err := p.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// check for errors that arose while stopping
	err = <-stopError
	if err != nil {
		return err
	}

	p.logger.Infof("stopped web server")
	return nil
}
