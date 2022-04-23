package main

import (
	"context"
	"flag"
	"os"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/database"
	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/migrate"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := jsonlog.New(os.Stdout)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	// open a database connection
	conn, err := database.Connect(cfg.DatabaseURI)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}
	defer conn.Close(context.Background())

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	err = migrate.Migrate(conn, logger)
	if err != nil {
		logger.PrintError(err, nil)
		return 1
	}

	return 0
}
