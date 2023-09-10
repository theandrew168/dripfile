package transfer

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/itinerary"
	"github.com/theandrew168/dripfile/backend/location"
)

type DomainService struct{}

func NewDomainService() *DomainService {
	srvc := DomainService{}
	return &srvc
}

// Perform the transfer on provided domain objects.
func (srvc *DomainService) Run(i *itinerary.Itinerary, from, to *location.Location) error {
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

type AppService struct {
	domainService *DomainService
	locationRepo  location.Repository
	itineraryRepo itinerary.Repository
}

func NewAppService(locationRepo location.Repository, itineraryRepo itinerary.Repository) *AppService {
	domainService := NewDomainService()

	srvc := AppService{
		domainService: domainService,
		locationRepo:  locationRepo,
		itineraryRepo: itineraryRepo,
	}
	return &srvc
}

// Lookup the domain objects from persistent storage and perform the transfer.
func (srvc *AppService) RunApp(itineraryID uuid.UUID) error {
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

	return srvc.domainService.Run(i, from, to)
}
