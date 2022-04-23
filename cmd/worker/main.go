package main

import (
	"encoding/hex"
	"flag"
	"os"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/database"
	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/mail"
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

	var mailer mail.Mailer
	if cfg.PostmarkAPIKey != "" {
		mailer = mail.NewPostmarkMailer(cfg.PostmarkAPIKey)
	} else {
		mailer = mail.NewMockMailer(logger)
	}

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

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the worker forever
	worker := task.NewWorker(logger, store, queue, box, billing, mailer)
	err = worker.Run()
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	return 0
}
