package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/pkg/app"
	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/postgres"
	"github.com/theandrew168/dripfile/pkg/secret"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/stripe"
	"github.com/theandrew168/dripfile/pkg/task"
)

func main() {
	os.Exit(run())
}

func run() int {
	infoLog := log.New(os.Stdout, "", 0)
	errorLog := log.New(os.Stderr, "error: ", 0)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	secretKeyBytes, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	// create secret.Box
	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	box := secret.NewBox(secretKey)

	// open a database connection pool
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		errorLog.Println(err)
		return 1
	}
	defer pool.Close()

	store := storage.New(pool)
	queue := task.NewQueue(pool)

	// init the stripe interface
	var stripeI stripe.Interface
	if cfg.StripeSecretKey != "" {
		stripeI = stripe.New(
			cfg.StripeSecretKey,
			cfg.SiteURL+"/stripe/success",
			cfg.SiteURL+"/stripe/cancel",
		)
	} else {
		stripeI = stripe.NewMock(infoLog)
	}

	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	handler := app.New(cfg, box, store, queue, stripeI, infoLog, errorLog)

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
		errorLog.Println(err)
		return 1
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	infoLog.Printf("started web server on %s\n", addr)

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
		infoLog.Println("stopping server")
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
		errorLog.Println(err)
		return 1
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	infoLog.Println("stopped server")
	return 0
}
