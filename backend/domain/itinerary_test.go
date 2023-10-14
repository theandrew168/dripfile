package domain_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNewItinerary(t *testing.T) {
	pattern := "*"

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(pattern, from, to)
	test.AssertNilError(t, err)

	test.AssertEqual(t, itinerary.Pattern(), pattern)
	test.AssertEqual(t, itinerary.FromLocationID(), from.ID())
	test.AssertEqual(t, itinerary.ToLocationID(), to.ID())
}

func TestNewItineraryInvalidPattern(t *testing.T) {
	pattern := ""

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	_, err = domain.NewItinerary(pattern, from, to)
	test.AssertErrorIs(t, err, domain.ErrItineraryInvalidPattern)
}

func TestNewItinerarySameLocation(t *testing.T) {
	pattern := "*"

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	_, err = domain.NewItinerary(pattern, from, from)
	test.AssertErrorIs(t, err, domain.ErrItinerarySameLocation)
}

func TestItineraryCanDelete(t *testing.T) {
	pattern := "*"

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(pattern, from, to)
	test.AssertNilError(t, err)

	test.AssertNilError(t, itinerary.CheckDelete())
}
