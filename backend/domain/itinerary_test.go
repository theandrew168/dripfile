package domain_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNewItinerary(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	test.AssertEqual(t, itinerary.Pattern(), "*")
	test.AssertEqual(t, itinerary.FromLocationID(), from.ID())
	test.AssertEqual(t, itinerary.ToLocationID(), to.ID())
}

func TestNewItineraryInvalidPattern(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	_, err = domain.NewItinerary(from, to, "")
	test.AssertErrorIs(t, err, domain.ErrItineraryInvalidPattern)
}

func TestNewItinerarySameLocation(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	_, err = domain.NewItinerary(from, from, "*")
	test.AssertErrorIs(t, err, domain.ErrItinerarySameLocation)
}

func TestItineraryCanDelete(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	test.AssertNilError(t, itinerary.CheckDelete())
}
