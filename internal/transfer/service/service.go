package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/internal/location"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Service struct {
	locationStorage location.Storage
	transferStorage transfer.Storage
}

func New(locationStorage location.Storage, transferStorage transfer.Storage) *Service {
	srvc := Service{
		locationStorage: locationStorage,
		transferStorage: transferStorage,
	}
	return &srvc
}

func (srvc *Service) GetByID(query transfer.GetByIDQuery) (*transfer.Transfer, error) {
	_, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, transfer.ErrInvalidUUID
	}

	return srvc.transferStorage.Read(query.ID)
}

func (srvc *Service) GetAll(query transfer.GetAllQuery) ([]*transfer.Transfer, error) {
	return srvc.transferStorage.List()
}

func (srvc *Service) Add(cmd transfer.AddCommand) error {
	t, err := transfer.New(
		cmd.ID,
		cmd.Pattern,
		cmd.FromLocationID,
		cmd.ToLocationID,
	)
	if err != nil {
		return err
	}

	return srvc.transferStorage.Create(t)
}

func (srvc *Service) Remove(cmd transfer.RemoveCommand) error {
	_, err := uuid.Parse(cmd.ID)
	if err != nil {
		return transfer.ErrInvalidUUID
	}

	return srvc.transferStorage.Delete(cmd.ID)
}
