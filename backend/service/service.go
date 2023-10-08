package service

import "github.com/theandrew168/dripfile/backend/repository"

type Service struct {
	Location  *LocationService
	Itinerary *ItineraryService
	Transfer  *TransferService
}

func New(repo *repository.Repository) *Service {
	srvc := Service{
		Location:  NewLocationService(repo),
		Itinerary: NewItineraryService(repo),
		Transfer:  NewTransferService(repo),
	}
	return &srvc
}
