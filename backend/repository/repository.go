package repository

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
