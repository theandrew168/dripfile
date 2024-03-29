package domain_test

import (
	"testing"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestNewTransfer(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)
	test.AssertEqual(t, transfer.Status(), domain.TransferStatusPending)
	test.AssertEqual(t, transfer.Progress(), 0)
}

func TestTransferCanDelete(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	test.AssertNilError(t, transfer.CheckDelete())
}

func TestTransferUpdateStatus(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = transfer.SetStatus(domain.TransferStatusSuccess)
	test.AssertNilError(t, err)

	test.AssertEqual(t, transfer.Status(), domain.TransferStatusSuccess)
}

func TestTransferUpdateProgress(t *testing.T) {
	t.Parallel()

	from, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	to, err := domain.NewMemoryLocation()
	test.AssertNilError(t, err)

	itinerary, err := domain.NewItinerary(from, to, "*")
	test.AssertNilError(t, err)

	transfer, err := domain.NewTransfer(itinerary)
	test.AssertNilError(t, err)

	err = transfer.SetProgress(100)
	test.AssertNilError(t, err)

	test.AssertEqual(t, transfer.Progress(), 100)
}
