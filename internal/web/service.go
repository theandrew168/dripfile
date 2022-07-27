package web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/service"
)

type webService struct {
	logger *jsonlog.Logger
	server *http.Server
}

func NewService(logger *jsonlog.Logger, addr string, handler http.Handler) service.Service {
	server := http.Server{
		Addr:    addr,
		Handler: handler,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	s := webService{
		logger: logger,
		server: &server,
	}
	return &s
}

func (s *webService) Start() error {
	s.logger.Info("starting server", map[string]string{
		"addr": s.server.Addr,
	})

	// listen and serve forever
	err := s.server.ListenAndServe()
	if err != nil {
		// ignore http.ErrServerClosed (expected upon stop)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}

func (s *webService) Stop() error {
	// give the web server 5 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.logger.Info("stopping web server", nil)
	s.server.SetKeepAlivesEnabled(false)
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	s.logger.Infof("stopped web server")
	return nil
}
