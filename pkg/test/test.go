package test

import (
	"testing"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/database"
	"github.com/theandrew168/dripfile/pkg/storage"
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

func Database(t *testing.T) (database.Interface, CloserFunc) {
	t.Helper()

	cfg := Config(t)
	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		t.Fatal(err)
	}

	return pool, pool.Close
}

func Storage(t *testing.T) (*storage.Storage, CloserFunc) {
	t.Helper()

	cfg := Config(t)
	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		t.Fatal(err)
	}

	store := storage.New(pool)
	return store, pool.Close
}
