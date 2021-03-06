package test

import (
	"context"
	"testing"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/storage"
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

func Database(t *testing.T) (postgresql.Conn, CloserFunc) {
	t.Helper()

	cfg := Config(t)
	conn, err := postgresql.Connect(cfg.PostgreSQLURL)
	if err != nil {
		t.Fatal(err)
	}

	return conn, func() {
		conn.Close(context.Background())
	}
}

func Storage(t *testing.T) (*storage.Storage, CloserFunc) {
	t.Helper()

	db, closer := Database(t)
	store := storage.New(db)
	return store, closer
}
