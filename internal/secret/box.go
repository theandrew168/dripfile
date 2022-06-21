// Nicer wrapper around the nacl/secretbox package.
package secret

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

type Box struct {
	key [32]byte
}

func NewBox(key [32]byte) *Box {
	b := Box{
		key: key,
	}
	return &b
}

func (b *Box) Nonce() ([24]byte, error) {
	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return [24]byte{}, err
	}

	return nonce, nil
}

func (b *Box) Encrypt(nonce [24]byte, message []byte) []byte {
	return secretbox.Seal(nonce[:], message, &nonce, &b.key)
}

func (b *Box) Decrypt(box []byte) ([]byte, error) {
	var nonce [24]byte
	copy(nonce[:], box[:24])

	message, ok := secretbox.Open(nil, box[24:], &nonce, &b.key)
	if !ok {
		return nil, fmt.Errorf("decryption error")
	}

	return message, nil
}
