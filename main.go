package main

import (
	"context"
	"embed"
	"encoding/hex"
	"flag"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/cli"
	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/location"
	locationService "github.com/theandrew168/dripfile/internal/location/service"
	locationStorage "github.com/theandrew168/dripfile/internal/location/storage"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/transfer"
	transferService "github.com/theandrew168/dripfile/internal/transfer/service"
	transferStorage "github.com/theandrew168/dripfile/internal/transfer/storage"
)

//go:embed migration
var migrationFS embed.FS

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

	locationStorage := locationStorage.New(pool, box)
	transferStorage := transferStorage.New(pool)

	locationService := locationService.New(locationStorage)
	transferService := transferService.New(locationStorage, transferStorage)

	fooID, _ := uuid.NewRandom()
	err = locationService.AddS3(location.AddS3Command{
		ID: fooID.String(),

		Endpoint:        "localhost:9000",
		Bucket:          "foo",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	})
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	barID, _ := uuid.NewRandom()
	err = locationService.AddS3(location.AddS3Command{
		ID: barID.String(),

		Endpoint:        "localhost:9000",
		Bucket:          "bar",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	})
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	tID, _ := uuid.NewRandom()
	err = transferService.Add(transfer.AddCommand{
		ID: tID.String(),

		Pattern:        "*.png",
		FromLocationID: fooID.String(),
		ToLocationID:   barID.String(),
	})
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	err = transferService.Run(transfer.RunCommand{
		ID: tID.String(),
	})
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	cli := cli.New(flag.Args(), locationService, transferService)
	err = cli.Run(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
}
