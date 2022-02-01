package test

import (
	"math/rand"

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

func RandomString(n int) string {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return string(buf)
}
