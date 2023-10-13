package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/model"
	"github.com/theandrew168/dripfile/backend/repository"
)

type ItineraryService struct {
	repo *repository.Repository
}

func NewItineraryService(repo *repository.Repository) *ItineraryService {
	srvc := ItineraryService{
		repo: repo,
	}
	return &srvc
}

// TODO: ever both with this? Or just let the DB error via FK constraint?
func (srvc *ItineraryService) Create(itinerary model.Itinerary) error {
	_, err := srvc.repo.Location.Read(itinerary.FromLocationID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotExist):
			return repository.ErrConflict
		default:
			return err
		}
	}

	_, err = srvc.repo.Location.Read(itinerary.ToLocationID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotExist):
			return repository.ErrConflict
		default:
			return err
		}
	}

	return srvc.repo.Itinerary.Create(itinerary)
}

func (srvc *ItineraryService) List() ([]model.Itinerary, error) {
	return srvc.repo.Itinerary.List()
}

func (srvc *ItineraryService) Read(id uuid.UUID) (model.Itinerary, error) {
	return srvc.repo.Itinerary.Read(id)
}

func (srvc *ItineraryService) Delete(id uuid.UUID) error {
	return srvc.repo.Itinerary.Delete(id)
}
