package main

import (
	"embed"
	"flag"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
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
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logger.Error(err.Error())
	}

	return 0
}
