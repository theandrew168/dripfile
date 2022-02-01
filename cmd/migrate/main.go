package main

import (
	"context"
	"flag"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/migrate"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger(os.Stdout)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Error(err)
		return 1
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), cfg.DatabaseURI)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer conn.Close()

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		logger.Error(err)
		return 1
	}

	ctx := context.Background()
	err = migrate.Migrate(ctx, conn, logger)
	if err != nil {
		logger.Error(err)
		return 1
	}

	return 0
}
