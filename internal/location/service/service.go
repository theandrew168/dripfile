package service

import (
	"encoding/json"
	"errors"

	"github.com/theandrew168/dripfile/internal/fileserver/s3"
	"github.com/theandrew168/dripfile/internal/location"
	locationRepo "github.com/theandrew168/dripfile/internal/location/repository"
	"github.com/theandrew168/dripfile/internal/secret"
)

var (
	ErrLocationInUse = errors.New("location: in use")
)

type Service struct {
	box          *secret.Box
	locationRepo *locationRepo.Repository
}

func New(box *secret.Box, locationRepo *locationRepo.Repository) *Service {
	s := Service{
		box:          box,
		locationRepo: locationRepo,
	}
	return &s
}

func (s *Service) CreateS3(info s3.Info) (location.Location, error) {
	fs, err := s3.New(info)
	if err != nil {
		return location.Location{}, err
	}

	err = fs.Ping()
	if err != nil {
		return location.Location{}, err
	}

	data, err := json.Marshal(info)
	if err != nil {
		return location.Location{}, err
	}

	encrypted, err := s.box.Encrypt(data)
	if err != nil {
		return location.Location{}, err
	}

	m := location.New(location.KindS3, encrypted)
	err = s.locationRepo.Create(&m)
	if err != nil {
		return location.Location{}, err
	}

	return m, nil
}

func (s *Service) Read(id string) (location.Location, error) {
	// TODO: decrypt info
	return s.locationRepo.Read(id)
}

func (s *Service) List() ([]location.Location, error) {
	// TODO: decrypt each location's info
	return s.locationRepo.List()
}

func (s *Service) Update(location location.Location) error {
	// TODO: decrypt info
	// TODO: make updates
	// TODO: encrypt info
	return s.locationRepo.Update(location)
}

func (s *Service) Delete(id string) error {
	// TODO: check for transfers that use this location
	return s.locationRepo.Delete(id)
}
