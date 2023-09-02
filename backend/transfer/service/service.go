package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/transfer"
)

// ensure Service interface is satisfied
var _ Service = (*service)(nil)

// service interface (other code depends on this)
type Service interface {
	RunDomain(transfer *transfer.Transfer, from, to *location.Location) error
	RunApplication(transferID uuid.UUID) error
}

// service implementation
type service struct {
	locationRepo location.Repository
	transferRepo transfer.Repository
}

func New(
	locationRepo location.Repository,
	transferRepo transfer.Repository,
) Service {
	srvc := service{
		locationRepo: locationRepo,
		transferRepo: transferRepo,
	}
	return &srvc
}

// Transfer Domain Service - run a transfer from a domain point of view
func (srvc *service) RunDomain(transfer *transfer.Transfer, from, to *location.Location) error {
	fromFS, err := from.Connect()
	if err != nil {
		return err
	}

	toFS, err := to.Connect()
	if err != nil {
		return err
	}

	files, err := fromFS.Search(transfer.Pattern())
	if err != nil {
		return err
	}

	var totalBytes int64
	for _, file := range files {
		r, err := fromFS.Read(file)
		if err != nil {
			return err
		}

		err = toFS.Write(file, r)
		if err != nil {
			return err
		}

		totalBytes += file.Size
	}

	return nil
}

// Transfer Application Service - run a transfer from an application point of view (repo lookups, notifications, etc)
func (srvc *service) RunApplication(transferID uuid.UUID) error {
	t, err := srvc.transferRepo.Read(transferID)
	if err != nil {
		return err
	}

	from, err := srvc.locationRepo.Read(t.FromLocationID())
	if err != nil {
		return err
	}

	to, err := srvc.locationRepo.Read(t.ToLocationID())
	if err != nil {
		return err
	}

	return srvc.RunDomain(t, from, to)
}
