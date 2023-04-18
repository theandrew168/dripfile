package transfer

import (
	"encoding/json"
	"time"

	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/history"
	"github.com/theandrew168/dripfile/internal/location"
)

func Run(
	transferID string,
	transferRepo Repository,
	locationRepo location.Repository,
	historyRepo history.Repository,
) error {
	transfer, err := transferRepo.Read(transferID)
	if err != nil {
		return err
	}

	fromLocation, err := locationRepo.Read(transfer.FromLocationID)
	if err != nil {
		return err
	}

	var fromLocationInfo fileserver.S3Info
	err = json.Unmarshal(fromLocation.Info, &fromLocationInfo)
	if err != nil {
		return err
	}

	fromLocationFileServer, err := fileserver.NewS3(fromLocationInfo)
	if err != nil {
		return err
	}

	toLocation, err := locationRepo.Read(transfer.ToLocationID)
	if err != nil {
		return err
	}

	var toLocationInfo fileserver.S3Info
	err = json.Unmarshal(toLocation.Info, &toLocationInfo)
	if err != nil {
		return err
	}

	toLocationFileServer, err := fileserver.NewS3(toLocationInfo)
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
	err = historyRepo.Create(&history)
	if err != nil {
		return err
	}

	return nil
}
