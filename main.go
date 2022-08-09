package main

import (
	"context"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/scheduler"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/web"
)

//go:embed migration
var migrationFS embed.FS

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

	pool, err := postgresql.ConnectPool(cfg.PostgreSQLURL)
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
		err := migrate.Migrate(logger, pool, migrationFS)
		if err != nil {
			logger.Error(err, nil)
			return 1
		}
		return 0
	}

	store := storage.New(pool)
	queue := task.NewQueue(pool)

	// init the mailer interface
	var mailer mail.Mailer
	if cfg.SMTPURL != "" {
		mailer, err = mail.NewSMTPMailer(cfg.SMTPURL)
	} else {
		logger.Infof("using mock mailer")
		mailer, err = mail.NewMockMailer(logger)
	}
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	// scheduler: run scheduler forever
	if action == "scheduler" {
		// let systemd know that we are good to go (no-op if not using systemd)
		daemon.SdNotify(false, daemon.SdNotifyReady)

		s := scheduler.New(logger, store, queue)
		ctx := newSignalHandlerContext()
		err := s.Run(ctx)
		if err != nil {
			logger.Error(err, nil)
			return 1
		}
		return 0
	}

	// worker: run worker forever
	if action == "worker" {
		// let systemd know that we are good to go (no-op if not using systemd)
		daemon.SdNotify(false, daemon.SdNotifyReady)

		w := task.NewWorker(logger, store, queue, box, mailer)
		ctx := newSignalHandlerContext()
		err := w.Run(ctx)
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

	// instantiate main web application
	app := web.NewApplication(logger, store, queue, box)

	// let port be overridable by an env var
	port := cfg.Port
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	ctx := newSignalHandlerContext()
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	err = app.Run(ctx, addr)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	return 0
}

func newSignalHandlerContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// idle until a signal is caught (must be a buffered channel)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		cancel()
	}()

	return ctx
}
