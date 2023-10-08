package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/repository"
)

type TransferService struct {
	repo *repository.Repository
}

func NewTransferService(repo *repository.Repository) *TransferService {
	srvc := TransferService{
		repo: repo,
	}
	return &srvc
}

func (srvc *TransferService) Run(itineraryID uuid.UUID) error {
	i, err := srvc.repo.Itinerary.Read(itineraryID)
	if err != nil {
		return err
	}

	from, err := srvc.repo.Location.Read(i.FromLocationID)
	if err != nil {
		return err
	}

	to, err := srvc.repo.Location.Read(i.ToLocationID)
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

	_, err = fileserver.Transfer(i.Pattern, fromFS, toFS)
	if err != nil {
		return err
	}

	return nil
}
