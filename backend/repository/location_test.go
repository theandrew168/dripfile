package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestLocationRepositoryCreate(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(location)
	test.AssertNilError(t, err)
}

func TestLocationRepositoryList(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(location)
	test.AssertNilError(t, err)

	locations, err := repo.Location.List()
	test.AssertNilError(t, err)

	var ids []uuid.UUID
	for _, location := range locations {
		ids = append(ids, location.ID())
	}

	test.AssertSliceContains(t, ids, location.ID())
}

func TestLocationRepositoryListUsedBy(t *testing.T) {
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

	locations, err := repo.Location.List()
	test.AssertNilError(t, err)

	var got []*domain.Location
	for _, location := range locations {
		if location.ID() == from.ID() || location.ID() == to.ID() {
			got = append(got, location)
		}
	}

	for _, location := range got {
		test.AssertSliceContains(t, location.UsedBy(), itinerary.ID())
	}
}

func TestLocationRepositoryRead(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(location)
	test.AssertNilError(t, err)

	got, err := repo.Location.Read(location.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), location.ID())
	test.AssertEqual(t, got.Kind(), location.Kind())
	test.AssertEqual(t, len(got.UsedBy()), 0)
}

func TestLocationRepositoryReadNotFound(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	_, err := repo.Location.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestLocationRepositoryReadUsedBy(t *testing.T) {
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

	gotFrom, err := repo.Location.Read(from.ID())
	test.AssertNilError(t, err)
	test.AssertSliceContains(t, gotFrom.UsedBy(), itinerary.ID())

	gotTo, err := repo.Location.Read(to.ID())
	test.AssertNilError(t, err)
	test.AssertSliceContains(t, gotTo.UsedBy(), itinerary.ID())
}

func TestLocationRepositoryDelete(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = repo.Location.Create(location)
	test.AssertNilError(t, err)

	_, err = repo.Location.Read(location.ID())
	test.AssertNilError(t, err)

	err = repo.Location.Delete(location)
	test.AssertNilError(t, err)

	_, err = repo.Location.Read(location.ID())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}
