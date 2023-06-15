// Package secret is a nicer wrapper around the nacl/secretbox package.
package secret

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	ErrEncrypt = errors.New("secret: encryption error")
	ErrDecrypt = errors.New("secret: decryption error")
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

func (b *Box) Encrypt(message []byte) ([]byte, error) {
	nonce, err := generateNonce()
	if err != nil {
		return nil, ErrEncrypt
	}

	encrypted := secretbox.Seal(nonce[:], message, &nonce, &b.key)
	return encrypted, nil
}

func (b *Box) Decrypt(encrypted []byte) ([]byte, error) {
	var nonce [24]byte
	copy(nonce[:], encrypted[:24])

	message, ok := secretbox.Open(nil, encrypted[24:], &nonce, &b.key)
	if !ok {
		return nil, ErrDecrypt
	}

	return message, nil
}

func generateNonce() ([24]byte, error) {
	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return [24]byte{}, err
	}

	return nonce, nil
}
