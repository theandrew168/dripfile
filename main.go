package main

import (
	"context"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/coreos/go-systemd/daemon"

	"github.com/theandrew168/dripfile/backend/config"
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/migrate"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/secret"
	"github.com/theandrew168/dripfile/backend/web"
	"github.com/theandrew168/dripfile/backend/worker"
)

//go:embed migration
var migrationFS embed.FS

//go:embed public
var publicFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	conf := flag.String("conf", "dripfile.conf", "app config file")
	migrateOnly := flag.Bool("migrate", false, "apply migrations and exit")
	flag.Parse()

	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}
	defer pool.Close()

	applied, err := migrate.Migrate(pool, migrationFS)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	for _, migration := range applied {
		logger.Info("applied migration", "name", migration)
	}

	if *migrateOnly {
		return 0
	}

	secretKey, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	box := secret.NewBox([32]byte(secretKey))

	repo := repository.NewPostgres(pool, box)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a context that cancels upon receiving an interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	app := web.NewApplication(
		logger,
		publicFS,
		repo,
	)

	// let port be overridden by an env var
	port := cfg.Port
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, port)

	// start the web server in the background
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := app.Run(ctx, addr)
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	w := worker.New(logger, repo)

	// start worker in the background (standalone mode by default)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := w.Run(ctx)
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	// wait for the worker and web server to stop
	wg.Wait()

	return 0
}
