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
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/internal/api"
	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/scheduler"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/web"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := jsonlog.New(os.Stdout)

	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	secretKeyBytes, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	box := secret.NewBox(secretKey)

	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}
	defer pool.Close()

	// check for action (default web)
	args := flag.Args()
	var action string
	if len(args) == 0 {
		action = "web"
	} else {
		action = args[0]
	}

	// migrate: apply migrations and exit
	if action == "migrate" {
		err := migrate.Migrate(logger, pool)
		if err != nil {
			logger.Error(err, nil)
			return 1
		}
		return 0
	}

	store := storage.New(pool)

	redis, err := asynq.ParseRedisURI(cfg.RedisURI)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	queue := asynq.NewClient(redis)
	defer queue.Close()

	// init the mailer interface
	var mailer mail.Mailer
	if cfg.PostmarkAPIKey != "" {
		mailer = mail.NewPostmarkMailer(cfg.PostmarkAPIKey)
	} else {
		logger.Infof("using mock mailer")
		mailer = mail.NewMockMailer(logger)
	}

	// scheduler: run scheduler forever
	if action == "scheduler" {
		sched := scheduler.New(cfg, logger, store, queue)
		err := sched.Run()
		if err != nil {
			logger.Error(err, nil)
			return 1
		}
		return 0
	}

	// worker: run worker forever
	if action == "worker" {
		worker := task.NewWorker(cfg, logger, store, box, mailer)
		err := worker.Run()
		if err != nil {
			logger.Error(err, nil)
			return 1
		}
		return 0
	}

	// web: run web server forever (default)
	if action != "web" {
		logger.Errorf("invalid action: %s", action)
		return 1
	}

	// instantiate applications (DI happens here)
	apiApp := api.NewApplication(logger)
	webApp := web.NewApplication(cfg, logger, store, queue, box)

	// nest the API handler under the main web app
	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: webApp.Handler(apiApp.Handler()),

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.Info("starting server", map[string]string{
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
		logger.Info("stopping server", nil)
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
		logger.Error(err, nil)
		return 1
	}

	// check for shutdown errors
	err = <-shutdownError
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	logger.Infof("stopped server")
	return 0
}
