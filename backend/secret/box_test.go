package secret_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/secret"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	rand := test.NewRandom()
	secretKey := rand.Bytes(32)
	box := secret.NewBox([32]byte(secretKey))

	inputs := []string{
		"",
		"hello world",
		"longer message still works when encrypted / decrypted",
	}

	for _, input := range inputs {
		encrypted, err := box.Encrypt([]byte(input))
		test.AssertNilError(t, err)

		decrypted, err := box.Decrypt(encrypted)
		test.AssertNilError(t, err)

		test.AssertEqual(t, string(decrypted), input)
	}
}
