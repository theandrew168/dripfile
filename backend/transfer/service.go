package transfer

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/itinerary"
	"github.com/theandrew168/dripfile/backend/location"
)

type Service struct {
	locationRepo  location.Repository
	itineraryRepo itinerary.Repository
}

func NewService(locationRepo location.Repository, itineraryRepo itinerary.Repository) *Service {
	srvc := Service{
		locationRepo:  locationRepo,
		itineraryRepo: itineraryRepo,
	}
	return &srvc
}

// Perform the transfer given domain objects.
func (srvc *Service) RunDomain(i *itinerary.Itinerary, from, to *location.Location) error {
	fromFS, err := from.Connect()
	if err != nil {
		return err
	}

	toFS, err := to.Connect()
	if err != nil {
		return err
	}

	_, err = Transfer(i.Pattern(), fromFS, toFS)
	if err != nil {
		return err
	}

	return nil
}

// Lookup the domain objects from persistent storage and perform the transfer.
func (srvc *Service) RunApp(itineraryID uuid.UUID) error {
	i, err := srvc.itineraryRepo.Read(itineraryID)
	if err != nil {
		return err
	}

	from, err := srvc.locationRepo.Read(i.FromLocationID())
	if err != nil {
		return err
	}

	to, err := srvc.locationRepo.Read(i.ToLocationID())
	if err != nil {
		return err
	}

	return srvc.RunDomain(i, from, to)
}
