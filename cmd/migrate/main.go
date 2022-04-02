package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/migrate"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

func main() {
	os.Exit(run())
}

func run() int {
	infoLog := log.New(os.Stdout, "", 0)
	errorLog := log.New(os.Stderr, "error: ", 0)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	// open a regular connection
	conn, err := postgres.Connect(cfg.DatabaseURI)
	if err != nil {
		errorLog.Println(err)
		return 1
	}
	defer conn.Close(context.Background())

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		errorLog.Println(err)
		return 1
	}

	err = migrate.Migrate(conn, infoLog)
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	return 0
}
