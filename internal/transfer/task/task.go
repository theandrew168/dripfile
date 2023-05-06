package task

import (
	"encoding/json"
	"time"

	"github.com/theandrew168/dripfile/internal/fileserver/s3"
	"github.com/theandrew168/dripfile/internal/history"
	"github.com/theandrew168/dripfile/internal/location"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Task struct {
	transferID   string
	transferRepo transfer.Repository
	locationRepo location.Repository
	historyRepo  history.Repository
}

func New(
	transferID string,
	transferRepo transfer.Repository,
	locationRepo location.Repository,
	historyRepo history.Repository,
) *Task {
	t := Task{
		transferID:   transferID,
		transferRepo: transferRepo,
		locationRepo: locationRepo,
		historyRepo:  historyRepo,
	}
	return &t
}

func (t *Task) Run() error {
	transfer, err := t.transferRepo.Read(t.transferID)
	if err != nil {
		return err
	}

	fromLocation, err := t.locationRepo.Read(transfer.FromLocationID)
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

	toLocation, err := t.locationRepo.Read(transfer.ToLocationID)
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

	files, err := fromLocationFileServer.Search(transfer.Pattern)
	if err != nil {
		return err
	}

	start := time.Now().UTC()

	// TODO: write all to temps, rename if success, else rollback
	// TODO: perform these in parallel?
	var totalBytes int64
	for _, file := range files {
		r, err := fromLocationFileServer.Read(file)
		if err != nil {
			return err
		}

		err = toLocationFileServer.Write(file, r)
		if err != nil {
			return err
		}

		totalBytes += file.Size
	}

	finish := time.Now().UTC()

	history := history.New(
		totalBytes,
		start,
		finish,
		transfer.ID,
	)
	err = t.historyRepo.Create(&history)
	if err != nil {
		return err
	}

	return nil
}
