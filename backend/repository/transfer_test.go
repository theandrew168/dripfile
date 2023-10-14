package repository_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each Repository impl

func TestTransferRepositoryCreate(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()
	transferRepo := repository.NewMemoryTransferRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	pattern := "*"

	itinerary, err := domain.NewItinerary(pattern, from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = transferRepo.Create(transfer)
	test.AssertNilError(t, err)
}

func TestTransferRepositoryList(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()
	transferRepo := repository.NewMemoryTransferRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	pattern := "*"

	itinerary, err := domain.NewItinerary(pattern, from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = transferRepo.Create(transfer)
	test.AssertNilError(t, err)

	transfers, err := transferRepo.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(transfers), 1)
}

func TestTransferRepositoryRead(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()
	transferRepo := repository.NewMemoryTransferRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	pattern := "*"

	itinerary, err := domain.NewItinerary(pattern, from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = transferRepo.Create(transfer)
	test.AssertNilError(t, err)

	got, err := transferRepo.Read(transfer.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), transfer.ID())
	test.AssertEqual(t, got.ItineraryID(), transfer.ItineraryID())
	test.AssertEqual(t, got.Status(), transfer.Status())
	test.AssertEqual(t, got.Progress(), transfer.Progress())
}

func TestTransferRepositoryReadNotFound(t *testing.T) {
	transferRepo := repository.NewMemoryTransferRepository()

	_, err := transferRepo.Read(uuid.New())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}

func TestTransferRepositoryDelete(t *testing.T) {
	locationRepo := repository.NewMemoryLocationRepository()
	itineraryRepo := repository.NewMemoryItineraryRepository()
	transferRepo := repository.NewMemoryTransferRepository()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(from)
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	err = locationRepo.Create(to)
	test.AssertNilError(t, err)

	pattern := "*"

	itinerary, err := domain.NewItinerary(pattern, from, to)
	test.AssertNilError(t, err)

	err = itineraryRepo.Create(itinerary)
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = transferRepo.Create(transfer)
	test.AssertNilError(t, err)

	_, err = transferRepo.Read(transfer.ID())
	test.AssertNilError(t, err)

	err = transferRepo.Delete(transfer.ID())
	test.AssertNilError(t, err)

	_, err = transferRepo.Read(transfer.ID())
	test.AssertErrorIs(t, err, repository.ErrNotExist)
}
