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
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	code := 0
	if err := run(logger); err != nil {
		logger.Error(err.Error())
		code = 1
	}

	os.Exit(code)
}

func run(logger *slog.Logger) error {
	conf := flag.String("conf", "dripfile.conf", "app config file")
	migrateOnly := flag.Bool("migrate", false, "apply migrations and exit")
	flag.Parse()

	cfg, err := config.ReadFile(*conf)
	if err != nil {
		return err
	}

	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		return err
	}
	defer pool.Close()

	applied, err := migrate.Migrate(pool, migrationFS)
	if err != nil {
		return err
	}

	for _, migration := range applied {
		logger.Info("applied migration", "name", migration)
	}

	if *migrateOnly {
		return nil
	}

	secretKey, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		return err
	}

	box := secret.NewBox([32]byte(secretKey))

	repo := repository.NewPostgres(pool, box)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// create a context that cancels upon receiving an interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// create a WaitGroup with an initial counter of two:
	// 1. web server
	// 2. worker
	var wg sync.WaitGroup
	wg.Add(2)

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
	go func() {
		defer wg.Done()
		app.Run(ctx, addr)
	}()

	w := worker.New(logger, repo)

	// start worker in the background (standalone mode by default)
	go func() {
		defer wg.Done()
		w.Run(ctx)
	}()

	// wait for the worker and web server to stop
	wg.Wait()

	return nil
}
