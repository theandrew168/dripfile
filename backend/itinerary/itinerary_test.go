package itinerary_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/itinerary"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNew(t *testing.T) {
	fromLocationID := uuid.New()
	toLocationID := uuid.New()
	pattern := "*.txt"

	i, err := itinerary.New(fromLocationID, toLocationID, pattern)
	test.AssertNilError(t, err)

	test.AssertEqual(t, i.Pattern(), pattern)
	test.AssertEqual(t, i.FromLocationID(), fromLocationID)
	test.AssertEqual(t, i.ToLocationID(), toLocationID)
}

func TestNewInvalidPattern(t *testing.T) {
	fromLocationID := uuid.New()
	toLocationID := uuid.New()
	pattern := ""

	_, err := itinerary.New(fromLocationID, toLocationID, pattern)
	test.AssertErrorIs(t, err, itinerary.ErrInvalidPattern)
}

func TestNewSameLocation(t *testing.T) {
	fromLocationID := uuid.New()
	toLocationID := fromLocationID
	pattern := "*.txt"

	_, err := itinerary.New(fromLocationID, toLocationID, pattern)
	test.AssertErrorIs(t, err, itinerary.ErrSameLocation)
}
