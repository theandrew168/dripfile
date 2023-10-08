package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
	"github.com/theandrew168/dripfile/backend/model"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestItineraryRepositoryCreate(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocation := model.NewMemoryLocation()
	toLocation := model.NewMemoryLocation()

	i := model.NewItinerary(pattern, fromLocation, toLocation)

	err := repo.Create(i)
	test.AssertNilError(t, err)
}

func TestItineraryRepositoryList(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocation := model.NewMemoryLocation()
	toLocation := model.NewMemoryLocation()

	i := model.NewItinerary(pattern, fromLocation, toLocation)

	err := repo.Create(i)
	test.AssertNilError(t, err)

	is, err := repo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(is), 1)
}

func TestItineraryRepositoryRead(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocation := model.NewMemoryLocation()
	toLocation := model.NewMemoryLocation()

	i := model.NewItinerary(pattern, fromLocation, toLocation)

	err := repo.Create(i)
	test.AssertNilError(t, err)

	got, err := repo.Read(i.ID)
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID, i.ID)
	test.AssertEqual(t, got.Pattern, i.Pattern)
	test.AssertEqual(t, got.FromLocationID, i.FromLocationID)
	test.AssertEqual(t, got.ToLocationID, i.ToLocationID)
}

func TestItineraryRepositoryReadNotFound(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	_, err := repo.Read(uuid.New())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}

func TestItineraryRepositoryDelete(t *testing.T) {
	repo := repository.NewMemoryItineraryRepository()

	pattern := "*.txt"
	fromLocation := model.NewMemoryLocation()
	toLocation := model.NewMemoryLocation()

	i := model.NewItinerary(pattern, fromLocation, toLocation)

	err := repo.Create(i)
	test.AssertNilError(t, err)

	_, err = repo.Read(i.ID)
	test.AssertNilError(t, err)

	err = repo.Delete(i.ID)
	test.AssertNilError(t, err)

	_, err = repo.Read(i.ID)
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}
