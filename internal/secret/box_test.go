package secret_test

import (
	"testing"
	"time"

	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/test"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	random := test.NewRandom(time.Now().Unix())
	secretKey := random.Bytes(32)
	box := secret.NewBox([32]byte(secretKey))

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
