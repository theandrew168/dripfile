package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/history"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/transfer"
)

type Service struct {
	locationRepo    location.Repository
	transferStorage transfer.Storage
	historyStorage  history.Storage
}

func New(
	locationRepo location.Repository,
	transferStorage transfer.Storage,
	historyStorage history.Storage,
) *Service {
	srvc := Service{
		locationRepo:    locationRepo,
		transferStorage: transferStorage,
		historyStorage:  historyStorage,
	}
	return &srvc
}

func (srvc *Service) Run(transferID string) error {
	_, err := uuid.Parse(transferID)
	if err != nil {
		return transfer.ErrInvalidUUID
	}

	start := time.Now().UTC()

	t, err := srvc.transferStorage.Read(transferID)
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

	fromFS, err := from.Connect()
	if err != nil {
		return err
	}

	toFS, err := to.Connect()
	if err != nil {
		return err
	}

	files, err := fromFS.Search(t.Pattern())
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

	finish := time.Now().UTC()

	hID, _ := uuid.NewRandom()
	h, err := history.New(hID.String(), totalBytes, start, finish, t.ID())
	if err != nil {
		return err
	}

	err = srvc.historyStorage.Create(h)
	if err != nil {
		return err
	}

	return nil
}
