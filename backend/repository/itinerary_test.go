package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestItineraryRepositoryCreate(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary("*", from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)
}

func TestItineraryRepositoryList(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary("*", from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	itineraries, err := itineraryRepo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(itineraries), 1)
}

func TestItineraryRepositoryRead(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary("*", from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	got, err := itineraryRepo.Read(itinerary.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), itinerary.ID())
	test.AssertEqual(t, got.Pattern(), itinerary.Pattern())
	test.AssertEqual(t, got.FromLocationID(), itinerary.FromLocationID())
	test.AssertEqual(t, got.ToLocationID(), itinerary.ToLocationID())
}

func TestItineraryRepositoryReadNotFound(t *testing.T) {
	itineraryRepo := repository.NewMemoryItineraryRepository()

	_, err := itineraryRepo.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestItineraryRepositoryDelete(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary("*", from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	_, err = itineraryRepo.Read(itinerary.ID())
	test.AssertNilError(t, err)

	err = itineraryRepo.Delete(itinerary)
	test.AssertNilError(t, err)

	_, err = itineraryRepo.Read(itinerary.ID())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}
