package service

import (
	"github.com/theandrew168/dripfile/internal/location"
	locationStorage "github.com/theandrew168/dripfile/internal/location/storage"
)

type Service struct {
	locationStore locationStorage.Storage
}

func New(locationStore locationStorage.Storage) *Service {
	s := Service{
		locationStore: locationStore,
	}
	return &s
}

type GetByIDQuery struct {
	ID string
}

func (s *Service) GetByID(query GetByIDQuery) (*location.Location, error) {
	return s.locationStore.Read(query.ID)
}

type AddS3Command struct {
	ID string

	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

func (s *Service) AddS3(cmd AddS3Command) error {
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

	return s.locationStore.Create(l)
}
