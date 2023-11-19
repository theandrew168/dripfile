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
}

func TestLocationRepositoryReadNotFound(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	_, err := repo.Location.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
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
