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

func TestLocationRepositoryCreate(t *testing.T) {
	repo := repository.NewMemoryLocationRepository()

	l := model.NewMemoryLocation()

	err := repo.Create(l)
	test.AssertNilError(t, err)
}

func TestLocationRepositoryList(t *testing.T) {
	repo := repository.NewMemoryLocationRepository()

	l := model.NewMemoryLocation()

	err := repo.Create(l)
	test.AssertNilError(t, err)

	ls, err := repo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(ls), 1)
}

func TestLocationRepositoryRead(t *testing.T) {
	repo := repository.NewMemoryLocationRepository()

	l := model.NewMemoryLocation()

	err := repo.Create(l)
	test.AssertNilError(t, err)

	got, err := repo.Read(l.ID)
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID, l.ID)
	test.AssertEqual(t, got.Kind, l.Kind)
}

func TestLocationRepositoryReadNotFound(t *testing.T) {
	repo := repository.NewMemoryLocationRepository()

	_, err := repo.Read(uuid.New())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}

func TestLocationRepositoryDelete(t *testing.T) {
	repo := repository.NewMemoryLocationRepository()

	l := model.NewMemoryLocation()

	err := repo.Create(l)
	test.AssertNilError(t, err)

	_, err = repo.Read(l.ID)
	test.AssertNilError(t, err)

	err = repo.Delete(l.ID)
	test.AssertNilError(t, err)

	_, err = repo.Read(l.ID)
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}
