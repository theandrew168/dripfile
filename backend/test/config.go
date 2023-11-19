package test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/theandrew168/dripfile/backend/config"
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/secret"
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

func Repository(t *testing.T) (*repository.Repository, CloserFunc) {
	t.Helper()

	conn, closer := Database(t)

	cfg := Config(t)
	secretKey, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		t.Fatal(err)
	}
	box := secret.NewBox([32]byte(secretKey))

	repo := repository.NewPostgres(conn, box)
	return repo, closer
}
