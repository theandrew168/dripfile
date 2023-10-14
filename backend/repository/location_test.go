package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestLocationRepositoryCreate(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(location)
	test.AssertNilError(t, err)
}

func TestLocationRepositoryList(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(location)
	test.AssertNilError(t, err)

	locations, err := locationRepo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(locations), 1)
}

func TestLocationRepositoryRead(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(location)
	test.AssertNilError(t, err)

	got, err := locationRepo.Read(location.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), location.ID())
	test.AssertEqual(t, got.Kind(), location.Kind())
}

func TestLocationRepositoryReadNotFound(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()

	_, err := locationRepo.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestLocationRepositoryDelete(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()

	location, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(location)
	test.AssertNilError(t, err)

	_, err = locationRepo.Read(location.ID())
	test.AssertNilError(t, err)

	err = locationRepo.Delete(location)
	test.AssertNilError(t, err)

	_, err = locationRepo.Read(location.ID())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}
