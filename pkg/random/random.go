package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Bytes(n int) []byte {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return buf
}

func String(n int) string {
	return string(Bytes(n))
}

func URL(n int) string {
	return "https://" + String(n)
}

func Time() time.Time {
	return time.Now().UTC()
}
