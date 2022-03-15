package main

import (
	"flag"
	"os"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/postgres"
	"github.com/theandrew168/dripfile/internal/task"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger(os.Stdout)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Error(err)
		return 1
	}

	// open a database connection pool
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer pool.Close()

	storage := database.NewPostgresStorage(pool)
	queue := task.NewPostgresQueue(pool)

	var mailer mail.Mailer
	if cfg.PostmarkAPIKey != "" {
		mailer = mail.NewPostmarkMailer(cfg.PostmarkAPIKey)
	} else {
		mailer = mail.NewLogMailer(logger)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the worker forever
	worker := task.NewWorker(queue, storage, mailer, logger)
	err = worker.Run()
	if err != nil {
		logger.Error(err)
		return 1
	}

	return 0
}
