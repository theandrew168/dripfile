package test

import (
	"testing"

	"github.com/theandrew168/dripfile/pkg/config"
)

func Config(t *testing.T) config.Config {
	t.Helper()

	// read the local development config file
	cfg, err := config.ReadFile("../../dripfile.conf")
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}
