package service

import (
	"encoding/json"
	"time"

	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/fileserver/s3"
	"github.com/theandrew168/dripfile/internal/history"
	locationService "github.com/theandrew168/dripfile/internal/location/service"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Service struct {
	transferRepo    transfer.Repository
	historyRepo     history.Repository
	locationService *locationService.Service
}

func New(
	transferRepo transfer.Repository,
	historyRepo history.Repository,
	locationService *locationService.Service,
) *Service {
	svc := Service{
		transferRepo:    transferRepo,
		historyRepo:     historyRepo,
		locationService: locationService,
	}
	return &svc
}

func (svc *Service) Create(pattern, fromLocationID, toLocationID string) (transfer.Transfer, error) {
	t := transfer.New(pattern, fromLocationID, toLocationID)
	err := svc.transferRepo.Create(&t)
	if err != nil {
		return transfer.Transfer{}, err
	}

	return t, nil
}

func (svc *Service) List() ([]transfer.Transfer, error) {
	return svc.transferRepo.List()
}

func (svc *Service) Execute(transferID string) error {
	transfer, err := svc.transferRepo.Read(transferID)
	if err != nil {
		return err
	}

	fromLocation, err := svc.locationService.Read(transfer.FromLocationID)
	if err != nil {
		return err
	}

	var fromLocationInfo s3.Info
	err = json.Unmarshal(fromLocation.Info, &fromLocationInfo)
	if err != nil {
		return err
	}

	fromLocationFileServer, err := s3.New(fromLocationInfo)
	if err != nil {
		return err
	}

	toLocation, err := svc.locationService.Read(transfer.ToLocationID)
	if err != nil {
		return err
	}

	var toLocationInfo s3.Info
	err = json.Unmarshal(toLocation.Info, &toLocationInfo)
	if err != nil {
		return err
	}

	toLocationFileServer, err := s3.New(toLocationInfo)
	if err != nil {
		return err
	}

	start := time.Now().UTC()

	totalBytes, err := execute(transfer.Pattern, fromLocationFileServer, toLocationFileServer)
	if err != nil {
		return err
	}

	finish := time.Now().UTC()

	history := history.New(
		totalBytes,
		start,
		finish,
		transfer.ID,
	)
	err = svc.historyRepo.Create(&history)
	if err != nil {
		return err
	}

	return nil
}

// TODO: where should this live?
func execute(pattern string, from, to fileserver.FileServer) (int64, error) {
	files, err := from.Search(pattern)
	if err != nil {
		return 0, err
	}

	// TODO: write all to temps, rename if success, else rollback
	// TODO: perform these in parallel?
	var totalBytes int64
	for _, file := range files {
		r, err := from.Read(file)
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
