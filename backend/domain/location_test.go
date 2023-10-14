package domain_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNewMemoryLocation(t *testing.T) {
	l, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)
	test.AssertEqual(t, l.Kind(), domain.LocationKindMemory)
}

func TestLocationCanDelete(t *testing.T) {
	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	test.AssertEqual(t, from.CanDelete(), true)
	test.AssertEqual(t, to.CanDelete(), true)

	_, err = domain.NewItinerary("*", from, to)
	test.AssertNilError(t, err)

	test.AssertEqual(t, from.CanDelete(), false)
	test.AssertEqual(t, to.CanDelete(), false)
}
