package location_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNewMemory(t *testing.T) {
	l, err := location.NewMemory()
	test.AssertNilError(t, err)
	test.AssertEqual(t, l.Kind(), location.KindMemory)
}
