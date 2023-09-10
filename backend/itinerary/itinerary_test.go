package itinerary_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/itinerary"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNew(t *testing.T) {
	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	i, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, i.Pattern(), pattern)
	test.AssertEqual(t, i.FromLocationID(), fromLocationID)
	test.AssertEqual(t, i.ToLocationID(), toLocationID)
}

func TestNewInvalidPattern(t *testing.T) {
	pattern := ""
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	_, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertErrorIs(t, err, itinerary.ErrInvalidPattern)
}

func TestNewSameLocation(t *testing.T) {
	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := fromLocationID

	_, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertErrorIs(t, err, itinerary.ErrSameLocation)
}
