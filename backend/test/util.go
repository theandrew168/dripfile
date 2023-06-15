package test

import (
	"context"
	"testing"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
)

type CloserFunc func()

func Config(t *testing.T) config.Config {
	t.Helper()

	// read the local development config file
	cfg, err := config.ReadFile("../../dripfile.conf")
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func Database(t *testing.T) (database.Conn, CloserFunc) {
	t.Helper()

	cfg := Config(t)
	conn, err := database.Connect(cfg.DatabaseURI)
	if err != nil {
		t.Fatal(err)
	}

	return conn, func() {
		conn.Close(context.Background())
	}
}
