package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/model"
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

func (srvc *TransferService) Create(transfer model.Transfer) error {
	return srvc.repo.Transfer.Create(transfer)
}

func (srvc *TransferService) List() ([]model.Transfer, error) {
	return srvc.repo.Transfer.List()
}

func (srvc *TransferService) Read(id uuid.UUID) (model.Transfer, error) {
	return srvc.repo.Transfer.Read(id)
}

func (srvc *TransferService) Delete(id uuid.UUID) error {
	return srvc.repo.Transfer.Delete(id)
}

func (srvc *TransferService) Run(id uuid.UUID) error {
	transfer, err := srvc.repo.Transfer.Read(id)
	if err != nil {
		return err
	}

	itinerary, err := srvc.repo.Itinerary.Read(transfer.ItineraryID)
	if err != nil {
		return err
	}

	fromLocation, err := srvc.repo.Location.Read(itinerary.FromLocationID)
	if err != nil {
		return err
	}

	toLocation, err := srvc.repo.Location.Read(itinerary.ToLocationID)
	if err != nil {
		return err
	}

	fromFS, err := fromLocation.Connect()
	if err != nil {
		return err
	}

	toFS, err := toLocation.Connect()
	if err != nil {
		return err
	}

	// TODO: read status from xfer channel and update DB periodically
	_, err = fileserver.Transfer(itinerary.Pattern, fromFS, toFS)
	if err != nil {
		return err
	}

	return nil
}
