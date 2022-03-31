package secret

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

type Box interface {
	Nonce() ([24]byte, error)
	Encrypt(nonce [24]byte, message []byte) []byte
	Decrypt(box []byte) ([]byte, error)
}

type box struct {
	key [32]byte
}

func NewBox(key [32]byte) Box {
	b := box{
		key: key,
	}
	return &b
}

func (b *box) Nonce() ([24]byte, error) {
	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return [24]byte{}, err
	}

	return nonce, nil
}

func (b *box) Encrypt(nonce [24]byte, message []byte) []byte {
	return secretbox.Seal(nonce[:], message, &nonce, &b.key)
}

func (b *box) Decrypt(box []byte) ([]byte, error) {
	var nonce [24]byte
	copy(nonce[:], box[:24])
	message, ok := secretbox.Open(nil, box[24:], &nonce, &b.key)
	if !ok {
		return nil, fmt.Errorf("decryption error")
	}

	return message, nil
}
