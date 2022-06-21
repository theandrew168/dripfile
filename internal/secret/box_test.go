package secret_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/random"
	"github.com/theandrew168/dripfile/internal/secret"
)

func TestEncryptDecrypt(t *testing.T) {
	secretKeyBytes := random.Bytes(32)

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
		if err != nil {
			t.Fatal(err)
		}

		encrypted := box.Encrypt(nonce, []byte(input))
		decrypted, err := box.Decrypt(encrypted)
		if err != nil {
			t.Fatal(err)
		}

		if string(decrypted) != input {
			t.Fatal("mismatched input after encrypt / decrypt")
		}
	}
}
