package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/model"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestItineraryRepositoryCreate(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	itinerary := model.NewItinerary(pattern, fromLocationID, toLocationID)

	err := repo.Create(itinerary)
	test.AssertNilError(t, err)
}

func TestItineraryRepositoryList(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	itinerary := model.NewItinerary(pattern, fromLocationID, toLocationID)

	err := repo.Create(itinerary)
	test.AssertNilError(t, err)

	is, err := repo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(is), 1)
}

func TestItineraryRepositoryRead(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	itinerary := model.NewItinerary(pattern, fromLocationID, toLocationID)

	err := repo.Create(itinerary)
	test.AssertNilError(t, err)

	got, err := repo.Read(itinerary.ID)
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID, itinerary.ID)
	test.AssertEqual(t, got.Pattern, itinerary.Pattern)
	test.AssertEqual(t, got.FromLocationID, itinerary.FromLocationID)
	test.AssertEqual(t, got.ToLocationID, itinerary.ToLocationID)
}

func TestItineraryRepositoryReadNotFound(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	_, err := repo.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestItineraryRepositoryDelete(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	itinerary := model.NewItinerary(pattern, fromLocationID, toLocationID)

	err := repo.Create(itinerary)
	test.AssertNilError(t, err)

	_, err = repo.Read(itinerary.ID)
	test.AssertNilError(t, err)

	err = repo.Delete(itinerary.ID)
	test.AssertNilError(t, err)

	_, err = repo.Read(itinerary.ID)
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}
