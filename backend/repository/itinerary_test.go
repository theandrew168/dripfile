package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestItineraryRepositoryCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	err = repo.Itinerary.Create(itinerary)
	test.AssertNilError(t, err)
}

func TestItineraryRepositoryList(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	err = repo.Itinerary.Create(itinerary)
	test.AssertNilError(t, err)

	itineraries, err := repo.Itinerary.List()
	test.AssertNilError(t, err)

	var ids []uuid.UUID
	for _, itinerary := range itineraries {
		ids = append(ids, itinerary.ID())
	}

	test.AssertSliceContains(t, ids, itinerary.ID())
}

func TestItineraryRepositoryRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	err = repo.Itinerary.Create(itinerary)
	test.AssertNilError(t, err)

	got, err := repo.Itinerary.Read(itinerary.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), itinerary.ID())
	test.AssertEqual(t, got.Pattern(), itinerary.Pattern())
	test.AssertEqual(t, got.FromLocationID(), itinerary.FromLocationID())
	test.AssertEqual(t, got.ToLocationID(), itinerary.ToLocationID())
}

func TestItineraryRepositoryReadNotFound(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	_, err := repo.Itinerary.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestItineraryRepositoryDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(to)
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	err = repo.Itinerary.Create(itinerary)
	test.AssertNilError(t, err)

	_, err = repo.Itinerary.Read(itinerary.ID())
	test.AssertNilError(t, err)

	err = repo.Itinerary.Delete(itinerary)
	test.AssertNilError(t, err)

	_, err = repo.Itinerary.Read(itinerary.ID())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}
