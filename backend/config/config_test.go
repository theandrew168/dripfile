package config_test

import (
	"fmt"
	"testing"

	"github.com/theandrew168/dripfile/backend/config"
	"github.com/theandrew168/dripfile/backend/test"
)

const (
	secretKey   = "secret"
	databaseURI = "postgresql://foo:bar@localhost:5432/postgres"
	smtpURI     = "smtp://foo:bar@localhost:587"
	port        = "5000"
)

func TestRead(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
		database_uri = "%s"
		smtp_uri = "%s"
		port = "%s"
	`, secretKey, databaseURI, smtpURI, port)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.SecretKey, secretKey)
	test.AssertEqual(t, cfg.DatabaseURI, databaseURI)
	test.AssertEqual(t, cfg.SMTPURI, smtpURI)
	test.AssertEqual(t, cfg.Port, port)
}

func TestOptional(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
		database_uri = "%s"
	`, secretKey, databaseURI)

	cfg, err := config.Read(data)
	test.AssertNilError(t, err)

	test.AssertEqual(t, cfg.SecretKey, secretKey)
	test.AssertEqual(t, cfg.DatabaseURI, databaseURI)
	test.AssertEqual(t, cfg.SMTPURI, "")
	test.AssertEqual(t, cfg.Port, config.DefaultPort)
}

func TestRequired(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
	`, secretKey)

	_, err := config.Read(data)
	test.AssertErrorContains(t, err, "missing")
	test.AssertErrorContains(t, err, "database_uri")
}

func TestExtra(t *testing.T) {
	t.Parallel()

	data := fmt.Sprintf(`
		secret_key = "%s"
		database_uri = "%s"
		foo = "bar"
	`, secretKey, databaseURI)

	_, err := config.Read(data)
	test.AssertErrorContains(t, err, "extra")
	test.AssertErrorContains(t, err, "foo")
}
