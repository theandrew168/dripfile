package main

import (
	"context"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/go-systemd/daemon"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/config"
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/history"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/migrate"
	"github.com/theandrew168/dripfile/backend/secret"
	"github.com/theandrew168/dripfile/backend/transfer"
	transferService "github.com/theandrew168/dripfile/backend/transfer/service"
	"github.com/theandrew168/dripfile/backend/web"
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

	secretKey, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	box := secret.NewBox([32]byte(secretKey))

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

	locationRepo := location.NewRepository(pool, box)
	transferRepo := transfer.NewRepository(pool)
	historyRepo := history.NewRepository(pool)

	transferService := transferService.New(locationRepo, transferRepo, historyRepo)

	fooID, _ := uuid.NewRandom()
	fooLoc, err := location.NewS3(
		fooID.String(),
		"localhost:9000",
		"foo",
		"minioadmin",
		"minioadmin",
	)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	err = locationRepo.Create(fooLoc)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	barID, _ := uuid.NewRandom()
	barLoc, err := location.NewS3(
		barID.String(),
		"localhost:9000",
		"bar",
		"minioadmin",
		"minioadmin",
	)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	err = locationRepo.Create(barLoc)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	tfID, _ := uuid.NewRandom()
	tf, err := transfer.New(
		tfID.String(),
		"*.png",
		fooID.String(),
		barID.String(),
	)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	err = transferRepo.Create(tf)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	err = transferService.Run(tf.ID())
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	app := web.NewApplication(
		logger,
		publicFS,
		locationRepo,
		transferRepo,
		historyRepo,
		transferService,
	)

	// let port be overridable by an env var
	port := cfg.Port
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	ctx := newSignalHandlerContext()
	err = app.Run(ctx, addr)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
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
