package repository

import (
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/secret"
)

type Repository struct {
	Location  LocationRepository
	Itinerary ItineraryRepository
	Transfer  TransferRepository
}

func NewMemory() *Repository {
	repo := Repository{
		Location:  NewMemoryLocationRepository(),
		Itinerary: NewMemoryItineraryRepository(),
		Transfer:  NewMemoryTransferRepository(),
	}
	return &repo
}

func NewPostgres(conn database.Conn, box *secret.Box) *Repository {
	repo := Repository{
		Location:  NewPostgresLocationRepository(conn, box),
		Itinerary: NewPostgresItineraryRepository(conn),
		Transfer:  NewMemoryTransferRepository(),
	}
	return &repo
}
