package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/internal/location"
)

type Service struct {
	locationStorage location.Storage
}

func New(locationStorage location.Storage) *Service {
	srvc := Service{
		locationStorage: locationStorage,
	}
	return &srvc
}

func (srvc *Service) GetByID(query location.GetByIDQuery) (*location.Location, error) {
	_, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, location.ErrInvalidUUID
	}

	return srvc.locationStorage.Read(query.ID)
}

func (srvc *Service) GetAll(query location.GetAllQuery) ([]*location.Location, error) {
	return srvc.locationStorage.List()
}

func (srvc *Service) AddMemory(cmd location.AddMemoryCommand) error {
	l, err := location.NewMemory(
		cmd.ID,
	)
	if err != nil {
		return err
	}

	return srvc.locationStorage.Create(l)
}

func (srvc *Service) AddS3(cmd location.AddS3Command) error {
	l, err := location.NewS3(
		cmd.ID,
		cmd.Endpoint,
		cmd.Bucket,
		cmd.AccessKeyID,
		cmd.SecretAccessKey,
	)
	if err != nil {
		return err
	}

	return srvc.locationStorage.Create(l)
}

func (srvc *Service) Remove(cmd location.RemoveCommand) error {
	_, err := uuid.Parse(cmd.ID)
	if err != nil {
		return location.ErrInvalidUUID
	}

	return srvc.locationStorage.Delete(cmd.ID)
}
