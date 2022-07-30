package secret_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/test"
)

func TestEncryptDecrypt(t *testing.T) {
	secretKeyBytes := test.RandomBytes(32)

	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	box := secret.NewBox(secretKey)

	inputs := []string{
		"",
		"hello world",
		"longer message still works when encrypted / decrypted",
	}

	for _, input := range inputs {
		nonce, err := box.Nonce()
		test.AssertNilError(t, err)

		encrypted := box.Encrypt(nonce, []byte(input))
		decrypted, err := box.Decrypt(encrypted)
		test.AssertNilError(t, err)

		test.AssertEqual(t, string(decrypted), input)
	}
}
