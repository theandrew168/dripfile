package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/config"
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/migrate"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/web"
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

	// secretKey, err := hex.DecodeString(cfg.SecretKey)
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	return 1
	// }

	// box := secret.NewBox([32]byte(secretKey))

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

	repo := repository.NewMemory()

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

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	ctx := newSignalHandlerContext()
	err = app.Run(ctx, addr)
	if err != nil {
		return err
	}

	return nil
}

// create a context that cancels upon receiving an exit signal
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
