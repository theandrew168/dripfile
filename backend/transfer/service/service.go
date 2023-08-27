package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/history"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/transfer"
)

// service interface (other code depends on this)
type Service interface {
	RunDomain(transfer *transfer.Transfer, from, to *location.Location) (*history.History, error)
	RunApplication(transferID string) error
}

// service implementation
type service struct {
	locationRepo location.Repository
	transferRepo transfer.Repository
	historyRepo  history.Repository
}

func New(
	locationRepo location.Repository,
	transferRepo transfer.Repository,
	historyRepo history.Repository,
) Service {
	srvc := service{
		locationRepo: locationRepo,
		transferRepo: transferRepo,
		historyRepo:  historyRepo,
	}
	return &srvc
}

// Transfer Domain Service - run a transfer from a domain point of view
func (srvc *service) RunDomain(transfer *transfer.Transfer, from, to *location.Location) (*history.History, error) {
	start := time.Now().UTC()

	fromFS, err := from.Connect()
	if err != nil {
		return nil, err
	}

	toFS, err := to.Connect()
	if err != nil {
		return nil, err
	}

	files, err := fromFS.Search(transfer.Pattern())
	if err != nil {
		return nil, err
	}

	var totalBytes int64
	for _, file := range files {
		r, err := fromFS.Read(file)
		if err != nil {
			return nil, err
		}

		err = toFS.Write(file, r)
		if err != nil {
			return nil, err
		}

		totalBytes += file.Size
	}

	finish := time.Now().UTC()

	hID, _ := uuid.NewRandom()
	h, err := history.New(hID.String(), totalBytes, start, finish, transfer.ID())
	if err != nil {
		return nil, err
	}

	return h, nil
}

// Transfer Application Service - run a transfer from an application point of view (repo lookups, notifications, etc)
func (srvc *service) RunApplication(transferID string) error {
	_, err := uuid.Parse(transferID)
	if err != nil {
		return transfer.ErrInvalidUUID
	}

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

	h, err := srvc.RunDomain(t, from, to)
	if err != nil {
		return err
	}

	err = srvc.historyRepo.Create(h)
	if err != nil {
		return err
	}

	return nil
}
