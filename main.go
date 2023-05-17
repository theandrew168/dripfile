package main

import (
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	locationService "github.com/theandrew168/dripfile/internal/location/service"
	"github.com/theandrew168/dripfile/internal/location/service/command"
	"github.com/theandrew168/dripfile/internal/location/service/query"
	locationStorage "github.com/theandrew168/dripfile/internal/location/storage/postgres"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/secret"
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
	locationService := locationService.New(locationStorage)

	id, _ := uuid.NewRandom()
	err = locationService.Command.CreateS3.Handle(command.CreateS3{
		ID:              id.String(),
		Endpoint:        "localhost:9000",
		Bucket:          "foo",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	})
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	l, err := locationService.Query.Read.Handle(query.Read{
		ID: id.String(),
	})
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	fmt.Printf("%+v\n", l)

	fs, err := l.Connect()
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	err = fs.Ping()
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
}
