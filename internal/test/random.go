package test

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomBytes(n int) []byte {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return buf
}

func RandomInt() int {
	return rand.Int()
}

func RandomUUID() string {
	return "2dbb189e-d392-43a1-bef6-724c63c0550a"
}

func RandomString(n int) string {
	return string(RandomBytes(n))
}

func RandomURL(n int) string {
	return "https://" + RandomString(n)
}

func RandomTime() time.Time {
	return time.Now().UTC()
}
