package test

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Random struct {
	rand *rand.Rand
}

func NewRandom() *Random {
	return NewRandomFromSeed(time.Now().UnixNano())
}

func NewRandomFromSeed(seed int64) *Random {
	r := Random{
		rand: rand.New(rand.NewSource(seed)),
	}
	return &r
}

func (r *Random) Bytes(n int) []byte {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[r.rand.Intn(len(valid))]
	}

	return buf
}

func (r *Random) Int() int {
	return r.rand.Int()
}

func (r *Random) UUID() string {
	id, err := uuid.NewRandomFromReader(r.rand)
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (r *Random) String(n int) string {
	return string(r.Bytes(n))
}

func (r *Random) URL() string {
	return "https://" + r.String(16) + "example.com"
}

func (r *Random) Time() time.Time {
	return time.Now().UTC()
}
