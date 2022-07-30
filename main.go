package main

import (
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/scheduler"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/service"
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
		s := scheduler.New(logger, store, queue)
		err := service.Run(s)
		if err != nil {
			logger.Error(err, nil)
			return 1
		}
		return 0
	}

	// worker: run worker forever
	if action == "worker" {
		w := task.NewWorker(logger, store, queue, box, mailer)
		err := service.Run(w)
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

	// nest the API handler under the main web app
	addr := fmt.Sprintf("127.0.0.1:%s", port)
	handler := app.Handler()

	s := web.NewService(logger, addr, handler)
	err = service.Run(s)
	if err != nil {
		logger.Error(err, nil)
		return 1
	}

	return 0
}
