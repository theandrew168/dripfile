package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestTransferRepositoryCreate(t *testing.T) {
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

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = repo.Transfer.Create(transfer)
	test.AssertNilError(t, err)
}

func TestTransferRepositoryList(t *testing.T) {
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

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = repo.Transfer.Create(transfer)
	test.AssertNilError(t, err)

	transfers, err := repo.Transfer.List()
	test.AssertNilError(t, err)

	var ids []uuid.UUID
	for _, transfer := range transfers {
		ids = append(ids, transfer.ID())
	}

	test.AssertSliceContains(t, ids, transfer.ID())
}

func TestTransferRepositoryRead(t *testing.T) {
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

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = repo.Transfer.Create(transfer)
	test.AssertNilError(t, err)

	got, err := repo.Transfer.Read(transfer.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), transfer.ID())
	test.AssertEqual(t, got.ItineraryID(), transfer.ItineraryID())
	test.AssertEqual(t, got.Status(), transfer.Status())
	test.AssertEqual(t, got.Progress(), transfer.Progress())
}

func TestTransferRepositoryReadNotFound(t *testing.T) {
	t.Parallel()

	repo, closer := test.Repository(t)
	defer closer()

	_, err := repo.Transfer.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestTransferRepositoryUpdate(t *testing.T) {
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

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = repo.Transfer.Create(transfer)
	test.AssertNilError(t, err)

	transfer.SetStatus(domain.TransferStatusSuccess)
	transfer.SetProgress(100)

	err = repo.Transfer.Update(transfer)
	test.AssertNilError(t, err)

	transfer, err = repo.Transfer.Read(transfer.ID())
	test.AssertNilError(t, err)

	test.AssertEqual(t, transfer.Status(), domain.TransferStatusSuccess)
	test.AssertEqual(t, transfer.Progress(), 100)
	test.AssertNotEqual(t, transfer.UpdatedAt(), transfer.CreatedAt())
}

func TestTransferRepositoryDelete(t *testing.T) {
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

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = repo.Transfer.Create(transfer)
	test.AssertNilError(t, err)

	_, err = repo.Transfer.Read(transfer.ID())
	test.AssertNilError(t, err)

	err = repo.Transfer.Delete(transfer)
	test.AssertNilError(t, err)

	_, err = repo.Transfer.Read(transfer.ID())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestTransferRepositoryAcquire(t *testing.T) {
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

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = repo.Transfer.Create(transfer)
	test.AssertNilError(t, err)

	transfer, err = repo.Transfer.Acquire()
	test.AssertNilError(t, err)
	test.AssertEqual(t, transfer.Status(), domain.TransferStatusRunning)
}
