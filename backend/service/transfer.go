package service

import (
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/model"
	"github.com/theandrew168/dripfile/backend/repository"
)

type TransferService struct {
	repo *repository.Repository
}

func NewTransferService(repo *repository.Repository) *TransferService {
	srvc := TransferService{
		repo: repo,
	}
	return &srvc
}

func (srvc *TransferService) Run(itineraryID uuid.UUID) error {
	i, err := srvc.repo.Itinerary.Read(itineraryID)
	if err != nil {
		return err
	}

	from, err := srvc.repo.Location.Read(i.FromLocationID)
	if err != nil {
		return err
	}

	to, err := srvc.repo.Location.Read(i.ToLocationID)
	if err != nil {
		return err
	}

	return Run(i, from, to)
}

func Run(i model.Itinerary, from, to model.Location) error {
	fromFS, err := from.Connect()
	if err != nil {
		return err
	}

	toFS, err := to.Connect()
	if err != nil {
		return err
	}

	_, err = Transfer(i.Pattern, fromFS, toFS)
	if err != nil {
		return err
	}

	return nil
}

// Transfer all files matching a given pattern from one FileServer to another.
// Returns the total number of bytes transferred or an error.
func Transfer(pattern string, from, to fileserver.FileServer) (int, error) {
	files, err := from.Search(pattern)
	if err != nil {
		return 0, err
	}

	// TODO: spawn a goro and return a progress channel

	var totalBytes int
	for _, file := range files {
		r, err := from.Read(file.Name)
		if err != nil {
			return 0, err
		}

		err = to.Write(file, r)
		if err != nil {
			return 0, err
		}

		totalBytes += file.Size
	}

	return totalBytes, nil
}
