package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	locationRepo "github.com/theandrew168/dripfile/internal/location/repository"
	"github.com/theandrew168/dripfile/internal/migrate"
)

//go:embed migration
var migrationFS embed.FS

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stdout))

	conf := flag.String("conf", "dripfile.conf", "app config file")
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

	locationRepo := locationRepo.NewPostgresRepository(pool)

	app := &cli.App{
		Name:  "dripfile",
		Usage: "Managed File Transfers as a Service",
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "Applies migrations and exits",
				Action: func(*cli.Context) error {
					return migrate.Migrate(logger, pool, migrationFS)
				},
			},
			{
				Name:  "location",
				Usage: "Options for managing locations",
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Creates a new location",
						Subcommands: []*cli.Command{
							{
								Name:  "s3",
								Usage: "Creates a new S3 location",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "endpoint",
										Usage:    "S3 endpoint URL",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "bucket",
										Usage:    "S3 bucket name",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "access_key_id",
										Usage:    "S3 access key id",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "secret_access_key",
										Usage:    "S3 secret access key",
										Required: true,
									},
								},
								Action: func(*cli.Context) error {
									location, err := locationRepo.Read("asdf")
									if err != nil {
										return err
									}

									fmt.Printf("%+v\n", location)
									return nil
								},
							},
						},
					},
					{
						Name:  "read",
						Usage: "Reads a location by its ID",
						Action: func(*cli.Context) error {
							location, err := locationRepo.Read("asdf")
							if err != nil {
								return err
							}

							fmt.Printf("%+v\n", location)
							return nil
						},
					},
					{
						Name:  "list",
						Usage: "Lists all locations",
						Action: func(*cli.Context) error {
							locations, err := locationRepo.List()
							if err != nil {
								return err
							}

							for _, location := range locations {
								fmt.Printf("%+v\n", location)
							}

							return nil
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logger.Error(err.Error())
	}

	return 0
}
