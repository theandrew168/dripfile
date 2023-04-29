package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/fileserver/s3"
	locationService "github.com/theandrew168/dripfile/internal/location/service"
)

type CLI struct {
	locationService *locationService.Service

	args []string
}

func New(locationService *locationService.Service, args []string) *CLI {
	c := CLI{
		locationService: locationService,

		args: args,
	}
	return &c
}

func (c *CLI) Run() error {
	app := &cli.App{
		Name:  "dripfile",
		Usage: "Managed File Transfers as a Service",
		Commands: []*cli.Command{
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
									endpoint, err := input("Endpoint: ")
									if err != nil {
										return err
									}

									bucket, err := input("Bucket: ")
									if err != nil {
										return err
									}

									accessKeyID, err := input("AccessKeyID: ")
									if err != nil {
										return err
									}

									secretAccessKey, err := input("SecretAccessKey: ")
									if err != nil {
										return err
									}

									info := s3.Info{
										Endpoint:        endpoint,
										Bucket:          bucket,
										AccessKeyID:     accessKeyID,
										SecretAccessKey: secretAccessKey,
									}
									location, err := c.locationService.CreateS3(info)
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

							location, err := c.locationService.Read(id)
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
							locations, err := c.locationService.List()
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

							err := c.locationService.Delete(id)
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

	args := append([]string{"dripfile"}, c.args...)
	return app.Run(args)
}

func input(prompt string) (string, error) {
	fmt.Print(prompt)

	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}
