package location_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/memorydb"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestRepositoryCreate(t *testing.T) {
	repo := location.NewMemoryRepository()

	l, err := location.NewMemory()
	test.AssertNilError(t, err)

	err = repo.Create(l)
	test.AssertNilError(t, err)
}

func TestRepositoryList(t *testing.T) {
	repo := location.NewMemoryRepository()

	l, err := location.NewMemory()
	test.AssertNilError(t, err)

	err = repo.Create(l)
	test.AssertNilError(t, err)

	ls, err := repo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(ls), 1)
}

func TestRepositoryRead(t *testing.T) {
	repo := location.NewMemoryRepository()

	l, err := location.NewMemory()
	test.AssertNilError(t, err)

	err = repo.Create(l)
	test.AssertNilError(t, err)

	got, err := repo.Read(l.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), l.ID())
	test.AssertEqual(t, got.Kind(), l.Kind())
}

func TestRepositoryReadNotFound(t *testing.T) {
	repo := location.NewMemoryRepository()

	_, err := repo.Read(uuid.New())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}

func TestRepositoryDelete(t *testing.T) {
	repo := location.NewMemoryRepository()

	l, err := location.NewMemory()
	test.AssertNilError(t, err)

	err = repo.Create(l)
	test.AssertNilError(t, err)

	_, err = repo.Read(l.ID())
	test.AssertNilError(t, err)

	err = repo.Delete(l.ID())
	test.AssertNilError(t, err)

	_, err = repo.Read(l.ID())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}
