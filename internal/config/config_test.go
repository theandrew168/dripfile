package config_test

import (
	"fmt"
	"testing"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/test"
)

const (
	secretKey     = "secret"
	postgreSQLURL = "postgresql://foo:bar@localhost:5432/postgres"
	smtpURL       = "smtp://foo:bar@localhost:587"
	port          = "5000"
)

func TestRead(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
		postgresql_url = "%s"
		smtp_url = "%s"
		port = "%s"
	`, secretKey, postgreSQLURL, smtpURL, port)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.SecretKey, secretKey)
	test.AssertEqual(t, cfg.PostgreSQLURL, postgreSQLURL)
	test.AssertEqual(t, cfg.SMTPURL, smtpURL)
	test.AssertEqual(t, cfg.Port, port)
}

func TestOptional(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
		postgresql_url = "%s"
	`, secretKey, postgreSQLURL)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.SecretKey, secretKey)
	test.AssertEqual(t, cfg.PostgreSQLURL, postgreSQLURL)
	test.AssertEqual(t, cfg.SMTPURL, "")
	test.AssertEqual(t, cfg.Port, config.DefaultPort)
}

func TestRequired(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
	`, secretKey)

	_, err := config.Read(data)
	test.AssertErrorContains(t, err, "missing")
	test.AssertErrorContains(t, err, "postgresql_url")
}

func TestExtra(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
		postgresql_url = "%s"
		foo = "bar"
	`, secretKey, postgreSQLURL)

	_, err := config.Read(data)
	test.AssertErrorContains(t, err, "extra")
	test.AssertErrorContains(t, err, "foo")
}
