package main

import (
	"embed"
	"encoding/hex"
	"flag"
	"os"

	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/cli"
	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	locationRepo "github.com/theandrew168/dripfile/internal/location/repository"
	locationService "github.com/theandrew168/dripfile/internal/location/service"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/secret"
	transferRepo "github.com/theandrew168/dripfile/internal/transfer/repository"
)

//go:embed migration
var migrationFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stdout))

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

	locationRepo := locationRepo.New(pool)
	transferRepo := transferRepo.New(pool)
	locationService := locationService.New(box, locationRepo, transferRepo)

	cli := cli.New(locationService, flag.Args())
	err = cli.Run()
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
}
