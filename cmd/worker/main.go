package main

import (
	"context"
	"flag"
	"os"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/postgres"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/worker"
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

	// open a regular connection (for listen / notify)
	conn, err := postgres.Connect(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer conn.Close(context.Background())

	// open a connection pool (for everything else)
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer pool.Close()

	storage := database.NewPostgresStorage(pool)
	queue := task.NewPostgresQueue(conn, pool)

	w := worker.New(storage, queue, logger)

	// run the worker forever
	err = w.Run()
	if err != nil {
		logger.Error(err)
		return 1
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	return 0
}
