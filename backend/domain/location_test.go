package domain_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNewMemoryLocation(t *testing.T) {
	t.Parallel()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)
	test.AssertEqual(t, location.Kind(), domain.LocationKindMemory)
}

func TestLocationCheckDelete(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	test.AssertNilError(t, from.CheckDelete())
	test.AssertNilError(t, to.CheckDelete())

	_, err = domain.NewItinerary("*", from, to)
	test.AssertNilError(t, err)

	test.AssertErrorIs(t, from.CheckDelete(), domain.ErrLocationInUse)
	test.AssertErrorIs(t, to.CheckDelete(), domain.ErrLocationInUse)
}
