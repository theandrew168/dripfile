package query

import (
	"github.com/theandrew168/dripfile/internal/location"
	locationStorage "github.com/theandrew168/dripfile/internal/location/storage"
)

type Read struct {
	ID string
}

type ReadHandler struct {
	locationStorage locationStorage.Storage
}

func NewReadHandler(locationStore locationStorage.Storage) *ReadHandler {
	h := ReadHandler{
		locationStorage: locationStore,
	}
	return &h
}

func (h *ReadHandler) Handle(query Read) (*location.Location, error) {
	return h.locationStorage.Read(query.ID)
}
