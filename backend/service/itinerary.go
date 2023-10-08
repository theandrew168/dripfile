package service

import (
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

func (srvc *ItineraryService) Create(itinerary model.Itinerary) error {
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
