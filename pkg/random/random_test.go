package random_test

import (
	"strings"
	"testing"

	"github.com/theandrew168/dripfile/pkg/random"
)

func TestString(t *testing.T) {
	if got := random.String(0); got != "" {
		t.Errorf("random string broken when empty")
	}
	if got := random.String(8); len(got) != 8 {
		t.Errorf("random string invalid length")
	}
}

func TestURL(t *testing.T) {
	if got := random.URL(8); !strings.HasPrefix(got, "https://") {
		t.Errorf("random url missing https:// prefix")
	}
}

func TestTime(t *testing.T) {
	random.Time()
}
