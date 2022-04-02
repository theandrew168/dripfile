package main

import (
	"encoding/hex"
	"flag"
	"log"
	"os"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/database"
	"github.com/theandrew168/dripfile/pkg/postgres"
	"github.com/theandrew168/dripfile/pkg/postmark"
	"github.com/theandrew168/dripfile/pkg/secret"
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

	storage := database.NewStorage(pool)
	queue := task.NewQueue(pool)

	var postmarkI postmark.Interface
	if cfg.PostmarkAPIKey != "" {
		postmarkI = postmark.New(cfg.PostmarkAPIKey)
	} else {
		postmarkI = postmark.NewMock(infoLog)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the worker forever
	worker := task.NewWorker(box, queue, storage, postmarkI, infoLog, errorLog)
	err = worker.Run()
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	return 0
}
