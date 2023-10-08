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

func TestTransferRepositoryCreate(t *testing.T) {
	repo := repository.NewMemoryTransferRepository()

	itineraryID := uuid.New()

	transfer := model.NewTransfer(itineraryID)

	err := repo.Create(transfer)
	test.AssertNilError(t, err)
}

func TestTransferRepositoryList(t *testing.T) {
	repo := repository.NewMemoryTransferRepository()

	itineraryID := uuid.New()

	transfer := model.NewTransfer(itineraryID)

	err := repo.Create(transfer)
	test.AssertNilError(t, err)

	is, err := repo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(is), 1)
}

func TestTransferRepositoryRead(t *testing.T) {
	repo := repository.NewMemoryTransferRepository()

	itineraryID := uuid.New()

	transfer := model.NewTransfer(itineraryID)

	err := repo.Create(transfer)
	test.AssertNilError(t, err)

	got, err := repo.Read(transfer.ID)
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID, transfer.ID)
	test.AssertEqual(t, got.ItineraryID, transfer.ItineraryID)
}

func TestTransferRepositoryReadNotFound(t *testing.T) {
	repo := repository.NewMemoryTransferRepository()

	_, err := repo.Read(uuid.New())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}

func TestTransferRepositoryDelete(t *testing.T) {
	repo := repository.NewMemoryTransferRepository()

	itineraryID := uuid.New()

	transfer := model.NewTransfer(itineraryID)

	err := repo.Create(transfer)
	test.AssertNilError(t, err)

	_, err = repo.Read(transfer.ID)
	test.AssertNilError(t, err)

	err = repo.Delete(transfer.ID)
	test.AssertNilError(t, err)

	_, err = repo.Read(transfer.ID)
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}
