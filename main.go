package main

import (
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/fileserver/s3"
	locationRepo "github.com/theandrew168/dripfile/internal/location/repository"
	locationService "github.com/theandrew168/dripfile/internal/location/service"
	"github.com/theandrew168/dripfile/internal/migrate"
	"github.com/theandrew168/dripfile/internal/secret"
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

	locationRepo := locationRepo.New(pool)
	locationService := locationService.New(box, locationRepo)

	app := &cli.App{
		Name:  "dripfile",
		Usage: "Managed File Transfers as a Service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf",
				Usage: "Path to config file",
				Value: "dripfile.conf",
			},
		},
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
								Name:      "s3",
								Usage:     "Creates a new S3 location",
								ArgsUsage: "endpoint bucket access_key_id secret_access_key",
								Action: func(ctx *cli.Context) error {
									if ctx.Args().Len() != 4 {
										return cli.ShowSubcommandHelp(ctx)
									}

									endpoint := ctx.Args().Get(0)
									bucket := ctx.Args().Get(1)
									accessKeyID := ctx.Args().Get(2)
									secretAccessKey := ctx.Args().Get(3)

									info := s3.Info{
										Endpoint:        endpoint,
										Bucket:          bucket,
										AccessKeyID:     accessKeyID,
										SecretAccessKey: secretAccessKey,
									}
									location, err := locationService.CreateS3(info)
									if err != nil {
										return err
									}

									fmt.Printf("location created: %s\n", location.ID)
									return nil
								},
							},
						},
					},
					{
						Name:      "read",
						Usage:     "Reads a location by its ID",
						ArgsUsage: "id",
						Action: func(ctx *cli.Context) error {
							if ctx.Args().Len() != 1 {
								return cli.ShowSubcommandHelp(ctx)
							}

							id := ctx.Args().Get(0)

							location, err := locationRepo.Read(id)
							if err != nil {
								return err
							}

							switch location.Kind {
							case "s3":
								var info s3.Info
								err := json.Unmarshal(location.Info, &info)
								if err != nil {
									return err
								}

								fmt.Printf(
									"%s %s %s %s\n",
									location.ID,
									location.Kind,
									info.Endpoint,
									info.Bucket,
								)
							}

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
								fmt.Println(location.ID)
							}

							return nil
						},
					},
					{
						Name:      "delete",
						Usage:     "Deletes a location by its ID",
						ArgsUsage: "id",
						Action: func(ctx *cli.Context) error {
							if ctx.Args().Len() != 1 {
								return cli.ShowSubcommandHelp(ctx)
							}

							id := ctx.Args().Get(0)

							err := locationRepo.Delete(id)
							if err != nil {
								if errors.Is(err, database.ErrNotExist) {
									fmt.Printf("location does not exist: %s\n", id)
									return nil
								}

								return err
							}

							fmt.Printf("location deleted: %s\n", id)
							return nil
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

	return 0
}
