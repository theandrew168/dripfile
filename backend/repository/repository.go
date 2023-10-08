package repository

type Repository struct {
	Location  LocationRepository
	Itinerary ItineraryRepository
}

func NewMemory() *Repository {
	repo := Repository{
		Location:  NewMemoryLocationRepository(),
		Itinerary: NewMemoryItineraryRepository(),
	}
	return &repo
}
