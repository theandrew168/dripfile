package service

import (
	"encoding/json"
	"errors"

	"github.com/theandrew168/dripfile/internal/fileserver/s3"
	"github.com/theandrew168/dripfile/internal/location"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/transfer"
)

var (
	ErrLocationInUse = errors.New("location: in use")
)

type Service struct {
	box          *secret.Box
	locationRepo location.Repository
	transferRepo transfer.Repository
}

func New(
	box *secret.Box,
	locationRepo location.Repository,
	transferRepo transfer.Repository,
) *Service {
	svc := Service{
		box:          box,
		locationRepo: locationRepo,
		transferRepo: transferRepo,
	}
	return &svc
}

func (svc *Service) CreateS3(info s3.Info) (location.Location, error) {
	fs, err := s3.New(info)
	if err != nil {
		return location.Location{}, err
	}

	err = fs.Ping()
	if err != nil {
		return location.Location{}, err
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		return location.Location{}, err
	}

	encryptedInfoJSON, err := svc.box.Encrypt(infoJSON)
	if err != nil {
		return location.Location{}, err
	}

	l := location.New(location.KindS3, encryptedInfoJSON)
	err = svc.locationRepo.Create(&l)
	if err != nil {
		return location.Location{}, err
	}

	return l, nil
}

func (svc *Service) Read(id string) (location.Location, error) {
	l, err := svc.locationRepo.Read(id)
	if err != nil {
		return location.Location{}, err
	}

	decryptedInfo, err := svc.box.Decrypt(l.Info)
	if err != nil {
		return location.Location{}, err
	}

	l.Info = decryptedInfo
	return l, nil
}

func (svc *Service) List() ([]location.Location, error) {
	encryptedLocations, err := svc.locationRepo.List()
	if err != nil {
		return nil, err
	}

	var ls []location.Location
	for _, l := range encryptedLocations {
		decryptedInfo, err := svc.box.Decrypt(l.Info)
		if err != nil {
			return nil, err
		}

		l.Info = decryptedInfo
		ls = append(ls, l)
	}

	return ls, nil
}

func (svc *Service) Update(location location.Location) error {
	encryptedInfo, err := svc.box.Encrypt(location.Info)
	if err != nil {
		return err
	}

	location.Info = encryptedInfo
	return svc.locationRepo.Update(location)
}

func (svc *Service) Delete(id string) error {
	transfers, err := svc.transferRepo.ListByLocationID(id)
	if err != nil {
		return err
	}

	if len(transfers) > 0 {
		return ErrLocationInUse
	}

	return svc.locationRepo.Delete(id)
}
