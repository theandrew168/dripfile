package main

import (
	"flag"
	"os"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/postgres"
	"github.com/theandrew168/dripfile/internal/pubsub"
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
	conn, err := postgres.Connect(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer conn.Close()

	storage := database.NewPostgresStorage(conn)
	queue := pubsub.NewPostgresQueue(conn, storage)

	// simulate a single job
	transfer, err := queue.Transfer.Subscribe()
	if err != nil {
		logger.Error(err)
		return 1
	}

	logger.Info("transfer %s start\n", transfer.ID)
	logger.Info("transfer %s end\n", transfer.ID)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	return 0
}
