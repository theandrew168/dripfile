package storage

import (
	"github.com/theandrew168/dripfile/internal/location"
)

type Storage interface {
	Create(l *location.Location) error
	Read(id string) (*location.Location, error)
	List() ([]*location.Location, error)
	Delete(id string) error
}
