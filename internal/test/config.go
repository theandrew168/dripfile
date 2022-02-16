package test

import (
	"github.com/theandrew168/dripfile/internal/config"
)

func Config() config.Config {
	// read the local development config file
	cfg, err := config.ReadFile("../../dripfile.conf")
	if err != nil {
		panic(err)
	}

	return cfg
}
