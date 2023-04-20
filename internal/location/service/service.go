package service

import (
	"encoding/json"

	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/location"
	locationRepo "github.com/theandrew168/dripfile/internal/location/repository"
)

type Service struct {
	locationRepo locationRepo.Repository
}

func New(locationRepo locationRepo.Repository) *Service {
	s := Service{
		locationRepo: locationRepo,
	}
	return &s
}

func (s *Service) CreateS3(name string, info fileserver.S3Info) (location.Location, error) {
	fs, err := fileserver.NewS3(info)
	if err != nil {
		return location.Location{}, nil
	}

	err = fs.Ping()
	if err != nil {
		return location.Location{}, nil
	}

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return location.Location{}, nil
	}

	m := location.New(location.KindS3, name, jsonInfo)
	err = s.locationRepo.Create(&m)
	if err != nil {
		return location.Location{}, nil
	}

	return m, nil
}
