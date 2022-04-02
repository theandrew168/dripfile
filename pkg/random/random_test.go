package random_test

import (
	"testing"

	"github.com/theandrew168/dripfile/pkg/random"
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		want := i
		if got := random.String(i); len(got) != i {
			t.Errorf("want %v, got %v", want, len(got))
		}
	}
}
