package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/model"
	"github.com/theandrew168/dripfile/backend/repository"
)

type LocationService struct {
	repo *repository.Repository
}

func NewLocationService(repo *repository.Repository) *LocationService {
	srvc := LocationService{
		repo: repo,
	}
	return &srvc
}

func (srvc *LocationService) Create(location model.Location) error {
	return srvc.repo.Location.Create(location)
}

func (srvc *LocationService) List() ([]model.Location, error) {
	return srvc.repo.Location.List()
}

func (srvc *LocationService) Read(id uuid.UUID) (model.Location, error) {
	return srvc.repo.Location.Read(id)
}

func (srvc *LocationService) Delete(id uuid.UUID) error {
	// TODO: ensure no itineraries / active xfers reference this location
	return srvc.repo.Location.Delete(id)
}
