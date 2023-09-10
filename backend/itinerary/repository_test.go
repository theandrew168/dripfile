package itinerary_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/itinerary"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestRepositoryCreate(t *testing.T) {
	repo := itinerary.NewMemoryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	i, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertNilError(t, err)

	err = repo.Create(i)
	test.AssertNilError(t, err)
}

func TestRepositoryList(t *testing.T) {
	repo := itinerary.NewMemoryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	i, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertNilError(t, err)

	err = repo.Create(i)
	test.AssertNilError(t, err)

	is, err := repo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(is), 1)
}

func TestRepositoryRead(t *testing.T) {
	repo := itinerary.NewMemoryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	i, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertNilError(t, err)

	err = repo.Create(i)
	test.AssertNilError(t, err)

	got, err := repo.Read(i.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), i.ID())
	test.AssertEqual(t, got.Pattern(), i.Pattern())
	test.AssertEqual(t, got.FromLocationID(), i.FromLocationID())
	test.AssertEqual(t, got.ToLocationID(), i.ToLocationID())
}

func TestRepositoryReadNotFound(t *testing.T) {
	repo := itinerary.NewMemoryRepository()

	_, err := repo.Read(uuid.New())
	test.AssertErrorIs(t, err, itinerary.ErrNotFound)
}

func TestRepositoryDelete(t *testing.T) {
	repo := itinerary.NewMemoryRepository()

	pattern := "*.txt"
	fromLocationID := uuid.New()
	toLocationID := uuid.New()

	i, err := itinerary.New(pattern, fromLocationID, toLocationID)
	test.AssertNilError(t, err)

	err = repo.Create(i)
	test.AssertNilError(t, err)

	_, err = repo.Read(i.ID())
	test.AssertNilError(t, err)

	err = repo.Delete(i.ID())
	test.AssertNilError(t, err)

	_, err = repo.Read(i.ID())
	test.AssertErrorIs(t, err, itinerary.ErrNotFound)
}
