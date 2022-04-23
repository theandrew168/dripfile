package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/pkg/app"
	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/database"
	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/secret"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/stripe"
	"github.com/theandrew168/dripfile/pkg/task"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := jsonlog.New(os.Stdout)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	secretKeyBytes, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	// create secret.Box
	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	box := secret.NewBox(secretKey)

	// open a database connection pool
	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}
	defer pool.Close()

	store := storage.New(pool)
	queue := task.NewQueue(pool)

	// init the stripe billing interface
	var billing stripe.Billing
	if cfg.StripeSecretKey != "" {
		billing = stripe.NewBilling(
			logger,
			cfg.StripeSecretKey,
			cfg.SiteURL+"/billing/success",
			cfg.SiteURL+"/billing/cancel",
		)
	} else {
		billing = stripe.NewMockBilling(logger)
	}

	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	handler := app.New(cfg, logger, store, queue, box, billing)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.PrintInfo("starting server", map[string]string{
		"addr": addr,
	})

	// kick off a goroutine to listen for SIGINT and SIGTERM
	shutdownError := make(chan error)
	go func() {
		// idle until a signal is caught
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// give the web server 5 seconds to shutdown gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown the web server and track any errors
		logger.PrintInfo("stopping server", nil)
		srv.SetKeepAlivesEnabled(false)
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	// serve the app, check for ErrServerClosed (expected after shutdown)
	err = srv.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		logger.PrintError(err, nil)
		return 1
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	logger.PrintInfo("stopped server", nil)
	return 0
}
