package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

const DefaultPort = "5000"

type Config struct {
	SiteURL        string `toml:"site_url"`
	SecretKey      string `toml:"secret_key"`
	DatabaseURI    string `toml:"database_uri"`
	RedisURI       string `toml:"redis_uri"`
	PostmarkAPIKey string `toml:"postmark_api_key"`
	Port           string `toml:"port"`
}

func Read(data string) (Config, error) {
	// init Config struct with default values
	cfg := Config{
		Port: DefaultPort,
	}
	meta, err := toml.Decode(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	// gather extra values
	extra := []string{}
	for _, keys := range meta.Undecoded() {
		key := keys[0]
		extra = append(extra, key)
	}

	// error upon extra values
	if len(extra) > 0 {
		msg := strings.Join(extra, ", ")
		return Config{}, fmt.Errorf("extra config values: %s", msg)
	}

	// build set of present config keys
	present := make(map[string]bool)
	for _, keys := range meta.Keys() {
		key := keys[0]
		present[key] = true
	}

	required := []string{
		"site_url",
		"secret_key",
		"database_uri",
		"redis_uri",
	}

	// gather missing values
	missing := []string{}
	for _, key := range required {
		if _, ok := present[key]; !ok {
			missing = append(missing, key)
		}
	}

	// error upon missing values
	if len(missing) > 0 {
		msg := strings.Join(missing, ", ")
		return Config{}, fmt.Errorf("missing config values: %s", msg)
	}

	return cfg, nil
}

func ReadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return Read(string(data))
}
