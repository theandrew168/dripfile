package main

import (
	"encoding/hex"
	"flag"
	"os"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/postgres"
	"github.com/theandrew168/dripfile/pkg/postmark"
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
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}
	defer pool.Close()

	store := storage.New(pool)
	queue := task.NewQueue(pool)

	var postmarkI postmark.Interface
	if cfg.PostmarkAPIKey != "" {
		postmarkI = postmark.New(cfg.PostmarkAPIKey)
	} else {
		postmarkI = postmark.NewMock(logger)
	}

	var stripeI stripe.Interface
	if cfg.StripeSecretKey != "" {
		stripeI = stripe.New(
			logger,
			cfg.StripeSecretKey,
			cfg.SiteURL+"/billing/success",
			cfg.SiteURL+"/billing/cancel",
		)
	} else {
		stripeI = stripe.NewMock(logger)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the worker forever
	worker := task.NewWorker(logger, store, queue, box, stripeI, postmarkI)
	err = worker.Run()
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	return 0
}
