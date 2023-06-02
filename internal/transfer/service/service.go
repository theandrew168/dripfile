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

func (srvc *Service) Run(cmd transfer.RunCommand) error {
	_, err := uuid.Parse(cmd.ID)
	if err != nil {
		return transfer.ErrInvalidUUID
	}

	t, err := srvc.transferStorage.Read(cmd.ID)
	if err != nil {
		return err
	}

	from, err := srvc.locationStorage.Read(t.FromLocationID())
	if err != nil {
		return err
	}

	to, err := srvc.locationStorage.Read(t.ToLocationID())
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

	return nil
}
